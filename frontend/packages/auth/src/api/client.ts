import createClient from "openapi-fetch";
import type { paths as AuthPaths } from "../types/generated/api.auth.ts";

const API_BASE = import.meta.env.VITE_AUTH_API_BASE || "http://localhost:8100/api/v1";

export const authClient = createClient<AuthPaths>({
	baseUrl: `${API_BASE}`,
	credentials: "include",
});
