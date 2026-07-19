<script setup lang="ts">
import type { Driver } from "driver/types/driver.ts";
import { $coords, $startPoint } from "@geomove/maps";
import { haversineDistance, displayDistance } from "@geomove/geo";
import { $selectedDriver } from "../../../stores/driverStore.ts";
import { computed } from "vue";
import { useStore } from "@nanostores/vue";

const props = defineProps<Driver>();

const profileImageDefaultPath = "tow_image.jpg";

const coords = useStore($coords);
const startPoint = useStore($startPoint);

const distanceText = computed(() => {
  let center = startPoint.value;
  if (!center) {
    const coordsCenter = coords.value;
    if (coordsCenter) {
      center = coordsCenter.center;
    }
  }
  if (!center) return "";
  return displayDistance(
    haversineDistance([props.lat, props.lon], [center.lat, center.lon]),
  );
});

function openPopup() {
  $selectedDriver.set(props);
}
</script>

<template>
  <div
    class="flex items-center gap-4 rounded-xl p-3 bg-gray-100 hover:bg-gray-200 active:bg-gray-300"
    @click="openPopup"
  >
    <img
      :src="profileImageDefaultPath"
      class="w-32 h-32 rounded-full object-cover border-2 border-gray-200 shrink-0"
    />
    <div class="flex flex-col gap-0.5 min-w-0">
      <div class="font-semibold text-lg text-gray-800 truncate">
        {{ props.name }}
      </div>
      <div v-if="props.rating" class="flex items-center gap-1 text-yellow-500">
        <span>&#9733;</span>
        {{ props.rating }}
      </div>
      <div class="text-sm text-gray-600">До вас: {{ distanceText }}</div>
    </div>
  </div>
</template>
