import type { components, operations } from "./generated/api.driver.ts";
export type Driver = components["schemas"]["Driver"];

export type GetFilteredDriversRequest = operations["getFilteredDrivers"]["requestBody"]["content"]["application/json"];
export type GetFilteredDriversResponse = operations["getFilteredDrivers"]["responses"]["200"]["content"]["application/json"];

