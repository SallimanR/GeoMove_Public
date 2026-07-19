export interface UsePushNotificationsOptions {
  notificationWorkerPath?: string
  deviceType?: string
}

export async function registerServiceWorker(
  notificationWorkerPath: string = "/notification-worker.js",
): Promise<ServiceWorkerRegistration | null> {
  if (!("serviceWorker" in navigator)) {
    console.warn("Service workers are not supported in this browser")
    return null
  }

  try {
    const registration = await navigator.serviceWorker.register(notificationWorkerPath)
    await navigator.serviceWorker.ready
    return registration
  } catch (err) {
    console.error("Service worker registration failed:", err)
    return null
  }
}

export async function requestNotificationPermission(): Promise<boolean> {
  if (!("Notification" in window)) {
    console.warn("Notifications are not supported in this browser")
    return false
  }

  const permission = await Notification.requestPermission()
  return permission === "granted"
}
