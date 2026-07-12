import { atom, map } from "nanostores";

import { Map as MaplibreMap } from "maplibre-gl";
import { type MapboxOverlay } from "@deck.gl/mapbox";
import type { GeoPoint } from "../types/geoPoint";
import { type SearchResult } from "geo";

export const $coords = map<{ center: { lat: number, lon: number } }>()
export const $mapCenterAddress = atom<SearchResult | null>(null);
export const $mapCenterAddressText = atom<string>("");

export const $mapInstance = atom<MaplibreMap | null>(null);

export const $deckOverlay = atom<MapboxOverlay | null>(null)

export const $locationPicking = atom<boolean>(false);

let pickCallback: ((point: GeoPoint, address: string) => void) | null = null;

export function setPickCallback(cb: (point: GeoPoint, address: string) => void) {
	pickCallback = cb;
}

export function invokePickCallback(point: GeoPoint, address: string) {
	pickCallback?.(point, address);
	pickCallback = null;
}

export function clearPickCallback() {
	pickCallback = null;
}
