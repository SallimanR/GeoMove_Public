<script setup lang="ts">
import { ref, computed, onUnmounted } from "vue";
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

const now = ref(Date.now());
const expiryMs = 30 * 60 * 1000;
const remainingSeconds = computed(() => {
  if (props.order?.status !== "pending") return null;
  const created = new Date(props.order.created_at).getTime();
  const remaining = expiryMs - (now.value - created);
  return Math.max(0, Math.floor(remaining / 1000));
});
const remainingDisplay = computed(() => {
  const s = remainingSeconds.value;
  if (s == null) return null;
  const min = Math.floor(s / 60);
  const sec = s % 60;
  return `${min}:${String(sec).padStart(2, "0")}`;
});

let timer: ReturnType<typeof setInterval> | null = null;
if (props.order?.status === "pending") {
  timer = setInterval(() => {
    now.value = Date.now();
  }, 1000);
}
onUnmounted(() => {
  if (timer) clearInterval(timer);
});

const statusLabels: Record<string, string> = {
  pending: "Ожидает водителя",
  accepted: "Принят",
  in_progress: "В пути",
  completed: "Завершён",
  cancelled: "Отменён",
};

const canEdit = props.order?.status === "pending";

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
      <span class="rounded-full bg-blue-100 px-2 py-0.5 text-sm text-blue-700">
        {{ statusLabels[order.status] ?? order.status }}
      </span>
    </div>

    <div
      v-if="remainingDisplay && remainingSeconds && remainingSeconds <= 300"
      class="rounded-lg bg-orange-100 p-2 text-center text-sm font-medium text-orange-700"
    >
      Заказ будет автоматически отменён через {{ remainingDisplay }}, если водитель не примет его
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

    <div class="flex flex-col gap-1 rounded-xl bg-gray-100 p-3 text-sm">
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
      class="max-h-48 w-full rounded-lg object-contain"
    />

    <div v-if="submitError" class="rounded-xl bg-red-100 p-3 text-center text-red-600">
      {{ submitError }}
    </div>

    <div
      v-if="order.status === 'cancelled' && (order as any).cancellation_reason"
      class="rounded-lg bg-red-50 p-3 text-sm text-red-700"
    >
      <strong>Причина отмены:</strong> {{ (order as any).cancellation_reason }}
    </div>

    <template v-if="canEdit">
      <Button @click="emit('edit')" class="w-full !border-blue-500 !bg-blue-500">
        Редактировать
      </Button>

      <Button
        :loading="isCancelling"
        @click="handleCancel"
        class="w-full !border-red-500 !bg-red-500"
      >
        {{ isCancelling ? "Отмена..." : "Отменить заказ" }}
      </Button>
    </template>
  </div>
</template>
