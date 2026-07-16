import type { components, operations } from "./generated/api.freely_available.ts";

export type Location = components["schemas"]["Location"];

export type FreelyAvailable = components["schemas"]["FreelyAvailable"];
export type FreelyAvailableResponse = components["schemas"]["FreelyAvailableResponse"];
export type FreelyAvailableDriver = components["schemas"]["FreelyAvailableDriverResponse"];

export type CreateFreelyAvailableRequest = operations["createFreelyAvailable"]["requestBody"]["content"]["application/json"];
export type UpdateFreelyAvailableRequest = operations["updateFreelyAvailable"]["requestBody"]["content"]["application/json"];
export type GetFreelyAvailableByUserIdResponse = operations["getFreelyAvailableByUserId"]["responses"]["200"]["content"]["application/json"];
export type GetFreelyAvailableDriversRequest = operations["getFreelyAvailableDrivers"]["requestBody"]["content"]["application/json"];
export type GetFreelyAvailableDriversResponse = operations["getFreelyAvailableDrivers"]["responses"]["200"]["content"]["application/json"];
