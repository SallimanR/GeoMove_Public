<script setup lang="ts">
import { ref } from "vue";
import { $orders } from "order/store/orderStore.ts";
import { orderClient } from "order/api/client.ts";
import type { Order } from "order";
import Button from "primevue/button";

const props = defineProps<{
  order: Order;
}>();

const emit = defineEmits<{
  edit: [];
}>();

const isCancelling = ref(false);
const submitError = ref<string | null>(null);

const statusLabels: Record<string, string> = {
  forming: "Формируется",
  pending: "Ожидает водителя",
  accepted: "Принят",
  in_progress: "В пути",
  completed: "Завершён",
  cancelled: "Отменён",
};

const canEdit = (props.order?.status === "forming" || props.order?.status === "pending");

async function handleCancel() {
  if (!props.order) return;
  isCancelling.value = true;
  submitError.value = null;
  try {
    const { error } = await orderClient.DELETE("/order/my/active");

    if (error) {
      submitError.value = error?.error ?? "Ошибка при отмене заказа";
    } else {
      $orders.set($orders.get().filter((o) => o.id !== props.order.id));
    }
  } catch {
    submitError.value = "Не удалось отменить заказ";
  } finally {
    isCancelling.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col gap-3">
    <div class="flex items-center gap-2">
      <h3 class="font-semibold">Заказ #{{ order.id }}</h3>
      <span class="text-sm px-2 py-0.5 bg-blue-100 text-blue-700 rounded-full">
        {{ statusLabels[order.status] ?? order.status }}
      </span>
    </div>

    <div class="flex gap-2 text-sm">
      <div class="flex-1 rounded-xl bg-gray-100 p-2">
        <div class="text-gray-500">Откуда</div>
        <div>{{ order.from_address ?? `${order.from_lat}, ${order.from_lon}` }}</div>
      </div>
      <div class="flex-1 rounded-xl bg-gray-100 p-2">
        <div class="text-gray-500">Куда</div>
        <div>{{ order.to_address ?? `${order.to_lat}, ${order.to_lon}` }}</div>
      </div>
    </div>

    <div class="rounded-xl bg-gray-100 p-3 text-sm flex flex-col gap-1">
      <div><strong>Авто:</strong> {{ (order as any).car_name }}</div>
      <div><strong>Тип:</strong> {{ (order as any).car_type }}</div>
      <div><strong>Вес:</strong> {{ (order as any).car_weight_kg }} кг</div>
      <div><strong>Длина:</strong> {{ (order as any).car_length_meters }} м</div>
      <div><strong>Колёс заблокировано:</strong> {{ order.how_many_wheels_blocked }}</div>
      <div v-if="(order as any).customer_message">
        <strong>Сообщение:</strong> {{ (order as any).customer_message }}
      </div>
    </div>

    <img
      v-if="(order as any).car_photo_url"
      :src="(order as any).car_photo_url"
      class="w-full max-h-48 object-contain rounded-lg"
    />

    <div
      v-if="submitError"
      class="rounded-xl bg-red-100 p-3 text-center text-red-600"
    >
      {{ submitError }}
    </div>

    <template v-if="canEdit">
      <Button
        @click="emit('edit')"
        class="w-full !bg-blue-500 !border-blue-500"
      >
        Редактировать
      </Button>

      <Button
        :loading="isCancelling"
        @click="handleCancel"
        class="w-full !bg-red-500 !border-red-500"
      >
        {{ isCancelling ? "Отмена..." : "Отменить заказ" }}
      </Button>
    </template>
  </div>
</template>
