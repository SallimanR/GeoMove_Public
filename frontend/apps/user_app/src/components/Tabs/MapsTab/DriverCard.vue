<script setup lang="ts">
import type { Driver } from "driver/types/driver.ts";
import { $coords, $startPoint } from "@geomove/maps";
import { haversineDistance, displayDistance } from "@geomove/geo";
import { $selectedDriver } from "../../../stores/driverStore.ts";
import { resolveImageUrl } from "../../../stores/driverStore.ts";
import { computed } from "vue";
import { useStore } from "@nanostores/vue";

const props = defineProps<Driver>();

const defaultImage = "tow_image.jpg";

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
  return displayDistance(haversineDistance([props.lat, props.lon], [center.lat, center.lon]));
});

const carSpecs = computed(() => {
  const parts: string[] = [];
  if (props.max_car_weight_kg) {
    parts.push(`${props.max_car_weight_kg} кг`);
  }
  if (props.max_car_length_meters) {
    parts.push(`${props.max_car_length_meters} м`);
  }
  return parts.join(" ⸱ ");
});

function openPopup() {
  $selectedDriver.set(props);
}
</script>

<template>
  <div
    class="flex cursor-pointer items-center gap-4 rounded-xl bg-gray-100 p-3 hover:bg-gray-200 active:bg-gray-300"
    @click="openPopup"
  >
    <img
      :src="resolveImageUrl(props.car_photo_main) || defaultImage"
      class="h-16 w-16 shrink-0 rounded-full border-2 border-gray-200 object-cover"
    />
    <div class="flex min-w-0 flex-col gap-0.5">
      <div class="truncate text-lg font-semibold text-gray-800">
        {{ props.name }}
      </div>
      <div v-if="props.rating" class="flex items-center gap-1 text-sm text-yellow-500">
        <span>&#9733;</span>
        {{ props.rating }}
      </div>
      <div v-if="carSpecs" class="text-sm text-gray-500">
        {{ carSpecs }}
      </div>
      <div class="text-sm text-gray-600">До вас: {{ distanceText }}</div>
    </div>
  </div>
</template>
