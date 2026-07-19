import { ref, onMounted } from "vue"
import {
	registerServiceWorker,
	requestNotificationPermission,
} from "../register.ts"
import { getVapidPublicKey, subscribe, unsubscribe } from "../api/api.ts"
import type { UsePushNotificationsOptions } from "../register.ts"

export function usePushNotifications(options: UsePushNotificationsOptions = {}) {
	const { notificationWorkerPath = "/notification-worker.js", deviceType = "web" } = options

	const isSupported = ref(false)
	const isSubscribed = ref(false)
	const permissionGranted = ref(false)
	const error = ref<string | null>(null)

	onMounted(async () => {
		isSupported.value = "serviceWorker" in navigator && "PushManager" in window
	})

	async function init(): Promise<void> {
		if (!isSupported.value) {
			error.value = "Push уведомления не поддерживаются"
			return
		}

		const permitted = await requestNotificationPermission()
		permissionGranted.value = permitted
		if (!permitted) {
			error.value = "Нет разрешения на уведомления"
			return
		}

		const registration = await registerServiceWorker(notificationWorkerPath)
		if (!registration) {
			error.value = "Не удалось инициализировать service worker"
			return
		}

		const existingSubscription =
			await registration.pushManager.getSubscription()
		if (existingSubscription) {
			isSubscribed.value = true
			return
		}

		const vapidKey = await getVapidPublicKey()
		if (!vapidKey) {
			error.value = "Не удалось получить публичный VAPID ключ"
			return
		}

		try {
			const pushSubscription = await registration.pushManager.subscribe({
				userVisibleOnly: true,
				applicationServerKey: vapidKey,
			})

			const ok = await subscribe(pushSubscription.toJSON(), deviceType)
			if (ok) {
				isSubscribed.value = true
				error.value = null
			} else {
				error.value = "Не удалось сохранить подписку на сервере"
			}
		} catch (err) {
			console.error("Push subscription failed:", err)
			error.value = "Ошибка push уведомления"
		}
	}

	async function unsubscribeFromPush(): Promise<void> {
		const registration = await navigator.serviceWorker.ready
		const pushSubscription = await registration.pushManager.getSubscription()
		if (pushSubscription) {
			await unsubscribe(pushSubscription.endpoint)
			await pushSubscription.unsubscribe()
			isSubscribed.value = false
		}
	}

	return {
		isSupported,
		isSubscribed,
		permissionGranted,
		error,
		init,
		unsubscribeFromPush,
	}
}

// function urlBase64ToUint8Array(base64String: string): Uint8Array {
// 	const padding = "=".repeat((4 - (base64String.length % 4)) % 4)
// 	const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/")
// 	const rawData = window.atob(base64)
// 	const outputArray = new Uint8Array(rawData.length)
// 	for (let i = 0; i < rawData.length; ++i) {
// 		outputArray[i] = rawData.charCodeAt(i)
// 	}
// 	return outputArray
// }
