import { ref, provide, inject, onUnmounted } from 'vue'
import { Map as MaplibreMap, type MapOptions, GeolocateControl } from "maplibre-gl"
import { MapboxOverlay } from '@deck.gl/mapbox';
import { type Layer, type LayersList } from "@deck.gl/core"
import { $coords, $deckOverlay, $mapInstance } from '../stores/mapsStore';
import { MAP_CONFIG, MAP_TILES_API, MapLayer_3dLayer } from '../mapConfig';

export type MapContext = {
	map: MaplibreMap
	deckOverlay: MapboxOverlay,
}

const mapContextKey = Symbol('mapContext')


export function useMaps() {
	const map = $mapInstance
	const deckOverlay = $deckOverlay
	// const map = ref<MaplibreMap | null>(null)
	// const deckOverlay = ref<MapboxOverlay | null>(null)
	// const deckLayers = ref<Layer[]>([])

	const initMap = (container: HTMLElement, options: Partial<MapOptions> = {}): void => {
		if (map.value) return

		const mapInstance = new MaplibreMap(
			{
				...MAP_CONFIG,
				...options,
				container,
			}
		)

		const center = mapInstance.getCenter();
		$coords.setKey("center", { lat: center.lat, lon: center.lng });

		mapInstance.once("load", () => {
			const layers = mapInstance.getStyle().layers;

			// Insert the layer beneath any symbol layer.
			let labelLayerId = "";
			// Find the index of the first symbol layer in the map style
			for (const layer of layers) {
				if (layer.type === "symbol" && layer.layout?.["text-field"]) {
					labelLayerId = layer.id;
					console.log("[DEBUG] labelLayerId", labelLayerId)
					break;
				}
			}

			mapInstance.addSource("map-tiles", {
				url: MAP_TILES_API,
				type: "vector",
			});
			mapInstance.addLayer(
				MapLayer_3dLayer,
				labelLayerId,
			);

		})
		mapInstance.on("moveend", () => {
			const center = mapInstance.getCenter();
			$coords.setKey("center", { lat: center.lat, lon: center.lng });
		});

		// GPS button on maps
		mapInstance.addControl(
			new GeolocateControl({
				positionOptions: {
					enableHighAccuracy: true
				},
				trackUserLocation: true
			})
		);

		const overlay = new MapboxOverlay({})
		mapInstance.addControl(overlay)

		map.set(mapInstance)
		deckOverlay.set(overlay)
	}

	onUnmounted(() => {
		if (map.value) {
			map.value.remove()
			map.set(null)
		}
		if (deckOverlay.value) {
			deckOverlay.value.finalize()
			deckOverlay.set(null)
		}
		// deckLayers.value = []
	})

	const provideMap = (): void => {
		if (!map.value || !deckOverlay.value) {
			throw new Error('Map not initialized')
		}
		provide(mapContextKey, {
			map: map.value,
			deckOverlay: deckOverlay.value,
		})
	}

	const injectMap = (): MapContext => {
		const context = inject<MapContext>(mapContextKey)
		if (!context) {
			throw new Error('useMap must be used within a Map component')
		}
		return context
	}

	return {
		initMap,
		provideMap,
		injectMap,
		map,
		deckOverlay,
	}
}
