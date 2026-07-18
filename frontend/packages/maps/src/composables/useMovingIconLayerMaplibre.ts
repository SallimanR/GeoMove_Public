import { $mapInstance } from "../stores/mapsStore";
import { addPopupToMap, type PopupEntry } from "../utils/mapPopup";
import type { MapLibreMap, GeoJSONSource } from "maplibre-gl";
import { onUnmounted } from "vue";
import { useMovingIconLayerCore } from "./useMovingIconLayerCore";
import type { MovingPosition, UseMovingIconLayerOptions } from "../types/movingIconLayerShared";

export function useMovingIconLayerMaplibre(options: UseMovingIconLayerOptions) {
	const ICON_ID = `${options.layerId ?? "moving-icons"}-icon`;
	const SOURCE_ID = `${options.layerId ?? "moving-icons"}-source`;
	const LAYER_ID = options.layerId ?? "moving-icons-maplibre";
	const POPUPS_GROUP = `${LAYER_ID}-popups`;

	let map = $mapInstance.get();
	let sourceReady = false;

	const popupTrackers = new Map<number, PopupEntry>();

	const core = useMovingIconLayerCore({
		pathsAtom: options.paths,
		adjustBearing: (b) => (90 - b + 360) % 360,
		isReady: () => !!map && sourceReady,
	});

	function initMapResources(m: MapLibreMap) {
		const onReady = () => {
			loadIcon(m)
				.then(() => ensureSourceAndLayer(m))
				.then(() => {
					sourceReady = true;
					core.start();
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

	core.setOnFrame((positions: MovingPosition[]) => {
		updateSource(positions);
		syncPopups(positions);
	});

	core.setOnStopCleanup(() => {
		cleanupMapResources();
		popupTrackers.forEach((t) => { t.destroy(); });
		popupTrackers.clear();
	});

	function cleanupMapResources() {
		if (!map) return;
		try { if (map.getLayer(LAYER_ID)) map.removeLayer(LAYER_ID); } catch (_) { }
		try { if (map.getSource(SOURCE_ID)) map.removeSource(SOURCE_ID); } catch (_) { }
		try { if (map.hasImage(ICON_ID)) map.removeImage(ICON_ID); } catch (_) { }
	}

	const unsubMap = $mapInstance.subscribe((v) => {
		map = v as MapLibreMap | null;
		if (map) initMapResources(map);
		else {
			sourceReady = false;
			core.fullStop();
		}
	});

	onUnmounted(() => {
		unsubMap();
	});

	return {
		start: () => core.start(),
		stop: () => core.fullStop(),
	};
}
