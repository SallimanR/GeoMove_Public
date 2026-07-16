import { $mapInstance, $deckOverlay } from "../stores/mapsStore";
import { createAnimatedIconLayer, type AnimatedIconConfig } from "../layers/animatedIconLayer";
import { Popup } from "maplibre-gl";
import type { MapLibreMap } from "maplibre-gl";
import type { MapboxOverlay } from "@deck.gl/mapbox";
import { createApp, onUnmounted, type Component, type App } from "vue";
import { haversineDistance } from "geo";
import type { ReadableAtom } from "nanostores";

export interface MovingPath {
	id: number;
	coordinates: [number, number][];
	time: number;
	distance: number;
}

export interface MovingPosition {
	id: number;
	position: [number, number];
	bearing: number;
}

interface CacheEntry {
	cumulativeDistances: number[];
	totalDistance: number;
	totalTimeSeconds: number;
}

export interface UseMovingIconLayerOptions {
	paths: ReadableAtom<MovingPath[]>;
	icon: AnimatedIconConfig;
	layerId?: string;
	iconSize?: number;
	popupComponent?: Component;
	onClick?: (id: number) => void;
	onHover?: (id: number) => void;
}

export function useMovingIconLayer(options: UseMovingIconLayerOptions) {
	let map = $mapInstance.get();
	let deck = $deckOverlay.get();
	let paths = options.paths.get();

	const cachePool = new Map<number, CacheEntry>();
	const popupTrackers = new Map<number, { popup: Popup; app?: App }>();

	let animationFrame: number | null = null;
	let startTime = 0;
	let isRunning = false;

	const unsubMap = $mapInstance.subscribe((v) => {
		map = v as MapLibreMap | null;
		maybeStart();
	});
	const unsubDeck = $deckOverlay.subscribe((v) => {
		deck = v as MapboxOverlay | null;
		maybeStart();
	});
	const unsubPaths = options.paths.subscribe((v) => {
		paths = v;
		if (paths.length === 0) {
			stop();
			return;
		}
		startTime = Date.now();
		buildCache();
		maybeStart();
	});

	function maybeStart() {
		if (isRunning) return;
		if (!map || !deck || paths.length === 0) return;
		startTime = Date.now();
		buildCache();
		isRunning = true;
		animate();
	}

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
			const bearing = (getBearing(s, e) - 90 + 360) % 360;

			result.push({ id: path.id, position: [lng, lat], bearing });
		}

		return result;
	}

	function syncPopups(positions: MovingPosition[]) {
		if (!map) return;

		for (const pos of positions) {
			let tracker = popupTrackers.get(pos.id);

			if (!tracker) {
				const popup = new Popup({
					closeButton: false,
					closeOnClick: false,
					anchor: "bottom",
					offset: [0, -25],
				})
					.setLngLat(pos.position)
					.addTo(map as MapLibreMap);

				if (options.popupComponent) {
					popup.setHTML('<div class="moving-popup-mount"></div>');
					const mountEl = popup.getElement().querySelector(".moving-popup-mount");
					if (mountEl) {
						const app = createApp(options.popupComponent, { id: pos.id });
						app.mount(mountEl);
						tracker = { popup, app };
					} else {
						popup.setHTML(`<div>#${pos.id}</div>`);
						tracker = { popup };
					}
				} else {
					tracker = { popup };
				}

				popupTrackers.set(pos.id, tracker);
			} else {
				tracker.popup.setLngLat(pos.position);
			}
		}
	}

	function animate() {
		if (!isRunning) return;
		if (!map || !deck) return;

		const positions = computePositions();
		const layer = createAnimatedIconLayer(positions, options.icon, {
			layerId: options.layerId,
			size: options.iconSize,
			onClick: options.onClick,
			onHover: options.onHover,
		});
		deck.setProps({ layers: [layer] });
		syncPopups(positions);

		animationFrame = requestAnimationFrame(animate);
	}

	function stop() {
		isRunning = false;
		if (animationFrame !== null) {
			cancelAnimationFrame(animationFrame);
			animationFrame = null;
		}
		if (deck) {
			deck.setProps({ layers: [] });
		}
		popupTrackers.forEach((t) => {
			t.app?.unmount();
			t.popup.remove();
		});
		popupTrackers.clear();
	}

	onUnmounted(() => {
		unsubMap();
		unsubDeck();
		unsubPaths();
		stop();
	});

	return {
		start: () => maybeStart(),
		stop,
	};
}
