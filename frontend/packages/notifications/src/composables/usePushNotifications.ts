import { ref, onMounted } from "vue";
import { registerServiceWorker, requestNotificationPermission } from "../register.ts";
import { getVapidPublicKey, subscribe, unsubscribe } from "../api/api.ts";
import type { UsePushNotificationsOptions } from "../register.ts";

export function usePushNotifications(options: UsePushNotificationsOptions = {}) {
	const { notificationWorkerPath = "/notification-worker.js", deviceType = "web" } = options;

	const isSupported = ref(false);
	const isSubscribed = ref(false);
	const permissionGranted = ref(false);
	const error = ref<string | null>(null);
	let initializing = false;

	onMounted(async () => {
		isSupported.value = "serviceWorker" in navigator && "PushManager" in window;
	});

	async function init(): Promise<void> {
		if (!isSupported.value) {
			error.value = "Push уведомления не поддерживаются";
			return;
		}

		if (initializing) return;
		initializing = true;

		const permitted = await requestNotificationPermission();
		permissionGranted.value = permitted;
		if (!permitted) {
			error.value = "Нет разрешения на уведомления";
			initializing = false;
			return;
		}

		const registration = await registerServiceWorker(notificationWorkerPath);
		if (!registration) {
			error.value = "Не удалось инициализировать service worker";
			initializing = false;
			return;
		}

		const existingSubscription = await registration.pushManager.getSubscription();
		if (existingSubscription) {
			isSubscribed.value = true;
			initializing = false;
			return;
		}

		const vapidKey = await getVapidPublicKey();
		if (!vapidKey) {
			error.value = "Не удалось получить публичный VAPID ключ";
			initializing = false;
			return;
		}

		try {
			const pushSubscription = await registration.pushManager.subscribe({
				userVisibleOnly: true,
				applicationServerKey: vapidKey,
			});

			const ok = await subscribe(pushSubscription.toJSON(), deviceType);
			if (ok) {
				isSubscribed.value = true;
				error.value = null;
			} else {
				error.value = "Не удалось сохранить подписку на сервере";
			}
		} catch (err) {
			console.error("Push subscription failed:", err);
			error.value = "Ошибка push уведомления";
		} finally {
			initializing = false;
		}
	}

	async function unsubscribeFromPush(): Promise<void> {
		const registration = await navigator.serviceWorker.ready;
		const pushSubscription = await registration.pushManager.getSubscription();
		if (pushSubscription) {
			await unsubscribe(pushSubscription.endpoint);
			await pushSubscription.unsubscribe();
			isSubscribed.value = false;
		}
	}

	return {
		isSupported,
		isSubscribed,
		permissionGranted,
		error,
		init,
		unsubscribeFromPush,
	};
}
