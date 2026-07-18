import { type AddLayerObject, type Map as MaplibreMap } from "maplibre-gl";
import { $endPoint, $isRouteLoading, $routePath, $startPoint } from "../stores/routeStore";
import { fetchRoute } from "geo";
import { GeoPoint } from "../types/geoPoint";

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
	const updateRouteLayer = (route: [number, number][]) => {
		removeRouteLayer()

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
			map.removeLayer(MapLayerID_Route);
		}
		if (map.getSource(MapSourceID_Route)) {
			map.removeSource(MapSourceID_Route)
		}

	}
	async function refetchRoute(startPoint: GeoPoint, endPoint: GeoPoint) {
		if (!startPoint || !endPoint) {
			$routePath.set(null);
			return;
		}
		$isRouteLoading.set(true);
		try {
			const route = await fetchRoute(startPoint, endPoint);
			$routePath.set(route);
		} catch (err) {
			$routePath.set(null);
			console.error("route watcher error: ", err);
		} finally {
			$isRouteLoading.set(false);
		}
	}

	$startPoint.subscribe(async (startPoint) => {
		const endPoint = $endPoint.get()
		refetchRoute(startPoint as GeoPoint, endPoint as GeoPoint)
	})

	$endPoint.subscribe(async (endPoint) => {
		const startPoint = $startPoint.get()
		refetchRoute(startPoint as GeoPoint, endPoint as GeoPoint)
	})

	$routePath.subscribe(routePath => {
		if (!routePath) {
			removeRouteLayer()
			return
		}
		if (routePath.paths.length > 0) {
			updateRouteLayer(routePath.paths[0].points.coordinates)
		}
	})

	return { removeRouteLayer }
}
