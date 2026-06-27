<script setup lang="ts">
import { onMounted, ref, watch, type Ref } from "vue";
import {
  $endPoint,
  $isRouteLoading,
  $routePath,
  $startPoint,
} from "../stores/routeStore";
import { fetchRoute } from "geo";
import { $mapInstance } from "../stores/mapsStore";
import { Marker, type MapMouseEvent } from "maplibre-gl";
import MapsSearchField from "./MapsSearchField.vue";
import { useStore } from "@nanostores/vue";
import type { GeoPoint } from "../types/geoPoint";

enum RouteInputMode {
  Deactivated = 0,
  SelectStartPoint = 1,
  SelectEndPoint = 2,
}

const inputState = ref(RouteInputMode.Deactivated);

function changeInputeState(state: RouteInputMode) {
  console.log("Click");
  inputState.value = state;
}

const startPoint = useStore($startPoint);
const endPoint = useStore($endPoint);

onMounted(() => {
  const map = $mapInstance.get();
  // FIXME: map is not ready
  if (!map) {
    console.log("Map is not ready");
    return;
  }

  const startMarker = new Marker({ draggable: true });
  startMarker.on("dragend", () => {
    const lngLat = startMarker.getLngLat();
    $startPoint.set({ lat: lngLat.lat, lon: lngLat.lng });
  });
  const endMarker = new Marker({ draggable: true });
  endMarker.on("dragend", () => {
    const lngLat = endMarker.getLngLat();
    $endPoint.set({ lat: lngLat.lat, lon: lngLat.lng });
  });

  map.on("click", (e: MapMouseEvent) => {
    switch (inputState.value) {
      case RouteInputMode.SelectStartPoint: {
        startMarker.setLngLat(e.lngLat).addTo(map);
        $startPoint.set({ lat: e.lngLat.lat, lon: e.lngLat.lng });
        console.log("start point: ", $startPoint.get());
      }
      case RouteInputMode.SelectEndPoint: {
        endMarker.setLngLat(e.lngLat).addTo(map);
        $endPoint.set({ lat: e.lngLat.lat, lon: e.lngLat.lng });
        console.log("end point: ", $endPoint.get());
      }
    }
  });
});

const startLocation = ref("");
const endLocation = ref("");

function onStartSearchResultClick(lat: number, lon: number): void {
  $startPoint.set({ lat: lat, lon: lon });
  console.log("onStartSearchResultClick: ", lat, lon);
  changeInputeState(RouteInputMode.SelectStartPoint);
}

function onEndSearchResultClick(lat: number, lon: number): void {
  $endPoint.set({ lat: lat, lon: lon });
  changeInputeState(RouteInputMode.SelectEndPoint);
}

function handleRemoveRoute() {
  changeInputeState(RouteInputMode.Deactivated);
  $startPoint.set(null);
  $endPoint.set(null);
  // $routePath.set(null);
}

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
      console.log(route);

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
</script>
<template>
  <div class="flex flex-col">
    <div>
      <MapsSearchField
        v-model="startLocation"
        :onSearchResultClick="onStartSearchResultClick"
        placeholder="Start"
      />
      <button
        @click="changeInputeState(RouteInputMode.SelectStartPoint)"
        class="rounded-xl bg-green-300 p-2"
      >
        on maps
      </button>
    </div>
    <div>
      <MapsSearchField
        v-model="endLocation"
        :onSearchResultClick="onEndSearchResultClick"
        placeholder="End"
      />
      <button
        @click="changeInputeState(RouteInputMode.SelectEndPoint)"
        class="rounded-xl bg-green-300 p-2"
      >
        on maps
      </button>
    </div>
    <button v-show="$routePath.get()" @click="handleRemoveRoute()">
      Cancel
    </button>
    <!-- <div> -->
    <!--   <button @click="handleFindRoute()" class="rounded-xl bg-orange-300 p-2"> -->
    <!--     Find route -->
    <!--   </button> -->
    <!--   <button @click="handleRemoveRoute()" class="rounded-xl bg-red-300 p-2"> -->
    <!--     Remove route -->
    <!--   </button> -->
    <!-- </div> -->
  </div>
</template>
