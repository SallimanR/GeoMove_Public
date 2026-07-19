# @geomove/maps

Vue 3 + Maplibre GL JS map components with routing, search, geocoding, and animated icon layers. Built on top of [@geomove/geo](https://www.npmjs.com/package/@geomove/geo).

## Install

```bash
pnpm add @geomove/maps @geomove/geo vue maplibre-gl @deck.gl/core @deck.gl/layers @deck.gl/mapbox nanostores @nanostores/vue pmtiles
```

## Import CSS

```ts
import '@geomove/maps/dist/maps.css'
```

This includes pre-compiled Tailwind utilities, Maplibre GL styles, and component styles. No Tailwind setup required on the consumer side.

## Quick start

```vue
<script setup lang="ts">
import { configureGeo } from '@geomove/geo'
import { Maps, MapsOverlayControls } from '@geomove/maps'

configureGeo({
  routingApi: 'https://your-api.example.com/routing',
  geocodingApi: 'https://your-api.example.com/geocoding',
})

const styleApi = 'https://api.protomaps.com/styles/v5/light/en.json?key=YOUR_KEY'
</script>

<template>
  <div style="width: 100%; height: 100vh; position: relative;">
    <Maps :styleApi="styleApi" />
    <MapsOverlayControls />
  </div>
</template>
```

## Components

| Component | Description |
|---|---|
| `Maps` | Root map container. Prop: `styleApi` — MapLibre style URL |
| `MapsOverlayControls` | Full overlay — search, route input, GPS, location picker |
| `MapsLocationPicker` | Map point picker with coordinates callback |

## Composables

| Export | Description |
|---|---|
| `useMovingIconLayer(options)` | Deck.gl animated icon layer for moving objects |
| `useMovingIconLayerMaplibre(options)` | Maplibre GL native animated symbol layer |

### useMovingIconLayer options

```ts
interface UseMovingIconLayerOptions {
  paths: ReadableAtom<MovingPath[]>
  iconUrl: string
  iconWidth: number
  iconHeight: number
  layerId?: string
  iconSize?: number
  popupComponent?: Component
  onClick?: (id: number) => void
  onHover?: (id: number) => void
}
```

## Stores (nanostores)

| Export | Description |
|---|---|
| `$coords` | Current map center `{ center: { lat, lon } }` |
| `$mapInstance` | MapLibreMap instance |
| `$deckOverlay` | Deck.gl MapboxOverlay instance |
| `$locationPicking` | Location picker active state |
| `$startPoint` / `$endPoint` | Route start/end coordinates |
| `$startAddress` / `$endAddress` | Route start/end address |
| `$routePath` | Current route response |

## Types

| Export | Description |
|---|---|
| `GeoPoint` | `{ lat: number; lon: number }` |
| `MovingPath` | Path with coordinates, time, distance, id |
| `MovingPosition` | Position with id and bearing |
| `UseMovingIconLayerOptions` | Config for moving icon layers |

## Utility functions

| Export | Description |
|---|---|
| `addPopupToMap(lat, lon, component, props, group, options)` | Mount a Vue component as a Maplibre popup |
| `removeAllPopups()` | Remove all popups |
| `removePopupsByGroup(group)` | Remove popups by group key |

## Peer dependencies

| Package | Version | Required |
|---|---|---|
| vue | ^3.5 | ✓ |
| maplibre-gl | ^5.0 | ✓ |
| @deck.gl/core | ^9.0 | ✓ |
| @deck.gl/layers | ^9.0 | ✓ |
| @deck.gl/mapbox | ^9.0 | ✓ |
| nanostores | ^1.0 | ✓ |
| @nanostores/vue | ^1.0 | ✓ |
| pmtiles | ^4.0 | ✓ |
| @capacitor/core | ^8.2 | optional (GPS) |
| @capacitor/geolocation | ^8.2 | optional (GPS) |
