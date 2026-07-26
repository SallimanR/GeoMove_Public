import { haversineDistance } from "../utils/geometry";
import { fetchRoute } from "./routing";

const CITY_AVG_SPEED_M_PER_MIN = 667;

export interface Point {
	lat: number;
	lon: number;
}

export async function shortestTimeToArrive(
	from: Point,
	candidates: Point[],
): Promise<number | null> {
	if (candidates.length === 0) return null;

	let closest = candidates[0];
	let closestDist = haversineDistance([closest.lon, closest.lat], [from.lon, from.lat]);
	for (const c of candidates) {
		const dist = haversineDistance([c.lon, c.lat], [from.lon, from.lat]);
		if (dist < closestDist) {
			closest = c;
			closestDist = dist;
		}
	}

	try {
		const route = await fetchRoute(closest, from);
		if (route.paths?.[0]?.distance) {
			return Math.round(route.paths[0].distance / CITY_AVG_SPEED_M_PER_MIN);
		}
	} catch {
		// fall through to straight-line estimate
	}

	return Math.round(closestDist / CITY_AVG_SPEED_M_PER_MIN);
}
