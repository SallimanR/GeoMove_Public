import { addPopupToMap, removePopupsByGroup, $mapInstance } from "@geomove/maps";
import { $driverStore, $freelyAvailableDriverStore } from "driver";
import DriverMapsPopup from "../components/Tabs/MapsTab/DriverMapsPopup.vue";
import FreelyAvailableDriverCard from "../components/Tabs/MapsTab/FreelyAvailableDriverCard.vue";
import DriverClusterPopup from "../components/Tabs/MapsTab/DriverClusterPopup.vue";

const CLUSTER_ZOOM = 11;
const CLUSTER_GROUP = "driver-clusters";
const FA_CLUSTER_GROUP = "fa-driver-clusters";

let zoomSub: (() => void) | undefined;
let zoomUnsub: (() => void) | undefined;
let currentZoom = 12;

let redrawCallback: (() => void) | undefined;

export function initDriverPopupZoom(onRedraw: () => void) {
	redrawCallback = onRedraw;
	if (zoomSub) return;
	zoomSub = $mapInstance.subscribe((map) => {
		if (zoomUnsub) zoomUnsub();
		if (!map) return;
		currentZoom = map.getZoom();
		map.on("zoom", () => {
			const z = map.getZoom();
			const wasClustering = currentZoom < CLUSTER_ZOOM;
			const nowClustering = z < CLUSTER_ZOOM;
			currentZoom = z;
			if (wasClustering !== nowClustering && redrawCallback) {
				redrawCallback();
			}
		});
		zoomUnsub = () => {
			map.off("zoom", () => { });
		};
		if (redrawCallback) redrawCallback();
	});
}

function isClustering() {
	return currentZoom < CLUSTER_ZOOM;
}

function zoomTo(lat: number, lon: number) {
	const map = $mapInstance.get();
	if (!map) return;
	map.flyTo({ center: [lon, lat], zoom: CLUSTER_ZOOM + 1 });
}

function clusterAndDisplay(
	group: string,
	clusterGroup: string,
	items: { lat: number; lon: number; data: any }[],
	Component: any,
	propsFn: (item: any) => any,
	clusterColor = "bg-blue-500",
) {
	removePopupsByGroup(group);
	removePopupsByGroup(clusterGroup);

	if (!isClustering()) {
		for (const item of items) {
			addPopupToMap(item.lat, item.lon, Component, propsFn(item.data), group);
		}
		return;
	}

	const map = $mapInstance.get();
	if (!map) return;

	const threshold = Math.max(40, 150 - (map.getZoom() - 5) * 22);

	const projected = items.map((item) => {
		const p = map.project([item.lon, item.lat]);
		return { item, x: p.x, y: p.y };
	});

	const clusters: { lat: number; lon: number; count: number; pts: typeof projected }[] = [];

	for (const pt of projected) {
		let found = false;
		for (const cl of clusters) {
			const cx = cl.pts.reduce((s, p) => s + p.x, 0) / cl.pts.length;
			const cy = cl.pts.reduce((s, p) => s + p.y, 0) / cl.pts.length;
			const dx = pt.x - cx;
			const dy = pt.y - cy;
			if (Math.sqrt(dx * dx + dy * dy) < threshold) {
				cl.lat = (cl.lat * cl.count + pt.item.lat) / (cl.count + 1);
				cl.lon = (cl.lon * cl.count + pt.item.lon) / (cl.count + 1);
				cl.count++;
				cl.pts.push(pt);
				found = true;
				break;
			}
		}
		if (!found) {
			clusters.push({
				lat: pt.item.lat,
				lon: pt.item.lon,
				count: 1,
				pts: [pt],
			});
		}
	}

	for (const cl of clusters) {
		addPopupToMap(
			cl.lat,
			cl.lon,
			DriverClusterPopup,
			{ count: cl.count, color: clusterColor, onZoomIn: () => zoomTo(cl.lat, cl.lon) },
			clusterGroup,
		);
	}
}

export function displayDriverPopups() {
	const drivers = $driverStore.get().map((d) => ({ lat: d.lat, lon: d.lon, data: d }));
	clusterAndDisplay("drivers", CLUSTER_GROUP, drivers, DriverMapsPopup, (d) => d);
}

export function displayFreelyAvailableDriverPopups() {
	const drivers = $freelyAvailableDriverStore.get().map((d) => ({
		lat: d.from_location.lat,
		lon: d.from_location.lon,
		data: d,
	}));
	clusterAndDisplay(
		"fa-drivers",
		FA_CLUSTER_GROUP,
		drivers,
		FreelyAvailableDriverCard,
		(d) => ({
			...d,
			asMapPopup: true,
		}),
		"bg-orange-500",
	);
}
