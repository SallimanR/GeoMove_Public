<script setup lang="ts">
import type { Order } from "order";
import { computed, ref } from "vue";
import { displayDistance } from "@geomove/geo";
import { $endPoint, $startPoint } from "@geomove/maps";

const props = withDefaults(
  defineProps<{
    order: Order;
    isMyOrder?: boolean;
    accepting?: boolean;
    declining?: boolean;
    restoring?: boolean;
    cancelling?: boolean;
    showAcceptButton?: boolean;
    showRestoreButton?: boolean;
    showCancelButton?: boolean;
    onAccept?: (orderId: number) => void;
    onDecline?: (orderId: number) => void;
    onRestore?: (orderId: number) => void;
    onCancel?: (orderId: number, reason: string) => void;
    onShowOnMap?: (lat: number, lon: number) => void;
    onClose?: () => void;
  }>(),
  {
    isMyOrder: false,
    accepting: false,
    showAcceptButton: true,
  },
);

const cancelReasons = [
  "Неисправность автомобиля",
  "Не могу добраться до клиента",
  "Клиент передумал",
  "Другая причина",
];

const showCancelReason = ref(false);
const selectedReason = ref("");
const customReason = ref("");

function getCancelReason(): string {
  if (selectedReason.value === "Другая причина") {
    return customReason.value.trim() || "Другая причина";
  }
  return selectedReason.value;
}

function handleCancel() {
  const reason = getCancelReason();
  if (!reason) return;
  if (props.onCancel) {
    props.onCancel(props.order.id, reason);
  } else {
    emit("cancel", props.order.id, reason);
  }
  showCancelReason.value = false;
  selectedReason.value = "";
  customReason.value = "";
}

function cancelAction() {
  showCancelReason.value = true;
}

const timeAgo = computed(() => {
  const diffMs = Date.now() - new Date(props.order.created_at).getTime();
  const diffSec = Math.floor(diffMs / 1000);
  if (diffSec < 60) return "только что";
  const diffMin = Math.floor(diffSec / 60);
  if (diffMin < 60) return `${diffMin} мин назад`;
  const diffHrs = Math.floor(diffMin / 60);
  if (diffHrs < 24) return `${diffHrs} ч назад`;
  const diffDays = Math.floor(diffHrs / 24);
  return `${diffDays} дн назад`;
});

const isFresh = computed(() => {
  const diffMs = Date.now() - new Date(props.order.created_at).getTime();
  return diffMs < 5 * 60 * 1000;
});

const emit = defineEmits<{
  accept: [orderId: number];
  decline: [orderId: number];
  restore: [orderId: number];
  cancel: [orderId: number, reason: string];
  showOnMap: [lat: number, lon: number];
  close: [];
}>();

function handleDecline() {
  if (props.onDecline) {
    props.onDecline(props.order.id);
  } else {
    emit("decline", props.order.id);
  }
}

function handleRestore() {
  if (props.onRestore) {
    props.onRestore(props.order.id);
  } else {
    emit("restore", props.order.id);
  }
}

function handleAccept() {
  if (props.onAccept) {
    props.onAccept(props.order.id);
  } else {
    emit("accept", props.order.id);
  }
}

function handleShowOnMap(lat: number, lon: number) {
  if (props.onShowOnMap) {
    props.onShowOnMap(lat, lon);
  } else {
    emit("showOnMap", lat, lon);
  }

  $startPoint.set({ lat: props.order.from_lat, lon: props.order.from_lon });
  $endPoint.set({ lat: props.order.to_lat, lon: props.order.to_lon });
}

function handleClose() {
  if (props.onClose) {
    props.onClose();
  } else {
    emit("close");
  }
}

const statusLabels: Record<string, string> = {
  pending: "Ожидает",
  accepted: "Принят",
  in_progress: "В пути",
  completed: "Завершён",
  cancelled: "Отменён",
};
</script>

<template>
  <div
    class="flex min-w-64 flex-col gap-2 rounded-xl border p-3"
    :class="[
      isMyOrder ? 'border-green-300 bg-green-50' : 'bg-white',
      isFresh && !isMyOrder ? 'border-l-4 border-l-blue-400' : '',
    ]"
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span class="font-medium">Заказ #{{ order.id }}</span>
        <span
          v-if="isFresh"
          class="rounded bg-blue-100 px-1.5 py-0.5 text-xs font-medium text-blue-600"
        >
          новый
        </span>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-sm text-gray-400">{{ timeAgo }}</span>
        <span
          class="rounded-full px-2 py-0.5 text-xs"
          :class="
            order.status === 'accepted' || order.status === 'in_progress'
              ? 'bg-green-100 text-green-700'
              : 'bg-blue-100 text-blue-700'
          "
        >
          {{ statusLabels[order.status] ?? order.status }}
        </span>
        <button
          v-if="onClose"
          type="button"
          class="flex h-6 w-6 items-center justify-center rounded-full bg-gray-200 text-gray-500 hover:bg-gray-300 hover:text-gray-700"
          @click="handleClose"
        >
          ✕
        </button>
      </div>
    </div>

    <div class="flex flex-col gap-1 text-sm">
      <div class="flex flex-row gap-4">
        <div class="flex items-center justify-between">
          <div>
            <strong>Откуда:</strong>
            {{ order.from_address ?? `${order.from_lat}, ${order.from_lon}` }}
          </div>
        </div>
        <button
          type="button"
          class="mt-1 h-full rounded-lg bg-gray-100 p-2 hover:bg-blue-100 hover:text-blue-600"
          @click="handleShowOnMap(order.from_lat, order.from_lon)"
        >
          На карте
        </button>
      </div>
      <div class="flex flex-row gap-4">
        <div class="flex items-center justify-between">
          <div>
            <strong>Куда:</strong>
            {{ order.to_address ?? `${order.to_lat}, ${order.to_lon}` }}
          </div>
        </div>
        <button
          type="button"
          class="mt-1 h-full rounded-lg bg-gray-100 p-2 hover:bg-blue-100 hover:text-blue-600"
          @click="handleShowOnMap(order.to_lat, order.to_lon)"
        >
          На карте
        </button>
      </div>
      <div v-if="order.total_distance_meters">
        <strong>Расстояние:</strong> {{ displayDistance(order.total_distance_meters) }}
      </div>
      <div v-if="order.price_rubles" class="font-semibold text-green-700">
        Приблизительно: {{ order.price_rubles.toLocaleString() }} ₽
      </div>
      <div>
        <strong>Авто:</strong> {{ (order as any).car_name }} ({{ (order as any).car_type }})
      </div>
      <div>
        <strong>Вес:</strong> {{ (order as any).car_weight_kg }} кг, <strong>длина:</strong>
        {{ (order as any).car_length_meters }} м
      </div>
      <div><strong>Колёс заблокировано:</strong> {{ order.how_many_wheels_blocked }}</div>
      <div v-if="(order as any).customer_message" class="text-gray-500 italic">
        {{ (order as any).customer_message }}
      </div>
    </div>

    <button
      v-if="!isMyOrder && showAcceptButton"
      :disabled="accepting"
      @click="handleAccept"
      class="w-full rounded-lg bg-green-500 py-2 text-white hover:bg-green-600 disabled:opacity-60"
    >
      {{ accepting ? "Принятие..." : "Принять заказ" }}
    </button>
    <button
      v-if="!isMyOrder && showAcceptButton"
      :disabled="declining"
      @click="handleDecline"
      class="w-full rounded-lg bg-red-100 py-1.5 text-red-600 hover:bg-red-200 disabled:opacity-60"
    >
      {{ declining ? "Отказ..." : "Отказаться" }}
    </button>
    <button
      v-if="showRestoreButton"
      :disabled="restoring"
      @click="handleRestore"
      class="w-full rounded-lg bg-blue-500 py-2 text-white hover:bg-blue-600 disabled:opacity-60"
    >
      {{ restoring ? "Восстановление..." : "Восстановить" }}
    </button>

    <div v-if="showCancelButton && !showCancelReason">
      <button
        :disabled="cancelling"
        @click="cancelAction"
        class="w-full rounded-lg bg-red-500 py-2 text-white hover:bg-red-600 disabled:opacity-60"
      >
        {{ cancelling ? "Отмена..." : "Отменить заказ" }}
      </button>
    </div>

    <div
      v-if="showCancelReason"
      class="flex flex-col gap-2 rounded-lg border border-red-200 bg-red-50 p-3"
    >
      <span class="text-sm font-medium text-red-700">Укажите причину отмены:</span>
      <label
        v-for="reason in cancelReasons"
        :key="reason"
        class="flex cursor-pointer items-center gap-2"
      >
        <input type="radio" v-model="selectedReason" :value="reason" class="accent-red-500" />
        <span class="text-sm">{{ reason }}</span>
      </label>
      <input
        v-if="selectedReason === 'Другая причина'"
        v-model="customReason"
        type="text"
        placeholder="Опишите причину..."
        class="w-full rounded-md border border-gray-300 px-2 py-1 text-sm"
      />
      <div class="flex gap-2">
        <button
          :disabled="!selectedReason || cancelling"
          @click="handleCancel"
          class="flex-1 rounded-lg bg-red-500 py-1.5 text-white hover:bg-red-600 disabled:opacity-60"
        >
          Подтвердить отмену
        </button>
        <button
          @click="
            showCancelReason = false;
            selectedReason = '';
            customReason = '';
          "
          class="rounded-lg bg-gray-200 px-3 py-1.5 text-gray-600 hover:bg-gray-300"
        >
          Назад
        </button>
      </div>
    </div>
  </div>
</template>
