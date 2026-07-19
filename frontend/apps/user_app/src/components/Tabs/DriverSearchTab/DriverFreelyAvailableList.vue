<script setup lang="ts">
import { $freelyAvailableDriverStore } from "driver/store/driverStore.ts";
import { FreelyAvailableDriver } from "driver/types/freelyAvailable.ts";
import { ref } from "vue";

const faDrivers = ref<FreelyAvailableDriver[]>([]);
$freelyAvailableDriverStore.subscribe((d) => {
  faDrivers.value = d;
});

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString("ru-RU");
}

function formatPrice(val: number | null): string {
  if (val == null) return "—";
  return `${val} ₽/км`;
}

function locationText(loc: {
  lat: number;
  lon: number;
  address?: string;
}): string {
  return loc.address || `${loc.lat.toFixed(5)}, ${loc.lon.toFixed(5)}`;
}
</script>

<template>
  <h2 class="text-lg font-medium text-center">Свободные эвакуаторы</h2>

  <div
    v-for="d in faDrivers"
    :key="'fa-' + d.user_id"
    class="bg-gray-100 rounded-lg p-3 flex flex-col gap-2"
  >
    <div class="flex items-center justify-between">
      <span class="font-medium text-gray-800">{{ d.name }}</span>
      <span v-if="d.rating" class="text-yellow-500 text-sm"
        >★ {{ d.rating?.toFixed(1) }}</span
      >
    </div>

    <div class="flex flex-col gap-1 bg-white/60 rounded-lg p-2">
      <div class="flex justify-between text-xs">
        <span class="text-gray-500">Период:</span>
        <span>{{ formatDate(d.from_date) }} — {{ formatDate(d.to_date) }}</span>
      </div>
      <div class="flex justify-between text-xs">
        <span class="text-gray-500">Тариф:</span>
        <span>{{ formatPrice(d.tariff_per_km) }}</span>
      </div>
      <div class="flex justify-between text-xs">
        <span class="text-gray-500">Попутный заказ:</span>
        <span>{{ d.en_route_order ? "Да" : "Нет" }}</span>
      </div>
    </div>

    <div class="flex flex-col gap-0.5">
      <p class="text-xs text-gray-500">Точка отправления:</p>
      <p class="text-xs">{{ locationText(d.from_location) }}</p>
    </div>

    <div class="text-xs text-gray-500">
      Расстояние: {{ d.distance.toFixed(0) }} м
    </div>
  </div>

  <div v-if="!faDrivers.length" class="text-center text-sm text-gray-400 py-2">
    Нет свободных эвакуаторов
  </div>
</template>
