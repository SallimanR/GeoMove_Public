<script setup lang="ts">
import { addressToText, getReverseGeocoding } from "geo";
import type { MapLibreMap } from "maplibre-gl";
import { ref } from "vue";
import {
  $mapCenterAddress,
  $mapCenterAddressText,
  $mapInstance,
} from "../stores/mapsStore";
import { useStore } from "@nanostores/vue";

const mapCenterAddressText = useStore($mapCenterAddressText);

$mapInstance.subscribe((map) => {
  if (map) {
    setupMapListeners(map);
  } else {
  }
});

function setupMapListeners(map: MapLibreMap) {
  async function changeCenter() {
    const mapCenter = map.getCenter();
    const req = await getReverseGeocoding(mapCenter.lat, mapCenter.lng);
    $mapCenterAddress.set(req);
    $mapCenterAddressText.set(addressToText(req));
  }

  changeCenter();

  map.on("moveend", async () => {
    changeCenter();
  });
}
</script>

<template>
  <div
    v-if="mapCenterAddressText"
    class="flex justify-center pt-4 pointer-events-none"
  >
    <div
      class="bg-white/90 backdrop-blur-sm px-4 py-2 rounded-full shadow-md text-sm font-medium text-gray-700"
    >
      {{ mapCenterAddressText }}
    </div>
  </div>
</template>
