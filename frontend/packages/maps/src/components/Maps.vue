<script setup lang="ts">
import { ref, onMounted } from "vue";
import { type Map as MaplibreMaps } from "maplibre-gl";

import "maplibre-gl/dist/maplibre-gl.css";
import { useMaps } from "../composables/useMaps.ts";

import { useRouteDisplay } from "../composables/useRouteDisplay.ts";

const props = defineProps(["styleApi", "tilesApi"]);

const { map, initMap } = useMaps();

const mapContainer = ref<HTMLElement>();

onMounted(() => {
  if (!mapContainer.value) return;
  initMap(mapContainer.value, {}, props.styleApi, props.tilesApi);

  const mapInstance = map.value as MaplibreMaps | null;
  if (!mapInstance) return;

  useRouteDisplay(mapInstance);
});
</script>
<template>
  <div id="mapContainer" ref="mapContainer" class="h-[100vh] w-full"></div>
  <div
    class="absolute inset-0 pointer-events-none flex items-center justify-center"
  >
    <div
      class="w-8 h-8 rounded-full border-2 border-blue-500 bg-blue-500/20 shadow-lg flex items-center justify-center"
    >
      <div class="w-2 h-2 rounded-full bg-blue-600"></div>
    </div>
  </div>
</template>
