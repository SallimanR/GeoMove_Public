import axios from "axios"

export interface RouteResponse {
	paths: Array<{
		distance: number
		points: { coordinates: [number, number][] }
	}>;
}

const API_BASE = "http://localhost:8989"

export async function fetchRoute(start: { lat: number, lon: number }, end: { lat: number, lon: number }): Promise<RouteResponse> {
	const result = (await axios.get<RouteResponse>(`${API_BASE}/route?point=${[start.lat, start.lon]}&point=${[end.lat, end.lon]}&profile=car&points_encoded=false`)).data
	return result
}

