<script setup lang="ts">
import { ref, onMounted } from "vue";
import "maplibre-gl/dist/maplibre-gl.css";

import { useMaps } from "../composables/useMaps.ts";

import MapsCentralMarker from "./MapsCentralMarker.vue";

const props = defineProps(["styleApi"]);

const { initMap } = useMaps();

const mapContainer = ref<HTMLElement>();

onMounted(() => {
  if (!mapContainer.value) return;
  initMap(mapContainer.value, {}, props.styleApi);
});
</script>

<template>
  <div id="mapContainer" ref="mapContainer" class="relative h-full w-full">
    <MapsCentralMarker />
  </div>
</template>

<style scoped>
:deep(.maplibregl-popup-content) {
  padding: 0 !important;
  background: transparent !important;
  border-radius: 0 !important;
  box-shadow: none !important;
}
</style>
