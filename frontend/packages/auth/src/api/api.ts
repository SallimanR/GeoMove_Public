import { authClient } from "./client.ts";

export async function logout(): Promise<void> {
	await authClient.POST("/auth/logout");
}
