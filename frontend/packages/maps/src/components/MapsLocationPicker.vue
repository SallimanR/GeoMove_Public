<script setup lang="ts">
import type { GeoPoint } from "../types/geoPoint";
import { $locationPicking, setPickCallback } from "../stores/mapsStore";

const emit = defineEmits<{
  (e: "pick", point: GeoPoint, address: string): void;
  (e: "click"): void;
}>();

function startPicking() {
  setPickCallback((point: GeoPoint, address: string) =>
    emit("pick", point, address),
  );
  $locationPicking.set(true);
  emit("click");
}
</script>

<template>
  <button
    @click="startPicking"
    type="button"
    class="p-2 rounded-lg bg-gray-200 hover:bg-gray-300"
  >
    На карте
  </button>
</template>
