import createClient from "openapi-fetch";
import type { paths } from "../types/generated/api.geolocation.ts";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8100/api/v1";

export const geolocationClient = createClient<paths>({
	baseUrl: `${API_BASE}`,
});
