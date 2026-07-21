<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { orderClient } from "order/api/client.ts";
import type { Order } from "order";
import Button from "primevue/button";
import { displayDistance } from "@geomove/geo";

const availableOrders = ref<Order[]>([]);
const myOrders = ref<Order[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const acceptingId = ref<number | null>(null);

async function fetchOrders() {
  loading.value = true;
  error.value = null;
  try {
    const [availableRes, myRes] = await Promise.all([
      orderClient.GET("/order/available"),
      orderClient.GET("/order/my", { params: { query: { role: "driver" } } }),
    ]);

    if (availableRes.error) {
      error.value = availableRes.error?.error ?? "Ошибка загрузки";
      return;
    }
    availableOrders.value = availableRes.data?.orders ?? [];

    if (!myRes.error) {
      myOrders.value = myRes.data?.orders ?? [];
    }
  } catch {
    error.value = "Не удалось загрузить заказы";
  } finally {
    loading.value = false;
  }
}

const allOrders = computed(() => {
  const myIds = new Set(myOrders.value.map((o) => o.id));
  const available = availableOrders.value.filter((o) => !myIds.has(o.id));
  return [...available, ...myOrders.value].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  );
});

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

const statusLabels: Record<string, string> = {
  forming: "Формируется",
  pending: "Ожидает",
  accepted: "Принят",
  in_progress: "В пути",
  completed: "Завершён",
  cancelled: "Отменён",
};

const isMyOrder = (order: Order) => myOrders.value.some((o) => o.id === order.id);

onMounted(fetchOrders);
</script>

<template>
  <div class="p-4 flex flex-col gap-3 h-full overflow-y-auto max-h-full">
    <h3 class="font-semibold text-center">Заказы</h3>

    <div v-if="loading" class="text-center py-8 text-gray-500">Загрузка...</div>

    <div v-else-if="error" class="text-center py-8 text-red-500">{{ error }}</div>

    <div v-else-if="allOrders.length === 0" class="text-center py-8 text-gray-500">
      Нет доступных заказов
    </div>

    <div
      v-for="order in allOrders"
      :key="order.id"
      class="rounded-xl border p-3 flex flex-col gap-2"
      :class="isMyOrder(order) ? 'bg-green-50 border-green-300' : 'bg-white'"
    >
      <div class="flex items-center justify-between">
        <span class="font-medium">Заказ #{{ order.id }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full"
          :class="order.status === 'accepted' || order.status === 'in_progress' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'"
        >
          {{ statusLabels[order.status] ?? order.status }}
        </span>
      </div>

      <div class="text-sm flex flex-col gap-1">
        <div><strong>Откуда:</strong> {{ order.from_address ?? `${order.from_lat}, ${order.from_lon}` }}</div>
        <div><strong>Куда:</strong> {{ order.to_address ?? `${order.to_lat}, ${order.to_lon}` }}</div>
        <div v-if="order.total_distance_meters">
          <strong>Расстояние:</strong> {{ displayDistance(order.total_distance_meters) }}
        </div>
        <div><strong>Авто:</strong> {{ (order as any).car_name }} ({{ (order as any).car_type }})</div>
        <div><strong>Вес:</strong> {{ (order as any).car_weight_kg }} кг, <strong>длина:</strong> {{ (order as any).car_length_meters }} м</div>
        <div><strong>Колёс заблокировано:</strong> {{ order.how_many_wheels_blocked }}</div>
        <div v-if="(order as any).customer_message" class="text-gray-500 italic">
          {{ (order as any).customer_message }}
        </div>
      </div>

      <Button
        v-if="!isMyOrder(order)"
        :loading="acceptingId === order.id"
        :disabled="acceptingId !== null"
        @click="acceptOrder(order.id)"
        class="w-full !bg-green-500 !border-green-500"
      >
        {{ acceptingId === order.id ? "Принятие..." : "Принять заказ" }}
      </Button>
    </div>

    <Button class="w-full !bg-gray-200 !border-gray-200 !text-gray-700" @click="fetchOrders">
      Обновить
    </Button>
  </div>
</template>
