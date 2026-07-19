<script setup lang="ts">
import { computed } from "vue";
import type { Order } from "order/types/order.ts";
import { displayDistance } from "@geomove/geo";

import EditOrder from "./EditOrder.vue";

const props = defineProps<{
  orders: Order[];
}>();

interface OrderWithAddresses {
  order: Order;
  fromAddress: string;
  toAddress: string;
}

const items = computed<OrderWithAddresses[]>(() =>
  props.orders.map((order) => ({
    order,
    fromAddress: order.from_address || `${order.from_lat.toFixed(5)}, ${order.from_lon.toFixed(5)}`,
    toAddress: order.to_address || `${order.to_lat.toFixed(5)}, ${order.to_lon.toFixed(5)}`,
  })),
);

const statusLabel: Record<string, string> = {
  forming: "Формируется",
  pending: "Ожидает",
  accepted: "Принят",
  in_progress: "В пути",
  completed: "Завершён",
  cancelled: "Отменён",
};

const statusColor: Record<string, string> = {
  forming: "bg-gray-300",
  pending: "bg-yellow-300",
  accepted: "bg-blue-400",
  in_progress: "bg-orange-400",
  completed: "bg-green-400",
  cancelled: "bg-red-400",
};
</script>

<template>
  <div class="flex flex-col gap-2">
    <h3 class="font-semibold text-lg text-gray-800">Мои заказы</h3>

    <div
      v-for="item in items"
      :key="item.order.id"
      class="rounded-xl bg-gray-100 p-3 flex flex-col gap-2"
    >
      <div class="flex items-center justify-between">
        <span class="text-sm text-gray-500">
          {{ new Date(item.order.created_at).toLocaleString("ru") }}
        </span>
        <span
          class="rounded-full px-3 py-0.5 text-xs font-medium"
          :class="statusColor[item.order.status] ?? 'bg-gray-300'"
        >
          {{ statusLabel[item.order.status] ?? item.order.status }}
        </span>
      </div>

      <div class="flex flex-col gap-1 text-sm">
        <div class="flex gap-1">
          <span class="text-gray-400 shrink-0">От:</span>
          <span class="text-gray-700 truncate">{{ item.fromAddress }}</span>
        </div>
        <div class="flex gap-1">
          <span class="text-gray-400 shrink-0">До:</span>
          <span class="text-gray-700 truncate">{{ item.toAddress }}</span>
        </div>
      </div>

      <div class="text-sm text-gray-600">
        Колёс заблокировано: {{ item.order.how_many_wheels_blocked }}
        <div v-if="item.order.total_distance_meters">
          Расстояние: {{ displayDistance(item.order.total_distance_meters) }}
        </div>
        <div v-if="item.order.price_rubles">
          Цена: {{ item.order.price_rubles }} ₽
        </div>
      </div>

      <EditOrder :order="item.order" />

      <div
        v-if="item.order.status === 'cancelled' && item.order.cancellation_reason"
        class="text-xs text-red-500"
      >
        Причина: {{ item.order.cancellation_reason }}
      </div>
    </div>
  </div>
</template>
