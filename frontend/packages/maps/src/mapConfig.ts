import type { AddLayerObject, MapOptions } from "maplibre-gl"

export const MAP_CENTER_LAT = 55.743293
export const MAP_CENTER_LON = 37.618037

export const getMapConfig = (
	styleApi: string
): Partial<MapOptions> => ({
	style: styleApi,
	center: { lat: MAP_CENTER_LAT, lon: MAP_CENTER_LON },
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
	minzoom: 13,
	paint: {
		"fill-extrusion-color": [
			"interpolate",
			["linear"],
			["coalesce", ["get", "render_height"], ["get", "height"], 8],
			0,
			"lightgray",
			100,
			"darkgray",
		],
		"fill-extrusion-height": ["coalesce", ["get", "render_height"], ["get", "height"], 8],
		"fill-extrusion-opacity": 0.7,
		"fill-extrusion-base": ["coalesce", ["get", "render_min_height"], ["get", "min_height"], 0],
	},
}
