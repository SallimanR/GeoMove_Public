import { atom, onMount } from "nanostores";
import type { Order } from "../types/order.ts";

export const $orders = atom<Order[]>([]);
export const $currentOrder = atom<Order | null>(null);

export interface LocationEntry {
	lat: number;
	lon: number;
	address: string;
}

export interface PendingEdit {
	orderId: number;
	wheels: number;
	price: number;
	from: LocationEntry;
	to: LocationEntry;
	carWeightKg: number;
	carLengthMeters: number;
	carType: string;
	carName: string;
	carPhotoUrl: string;
	customerMessage: string;
}
export const $pendingEdit = atom<PendingEdit | null>(null);
export const $editDialogOpen = atom(false);

export interface CreateOrderForm {
	wheels: number;
	carWeightKg: number;
	carLengthMeters: number;
	carType: string;
	carName: string;
	carPhotoUrl: string;
	customerMessage: string;
}

function loadForm(): CreateOrderForm {
	try {
		const saved = localStorage.getItem("createOrderForm");
		if (saved) return JSON.parse(saved);
	} catch { /* ignore */ }
	return {
		wheels: 0,
		carWeightKg: 0,
		carLengthMeters: 0,
		carType: "Другое",
		carName: "",
		carPhotoUrl: "",
		customerMessage: "",
	};
}

export const $createOrderForm = atom<CreateOrderForm>(loadForm());

onMount($createOrderForm, () => {
	$createOrderForm.subscribe((val) => {
		const { carPhotoUrl, ...rest } = val;
		localStorage.setItem("createOrderForm", JSON.stringify(rest));
	});
});

interface OrderRoute {
	fromLat: number;
	fromLon: number;
	fromText: string;
	toLat: number;
	toLon: number;
	toText: string;
}

function loadRoute(): OrderRoute | null {
	try {
		const saved = localStorage.getItem("orderRoute");
		if (saved) return JSON.parse(saved);
	} catch { /* ignore */ }
	return null;
}

export const $orderRoute = atom<OrderRoute | null>(loadRoute());

onMount($orderRoute, () => {
	$orderRoute.subscribe((val) => {
		if (val) {
			localStorage.setItem("orderRoute", JSON.stringify(val));
		} else {
			localStorage.removeItem("orderRoute");
		}
	});
});

export type ActiveTabSetter = (tab: string) => void;
export const $activeTabSetter = atom<ActiveTabSetter | null>(null);

export function setOrders(orders: Order[]): void {
	$orders.set(orders);
}

export function addOrder(order: Order): void {
	const current = $orders.get();
	$orders.set([order, ...current]);
}

export function updateOrderInStore(updated: Order): void {
	const current = $orders.get();
	$orders.set(current.map((o) => (o.id === updated.id ? updated : o)));
}
