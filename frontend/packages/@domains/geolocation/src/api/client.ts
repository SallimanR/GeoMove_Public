import createClient from "openapi-fetch";
import type { paths } from "../types/generated/api.geolocation.ts";

export const geolocationClient = createClient<paths>({
	baseUrl: "http://localhost:8100/api/v1",
});
