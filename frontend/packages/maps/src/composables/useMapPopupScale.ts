import { ref, computed, onUnmounted } from "vue";
import { $mapInstance } from "../stores/mapsStore";

export function useMapPopupScale() {
	const zoom = ref(12);
	let unsub: (() => void) | undefined;

	const cancel = $mapInstance.subscribe((map) => {
		if (unsub) unsub();
		if (map) {
			zoom.value = map.getZoom();
			map.on("zoom", () => {
				zoom.value = map.getZoom();
			});
			unsub = () => {
				map.off("zoom", () => {});
			};
		}
	});

	const scale = computed(() => {
		if (zoom.value >= 18) return 1;
		if (zoom.value <= 5) return 0.3;
		return 0.3 + (zoom.value - 5) * (0.7 / 13);
	});

	onUnmounted(() => {
		if (unsub) unsub();
	});

	return { scale };
}
