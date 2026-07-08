import { $deckOverlay, $mapInstance } from "../stores/mapsStore"
import { type Layer, type LayersList } from "deck.gl"
import { ref } from "vue"

export function useMapsActions() {
	const deckLayers = ref<Layer[]>([])
	const updateDeckOverlayLayers = (): void => {
		const deck = $deckOverlay.get()
		if (deck) {
			deck.setProps({ layers: deckLayers.value as LayersList })
		}
	}

	const addDeckLayer = (layer: Layer): void => {
		deckLayers.value = [...deckLayers.value, layer]
		updateDeckOverlayLayers()
	}

	const removeDeckLayer = (layerId: string): void => {
		deckLayers.value = deckLayers.value.filter(l => l.id !== layerId)
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
