import type { AddLayerObject, MapOptions } from "maplibre-gl"

export const MAP_STYLE_API = import.meta.env.PUBLIC_MAP_STYLE_API
export const MAP_TILES_API = import.meta.env.PUBLIC_MAP_TILES_API

export const MAP_CONFIG = <MapOptions>{
	style: MAP_STYLE_API,
	center: [37.618037, 55.743293],
	zoom: 15,
	pitch: 0,
	// cooperativeGestures: true,
	// scrollZoom: {
	// 	around: "center"
	// },
	canvasContextAttributes: { antialias: true },
}

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
