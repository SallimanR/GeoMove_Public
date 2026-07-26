<script setup lang="ts">
import { inject } from "vue";
import { useStore } from "@nanostores/vue";
import { $selectedDriver, resolveImageUrl, $driverDropdownOpen } from "src/stores/driverStore.ts";
import { $mapInstance } from "@geomove/maps";
import { ACTIVE_TAB_KEY } from "src/injectionKeys.ts";

const driver = useStore($selectedDriver);
const activeTab = inject(ACTIVE_TAB_KEY)!;

function formatTime(iso: string | undefined | null): string {
  if (!iso) return "—";
  return iso.split("T")[1]?.slice(0, 5) ?? "—";
}

function close() {
  $selectedDriver.set(null);
}

function showOnMap() {
  const d = $selectedDriver.get();
  if (!d) return;
  const map = $mapInstance.get();
  if (map) {
    activeTab.value = "mapsTab";
    map.flyTo({ center: [d.lon, d.lat] });
  }
  $driverDropdownOpen.set(false);
  close();
}
</script>

<template>
  <Teleport to="body">
    <div
      v-if="driver"
      class="fixed inset-0 z-200 flex items-center justify-center p-4"
      @click.self="close"
    >
      <div class="absolute inset-0 bg-black/40" />

      <div
        class="relative max-h-[90vh] w-full max-w-sm overflow-y-auto rounded-2xl bg-white p-5 shadow-xl"
      >
        <button
          class="absolute top-3 right-3 flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 text-gray-500 transition hover:bg-gray-200 hover:text-gray-700"
          @click="close"
        >
          &#10005;
        </button>

        <div class="flex flex-col items-center gap-1">
          <div class="text-center">
            <div class="text-xl font-bold text-gray-800">{{ driver.name }}</div>
            <div
              v-if="driver.rating"
              class="mt-0.5 flex items-center justify-center gap-1 text-yellow-500"
            >
              <span class="text-lg">&#9733;</span>
              <span class="text-gray-600">{{ driver.rating }}</span>
            </div>
          </div>
        </div>

        <div class="mt-4 flex flex-col gap-2 text-sm">
          <div
            v-if="driver.work_starts || driver.work_ends"
            class="flex items-center justify-between px-1"
          >
            <span class="text-gray-500">Часы работы</span>
            <span class="font-medium text-gray-800">
              {{ formatTime(driver.work_starts) }} –
              {{ formatTime(driver.work_ends) }}
            </span>
          </div>

          <div v-if="driver.phone" class="flex items-center justify-between px-1">
            <span class="text-gray-500">Телефон</span>
            <span class="font-medium text-gray-800">{{ driver.phone }}</span>
          </div>

          <div v-if="driver.max_car_weight_kg" class="flex items-center justify-between px-1">
            <span class="text-gray-500">Макс. вес</span>
            <span class="font-medium text-gray-800">{{ driver.max_car_weight_kg }} кг</span>
          </div>

          <div v-if="driver.max_car_length_meters" class="flex items-center justify-between px-1">
            <span class="text-gray-500">Макс. длина</span>
            <span class="font-medium text-gray-800">{{ driver.max_car_length_meters }} м</span>
          </div>

          <div v-if="driver.address" class="flex items-center justify-between px-1">
            <span class="text-gray-500">Адрес</span>
            <span class="max-w-50 text-right font-medium text-gray-800">{{ driver.address }}</span>
          </div>
        </div>

        <div v-if="driver.car_photo_main" class="mt-4">
          <img
            :src="resolveImageUrl(driver.car_photo_main)"
            class="h-40 w-full rounded-xl object-cover"
          />
        </div>

        <div class="mt-4 flex w-full gap-3">
          <button
            class="flex-1 rounded-xl bg-blue-500 p-3 text-center font-medium text-white transition hover:bg-blue-600"
            @click="showOnMap"
          >
            На карте
          </button>
          <button
            class="flex-1 rounded-xl bg-green-500 p-3 text-center font-medium text-white transition hover:bg-green-600"
            @click="close"
          >
            Выбрать
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
