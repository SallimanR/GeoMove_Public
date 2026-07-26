import createClient from "openapi-fetch";
import type { paths as PriceCalcPaths } from "../types/generated/api.price_calculation.ts";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8100/api/v1";

export const priceCalcClient = createClient<PriceCalcPaths>({
	baseUrl: `${API_BASE}`,
	credentials: "include",
});
