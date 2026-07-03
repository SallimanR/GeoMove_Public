<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import {
  $endPoint,
  $isRouteLoading,
  $routePath,
  $startPoint,
} from "../stores/routeStore";
import { fetchRoute, getReverseGeocoding } from "geo";
import { $mapInstance } from "../stores/mapsStore";
import { LngLat, Marker, type MapMouseEvent } from "maplibre-gl";
import MapsSearchField from "./MapsSearchField.vue";
import { useStore } from "@nanostores/vue";
import { Map as MaplibreMap } from "maplibre-gl";

enum RouteInputMode {
  SelectStartPoint,
  SelectEndPoint,
  Ready,
}

const inputState = ref(RouteInputMode.SelectStartPoint);

function changeInputeState(state: RouteInputMode) {
  inputState.value = state;
}

enum PanelState {
  Inactive,
  Active,
  WaitingForClick,
}

const showPanel = ref(PanelState.Inactive);

function openPanel() {
  showPanel.value = PanelState.Active;
}
function closePanel() {
  showPanel.value = PanelState.Inactive;
}

const startPoint = useStore($startPoint);
const endPoint = useStore($endPoint);
const routePath = useStore($routePath);

const startMarker = new Marker({ draggable: true });
const endMarker = new Marker({ draggable: true });

async function changeStartPoint(point: LngLat, map: MaplibreMap) {
  startMarker.setLngLat(point).addTo(map);
  $startPoint.set({ lat: point.lat, lon: point.lng });
  startLocation.value = await getLocationText(point);
}

async function changeEndPoint(point: LngLat, map: MaplibreMap) {
  endMarker.setLngLat(point).addTo(map);
  $endPoint.set({ lat: point.lat, lon: point.lng });
  endLocation.value = await getLocationText(point);
  inputState.value = RouteInputMode.Ready;
}

async function getLocationText(point: LngLat): Promise<string> {
  const req = (await getReverseGeocoding(point.lat, point.lng)).features[0];
  if (!req) return "";
  return [
    req.properties.name,
    req.properties.street,
    req.properties.housenumber,
  ]
    .filter(Boolean)
    .join(" ");
}

const startLocation = ref("");
const endLocation = ref("");

function onStartSearchResultClick(lat: number, lon: number) {
  $startPoint.set({ lat, lon });
  changeInputeState(RouteInputMode.SelectStartPoint);
}
function onEndSearchResultClick(lat: number, lon: number) {
  $endPoint.set({ lat, lon });
  changeInputeState(RouteInputMode.SelectEndPoint);
}

function handleRemoveRoute() {
  changeInputeState(RouteInputMode.SelectStartPoint);
  $startPoint.set(null);
  $endPoint.set(null);
  startMarker.remove();
  endMarker.remove();
  startLocation.value = "";
  endLocation.value = "";
  inputState.value = RouteInputMode.SelectStartPoint;
}

onMounted(() => {
  const map = $mapInstance.get();
  if (!map) {
    console.log("Map is not ready");
    return;
  }

  map.on("click", (e: MapMouseEvent) => {
    showPanel.value = PanelState.Active;
    switch (inputState.value) {
      case RouteInputMode.SelectStartPoint:
        changeStartPoint(e.lngLat, map);
        break;
      case RouteInputMode.SelectEndPoint:
        changeEndPoint(e.lngLat, map);
        break;
    }
  });

  // TODO: disable all ui
  map.on("dragstart", () => {});

  map.on("moveend", () => {
    if (inputState.value === RouteInputMode.SelectStartPoint) {
      const center = map.getCenter();
      changeStartPoint(center, map);
    }
  });

  startMarker.on("dragend", () => {
    const lngLat = startMarker.getLngLat();
    changeStartPoint(lngLat, map);
  });

  endMarker.on("dragend", () => {
    const lngLat = endMarker.getLngLat();
    changeEndPoint(lngLat, map);
  });
});

watch(
  [startPoint, endPoint],
  async ([newStart, newEnd]) => {
    if (!newStart || !newEnd) {
      $routePath.set(null);
      return;
    }
    $isRouteLoading.set(true);
    try {
      const route = await fetchRoute(newStart, newEnd);
      $routePath.set(route);
    } catch (err) {
      $routePath.set(null);
      console.error(err);
    } finally {
      $isRouteLoading.set(false);
    }
  },
  { deep: true },
);

function handleOnMaps(state: RouteInputMode) {
  inputState.value = state;
  showPanel.value = PanelState.WaitingForClick;
}
</script>

<template>
  <div
    v-if="startLocation"
    class="flex justify-center pt-4 pointer-events-none"
  >
    <div
      class="bg-white/90 backdrop-blur-sm px-4 py-2 rounded-full shadow-md text-sm font-medium text-gray-700"
    >
      {{ startLocation }}
    </div>
  </div>

  <div
    v-if="showPanel == PanelState.Inactive"
    class="absolute bottom-4 left-4 right-4"
  >
    <div
      @click="openPanel"
      class="bg-white rounded-xl p-3 shadow-lg cursor-pointer text-gray-500 hover:bg-gray-50 transition text-center"
    >
      Where are we going?
    </div>
  </div>

  <div v-if="showPanel == PanelState.Active" class="absolute inset-0 z-100">
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black/30" @click="closePanel"></div>

    <!-- Panel Content – takes 90% of viewport height from bottom -->
    <div
      class="absolute bottom-0 left-0 right-0 h-[90vh] bg-white rounded-t-2xl shadow-2xl overflow-y-auto p-4"
    >
      <button
        @click="closePanel"
        class="absolute top-3 right-4 text-2xl text-gray-500 bg-red-300 hover:text-gray-700"
      >
        ×
      </button>

      <div class="flex flex-col space-y-2 mt-4">
        <div class="flex items-center gap-2">
          <div class="flex-1">
            <MapsSearchField
              v-model="startLocation"
              :onSearchResultClick="onStartSearchResultClick"
              placeholder="Start"
            />
          </div>
          <button
            @click="handleOnMaps(RouteInputMode.SelectStartPoint)"
            class="rounded-xl bg-green-300 p-2 whitespace-nowrap"
          >
            on maps
          </button>
        </div>

        <div class="flex items-center gap-2">
          <div class="flex-1">
            <MapsSearchField
              v-model="endLocation"
              :onSearchResultClick="onEndSearchResultClick"
              placeholder="End"
            />
          </div>
          <button
            @click="handleOnMaps(RouteInputMode.SelectEndPoint)"
            class="rounded-xl bg-green-300 p-2 whitespace-nowrap"
          >
            on maps
          </button>
        </div>

        <div class="flex gap-2">
          <button v-show="routePath" class="rounded-xl bg-green-300 p-2">
            Continue
          </button>
          <button
            v-show="routePath"
            @click="handleRemoveRoute"
            class="rounded-xl bg-red-300 p-2"
          >
            Remove route
          </button>
        </div>

        <div v-if="routePath">
          <div>Route distance: {{ $routePath.value?.paths[0].distance }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
