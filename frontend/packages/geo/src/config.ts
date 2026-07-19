let _routingApi = ""
let _geocodingApi = ""

export interface GeoModuleConfig {
  routingApi: string
  geocodingApi: string
}

export function configureGeo(config: Partial<GeoModuleConfig>) {
  if (config.routingApi !== undefined) _routingApi = config.routingApi
  if (config.geocodingApi !== undefined) _geocodingApi = config.geocodingApi
}

export function getGeoConfig(): GeoModuleConfig {
  return { routingApi: _routingApi, geocodingApi: _geocodingApi }
}
