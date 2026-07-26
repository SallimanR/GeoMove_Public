<script setup lang="ts">
import { computed } from "vue";
import { type FreelyAvailableDriver, $selectedFreelyAvailableDriver } from "driver";
import { useMapPopupScale } from "@geomove/maps";
import { displayDistance } from "@geomove/geo";

const props = defineProps<FreelyAvailableDriver & { asMapPopup?: boolean }>();

const distanceText = computed(() => {
  if (!props.distance) return "";
  return displayDistance(props.distance);
});

const { scale } = props.asMapPopup ? useMapPopupScale() : { scale: computed(() => 1) };

function openPopup() {
  $selectedFreelyAvailableDriver.set(props);
}
</script>

<template>
  <div
    class="flex cursor-pointer items-center gap-4 rounded-xl border border-dashed border-orange-300 bg-orange-50 p-3 hover:bg-orange-100 active:bg-orange-200"
    :style="asMapPopup ? { transform: `scale(${scale})`, transformOrigin: 'bottom center' } : {}"
    @click="openPopup"
  >
    <div
      class="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-orange-500 text-lg font-bold text-white"
    >
      {{ props.name.charAt(0).toUpperCase() }}
    </div>
    <div class="flex min-w-0 flex-1 flex-col gap-0.5">
      <div class="truncate text-lg font-semibold text-gray-800">
        {{ props.name }}
      </div>
      <div v-if="props.rating" class="flex items-center gap-1 text-sm text-yellow-500">
        <span>&#9733;</span>
        {{ props.rating }}
      </div>
      <div v-if="props.tariff_per_km" class="text-sm font-medium text-orange-600">
        {{ props.tariff_per_km }} ₽/км
      </div>
      <div v-if="props.en_route_order" class="text-xs text-green-600">Попутный заказ</div>
      <div class="text-sm text-gray-600">До вас: {{ distanceText }}</div>
    </div>
  </div>
</template>
