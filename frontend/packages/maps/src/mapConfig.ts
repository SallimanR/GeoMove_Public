import type { AddLayerObject, LngLatLike, MapOptions } from "maplibre-gl"

const MAP_CENTER_COORDINATES: LngLatLike = [37.618037, 55.743293]

export const getMapConfig = (
	styleApi: string
): Partial<MapOptions> => ({
	style: styleApi,
	center: MAP_CENTER_COORDINATES,
	zoom: 15,
	pitch: 0,
	// cooperativeGestures: true,
	// scrollZoom: {
	// 	around: "center"
	// },
	canvasContextAttributes: { antialias: true },
});

export const MapLayer_3dLayer = <AddLayerObject>{
	id: "3d-buildings",
	source: "protomaps",
	"source-layer": "buildings",
	type: "fill-extrusion",
	minzoom: 15,
	filter: [">", ["get", "height"], 0],
	paint: {
		"fill-extrusion-color": [
			"interpolate",
			["linear"],
			["get", "height"],
			0,
			"lightgray",
		],
		"fill-extrusion-height": ["get", "height"],
		"fill-extrusion-opacity": 0.7,
		"fill-extrusion-base": ["coalesce", ["get", "min_height"], 0],
	},
}
