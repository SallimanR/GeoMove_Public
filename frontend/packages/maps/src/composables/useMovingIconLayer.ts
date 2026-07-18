import { $mapInstance, $deckOverlay } from "../stores/mapsStore";
import { createAnimatedIconLayer } from "../layers/animatedIconLayer";
import { Popup } from "maplibre-gl";
import type { MapLibreMap } from "maplibre-gl";
import type { MapboxOverlay } from "@deck.gl/mapbox";
import { createApp, onUnmounted, type App } from "vue";
import { useMovingIconLayerCore } from "./useMovingIconLayerCore";
import type { MovingPosition, UseMovingIconLayerOptions } from "../types/movingIconLayerShared";

export function useMovingIconLayer(options: UseMovingIconLayerOptions) {
	let map = $mapInstance.get();
	let deck = $deckOverlay.get();

	const popupTrackers = new Map<number, { popup: Popup; app?: App }>();

	const core = useMovingIconLayerCore({
		pathsAtom: options.paths,
		adjustBearing: (b) => (b - 90 + 360) % 360,
		isReady: () => !!map && !!deck,
	});

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

	core.setOnFrame((positions: MovingPosition[]) => {
		const layer = createAnimatedIconLayer(positions,
			options.iconUrl,
			options.iconWidth,
			options.iconHeight
			, {
				layerId: options.layerId,
				size: options.iconSize ?? options.iconWidth,
				onClick: options.onClick,
				onHover: options.onHover,
			});
		deck!.setProps({ layers: [layer] });
		syncPopups(positions);
	});

	core.setOnStopCleanup(() => {
		if (deck) {
			deck.setProps({ layers: [] });
		}
		popupTrackers.forEach((t) => {
			t.app?.unmount();
			t.popup.remove();
		});
		popupTrackers.clear();
	});

	const unsubMap = $mapInstance.subscribe((v) => {
		map = v as MapLibreMap | null;
		core.start();
	});
	const unsubDeck = $deckOverlay.subscribe((v) => {
		deck = v as MapboxOverlay | null;
		core.start();
	});

	onUnmounted(() => {
		unsubMap();
		unsubDeck();
	});

	return {
		start: () => core.start(),
		stop: () => core.fullStop(),
	};
}
