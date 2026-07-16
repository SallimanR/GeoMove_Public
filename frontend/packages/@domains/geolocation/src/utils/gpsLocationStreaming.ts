import { ref, type Ref } from 'vue'
import { GPSUpdate } from '../api/proto/gps_update.ts'
import { Request, Channel, Response } from '../api/proto/ws_message.ts'

export interface GPSPosition {
	lat: number
	lng: number
	accuracy: number
	timestamp: number
}

export interface GpsProvider {
	start(onPosition: (pos: GPSPosition) => void, onError: (err: string) => void): void
	stop(): void
}

export interface GpsStreamState {
	isStreaming: Ref<boolean>
	currentPosition: Ref<GPSPosition | null>
	error: Ref<string | null>
}

class BrowserGpsProvider implements GpsProvider {
	private timer: ReturnType<typeof setInterval> | null = null
	private readonly POLL_INTERVAL = 500

	start(onPosition: (pos: GPSPosition) => void, onError: (err: string) => void): void {
		const poll = () => {
			navigator.geolocation.getCurrentPosition(
				(pos) => {
					onPosition({
						lat: pos.coords.latitude,
						lng: pos.coords.longitude,
						accuracy: pos.coords.accuracy,
						timestamp: Date.now(),
					})
				},
				(err) => onError(err.message),
				{ enableHighAccuracy: true, timeout: 10000, maximumAge: 0 },
			)
		}
		poll()
		this.timer = setInterval(poll, this.POLL_INTERVAL)
	}

	stop(): void {
		if (this.timer !== null) {
			clearInterval(this.timer)
			this.timer = null
		}
	}
}

class CapacitorGpsProvider implements GpsProvider {
	private watchId: string | null = null

	async start(onPosition: (pos: GPSPosition) => void, onError: (err: string) => void): Promise<void> {
		const { Geolocation } = await import('@capacitor/geolocation')
		try {
			const perm = await Geolocation.checkPermissions()
			if (perm.location !== 'granted') {
				const req = await Geolocation.requestPermissions()
				if (req.location !== 'granted') {
					onError('Доступ к геолокации отклонён')
					return
				}
			}
		} catch {
			onError('Ошибка запроса доступа к геолокации')
			return
		}

		this.watchId = await Geolocation.watchPosition(
			{ enableHighAccuracy: true, timeout: 10000 },
			(pos, err) => {
				if (err) { onError(err.message); return }
				if (pos) {
					onPosition({
						lat: pos.coords.latitude,
						lng: pos.coords.longitude,
						accuracy: pos.coords.accuracy,
						timestamp: Date.now(),
					})
				}
			},
		)
	}

	stop(): void {
		if (this.watchId !== null) {
			import('@capacitor/geolocation').then(({ Geolocation }) => {
				Geolocation.clearWatch({ id: this.watchId! })
				this.watchId = null
			})
		}
	}
}

class RingBufferGpsProvider implements GpsProvider {
	private timer: ReturnType<typeof setInterval> | null = null
	private index = 0
	private readonly POLL_INTERVAL = 500

	constructor(private readonly points: GPSPosition[]) { }

	start(onPosition: (pos: GPSPosition) => void, _onError: (err: string) => void): void {
		this.timer = setInterval(() => {
			const point = this.points[this.index]
			onPosition({
				lat: point.lat,
				lng: point.lng,
				accuracy: point.accuracy,
				timestamp: Date.now(),
			})
			this.index = (this.index + 1) % this.points.length
		}, this.POLL_INTERVAL)
	}

	stop(): void {
		if (this.timer !== null) {
			clearInterval(this.timer)
			this.timer = null
		}
	}
}

export const GpsProviders = {
	browser: () => new BrowserGpsProvider(),
	capacitor: () => new CapacitorGpsProvider(),
	ringBuffer: (points: GPSPosition[]) => new RingBufferGpsProvider(points),
}

export function streamGPS(wsUrl: string, provider: GpsProvider): GpsStreamState & { start: () => void; stop: () => void } {
	const FLUSH_INTERVAL = 5000
	const isStreaming = ref(false)
	const currentPosition = ref<GPSPosition | null>(null)
	const error = ref<string | null>(null)

	let socket: WebSocket | null = null
	let flushTimer: ReturnType<typeof setInterval> | null = null
	let updateSeq = 0

	const coordBuffer: [number, number][] = []
	const tsBuffer: number[] = []

	function cleanup() {
		provider.stop()
		if (flushTimer !== null) {
			clearInterval(flushTimer)
			flushTimer = null
		}
		socket?.close()
		socket = null
		isStreaming.value = false
		coordBuffer.length = 0
		tsBuffer.length = 0
	}

	function connectWebSocket(): WebSocket {
		const ws = new WebSocket(wsUrl)
		ws.binaryType = 'arraybuffer'
		ws.onopen = () => { error.value = null }
		ws.onerror = () => {
			cleanup()
			error.value = 'Ошибка соединения с сервером'
			console.error('[gps-stream] WebSocket onerror')
		}
		ws.onmessage = (event: MessageEvent) => {
			try {
				const resp = Response.decode(new Uint8Array(event.data))
				if (resp.response && resp.response.statusCode !== 200) {
					const msg = resp.response.errorMessage || `Ошибка сервера: ${resp.response.statusCode}`
					console.error('[gps-stream] server error:', msg)
					error.value = msg
				}
			} catch { /* non-response messages (broadcasts) are expected */ }
		}
		ws.onclose = (event: CloseEvent) => {
			if (isStreaming.value) {
				cleanup()
				error.value = `Соединение с сервером потеряно (код ${event.code})`
				console.error('[gps-stream] WebSocket closed', event.code)
			}
		}
		return ws
	}

	function flushBuffer() {
		if (coordBuffer.length === 0 || !socket || socket.readyState !== WebSocket.OPEN) return

		const coordinates = coordBuffer.map(([lat, lng]) => ({ latitude: lat, longitude: lng }))
		const timestamps = tsBuffer.slice()
		console.log('[gps-stream] sending', { coords: coordinates.length, ts: timestamps[0], sample: coordinates[0] })

		coordBuffer.length = 0
		tsBuffer.length = 0

		const gpsUpdate = GPSUpdate.encode({ coordinates, timestamps }).finish()
		const msg = Request.encode({
			channel: Channel.GPS_REALTIME,
			requestId: String(updateSeq++),
			publish: { data: gpsUpdate },
		}).finish()

		socket.send(msg)
	}

	function start() {
		error.value = null
		socket = connectWebSocket()
		flushTimer = setInterval(flushBuffer, FLUSH_INTERVAL)

		provider.start(
			(pos) => {
				console.log('[gps-stream] push', { now: pos.timestamp, epochSec: Math.floor(pos.timestamp / 1000), lat: pos.lat.toFixed(5), lng: pos.lng.toFixed(5) })
				currentPosition.value = pos
				coordBuffer.push([pos.lat, pos.lng])
				tsBuffer.push(Math.floor(pos.timestamp / 1000))
			},
			(err) => {
				cleanup()
				error.value = `Ошибка геолокации: ${err}`
			},
		)

		isStreaming.value = true
	}

	function stop() {
		cleanup()
		error.value = null
	}

	return { isStreaming, currentPosition, error, start, stop }
}
