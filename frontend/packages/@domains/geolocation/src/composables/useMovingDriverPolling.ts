import { $movingDrivers, $lastFetchTime } from "../stores/movingDriversStore.ts"
import type { MovingDriver } from "../types/geolocationTypes.ts"
import type { MovingPath } from "maps/composables/useMovingIconLayer.ts"
import { geolocationClient } from "../api/client.ts"

const POLL_INTERVAL = 5000
const MOCK_ADVANCE_INTERVAL = 5000

export interface MovingDriverProvider {
	getClosestMovingDrivers(lat: number, lon: number): Promise<MovingDriver[]>
	getMovingDriversByIDs(ids: number[]): Promise<MovingDriver[]>
}

class HttpMovingDriverProvider implements MovingDriverProvider {
	async getClosestMovingDrivers(lat: number, lon: number): Promise<MovingDriver[]> {
		const { data, error } = await geolocationClient.GET("/drivers/moving/closest", {
			params: { query: { lat, lon, radius_meters: 30_000 } },
		})
		if (error) throw new Error(`getClosestMovingDrivers failed: ${error}`)
		return data ?? []
	}

	async getMovingDriversByIDs(ids: number[]): Promise<MovingDriver[]> {
		const { data, error } = await geolocationClient.POST("/drivers/moving", {
			body: { ids },
		})
		if (error) throw new Error(`getMovingDriversByIDs failed: ${error}`)
		return data ?? []
	}
}

class MockMovingDriverProvider implements MovingDriverProvider {
	private ringBuffer: MovingDriver[][]
	private ringBufferIdx: number
	private lastRequestTime: number

	constructor(snapshots: MovingDriver[][]) {
		this.ringBuffer = snapshots
		this.ringBufferIdx = 0
		this.lastRequestTime = Date.now()
	}

	private maybeAdvance() {
		const now = Date.now()
		if (now - this.lastRequestTime >= MOCK_ADVANCE_INTERVAL) {
			this.ringBufferIdx = (this.ringBufferIdx + 1) % this.ringBuffer.length
			this.lastRequestTime = now
		}
	}

	async getClosestMovingDrivers(_lat: number, _lon: number): Promise<MovingDriver[]> {
		this.maybeAdvance()
		return this.ringBuffer[this.ringBufferIdx]
	}

	async getMovingDriversByIDs(ids: number[]): Promise<MovingDriver[]> {
		this.maybeAdvance()
		return this.ringBuffer[this.ringBufferIdx].filter((d) => ids.includes(d.driver_id))
	}
}

export const MovingDriverProviders = {
	http: () => new HttpMovingDriverProvider(),
	mock: (snapshots: MovingDriver[][]) => new MockMovingDriverProvider(snapshots),
}

function buildMovingDrivers(drivers: MovingDriver[]): MovingPath[] {
	const output: MovingPath[] = []
	for (const resp of drivers) {
		if (!resp || !resp.points) continue
		output.push({
			id: resp.driver_id,
			coordinates: resp.points as [number, number][],
			time: resp.travel_time,
			distance: resp.path_meters,
		})
	}
	return output
}

export function useMovingDriverPolling(
	provider: MovingDriverProvider,
	initialLat: number,
	initialLon: number,
) {
	let currentIds: number[] = []
	let pollTimer: ReturnType<typeof setInterval> | null = null

	async function poll() {
		try {
			let resp: MovingDriver[]
			if (currentIds.length === 0) {
				resp = await provider.getClosestMovingDrivers(initialLat, initialLon)
				currentIds = resp.map((d) => d.driver_id)
			} else {
				resp = await provider.getMovingDriversByIDs(currentIds)
			}

			$movingDrivers.set(buildMovingDrivers(resp))
			$lastFetchTime.set(Date.now())
		} catch (err) {
			console.error("[moving-drivers] poll error:", err)
		}
	}

	function start() {
		poll()
		pollTimer = setInterval(poll, POLL_INTERVAL)
	}

	function stop() {
		if (pollTimer !== null) {
			clearInterval(pollTimer)
			pollTimer = null
		}
	}

	return { start, stop }
}
