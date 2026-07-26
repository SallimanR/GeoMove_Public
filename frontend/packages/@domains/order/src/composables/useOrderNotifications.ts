import { onMounted, onUnmounted } from "vue";

export function useOrderNotifications(onNotification: () => void) {
	const bc = new BroadcastChannel("order-updates");

	onMounted(() => {
		try {
			bc.onmessage = () => {
				onNotification();
			};
		} catch {
			/* noop */
		}
	});

	onUnmounted(() => {
		bc?.close();
	});
}
