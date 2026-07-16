import { useMovingIconLayer } from "maps/composables/useMovingIconLayer.ts"
import { useMovingIconLayerMaplibre } from "maps/composables/useMovingIconLayerMaplibre.js"
import { $movingDrivers } from "../stores/movingDriversStore.ts"
import { useMovingDriverPolling, type MovingDriverProvider, MovingDriverProviders } from "./useMovingDriverPolling.ts"
import { mockMovingDrivers } from "./test/test_data/moving_driver_data.ts"
import { onMounted, onUnmounted, type Component } from "vue"

function carIconUrl(): string {
	return `data:image/svg+xml;charset=utf-8,${encodeURIComponent('\<svg width="112" height="116" viewBox="0 0 28 29" xmlns="http://www.w3.org/2000/svg"><g filter="url(#4510d74bbe320d84ab4638c9ccb2cfaf__a)"><path d="M6.932 9.334l-.008.008-.008.008c-3.89 3.993-3.888 10.454.004 14.446a9.842 9.842 0 0 0 14.16 0c3.892-3.992 3.893-10.453.004-14.446l-.008-.008-.008-.008L14.7 3.088 14 2.4l-.7.687-6.368 6.246z" stroke="#fff" stroke-width="2"/></g><path d="M7.636 23.098c-3.514-3.604-3.515-9.445-.003-13.05L14 3.801l6.367 6.245c3.512 3.605 3.51 9.447-.003 13.05a8.842 8.842 0 0 1-12.728 0z" fill="#3CB300"/><defs><filter id="4510d74bbe320d84ab4638c9ccb2cfaf__a" x="2" y="0" width="24" height="28.802" filterUnits="userSpaceOnUse" color-interpolation-filters="sRGB"><feFlood flood-opacity="0" result="BackgroundImageFix"/><feColorMatrix in="SourceAlpha" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/><feOffset/><feGaussianBlur stdDeviation=".5"/><feColorMatrix values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.3 0"/><feBlend in2="BackgroundImageFix" result="effect1_dropShadow_23120_22065"/><feBlend in="SourceGraphic" in2="effect1_dropShadow_23120_22065" result="shape"/></filter></defs></svg>')}`;
}

export function useMovingDriverLayer(
	provider: MovingDriverProvider = MovingDriverProviders.mock(mockMovingDrivers),
	initialLat = 55.751,
	initialLon = 37.618,
	popupComponent?: Component,
) {
	const { start: startAnimation } = useMovingIconLayer({
		paths: $movingDrivers,
		icon: { url: carIconUrl(), width: 112, height: 116 },
		layerId: "geolocation-moving",
		popupComponent,
		onClick(id) { console.log("moving driver:", id) },
	})

	const { start: startPolling, stop } = useMovingDriverPolling(provider, initialLat, initialLon)

	onMounted(() => startPolling())
	onUnmounted(() => stop())

	return { start: startAnimation }
}

export function useMovingDriverLayerMaplibre(
	provider: MovingDriverProvider = MovingDriverProviders.mock(mockMovingDrivers),
	initialLat = 55.751,
	initialLon = 37.618,
	popupComponent?: Component,
) {
	const { start: startAnimation } = useMovingIconLayerMaplibre({
		paths: $movingDrivers,
		iconUrl: carIconUrl(),
		iconWidth: 28,
		iconHeight: 29,
		layerId: "geolocation-moving",
		popupComponent,
		onClick(id) { console.log("moving driver:", id) },
	})

	const { start: startPolling, stop } = useMovingDriverPolling(provider, initialLat, initialLon)

	onMounted(() => startPolling())
	onUnmounted(() => stop())

	return { start: startAnimation }
}
