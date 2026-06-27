<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from "vue";
import { type MapMouseEvent, type Map as MaplibreMaps } from "maplibre-gl";

import "maplibre-gl/dist/maplibre-gl.css";
import { useMaps } from "../composables/useMaps.ts";

import { useReverseGeocoding } from "../composables/useReverseGeocoding.ts";
import { $coords } from "../stores/mapsStore.ts";
import { useRouteDisplay } from "../composables/useRouteDisplay.ts";

import { RealtimeDriversAnimator } from "geolocation/animation";

const {
  map,
  initMap,
  // injectMap,
  // provideMap,
  deckOverlay,
  // addDeckLayer,
  // removeDeckLayer,
} = useMaps();

const mapContainer = ref<HTMLElement>();

const animator = ref<RealtimeDriversAnimator>();

function handleRealtimeDriversStart() {
  const mapInstance = map.get() as MaplibreMaps | null;
  if (!mapInstance) return;
  if (deckOverlay.value) {
    animator.value = new RealtimeDriversAnimator(
      mapInstance,
      deckOverlay.value,
    );
    animator.value.start();
  }
}

function handleRealtimeDriversStop() {
  if (animator.value) {
    animator.value.stop();
  }
}

onMounted(() => {
  if (!mapContainer.value) return;
  initMap(mapContainer.value, {});
  // provideMap();

  const mapInstance = map.value as MaplibreMaps | null;
  if (!mapInstance) return;

  useRouteDisplay(mapInstance);

  // TODO:
  // const { showAddressPopup } = useReverseGeocoding(mapInstance);
  // mapInstance.on("click", async (e: MapMouseEvent) => {
  //   console.log(`lan/lot: ${e.lngLat.lat} ${e.lngLat.lng}`);
  //   console.log("bounds:", map.value?.getBounds().toArray());
  //
  //   const features = mapInstance.queryRenderedFeatures(e.point, {
  //     // Defined in style.json of maps.
  //     // style.json is located in /data/maps/
  //     layers: ["address_label"],
  //   });
  //
  //   console.log(features);
  //   // if (features.length == 0) return;
  //   // const feature = features[0];
  //   // console.log("feature:", feature.layer);
  //   // const address = feature.properties.addr_housenumber;
  //   // console.log("address:", address);
  //   // await showAddressPopup(e.lngLat);
  // });
});

onUnmounted(() => {
  handleRealtimeDriversStop();
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
  <div class="flex flex-col">
    <button class="p-2" @click="handleRealtimeDriversStart()">
      Start realtime drivers
    </button>
    <button class="p-2" @click="handleRealtimeDriversStop()">
      Stop realtime drivers
    </button>
  </div>
</template>
