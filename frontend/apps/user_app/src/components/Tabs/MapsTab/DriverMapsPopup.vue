<script setup lang="ts">
import type { Driver } from "driver/types/driver.ts";
import { $selectedDriver, resolveImageUrl } from "src/stores/driverStore.ts";
import { useMapPopupScale } from "@geomove/maps";

const props = defineProps<Driver>();

const defaultImage = "tow_image.jpg";

const { scale } = useMapPopupScale();

function openCardPopup() {
  $selectedDriver.set(props);
}
</script>

<template>
  <div
    class="flex min-w-20 cursor-pointer flex-col items-center rounded-2xl border border-gray-200 bg-white p-2 shadow-lg"
    :style="{ transform: `scale(${scale})`, transformOrigin: 'bottom center' }"
    @click="openCardPopup"
  >
    <img
      :src="resolveImageUrl(props.car_photo_main) || defaultImage"
      class="h-12 w-12 rounded-xl border-2 border-gray-100 object-cover"
    />
    <div
      class="mt-1 line-clamp-2 max-w-28 text-center text-sm leading-tight font-medium break-words text-gray-800"
    >
      {{ props.name }}
    </div>
    <div v-if="props.rating" class="mt-0.5 flex items-center gap-0.5 text-xs text-yellow-500">
      <span>&#9733;</span> {{ props.rating }}
    </div>
  </div>
</template>
