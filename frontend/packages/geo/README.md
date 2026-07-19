# @geomove/geo

Geospatial utility library — geocoding, reverse geocoding, routing, and distance helpers.

## Install

```bash
pnpm add @geomove/geo
```

## Quick start

```ts
import { configureGeo, haversineDistance, displayDistance } from '@geomove/geo'

configureGeo({
  geocodingApi: 'https://your-api.example.com/geocoding',
  routingApi: 'https://your-api.example.com/routing',
})

console.log(haversineDistance([37.62, 55.75], [30.31, 59.93]))
console.log(displayDistance(1500))
```

## API

### Configuration

| Export | Description |
|---|---|
| `configureGeo(config)` | Set API endpoints before using geocoding/routing |
| `GeoModuleConfig` | `{ routingApi: string; geocodingApi: string }` |

### Geocoding

| Export | Type |
|---|---|
| `getMapSearch(query, lat, lon)` | `Promise<SearchResultList>` |
| `getReverseGeocoding(lat, lon)` | `Promise<SearchResult>` |
| `SearchResult` | Interface — geometry + address properties |
| `SearchResultList` | `{ features: SearchResult[] }` |

### Routing

| Export | Type |
|---|---|
| `fetchRoute(start, end)` | `Promise<RouteResponse>` |
| `RouteResponse` | `{ paths: Array<{ distance: number; points: { coordinates: [number, number][] } }> }` |
| | `start` / `end` are `{ lat: number; lon: number }` |

### Geometry

| Export | Type |
|---|---|
| `haversineDistance(coord1, coord2)` | `(number, number) => number` — distance in meters |
| `degreesToRadians(degrees)` | `(number) => number` |

### Utilities

| Export | Description |
|---|---|
| `displayDistance(distance)` | `(number) => string` — `"500 м"` / `"1.5 км"` |
| `addressToText(address)` | `(SearchResult) => string` — "Ленина, 1, Москва" |
