import { onUnmounted } from "vue";
import { haversineDistance } from "geo";
import type { CacheEntry, MovingPath, MovingPosition } from "../types/movingIconLayerShared";
import type { ReadableAtom } from "nanostores";

export interface UseMovingIconLayerCoreOptions {
	pathsAtom: ReadableAtom<MovingPath[]>;
	adjustBearing: (rawBearing: number) => number;
	isReady: () => boolean;
}

export function useMovingIconLayerCore(options: UseMovingIconLayerCoreOptions) {
	let paths: MovingPath[] = options.pathsAtom.get() as MovingPath[];

	const cachePool = new Map<number, CacheEntry>();

	let animationFrame: number | null = null;
	let startTime = 0;
	let isRunning = false;

	let onFrame: ((positions: MovingPosition[]) => void) | null = null;
	let onStopCleanup: (() => void) | null = null;

	const unsubPaths = options.pathsAtom.subscribe((v) => {
		paths = v as MovingPath[];
		if (paths.length === 0) {
			fullStop();
			return;
		}
		buildCache();
		start();
	});

	function buildCache() {
		cachePool.clear();
		for (const path of paths) {
			const coords = path.coordinates;
			if (coords.length < 2) continue;

			const cumulativeDistances = new Array<number>(coords.length);
			cumulativeDistances[0] = 0;
			for (let i = 0; i < coords.length - 1; i++) {
				cumulativeDistances[i + 1] =
					cumulativeDistances[i] + haversineDistance(coords[i], coords[i + 1]);
			}

			cachePool.set(path.id, {
				cumulativeDistances,
				totalDistance: cumulativeDistances[coords.length - 1],
				totalTimeSeconds: path.time / 1000,
			});
		}
	}

	function getBearing(a: [number, number], b: [number, number]): number {
		return (Math.atan2(b[1] - a[1], b[0] - a[0]) * 180) / Math.PI;
	}

	function computePositions(): MovingPosition[] {
		const elapsed = (Date.now() - startTime) / 1000;
		const result: MovingPosition[] = [];

		for (const path of paths) {
			const cache = cachePool.get(path.id);
			if (!cache) continue;

			const totalTime = cache.totalTimeSeconds;
			let t = totalTime > 0 ? elapsed % totalTime : 0;

			if (t <= 0) {
				result.push({ id: path.id, position: path.coordinates[0], bearing: 0 });
				continue;
			}
			if (t >= totalTime) {
				const last = path.coordinates[path.coordinates.length - 1];
				result.push({ id: path.id, position: last, bearing: 0 });
				continue;
			}

			const targetDist = (t / totalTime) * path.distance;
			const clampedDist = Math.min(targetDist, cache.totalDistance);
			const cum = cache.cumulativeDistances;

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
			const segIndex = low - 1;
			const s = path.coordinates[segIndex];
			const e = path.coordinates[segIndex + 1];
			const segDist = cum[segIndex + 1] - cum[segIndex];
			const ratio = (clampedDist - cum[segIndex]) / segDist;

			const lng = s[0] + ratio * (e[0] - s[0]);
			const lat = s[1] + ratio * (e[1] - s[1]);
			const bearing = options.adjustBearing(getBearing(s, e));

			result.push({ id: path.id, position: [lng, lat], bearing });
		}

		return result;
	}

	function start() {
		if (isRunning) return;
		if (!options.isReady() || paths.length === 0) return;
		startTime = Date.now();
		buildCache();
		isRunning = true;
		animationFrame = requestAnimationFrame(loop);
	}

	function stop() {
		isRunning = false;
		if (animationFrame !== null) {
			cancelAnimationFrame(animationFrame);
			animationFrame = null;
		}
	}

	function loop() {
		if (!isRunning) return;
		if (!options.isReady()) return;
		const positions = computePositions();
		onFrame?.(positions);
		animationFrame = requestAnimationFrame(loop);
	}

	function setOnFrame(fn: (positions: MovingPosition[]) => void) {
		onFrame = fn;
	}

	function setOnStopCleanup(fn: () => void) {
		onStopCleanup = fn;
	}

	function fullStop() {
		stop();
		onStopCleanup?.();
	}

	onUnmounted(() => {
		unsubPaths();
		fullStop();
	});

	return {
		cachePool,
		buildCache,
		computePositions,
		getBearing,
		start,
		fullStop,
		setOnFrame,
		setOnStopCleanup,
	};
}
