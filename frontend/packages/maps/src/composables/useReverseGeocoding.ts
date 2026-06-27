import { ref, createApp } from "vue";
import { type Map as MaplibreMap, Popup, LngLat } from "maplibre-gl";
import { getReverseGeocoding } from "geo";
import ReverseGeocodingPopup from "../components/ReverseGeocodingPopup.vue";

export function useReverseGeocoding(map: MaplibreMap) {
	const popup = ref<Popup | null>(null);

	const showAddressPopup = async (lngLat: LngLat) => {
		try {
			const result = await getReverseGeocoding(
				lngLat.lat,
				lngLat.lng,
			);
			const feature = result.features[0]
			if (!feature) return;

			const { street, housenumber } = feature.properties;

			const clickPopup = new Popup({
				closeButton: true,
				closeOnClick: false,
			});
			clickPopup
				.setLngLat([lngLat.lng, lngLat.lat])
				.setHTML("<div id='reverse-geocoding-popup'></div>")
				.addTo(map);
			const mountPoint = clickPopup
				.getElement()
				.querySelector("#reverse-geocoding-popup");

			if (mountPoint) {
				const popupVueApp = createApp(ReverseGeocodingPopup, {
					street: street,
					housenumber: housenumber,
				});
				popupVueApp.mount(mountPoint);

				// TODO:
				clickPopup.on('close', () => {
					if (popupVueApp) {
						popupVueApp.unmount();
					}
				})
			}
		} catch (error) {
		}
	}
	return { showAddressPopup }
}
