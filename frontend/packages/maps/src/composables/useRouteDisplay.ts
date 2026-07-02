import { watch } from "vue";
import { type AddLayerObject, type Map as MaplibreMap } from "maplibre-gl";
import { $endPoint, $routePath, $startPoint } from "../stores/routeStore";
import { useStore } from '@nanostores/vue';

const MapSourceID_StartPoint = "source-start-point"
const MapLayerID_StartPoint = "layer-start-point"
const MapLayer_StartPoint = <AddLayerObject>{
	id: MapLayerID_StartPoint,
	source: MapSourceID_StartPoint,
	type: 'circle',
	paint: {
		'circle-radius': 8,
		'circle-color': '#4CAF50',
		'circle-stroke-width': 2,
		'circle-stroke-color': '#ffffff'
	}
}

const MapSourceID_EndPoint = "end-start-point"
const MapLayerID_EndPoint = "layer-end-point"
const MapLayer_EndPoint = <AddLayerObject>{
	id: MapLayerID_EndPoint,
	source: MapSourceID_EndPoint,
	type: 'circle',
	paint: {
		'circle-radius': 8,
		'circle-color': '#F44336',
		'circle-stroke-width': 2,
		'circle-stroke-color': '#ffffff'
	}
}

const MapSourceID_Route = "source-route"
const MapLayerID_Route = "layer-route"
const MapLayer_Route = <AddLayerObject>{
	id: MapLayerID_Route,
	source: MapSourceID_Route,
	type: 'line',
	layout: {
		'line-join': 'round',
		'line-cap': 'round'
	},
	// paint: {
	// 	'line-color': '#3887be',
	// 	'line-width': 5,
	// 	'line-opacity': 0.75
	// }
	// paint: {
	// 	"line-width": ["interpolate", ["exponential", 1.5], ["zoom"], 5, 2, 18, 3],
	// 	"line-color": "#4D93E3",
	// 	// "line-gap-width": ["interpolate", ["exponential", 1.5], ["zoom"], 5, 3, 18, 8],
	// },
	paint: {
		'line-color': [ // => using feature-state expression, that checks feature.properties.road_class value
			'match',
			['get', 'road_class'],
			'motorway',
			'#009933',
			'trunk',
			'#00cc99',
			'primary',
			'#009999',
			'secondary',
			'#00ccff',
			'tertiary',
			'#9999ff',
			'residential',
			'#9933ff',
			'service_other',
			'#ffcc66',
			'unclassified',
			'#666699',
			/* other */
			'#666699'
		],
		'line-width': 8
	}
}


export function useRouteDisplay(map: MaplibreMap) {
	const updateRouteLayer = (start: { lat: number, lon: number }, end: { lat: number, lon: number }, route: [number, number][]) => {
		removeRouteLayer()

		map.addSource(MapSourceID_StartPoint, {
			type: 'geojson',
			data: {
				type: 'Feature',
				geometry: {
					type: 'Point',
					coordinates: [start.lat, start.lon]
				},
				properties: {
					title: 'Start'
				}
			}
		});
		map.addLayer(MapLayer_StartPoint);

		map.addSource(MapSourceID_EndPoint, {
			type: 'geojson',
			data: {
				type: 'Feature',
				geometry: {
					type: 'Point',
					coordinates: [end.lat, end.lon]
				},
				properties: {
					title: 'End'
				}
			}
		});
		map.addLayer(MapLayer_EndPoint);

		map.addSource(MapSourceID_Route, {
			type: "geojson",
			data: {
				type: "FeatureCollection",
				features: [
					{
						type: 'Feature',
						properties: {},
						geometry: {
							type: 'LineString',
							coordinates: route
						}
					}
				]
			}
		})
		map.addLayer(MapLayer_Route);
	}

	const removeRouteLayer = () => {
		if (map.getLayer(MapLayerID_Route)) {
			map.removeLayer(MapLayerID_StartPoint)
			map.removeLayer(MapLayerID_EndPoint)
			map.removeLayer(MapLayerID_Route);
		}
		if (map.getSource(MapSourceID_Route)) {
			map.removeSource(MapSourceID_StartPoint)
			map.removeSource(MapSourceID_EndPoint)
			map.removeSource(MapSourceID_Route)
		}

	}

	const routePath = useStore($routePath);
	watch(
		routePath,
		(newRoute) => {
			console.log("[DEBUG] new route useRouteDisplay")
			if (!newRoute) {
				removeRouteLayer()
				return
			}
			const start = $startPoint.get()
			const end = $endPoint.get()
			if (start && end && newRoute && newRoute.paths.length > 0) {
				updateRouteLayer({ lat: start.lat, lon: start.lon }, { lat: end.lat, lon: end.lon }, newRoute.paths[0].points.coordinates)
			}
		},
		{ deep: true }
	)
	return { removeRouteLayer }
}
