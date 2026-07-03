<script setup lang="ts">
import { ref, onMounted } from "vue";
import { type Map as MaplibreMaps } from "maplibre-gl";
import "maplibre-gl/dist/maplibre-gl.css";

import { useMaps } from "../composables/useMaps.ts";
import { useRouteDisplay } from "../composables/useRouteDisplay.ts";

import MapsCentralMarker from "./MapsCentralMarker.vue";

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
  <div id="mapContainer" ref="mapContainer" class="relative h-full w-full">
    <MapsCentralMarker />
  </div>
</template>
