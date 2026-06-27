import { ref, onMounted, onUnmounted } from 'vue'
import type { RealtimeDriver } from '../types/realtimeDriver.ts'

type RealtimeUpdate = {
	drivers: RealtimeDriver[]
}

export function useRealtimeDrivers() {
	const drivers = ref<RealtimeDriver[]>([])
	const socket = ref<WebSocket | null>(null)

	const connect = () => {
		socket.value = new WebSocket('wss://localhost:8080/v1/gps/')
		socket.value.onmessage = (event: MessageEvent) => {
			try {
				const update: RealtimeUpdate = JSON.parse(event.data)
				drivers.value = update.drivers
			} catch (err) {
				console.error('Failed to parse WebSocket message:', err)
			}
		}
	}

	const disconnect = () => {
		socket.value?.close()
		socket.value = null
	}

	onMounted(() => connect())
	onUnmounted(() => disconnect())

	return { drivers }
}
