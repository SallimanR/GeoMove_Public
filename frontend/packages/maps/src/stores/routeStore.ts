import { atom } from "nanostores";
import type { GeoPoint } from "../types/geoPoint";
import type { RouteResponse } from "geo";

export const $startPoint = atom<GeoPoint | null>(null);
export const $endPoint = atom<GeoPoint | null>(null);
export const $routePath = atom<RouteResponse | null>(null);
export const $isRouteLoading = atom<boolean>(false);

export function setRoutePoints(start: GeoPoint, end: GeoPoint) {
	$startPoint.set(start);
	$endPoint.set(end)
}

export function clearRouteStores() {
	$startPoint.set(null)
	$endPoint.set(null)
	$routePath.set(null)
}
