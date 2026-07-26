<script setup lang="ts">
import { inject } from "vue";
import { useStore } from "@nanostores/vue";
import { $selectedFreelyAvailableDriver } from "driver/store/driverStore.ts";
import { $driverDropdownOpen } from "src/stores/driverStore.ts";
import { $mapInstance, flyToPointOnMap } from "@geomove/maps";
import { displayDistance } from "@geomove/geo";

const driver = useStore($selectedFreelyAvailableDriver);

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString("ru-RU");
}

function close() {
  $selectedFreelyAvailableDriver.set(null);
}

function showOnMap(lat: number, lon: number) {
  flyToPointOnMap(lat, lon);
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

        <div class="flex flex-col items-center gap-3">
          <div
            class="flex h-16 w-16 items-center justify-center rounded-xl bg-orange-500 text-2xl font-bold text-white"
          >
            {{ driver.name.charAt(0).toUpperCase() }}
          </div>

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
          <div class="flex items-center justify-between px-1">
            <span class="text-gray-500">Период</span>
            <span class="text-right text-xs font-medium text-gray-800">
              {{ formatDate(driver.from_date) }}<br />
              {{ formatDate(driver.to_date) }}
            </span>
          </div>

          <div v-if="driver.tariff_per_km" class="flex items-center justify-between px-1">
            <span class="text-gray-500">Тариф</span>
            <span class="font-medium text-gray-800">{{ driver.tariff_per_km }} ₽/км</span>
          </div>

          <div v-if="driver.from_location.address" class="flex items-center justify-between px-1">
            <span class="font-medium text-gray-500">Отправление</span>
            <span class="max-w-50 text-right text-gray-800">{{
              driver.from_location.address
            }}</span>
          </div>

          <div class="flex flex-col gap-1 px-1">
            <span class="font-medium text-gray-500">Точки Прибытия</span>
            <template v-if="driver.to_locations && driver.to_locations.length">
              <div
                v-for="(loc, i) in driver.to_locations"
                :key="i"
                class="text flex flex-row justify-between rounded-lg bg-gray-50 text-gray-800"
              >
                <div class="p-2">
                  {{ loc.address || `${loc.lat.toFixed(5)}, ${loc.lon.toFixed(5)}` }}
                </div>
                <button
                  class="h-full rounded-lg bg-gray-100 p-2 hover:bg-blue-100 hover:text-blue-600"
                  @click="showOnMap(loc.lat, loc.lon)"
                >
                  На карте
                </button>
              </div>
            </template>
            <span v-else class="text-xs text-gray-400">—</span>
          </div>

          <div v-if="driver.en_route_order" class="flex items-center justify-between px-1">
            <span class="text-gray-500">Попутный заказ</span>
            <span class="font-medium text-green-600">Да</span>
          </div>

          <div class="flex items-center justify-between px-1">
            <span class="text-gray-500">Расстояние</span>
            <span class="font-medium text-gray-800">{{ displayDistance(driver.distance) }}</span>
          </div>
        </div>

        <div class="mt-4 flex w-full gap-3">
          <button
            class="flex-1 rounded-xl bg-blue-500 p-3 text-center font-medium text-white transition hover:bg-blue-600"
            @click="showOnMap(driver.from_location.lat, driver.from_location.lon)"
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
