import createClient from "openapi-fetch";
import type { paths as OrderPaths } from "../types/generated/api.order.ts";

const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8100/api/v1";

export const orderClient = createClient<OrderPaths>({
	baseUrl: `${API_BASE}`,
	credentials: "include",
});
