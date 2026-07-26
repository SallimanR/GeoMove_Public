export type { UsePushNotificationsOptions } from "./register.ts";
export { registerServiceWorker } from "./register.ts";
export { usePushNotifications } from "./composables/usePushNotifications.ts";
export { notificationClient, getVapidPublicKey, subscribe, unsubscribe } from "./api/api.ts";
