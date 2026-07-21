import createClient from "openapi-fetch";
import type { paths as DriverPaths } from "../types/generated/api.driver.ts";
import type { paths as FaPaths } from "../types/generated/api.freely_available.ts";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8100/api/v1";

export const driverClient = createClient<DriverPaths>({
	baseUrl: `${API_BASE}`,
	credentials: "include",
});

export const freelyAvailableDriverClient = createClient<FaPaths>({
	baseUrl: `${API_BASE}`,
	credentials: "include",
});
