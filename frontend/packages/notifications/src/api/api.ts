import createClient from "openapi-fetch"
import type { paths as NotificationPaths } from "../types/generated/api.notifications.ts"

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8100/api/v1"

export const notificationClient = createClient<NotificationPaths>({
	baseUrl: API_BASE,
	credentials: "include",
})

export async function getVapidPublicKey(): Promise<string | null> {
	const { data, error } = await notificationClient.GET(
		"/notifications/vapid-public-key",
	)
	if (error || !data?.publicKey) {
		console.error("Failed to get VAPID public key:", error)
		return null
	}
	return data.publicKey
}

export async function subscribe(
	subscription: PushSubscriptionJSON,
	deviceType: string = "web",
): Promise<boolean> {
	const { error } = await notificationClient.POST("/notifications/subscribe", {
		body: {
			endpoint: subscription.endpoint ?? "",
			devicePublicKey: subscription.keys?.p256dh ?? "",
			authSecret: subscription.keys?.auth ?? "",
			deviceType,
		},
	})
	if (error) {
		console.error("Failed to save push subscription:", error)
		return false
	}
	return true
}

export async function unsubscribe(endpoint: string): Promise<boolean> {
	const { error } = await notificationClient.POST(
		"/notifications/unsubscribe",
		{
			body: { endpoint },
		},
	)
	if (error) {
		console.error("Failed to unsubscribe:", error)
		return false
	}
	return true
}
