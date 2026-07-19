import axios from "axios"
import { getGeoConfig } from "../config.ts"

export interface RouteResponse {
	paths: Array<{
		distance: number
		points: { coordinates: [number, number][] }
	}>;
}

export async function fetchRoute(start: { lat: number, lon: number }, end: { lat: number, lon: number }): Promise<RouteResponse> {
	const apiBase = getGeoConfig().routingApi
	const result = (await axios.get<RouteResponse>(`${apiBase}/route?point=${[start.lat, start.lon]}&point=${[end.lat, end.lon]}&profile=car&points_encoded=false`)).data
	return result
}
