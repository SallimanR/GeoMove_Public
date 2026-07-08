import { atom, map } from "nanostores";

import { Map as MaplibreMap } from "maplibre-gl";
import { type MapboxOverlay } from "@deck.gl/mapbox";

export const $coords = map<{ center: { lat: number, lon: number } }>()

export const $mapInstance = atom<MaplibreMap | null>(null);

export const $deckOverlay = atom<MapboxOverlay | null>(null)
