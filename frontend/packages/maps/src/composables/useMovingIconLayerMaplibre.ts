import { $mapInstance } from "../stores/mapsStore";
import { addPopupToMap, type PopupEntry } from "../utils/mapPopup";
import type { MapLibreMap, GeoJSONSource } from "maplibre-gl";
import { onUnmounted, type Component } from "vue";
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

export interface UseMovingIconLayerMaplibreOptions {
	paths: ReadableAtom<MovingPath[]>;
	iconUrl: string;
	iconWidth: number;
	iconHeight: number;
	layerId?: string;
	popupComponent?: Component;
	onClick?: (id: number) => void;
	onHover?: (id: number) => void;
}

export function useMovingIconLayerMaplibre(options: UseMovingIconLayerMaplibreOptions) {
	const ICON_ID = `${options.layerId ?? "moving-icons"}-icon`;
	const SOURCE_ID = `${options.layerId ?? "moving-icons"}-source`;
	const LAYER_ID = options.layerId ?? "moving-icons-maplibre";
	const POPUPS_GROUP = `${LAYER_ID}-popups`;

	let map = $mapInstance.get();
	let paths = options.paths.get();

	const cachePool = new Map<number, CacheEntry>();
	const popupTrackers = new Map<number, PopupEntry>();

	let animationFrame: number | null = null;
	let startTime = 0;
	let isRunning = false;
	let sourceReady = false;

	const unsubMap = $mapInstance.subscribe((v) => {
		map = v as MapLibreMap | null;
		if (map) initMapResources(map);
		else {
			sourceReady = false;
			stop();
		}
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

	function initMapResources(m: MapLibreMap) {
		const onReady = () => {
			loadIcon(m)
				.then(() => ensureSourceAndLayer(m))
				.then(() => {
					sourceReady = true;
					maybeStart();
				});
		};

		if (m.loaded()) {
			onReady();
		} else {
			m.once("load", onReady);
		}
	}

	function loadIcon(m: MapLibreMap): Promise<void> {
		return new Promise((resolve) => {
			if (m.hasImage(ICON_ID)) {
				resolve();
				return;
			}
			const img = new Image(options.iconWidth, options.iconHeight);
			img.onload = () => {
				if (!m.hasImage(ICON_ID)) m.addImage(ICON_ID, img);
				resolve();
			};
			img.onerror = () => resolve();
			img.src = options.iconUrl;
		});
	}

	function ensureSourceAndLayer(m: MapLibreMap) {
		if (m.getSource(SOURCE_ID)) return;

		m.addSource(SOURCE_ID, {
			type: "geojson",
			data: { type: "FeatureCollection", features: [] },
		});

		m.addLayer({
			id: LAYER_ID,
			type: "symbol",
			source: SOURCE_ID,
			layout: {
				"icon-image": ICON_ID,
				"icon-size": 1.0,
				"icon-rotate": ["get", "bearing"],
				"icon-rotation-alignment": "map",
				"icon-allow-overlap": true,
				"icon-ignore-placement": true,
			},
		});
	}

	function maybeStart() {
		if (isRunning) return;
		if (!map || !sourceReady || paths.length === 0) return;
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
			const bearing = (90 - getBearing(s, e) + 360) % 360;

			result.push({ id: path.id, position: [lng, lat], bearing });
		}

		return result;
	}

	function updateSource(positions: MovingPosition[]) {
		if (!map) return;
		const source = map.getSource(SOURCE_ID) as GeoJSONSource | undefined;
		if (!source) return;

		source.setData({
			type: "FeatureCollection",
			features: positions.map((p) => ({
				type: "Feature",
				geometry: { type: "Point", coordinates: p.position },
				properties: { id: p.id, bearing: p.bearing },
			})),
		});
	}

	function syncPopups(positions: MovingPosition[]) {
		if (!map) return;

		for (const pos of positions) {
			let entry = popupTrackers.get(pos.id);

			if (!entry && options.popupComponent) {
				const result = addPopupToMap(
					pos.position[1], pos.position[0],
					options.popupComponent,
					{ id: pos.id },
					POPUPS_GROUP,
					{ offset: [0, -25] },
				);
				if (result) popupTrackers.set(pos.id, result);
			} else if (entry) {
				entry.popup.setLngLat(pos.position);
			}
		}
	}

	function animate() {
		if (!isRunning) return;
		if (!map || !sourceReady) return;

		const positions = computePositions();
		updateSource(positions);
		syncPopups(positions);

		animationFrame = requestAnimationFrame(animate);
	}

	function stop() {
		isRunning = false;
		if (animationFrame !== null) {
			cancelAnimationFrame(animationFrame);
			animationFrame = null;
		}
		cleanupMapResources();
		popupTrackers.forEach((t) => { t.destroy(); });
		popupTrackers.clear();
	}

	function cleanupMapResources() {
		if (!map) return;
		try { if (map.getLayer(LAYER_ID)) map.removeLayer(LAYER_ID); } catch (_) { }
		try { if (map.getSource(SOURCE_ID)) map.removeSource(SOURCE_ID); } catch (_) { }
		try { if (map.hasImage(ICON_ID)) map.removeImage(ICON_ID); } catch (_) { }
	}

	onUnmounted(() => {
		unsubMap();
		unsubPaths();
		stop();
	});

	return {
		start: () => maybeStart(),
		stop,
	};
}
