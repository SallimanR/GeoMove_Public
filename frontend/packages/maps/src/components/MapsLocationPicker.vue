<script setup lang="ts">
import { ref, onUnmounted, toRaw } from "vue";
import type { GeoPoint } from "../types/geoPoint";
import { $mapInstance, $coords } from "../stores/mapsStore";
import { MapLibreMap, Marker } from "maplibre-gl";
import { getReverseGeocoding, addressToText } from "@geomove/geo";

const emit = defineEmits<{
  pick: [point: GeoPoint, address: string];
  click: [];
  cancel: [];
}>();

const active = ref(false);
const address = ref("");
let marker: Marker | null = null;

function getMap(): MapLibreMap | null {
  const raw = toRaw($mapInstance.get());
  return raw ? (raw as unknown as MapLibreMap) : null;
}

function updateAddressFromMap() {
  const map = getMap();
  if (!map) return;
  const center = map.getCenter();
  getReverseGeocoding(center.lat, center.lng).then((r) => {
    address.value = addressToText(r);
  });
}

function startPicking() {
  const map = getMap();
  if (!map) return;

  const center = map.getCenter();
  marker = new Marker({ color: "#3b82f6", draggable: false })
    .setLngLat(center)
    .addTo(map);

  address.value = "";
  active.value = true;
  updateAddressFromMap();
  emit("click");
}

function confirm() {
  const map = getMap();
  if (!map) return;
  const center = map.getCenter();

  emit("pick", { lat: center.lat, lon: center.lng }, address.value || `${center.lat}, ${center.lng}`);
  cleanup();
}

function cancel() {
  emit("cancel");
  cleanup();
}

function cleanup() {
  marker?.remove();
  marker = null;
  address.value = "";
  active.value = false;
}

const unsubCoords = $coords.subscribe((c) => {
  if (!active.value || !c?.center) return;
  const { lat, lon } = c.center;
  marker?.setLngLat([lon, lat]);
  updateAddressFromMap();
});

onUnmounted(() => {
  unsubCoords();
  cleanup();
});
</script>

<template>
  <button
    @click="startPicking"
    type="button"
    class="p-2 rounded-lg bg-gray-200 hover:bg-gray-300"
  >
    На карте
  </button>

  <template v-if="active">
    <Teleport to="body">
      <div class="fixed bottom-6 left-4 right-4 z-50 flex gap-2">
        <button
          @click="cancel"
          type="button"
          class="flex-1 p-2 bg-red-300 rounded-lg shadow-md hover:bg-red-400 transition-colors"
        >
          Отменить
        </button>
        <button
          @click="confirm"
          type="button"
          class="flex-1 p-2 bg-green-300 rounded-lg shadow-md hover:bg-green-400 transition-colors"
        >
          Выбрать
        </button>
      </div>
    </Teleport>
  </template>
</template>
