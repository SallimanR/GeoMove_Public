import { atom } from "nanostores";
import type { GeoPoint } from "../types/geoPoint";
import type { RouteResponse, SearchResult } from "geo";

export const $startPoint = atom<GeoPoint | null>(null);
export const $endPoint = atom<GeoPoint | null>(null);
export const $routePath = atom<RouteResponse | null>(null);
export const $isRouteLoading = atom<boolean>(false);

export const $startAddress = atom<SearchResult | null>(null);
export const $endAddress = atom<SearchResult | null>(null);
