import { atom } from "nanostores";
import { authClient } from "../api/client.ts";
import type { AuthUser } from "../types/auth.ts";

export const $user = atom<AuthUser | null>(null);
export const $isAuthenticated = atom(false);
export const $loading = atom(false);
export const $error = atom<string | null>(null);

export async function checkAuth() {
	$loading.set(true);
	$error.set(null);
	try {
		const { data, error: apiError } = await authClient.GET("/auth/me");
		if (apiError || !data?.user) {
			$user.set(null);
			$isAuthenticated.set(false);
			return;
		}
		$user.set(data.user as AuthUser);
		$isAuthenticated.set(true);
	} catch (err) {
		$user.set(null);
		$isAuthenticated.set(false);
		$error.set(err instanceof Error ? err.message : "Failed to check auth");
	} finally {
		$loading.set(false);
	}
}

export function setUser(u: NonNullable<AuthUser>) {
	$user.set(u);
	$isAuthenticated.set(true);
}

export function clearUser() {
	$user.set(null);
	$isAuthenticated.set(false);
}
