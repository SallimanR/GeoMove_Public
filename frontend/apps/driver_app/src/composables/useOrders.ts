import { ref, computed } from "vue";
import { orderClient } from "order/api/client.ts";
import type { Order } from "order";

function createOrdersStore() {
	const availableOrders = ref<Order[]>([]);
	const myOrders = ref<Order[]>([]);
	const declinedOrders = ref<Order[]>([]);
	const loading = ref(false);
	const error = ref<string | null>(null);
	const acceptingId = ref<number | null>(null);
	const decliningId = ref<number | null>(null);
	const restoringId = ref<number | null>(null);
	const cancellingId = ref<number | null>(null);

	const allOrders = computed(() => {
		const myIds = new Set(myOrders.value.map((o) => o.id));
		const available = availableOrders.value.filter((o) => !myIds.has(o.id));
		return [...available, ...myOrders.value].sort(
			(a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
		);
	});

	function isMyOrder(order: Order): boolean {
		return myOrders.value.some((o) => o.id === order.id);
	}

	function isDeclined(order: Order): boolean {
		return declinedOrders.value.some((o) => o.id === order.id);
	}

	async function fetchOrders() {
		loading.value = true;
		error.value = null;
		try {
			const [availableRes, myRes, declinedRes] = await Promise.all([
				orderClient.GET("/order/available"),
				orderClient.GET("/order/my", { params: { query: { role: "driver" } } }),
				orderClient.GET("/order/declined"),
			]);

			if (availableRes.error) {
				error.value = availableRes.error?.error ?? "Ошибка загрузки";
			} else {
				availableOrders.value = availableRes.data?.orders ?? [];
			}

			if (!myRes.error) {
				myOrders.value = myRes.data?.orders ?? [];
			}

			if (!declinedRes.error) {
				declinedOrders.value = declinedRes.data?.orders ?? [];
			}
		} catch {
			error.value = "Не удалось загрузить заказы";
		} finally {
			loading.value = false;
		}
	}

	async function acceptOrder(orderId: number) {
		acceptingId.value = orderId;
		error.value = null;
		try {
			const { data, error: err } = await orderClient.PATCH("/order/{order_id}/status", {
				params: { path: { order_id: orderId } },
				body: { status: "accepted" },
			});
			if (err) {
				error.value = err?.error ?? "Ошибка при принятии заказа";
			} else {
				const accepted = availableOrders.value.find((o) => o.id === orderId);
				if (accepted && data) {
					myOrders.value.push({ ...accepted, ...data, status: "accepted" });
				}
				availableOrders.value = availableOrders.value.filter((o) => o.id !== orderId);
			}
		} catch {
			error.value = "Не удалось принять заказ";
		} finally {
			acceptingId.value = null;
		}
	}

	async function declineOrder(orderId: number) {
		decliningId.value = orderId;
		error.value = null;
		try {
			const { error: err } = await orderClient.PATCH("/order/{order_id}/decline", {
				params: { path: { order_id: orderId } },
			});
			if (err) {
				error.value = err?.error ?? "Ошибка при отказе от заказа";
			} else {
				const declined = availableOrders.value.find((o) => o.id === orderId);
				if (declined) {
					declinedOrders.value.unshift(declined);
				}
				availableOrders.value = availableOrders.value.filter((o) => o.id !== orderId);
			}
		} catch {
			error.value = "Не удалось отказаться от заказа";
		} finally {
			decliningId.value = null;
		}
	}

	async function undoDeclineOrder(orderId: number) {
		restoringId.value = orderId;
		error.value = null;
		try {
			const { error: err } = await orderClient.DELETE("/order/{order_id}/decline", {
				params: { path: { order_id: orderId } },
			});
			if (err) {
				error.value = err?.error ?? "Ошибка при восстановлении заказа";
			} else {
				const restored = declinedOrders.value.find((o) => o.id === orderId);
				if (restored) {
					availableOrders.value.unshift(restored);
				}
				declinedOrders.value = declinedOrders.value.filter((o) => o.id !== orderId);
			}
		} catch {
			error.value = "Не удалось восстановить заказ";
		} finally {
			restoringId.value = null;
		}
	}

	async function cancelOrder(orderId: number, reason: string) {
		cancellingId.value = orderId;
		error.value = null;
		try {
			const { data, error: err } = await orderClient.PATCH("/order/{order_id}/status", {
				params: { path: { order_id: orderId } },
				body: { status: "pending", cancellation_reason: reason },
			});
			if (err) {
				error.value = err?.error ?? "Ошибка при отмене заказа";
			} else if (data) {
				myOrders.value = myOrders.value.filter((o) => o.id !== orderId);
			}
		} catch {
			error.value = "Не удалось отменить заказ";
		} finally {
			cancellingId.value = null;
		}
	}

	return {
		availableOrders,
		myOrders,
		declinedOrders,
		allOrders,
		loading,
		error,
		acceptingId,
		decliningId,
		restoringId,
		cancellingId,
		isMyOrder,
		isDeclined,
		fetchOrders,
		acceptOrder,
		declineOrder,
		undoDeclineOrder,
		cancelOrder,
	};
}

let _instance: ReturnType<typeof createOrdersStore> | null = null;

export function useOrders() {
	if (!_instance) {
		_instance = createOrdersStore();
	}
	return _instance;
}
