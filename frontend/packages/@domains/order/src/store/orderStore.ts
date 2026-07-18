import { atom } from "nanostores";
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
}
export const $pendingEdit = atom<PendingEdit | null>(null);
export const $editDialogOpen = atom(false);

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
