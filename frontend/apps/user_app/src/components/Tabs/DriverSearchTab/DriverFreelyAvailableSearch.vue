<script setup lang="ts">
import { onMounted, ref } from "vue";
import { freelyAvailableDriverClient } from "driver/api/client.ts";
import { $freelyAvailableDriverStore } from "driver/store/driverStore.ts";
import type {
  GetFreelyAvailableDriversRequest,
  FreelyAvailableDriver,
} from "driver/types/freelyAvailable.ts";
import { $coords, $startPoint } from "@geomove/maps";
import Checkbox from "primevue/checkbox";
import InputNumber from "primevue/inputnumber";

const center = getCenter();
const faReq = ref<GetFreelyAvailableDriversRequest>({
  user_lat: center.lat,
  user_lon: center.lon,
});

const faDrivers = ref<FreelyAvailableDriver[]>([]);
$freelyAvailableDriverStore.subscribe((d) => {
  faDrivers.value = d as FreelyAvailableDriver[];
});

const loading = ref(false);
const expanded = ref(false);

function getCenter() {
  const c = $startPoint.value || $coords.get().center;
  return c || { lat: 55, lon: 37 };
}

async function search() {
  loading.value = true;
  try {
    const center = getCenter();
    faReq.value.user_lat = center.lat;
    faReq.value.user_lon = center.lon;
    const { data, error } = await freelyAvailableDriverClient.POST(
      "/driver/freely-available/search",
      {
        body: faReq.value,
      },
    );
    if (!error && data?.drivers) {
      $freelyAvailableDriverStore.set(data.drivers);
    }
  } finally {
    loading.value = false;
  }
}

function toggle() {
  expanded.value = !expanded.value;
  if (expanded.value) {
    search();
  }
}

onMounted(() => {
  search();
});
</script>

<template>
  <div class="bg-gray-200 rounded-xl overflow-hidden">
    <div
      class="flex items-center justify-between p-4 cursor-pointer hover:bg-gray-300 transition-colors"
      @click="toggle"
    >
      <div class="flex items-center gap-2">
        <span class="text-sm">{{ expanded ? "▼" : "▶" }}</span>
        <span class="font-medium text-gray-800">Свободные эвакуаторы</span>
      </div>
      <span
        v-if="faDrivers.length"
        class="text-xs text-gray-500 bg-gray-300 px-2 py-0.5 rounded-full"
      >
        {{ faDrivers.length }}
      </span>
    </div>

    <div v-if="expanded" class="px-4 pb-4 flex flex-col gap-3">
      <div class="flex items-center gap-2">
        <Checkbox v-model="faReq.en_route_order" binary input-id="faEnRoute" />
        <label for="faEnRoute" class="text-sm text-gray-700"
          >Попутный заказ</label
        >
      </div>

      <div class="flex gap-2">
        <InputNumber
          v-model="faReq.min_tariff"
          :min="0"
          placeholder="Мин. ₽/км"
          class="w-full"
        />
        <InputNumber
          v-model="faReq.max_tariff"
          :min="0"
          placeholder="Макс. ₽/км"
          class="w-full"
        />
      </div>

      <button
        @click="search"
        :disabled="loading"
        class="w-full rounded-xl text-center p-2 bg-green-300 hover:bg-green-400 disabled:opacity-50 transition-colors text-sm font-medium"
      >
        {{ loading ? "Поиск..." : "Найти" }}
      </button>
    </div>
  </div>
</template>
