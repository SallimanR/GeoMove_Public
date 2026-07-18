<script setup lang="ts">
import { ref } from "vue";
import {
  $endAddress,
  $endPoint,
  $routePath,
  $startAddress,
  $startPoint,
} from "../stores/routeStore";
import { Map as MaplibreMap } from "maplibre-gl";

import { displayDistance, getReverseGeocoding, addressToText } from "geo";
import {
  $mapCenterAddress,
  $mapCenterAddressText,
  $mapInstance,
} from "../stores/mapsStore";
import { LngLat, Marker } from "maplibre-gl";
import MapsSearchField from "./MapsSearchField.vue";
import { useStore } from "@nanostores/vue";
import { useRouteDisplay } from "../composables/useRouteDisplay";

enum RouteInputMode {
  SelectStartPoint,
  SelectEndPoint,
  Ready,
}

const inputState = ref(RouteInputMode.SelectStartPoint);
const needsConfirm = ref(false);

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
  if (!mapInstance.value) return;

  if (!$startPoint.value) {
    const mapCenter = mapInstance.value?.getCenter();
    changeStartPoint(mapCenter, mapInstance.value);
  }
  showPanel.value = PanelState.Active;
}

function closePanel() {
  showPanel.value = PanelState.Inactive;
}

const routePath = useStore($routePath);

const startMarker = new Marker({ draggable: true, color: "#40fc0c" });
const endMarker = new Marker({ draggable: true, color: "#fc5b55" });

const startAddressText = ref("");
const endAddressText = ref("");

async function changeStartPoint(point: LngLat, map: MaplibreMap) {
  startMarker.setLngLat(point).addTo(map);
  $startPoint.set({ lat: point.lat, lon: point.lng });

  $startAddress.set($mapCenterAddress.get());
  startAddressText.value = $mapCenterAddressText.get();
}

async function changeEndPoint(point: LngLat, map: MaplibreMap) {
  endMarker.setLngLat(point).addTo(map);
  $endPoint.set({ lat: point.lat, lon: point.lng });

  $endAddress.set($mapCenterAddress.get());
  endAddressText.value = $mapCenterAddressText.get();

  inputState.value = RouteInputMode.Ready;
}

function clearStart() {
  startMarker.remove();
  startAddressText.value = "";
  $startPoint.set(null);
  $startAddress.set(null);
}

function clearEnd() {
  endMarker.remove();
  endAddressText.value = "";
  $endPoint.set(null);
  $endAddress.set(null);
}

function handleRemoveRoute() {
  clearStart();
  clearEnd();

  changeInputeState(RouteInputMode.SelectStartPoint);
}

function onStartSearchResultClick(lat: number, lon: number) {
  $startPoint.set({ lat, lon });
  changeInputeState(RouteInputMode.SelectStartPoint);
}
function onEndSearchResultClick(lat: number, lon: number) {
  $endPoint.set({ lat, lon });
  changeInputeState(RouteInputMode.SelectEndPoint);
}

const mapInstance = ref<MaplibreMap>();
const isMapMoving = ref(false);

$mapInstance.subscribe((map) => {
  if (!map) return;
  mapInstance.value = map as unknown as MaplibreMap;
  setupMapListeners(mapInstance.value);
  useRouteDisplay(map as MaplibreMap);
});

function setupMapListeners(map: MaplibreMap) {
  map.on("movestart", () => {
    isMapMoving.value = true;
  });
  map.on("moveend", async () => {
    isMapMoving.value = false;
  });

  startMarker.on("dragend", async () => {
    const point = startMarker.getLngLat();
    startMarker.setLngLat(point).addTo(map);
    $startPoint.set({ lat: point.lat, lon: point.lng });

    const req = await getReverseGeocoding(point.lat, point.lng);
    $startAddress.set(req);
    startAddressText.value = addressToText(req);
  });

  endMarker.on("dragend", async () => {
    const point = endMarker.getLngLat();
    endMarker.setLngLat(point).addTo(map);
    $endPoint.set({ lat: point.lat, lon: point.lng });

    const req = await getReverseGeocoding(point.lat, point.lng);
    $endAddress.set(req);
    endAddressText.value = addressToText(req);

    inputState.value = RouteInputMode.Ready;
  });
}

function handleOnMaps(state: RouteInputMode) {
  inputState.value = state;
  needsConfirm.value = true;
  showPanel.value = PanelState.WaitingForClick;
}

function handleConfirm() {
  if (!mapInstance.value) return;

  const mapCenter = mapInstance.value?.getCenter();

  switch (inputState.value) {
    case RouteInputMode.SelectStartPoint:
      changeStartPoint(mapCenter, mapInstance.value);
      break;
    case RouteInputMode.SelectEndPoint:
      changeEndPoint(mapCenter, mapInstance.value);
      break;
  }

  showPanel.value = PanelState.Active;
  needsConfirm.value = false;
}

function handleContinue() {
  inputState.value = RouteInputMode.Ready;
  showPanel.value = PanelState.Inactive;
}
</script>

<template>
  <div v-if="needsConfirm && !isMapMoving" class="">
    <button
      @click="handleConfirm()"
      class="absolute bottom-4 left-1/2 -translate-x-1/2 max-w-40 rounded-xl bg-green-400 mb-2 p-2 whitespace-nowrap text-center"
    >
      Выбрать
    </button>
  </div>

  <div
    v-if="showPanel == PanelState.Inactive"
    class="absolute bottom-4 left-1/2 -translate-x-1/2 w-full max-w-100"
  >
    <div
      @click="openPanel"
      class="bg-white rounded-xl p-3 shadow-lg cursor-pointer text-gray-500 hover:bg-gray-50 transition text-center"
    >
      Куда отвезти?
    </div>
  </div>

  <div v-if="showPanel == PanelState.Active" class="absolute inset-0 z-100">
    <!-- Backdrop -->
    <div
      class="absolute inset-0 bg-black/30 flex items-center justify-center"
      @click="closePanel"
    >
      <span class="text-white text-lg font-medium select-none">
        Нажмите чтобы закрыть
      </span>
    </div>

    <!-- Panel -->
    <div
      class="absolute bottom-0 left-1/2 -translate-x-1/2 w-full max-w-300 rounded-t-2xl bg-white overflow-y-auto p-4"
    >
      <div class="flex flex-col space-y-2 mt-4">
        <div class="flex items-center gap-2">
          <MapsSearchField
            class="flex-1"
            v-model="startAddressText"
            :onSearchResultClick="onStartSearchResultClick"
            :onClear="clearStart"
            placeholder="Откуда"
          />
          <button
            @click="handleOnMaps(RouteInputMode.SelectStartPoint)"
            class="rounded-xl bg-green-300 p-2 whitespace-nowrap"
          >
            на карте
          </button>
        </div>

        <div class="flex items-center gap-2">
          <MapsSearchField
            class="flex-1"
            v-model="endAddressText"
            :onSearchResultClick="onEndSearchResultClick"
            :onClear="clearEnd"
            placeholder="Куда"
          />
          <button
            @click="handleOnMaps(RouteInputMode.SelectEndPoint)"
            class="rounded-xl bg-green-300 p-2 whitespace-nowrap"
          >
            на карте
          </button>
        </div>

        <div class="flex gap-2">
          <button
            v-show="routePath"
            @click="handleContinue()"
            class="rounded-xl bg-green-300 p-2"
          >
            Продолжить
          </button>
          <button
            v-show="routePath"
            @click="handleRemoveRoute"
            class="rounded-xl bg-red-300 p-2"
          >
            Отменить путь
          </button>
        </div>

        <div v-if="routePath && routePath.paths[0].distance">
          Расстояние: {{ displayDistance(routePath.paths[0].distance) }}
        </div>
      </div>
    </div>
  </div>
</template>
