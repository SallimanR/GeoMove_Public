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
	source: "map-tiles",
	"source-layer": "building",
	type: "fill-extrusion",
	minzoom: 15,
	filter: ["!=", ["get", "hide_3d"], true],
	paint: {
		"fill-extrusion-color": [
			"interpolate",
			["linear"],
			["get", "render_height"],
			0,
			"lightgray",
			// 200,
			// "gray",
			// 400,
			// "gray",
		],
		"fill-extrusion-height": [
			"interpolate",
			["linear"],
			["zoom"],
			15,
			0,
			16,
			["get", "render_height"],
		],
		"fill-extrusion-opacity": 0.7,
		// "fill-extrusion-base": [
		// 	"case",
		// 	[">=", ["get", "zoom"], 16],
		// 	["get", "render_min_height"],
		// 	0,
		// ],
	},
}
