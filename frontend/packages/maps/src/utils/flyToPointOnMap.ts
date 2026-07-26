import { $mapInstance } from "../build-entry";

export function flyToPointOnMap(lat: number, lon: number) {
	const map = $mapInstance.get();
	if (!map) return;

	map.flyTo({ center: [lon, lat], zoom: 15 });
}
