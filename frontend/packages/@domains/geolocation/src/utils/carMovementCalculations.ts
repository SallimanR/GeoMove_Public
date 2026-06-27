import { rhumbDistance } from '@turf/turf';
import { type RealtimeDriver } from "../types/realtimeDriver.ts";
import { haversineDistance } from 'geo';

export function calculateSegmentsPercentages(segments: [number, number][], total_distance: number): number[] {
	const segmentsPercentages = new Array<number>(segments.length - 2)
	for (let i = 0; i < segments.length - 1; i++) {
		segmentsPercentages[i] = (rhumbDistance(segments[i], segments[i + 1], { units: "meters" }) / total_distance)
	}
	return segmentsPercentages
}

export function getDriverPosition(start: [number, number], end: number[], progress: number): [number, number] {
	// const distance = rhumbDistance(start, end, { units: "meters" })
	return [
		start[0] + (end[0] - start[0]) * progress,
		start[1] + (end[1] - start[1]) * progress
	]
}

export function getDriverBearing(start: [number, number], end: [number, number]): number {
	const dx = end[0] - start[0];
	const dy = end[1] - start[1];
	return Math.atan2(dy, dx) * 180 / Math.PI;
}

/**
 * Returns the interpolated position (lng, lat) on a driver's path at a given elapsed time.
 * @param driver The driver object containing coordinates, total time (ms), and total distance (m).
 * @param elapsedSeconds Time elapsed since start, in seconds.
 * @returns A [lng, lat] coordinate pair.
 */
export function getPositionAtTime(
	driver: RealtimeDriver,
	elapsedSeconds: number
): [number, number] {
	const coords = driver.coordinates;
	const totalTimeMs = driver.time;
	const totalDistanceM = driver.distance;

	// Handle edge cases
	if (coords.length === 0) {
		throw new Error("Driver has no coordinates");
	}
	if (elapsedSeconds <= 0) {
		return coords[0];
	}
	const totalTimeSeconds = totalTimeMs / 1000;
	if (elapsedSeconds >= totalTimeSeconds) {
		return coords[coords.length - 1];
	}

	// Distance that should have been covered by elapsedSeconds
	const targetDistance = totalDistanceM * (elapsedSeconds / totalTimeSeconds);

	// Walk along the polyline to find the containing segment
	let cumulativeDistance = 0;
	for (let i = 0; i < coords.length - 1; i++) {
		const p1 = coords[i];
		const p2 = coords[i + 1];
		const segmentDist = haversineDistance(p1, p2);

		if (cumulativeDistance + segmentDist >= targetDistance) {
			// Interpolate linearly between p1 and p2
			const remaining = targetDistance - cumulativeDistance;
			const ratio = remaining / segmentDist;
			const lng = p1[0] + ratio * (p2[0] - p1[0]);
			const lat = p1[1] + ratio * (p2[1] - p1[1]);
			return [lng, lat];
		}
		cumulativeDistance += segmentDist;
	}

	// Fallback (should never be reached if targetDistance <= totalDistance)
	return coords[coords.length - 1];
}
