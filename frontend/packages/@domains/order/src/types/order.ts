import type { components, operations } from "./generated/api.order.ts";

export type Order = components["schemas"]["Order"];

export type CreateOrderRequest = operations["createOrder"]["requestBody"]["content"]["application/json"];
export type CreateOrderResponse = operations["createOrder"]["responses"]["201"]["content"]["application/json"];

export type ListMyOrdersResponse = operations["listMyOrders"]["responses"]["200"]["content"]["application/json"];

export type UpdateOrderStatusRequest = operations["updateOrderStatus"]["requestBody"]["content"]["application/json"];
export type UpdateOrderStatusResponse = operations["updateOrderStatus"]["responses"]["200"]["content"]["application/json"];
