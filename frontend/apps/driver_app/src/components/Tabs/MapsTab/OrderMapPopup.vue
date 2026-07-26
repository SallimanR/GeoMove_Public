<script setup lang="ts">
import { computed } from "vue";
import type { Order } from "order";
import { displayDistance } from "@geomove/geo";
import { useMapPopupScale } from "@geomove/maps";

const props = defineProps<{
  order: Order;
  isMyOrder?: boolean;
  isDeclined?: boolean;
  onClick?: () => void;
}>();

const { scale } = useMapPopupScale();

const diffMin = computed(() => {
  const diffMs = Date.now() - new Date(props.order.created_at).getTime();
  return Math.floor(diffMs / 60000);
});

const timeAgo = computed(() => {
  const m = diffMin.value;
  if (m < 1) return "только что";
  if (m < 60) return `${m} мин назад`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h} ч назад`;
  const d = Math.floor(h / 24);
  return `${d} дн назад`;
});

const popupBg = computed(() => {
  if (props.isMyOrder) return "bg-green-50";
  if (props.isDeclined) return "bg-red-50";
  return "bg-white";
});

const bannerBg = computed(() => {
  if (props.isMyOrder || props.isDeclined) return "";
  const m = diffMin.value;
  if (m < 1) return "bg-green-400";
  if (m < 3) return "bg-green-300";
  if (m < 7) return "bg-green-200";
  if (m < 15) return "bg-lime-200";
  if (m < 30) return "bg-yellow-100";
  return "bg-gray-100";
});

const bannerText = computed(() => {
  if (props.isMyOrder) return "text-green-700";
  if (props.isDeclined) return "text-red-700";
  const m = diffMin.value;
  if (m < 7) return "text-green-800";
  if (m < 30) return "text-yellow-800";
  return "text-gray-500";
});

const bannerLabel = computed(() => {
  if (props.isMyOrder) return "Принят";
  if (props.isDeclined) return "Отклонён";
  return timeAgo.value;
});
</script>

<template>
  <div
    class="flex min-w-52 cursor-pointer flex-col overflow-hidden rounded-xl shadow-lg"
    :class="popupBg"
    :style="{ transform: `scale(${scale})`, transformOrigin: 'bottom center' }"
    @click="onClick?.()"
  >
    <div class="px-3 py-1 text-center font-medium" :class="[bannerBg, bannerText]">
      {{ bannerLabel }}
    </div>
    <div class="flex flex-col gap-1 p-3 text-sm">
      <div class="text-xs text-gray-400">Заказ #{{ order.id }}</div>
      <div><strong>Куда:</strong> {{ order.to_address ?? `${order.to_lat}, ${order.to_lon}` }}</div>
      <div v-if="order.total_distance_meters">
        <strong>Расстояние:</strong> {{ displayDistance(order.total_distance_meters) }}
      </div>
      <div>
        <strong>Авто:</strong> {{ (order as any).car_name }} ({{ (order as any).car_type }})
      </div>
      <div>
        <strong>Вес:</strong> {{ (order as any).car_weight_kg }} кг, <strong>длина:</strong>
        {{ (order as any).car_length_meters }} м
      </div>
      <div><strong>Колёс заблокировано:</strong> {{ order.how_many_wheels_blocked }}</div>
    </div>
  </div>
</template>
