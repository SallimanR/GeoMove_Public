<script setup lang="ts">
import { inject } from "vue";
import { useStore } from "@nanostores/vue";
import { $selectedDriver } from "src/stores/driverStore.ts";
import { $mapInstance } from "@geomove/maps";
import { ACTIVE_TAB_KEY } from "src/injectionKeys.ts";

const driver = useStore($selectedDriver);
const activeTab = inject(ACTIVE_TAB_KEY)!;
const imagePath = "tow_image.jpg";

function close() {
  $selectedDriver.set(null);
}

function showOnMap() {
  const d = $selectedDriver.get();
  if (!d) return;
  const map = $mapInstance.get();
  if (map) {
    activeTab.value = "mapsTab";
    map.jumpTo({ center: [d.lon, d.lat] });
  }
  close();
}
</script>

<template>
  <div
    v-if="driver"
    class="absolute bottom-0 left-0 right-0 z-100 rounded-t-2xl bg-gray-200 overflow-y-auto p-4 flex flex-col gap-4 border-4"
    @click="close"
  >
    <button
      class="absolute top-3 right-3 w-8 h-8 rounded-full bg-gray-200 hover:bg-gray-300 flex items-center justify-center text-gray-600 font-bold transition"
      @click="close"
    >
      &#10005;
    </button>

    <img
      :src="imagePath"
      class="w-24 h-24 rounded-full object-cover border-4 border-gray-100"
    />

    <div class="text-center">
      <div class="text-xl font-bold">{{ driver.name }}</div>
      <div
        class="flex items-center justify-center gap-1 text-sm text-gray-600 mt-1"
      >
        <span class="text-yellow-500 text-lg">&#9733;</span>
        {{ driver.rating ?? "—" }}
      </div>
    </div>

    <div class="w-full flex flex-col gap-2 text-sm text-gray-500">
      <div
        v-if="driver.work_starts || driver.work_ends"
        class="flex justify-between"
      >
        <span>Часы работы</span>
        <span
          >{{ driver.work_starts ?? "—" }} – {{ driver.work_ends ?? "—" }}</span
        >
      </div>
    </div>

    <div class="flex gap-3 w-full mt-2">
      <button
        class="flex-1 rounded-xl bg-green-300 p-3 text-center font-medium text-gray-800 hover:bg-green-400 transition"
        @click="showOnMap"
      >
        На карте
      </button>
      <button
        class="flex-1 rounded-xl bg-green-400 p-3 text-center font-medium text-white hover:bg-green-500 transition"
        @click="close"
      >
        Выбрать
      </button>
    </div>
  </div>
</template>
