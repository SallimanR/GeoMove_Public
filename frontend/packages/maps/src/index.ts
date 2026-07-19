export { MAP_CENTER_LAT, MAP_CENTER_LON } from './mapConfig'

export { default as Maps } from './components/Maps.vue'
export { default as MapsOverlayControls } from './components/MapsOverlayControls.vue'
export { default as MapsLocationPicker } from './components/MapsLocationPicker.vue'

export { type UseMovingIconLayerOptions, type MovingPath, type MovingPosition } from './types/movingIconLayerShared.ts'
export type { GeoPoint } from './types/geoPoint.ts'

export { useMovingIconLayer } from './composables/useMovingIconLayer.ts'
export { useMovingIconLayerMaplibre } from './composables/useMovingIconLayerMaplibre.ts'

export {
  $coords,
  $mapInstance,
  $deckOverlay,
  $mapCenterAddress,
  $mapCenterAddressText,
  $locationPicking,
  setPickCallback,
  invokePickCallback,
  clearPickCallback,
} from './stores/mapsStore'

export {
  $startPoint,
  $endPoint,
  $startAddress,
  $endAddress,
  $routePath,
  $isRouteLoading,
} from './stores/routeStore'

export { addPopupToMap, removeAllPopups, removePopupsByGroup } from './utils/mapPopup'
