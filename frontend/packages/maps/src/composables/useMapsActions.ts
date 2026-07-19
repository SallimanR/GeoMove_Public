import { $deckOverlay, $mapInstance } from "../stores/mapsStore"
import type { Layer } from "@deck.gl/core"
import type { LayersList } from "@deck.gl/core"
import { ref, type Ref } from "vue"

export function useMapsActions() {
	const deckLayers: Ref<Layer[]> = ref([])
	const updateDeckOverlayLayers = (): void => {
		const deck = $deckOverlay.get()
		if (deck) {
			deck.setProps({ layers: deckLayers.value as LayersList })
		}
	}

	const addDeckLayer = (layer: Layer): void => {
		const current = deckLayers.value
		deckLayers.value = current.concat(layer)
		updateDeckOverlayLayers()
	}

	const removeDeckLayer = (layerId: string): void => {
		const current = deckLayers.value
		deckLayers.value = current.filter((l: any) => l.id !== layerId)
		updateDeckOverlayLayers()
	}

	const flyTo = (lat: number, lon: number): void => {
		const map = $mapInstance.get()
		if (!map) {
			return
		}
		map.flyTo({
			center: { lat: lat, lon: lon },
			duration: 500,
		})
	}

	return {
		addDeckLayer,
		removeDeckLayer,
		flyTo,
	}
}
