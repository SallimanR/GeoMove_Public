<script setup lang="ts">
import { ref, onUnmounted } from "vue";
import { useStore } from "@nanostores/vue";
import { Map as MaplibreMap, Marker } from "maplibre-gl";
import { getReverseGeocoding, addressToText } from "geo";
import {
  $mapInstance,
  $coords,
  $locationPicking,
  invokePickCallback,
  clearPickCallback,
} from "../stores/mapsStore";
import type { GeoPoint } from "../types/geoPoint";
import MapsRouteInput from "./MapsRouteInput.vue";
import MapsSearchInput from "./MapsSearchInput.vue";
import MapsCurrentLocationBox from "./MapsCurrentLocationBox.vue";
import MapsGPSLocation from "./MapsGPSLocation.vue";

const isPicking = useStore($locationPicking);
const address = ref("");
let marker: Marker | null = null;

async function updateAddress(lat: number, lng: number) {
  try {
    const result = await getReverseGeocoding(lat, lng);
    address.value = addressToText(result);
  } catch {
    address.value = `${lat.toFixed(5)}, ${lng.toFixed(5)}`;
  }
}

function cancelPicking() {
  clearPickCallback();
  $locationPicking.set(false);
  marker?.remove();
  marker = null;
  address.value = "";
}

function confirmPick() {
  const map = $mapInstance.get();
  if (!map) return;

  const center = (map as unknown as MaplibreMap).getCenter();
  const point: GeoPoint = { lat: center.lat, lon: center.lng };

  invokePickCallback(point, address.value || `${point.lat}, ${point.lon}`);
  $locationPicking.set(false);
  marker?.remove();
  marker = null;
  address.value = "";
}

const unsubPicking = $locationPicking.subscribe((active) => {
  if (!active) {
    marker?.remove();
    marker = null;
    address.value = "";
    return;
  }

  const map = $mapInstance.get();
  if (!map) return;

  marker = new Marker({ color: "#3b82f6", draggable: false })
    .setLngLat(map.getCenter())
    .addTo(map as unknown as MaplibreMap);

  const center = (map as unknown as MaplibreMap).getCenter();
  updateAddress(center.lat, center.lng);
});

const unsubCoords = $coords.subscribe((coords) => {
  if (!$locationPicking.get() || !coords?.center) return;
  const { lat, lon } = coords.center;
  marker?.setLngLat([lon, lat]);
  updateAddress(lat, lon);
});

onUnmounted(() => {
  unsubPicking();
  unsubCoords();
  marker?.remove();
});
</script>

<template>
  <div class="absolute inset-0 flex flex-col pointer-events-none">
    <div class="flex gap-2 mt-4 mr-4 ml-4">
      <MapsSearchInput class="flex-1 pointer-events-auto" />
      <MapsGPSLocation class="pointer-events-auto" />
    </div>
    <MapsCurrentLocationBox />
    <template v-if="isPicking">
      <div
        class="fixed bottom-6 left-4 right-4 z-50 flex gap-2 pointer-events-auto"
      >
        <button
          @click="cancelPicking"
          type="button"
          class="flex-1 p-2 bg-red-300 rounded-lg shadow-md hover:bg-red-400 transition-colors"
        >
          Отменить
        </button>
        <button
          @click="confirmPick"
          type="button"
          class="flex-1 p-2 bg-green-300 rounded-lg shadow-md hover:bg-green-400 transition-colors"
        >
          Выбрать
        </button>
      </div>
    </template>

    <div v-else class="pointer-events-auto">
      <MapsRouteInput />
    </div>
  </div>
</template>
