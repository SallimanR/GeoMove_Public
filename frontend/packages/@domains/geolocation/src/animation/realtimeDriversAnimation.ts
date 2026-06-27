import { type Map as MaplibreMap } from 'maplibre-gl';
import { type MapboxOverlay } from '@deck.gl/mapbox';
import { $realtimeDrivers, $lastFetchTime } from '../stores/realtimeDriversStore.ts';
import {
	// calculateSegmentsPercentages,
	// getDriverPosition,
	getDriverBearing,
} from '../utils/carMovementCalculations.ts';
import { newRealtimeDriversLayer } from '../maps_display/realtimeDriversLayer.ts';
import type {
	RealtimeDriver,
	RealtimeDriverPosition
} from '../types/realtimeDriver.ts';
import { FetchedRealtimeDrivers } from './driver_data.ts';
import { haversineDistance } from 'geo';

interface DriverAnimationInfo {
	cumulativeDistances: number[];   // distance from start to each vertex (meters)
	totalDistance: number;
	totalTimeSeconds: number;
}

const ICON_BEARING_OFFSET = 90;

export class RealtimeDriversAnimator {
	private map: MaplibreMap;
	private deck: MapboxOverlay;
	private animationFrame: number | null = null;
	private removeSubscriptionListener: (() => void) | null = null;
	private cachePool = new Map<number, DriverAnimationInfo>();
	private drivers: readonly RealtimeDriver[] = [];
	private lastFetchTime = 0;

	constructor(map: MaplibreMap, deck: MapboxOverlay) {
		this.map = map;
		this.deck = deck;
	}

	start() {
		$realtimeDrivers.set(FetchedRealtimeDrivers);
		$lastFetchTime.set(Date.now());

		// Subscribe to store changes – rebuild cache when drivers are updated
		this.removeSubscriptionListener = $realtimeDrivers.subscribe((drivers) => {
			this.drivers = drivers;
			this.lastFetchTime = $lastFetchTime.get();
			this.rebuildCachePool();
		});

		this.animate();
	}

	stop() {
		$realtimeDrivers.set([]);
		if (this.animationFrame) {
			cancelAnimationFrame(this.animationFrame);
			this.animationFrame = null;
		}
		if (this.removeSubscriptionListener) {
			this.removeSubscriptionListener();
			this.removeSubscriptionListener = null;
		}
		if (this.deck) {
			this.deck.setProps({ layers: [] });
		}
	}

	private rebuildCachePool() {
		this.cachePool.clear();
		for (const driver of this.drivers) {
			const coords = driver.coordinates;
			const numCoords = coords.length;
			if (numCoords < 2) continue; // not enough points to animate

			const cumulativeDistances = new Array<number>(numCoords);
			cumulativeDistances[0] = 0;
			for (let i = 0; i < numCoords - 1; i++) {
				const dist = haversineDistance(coords[i], coords[i + 1]);
				cumulativeDistances[i + 1] = cumulativeDistances[i] + dist
			}

			this.cachePool.set(driver.id, {
				cumulativeDistances: cumulativeDistances,
				totalDistance: cumulativeDistances[numCoords - 1],
				totalTimeSeconds: driver.time / 1000,
			});
		}
	}

	private animate = () => {
		const now = Date.now();
		const elapsedSeconds = (now - this.lastFetchTime) / 1000;

		const updates = Array<RealtimeDriverPosition>(this.drivers.length - 1);

		let count = 0
		for (const driver of this.drivers) {
			const cache = this.cachePool.get(driver.id);
			if (!cache) continue;

			const totalTime = cache.totalTimeSeconds;
			// Loop the animation by taking elapsed modulo total time.
			// If totalTime is zero, keep the driver at its start.
			let t = totalTime > 0 ? elapsedSeconds % totalTime : 0;

			// Handle edge cases quickly.
			if (t <= 0) {
				updates[count++] = {
					id: driver.id,
					position: driver.coordinates[0],
					bearing: 0, // or compute initial bearing from first segment
				};
				continue;
			}
			if (t >= totalTime) {
				const lastIdx = driver.coordinates.length - 1;
				updates[count++] = {
					id: driver.id,
					position: driver.coordinates[lastIdx],
					bearing: 0, // or final bearing
				};
				continue;
			}

			// Distance that should have been covered according to the provided total distance.
			const targetDist = (t / totalTime) * driver.distance;
			// Clamp to the actual computed path length (avoid overshoot due to rounding).
			const clampedDist = Math.min(targetDist, cache.totalDistance);
			const cum = cache.cumulativeDistances;

			// Binary search to find the segment containing clampedDist.
			let low = 0;
			let high = cum.length - 1;
			while (low < high) {
				const mid = Math.floor((low + high) / 2);
				if (cum[mid] < clampedDist) {
					low = mid + 1;
				} else {
					high = mid;
				}
			}
			const segIndex = low - 1; // segment between segIndex and segIndex+1
			const startPoint = driver.coordinates[segIndex];
			const endPoint = driver.coordinates[segIndex + 1];
			const segDist = cum[segIndex + 1] - cum[segIndex];
			const ratio = (clampedDist - cum[segIndex]) / segDist;

			// Linear interpolation
			const lng = startPoint[0] + ratio * (endPoint[0] - startPoint[0]);
			const lat = startPoint[1] + ratio * (endPoint[1] - startPoint[1]);

			const bearing = getDriverBearing(
				startPoint,
				endPoint
			);
			const adjustedBearing = (bearing - ICON_BEARING_OFFSET + 360) % 360;

			updates[count++] = {
				id: driver.id,
				position: [lng, lat],
				bearing: adjustedBearing,
			};
		}

		const layer = newRealtimeDriversLayer(updates);
		this.deck.setProps({ layers: [layer] });

		// setTimeout(() => {
		// 	this.animationFrame = requestAnimationFrame(this.animate)
		// }, 125)

		this.animationFrame = requestAnimationFrame(this.animate)
	};
}

// interface DriverAnimationInfo {
// 	segmentsPercentages: number[];
// 	segments: [number, number][];
// 	currentSegment: number;
// 	segmentPercentageElapsed: number;
// }

// export class RealtimeDriversAnimator {
// 	private map: MaplibreMap;
// 	private deck: MapboxOverlay;
// 	private animationFrame: number | null = null;
// 	private removeSubscriptionListener: (() => void) | null = null;
// 	private cachePool = new Map<number, DriverAnimationInfo>();
// 	private drivers: readonly RealtimeDriver[] = [];
// 	private lastFetchTime = 0;
// 	private pathTime = 5000;
//
// 	constructor(map: MaplibreMap, deck: MapboxOverlay) {
// 		this.map = map;
// 		this.deck = deck;
// 	}
//
// 	start() {
// 		$realtimeDrivers.set(FetchedRealtimeDrivers)
// 		$lastFetchTime.set(Date.now())
// 		// Subscribe to store updates
// 		this.removeSubscriptionListener = $realtimeDrivers.subscribe((drivers) => {
// 			this.drivers = drivers;
// 			this.lastFetchTime = $lastFetchTime.get();
// 			this.rebuildCachePool();
// 		});
//
// 		this.animate();
// 	}
//
// 	stop() {
// 		$realtimeDrivers.set([]);
// 		if (this.animationFrame) {
// 			cancelAnimationFrame(this.animationFrame)
// 			this.animationFrame = null
// 		};
// 		if (this.removeSubscriptionListener) {
// 			this.removeSubscriptionListener()
// 			this.removeSubscriptionListener = null
// 		};
// 		if (this.deck) {
// 			this.deck.setProps({ layers: [] });
// 		}
// 	}
//
// 	private rebuildCachePool() {
// 		this.cachePool.clear();
// 		for (const driver of this.drivers) {
// 			this.cachePool.set(driver.id, {
// 				segmentsPercentages: calculateSegmentsPercentages(
// 					driver.coordinates,
// 					driver.distance
// 				),
// 				segments: driver.coordinates,
// 				currentSegment: 0,
// 				segmentPercentageElapsed: 0,
// 			});
//
// 			const cache = this.cachePool.get(driver.id)
// 			let segmentsPercentages = 0
// 			cache?.segmentsPercentages.forEach((p) => {
// 				segmentsPercentages += p
// 			})
// 			console.log(cache)
// 			console.log(segmentsPercentages)
// 		}
// 	}
//
// 	private animate = () => {
// 		if (!this.drivers.length) {
// 			this.animationFrame = requestAnimationFrame(this.animate);
// 			return;
// 		}
//
// 		const elapsedTime = Date.now() - this.lastFetchTime;
// 		const progress = Math.min(elapsedTime / this.pathTime, 1);
//
// 		const updates = Array<RealtimeDriverPosition>(this.drivers.length - 1);
//
// 		for (const driver of this.drivers) {
// 			const cache = this.cachePool.get(driver.id);
// 			if (!cache) continue;
//
// 			const [currentSegment, segmentElapsed] = this.getCurrentSegment(
// 				cache.currentSegment,
// 				cache.segmentPercentageElapsed,
// 				cache.segmentsPercentages,
// 				progress
// 			);
//
// 			console.log(progress)
// 			console.log("currentSegment", currentSegment)
// 			console.log("segmentElapsed", segmentElapsed)
// 			const position = getDriverPosition(
// 				cache.segments[currentSegment],
// 				cache.segments[currentSegment + 1],
// 				progress
// 			);
// 			const bearing = getDriverBearing(
// 				cache.segments[currentSegment],
// 				cache.segments[currentSegment + 1]
// 			);
//
// 			this.cachePool.set(driver.id, {
// 				...cache,
// 				currentSegment,
// 				segmentPercentageElapsed: segmentElapsed,
// 			});
//
// 			updates.push({ id: driver.id, position, bearing });
// 		}
//
// 		const layer = newRealtimeDriversLayer(updates);
// 		this.deck.setProps({ layers: [layer] });
//
// 		this.animationFrame = requestAnimationFrame(this.animate);
// 	};
//
// 	private getCurrentSegment(
// 		currentSegment: number,
// 		segmentElapsed: number,
// 		segmentsPercentages: number[],
// 		progress: number
// 	): [number, number] {
// 		for (currentSegment; currentSegment < segmentsPercentages.length - 1; currentSegment++) {
// 			console.log(currentSegment)
// 			if (progress <= segmentElapsed) {
// 				return [currentSegment, segmentElapsed]
// 			}
// 			segmentElapsed += segmentsPercentages[currentSegment + 1]
// 		}
//
// 		return [currentSegment, segmentElapsed];
// 	}
// }
