import createClient from "openapi-fetch";
import type { paths as DriverPaths } from "../types/generated/api.driver.ts";
import type { paths as FaPaths } from "../types/generated/api.freely_available.ts";

export const driverClient = createClient<DriverPaths>({
	baseUrl: "http://localhost:8100/api/v1",
});

export const freelyAvailableDriverClient = createClient<FaPaths>({
	baseUrl: "http://localhost:8100/api/v1",
});
