<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { driverClient } from "driver/api/client.ts";
import type { GetFilteredDriversRequest } from "driver/types/driver.ts";
import { $driverStore } from "driver/store/driverStore.ts";

import InputNumber from "primevue/inputnumber";

import { TimePicker } from "ui";

import { $coords, $mapInstance, $startAddress, $startPoint, removeAllPopups, MAP_CENTER_LAT, MAP_CENTER_LON } from "@geomove/maps";

import { addressToText, getReverseGeocoding, SearchResult } from "@geomove/geo";

import { displayDriverPopups } from "src/maps/displayDrivers.ts";
import { useStore } from "@nanostores/vue";

const defaultReqParams = <GetFilteredDriversRequest>{
  user_lat: MAP_CENTER_LAT,
  user_lon: MAP_CENTER_LON,
};
const reqParams = ref(defaultReqParams);

const mapInstance = useStore($mapInstance);

const fetchingError = ref(false);

async function handleFind() {
  removeAllPopups();

  let center = $startPoint.value;
  if (!center) {
    const coordsCenter = $coords.get().center;
    if (coordsCenter) {
      center = coordsCenter;
    } else {
      center = { lat: 55, lon: 37 };
    }
  }

  const req = reqParams.value;
  req.user_lat = center.lat;
  req.user_lon = center.lon;

  try {
    const { data, error } = await driverClient.POST("/driver/filter", {
      body: req,
    });
    if (!error && data?.drivers) {
      $driverStore.set(data.drivers);
      displayDriverPopups();
      fetchingError.value = false;
    } else {
      fetchingError.value = true;
    }
  } catch {
    fetchingError.value = true;
  }
}

function handleClear() {
  reqParams.value = defaultReqParams;
  handleFind();
}

const currentAddressText = ref("");

const startAddress = useStore($startAddress);

async function updateAddress() {
  try {
    if (startAddress.value) {
      currentAddressText.value = addressToText(
        startAddress.value as SearchResult,
      );
      return;
    }
    if (mapInstance.value) {
      const center = mapInstance.value.getCenter();
      const address = await getReverseGeocoding(center.lat, center.lng);
      currentAddressText.value = addressToText(address);
    } else {
      currentAddressText.value = "Карты ещё не готовы";
    }
  } catch (error) {
    currentAddressText.value = "Неудалось получить данные";
    console.error(error);
  }
}

watch(startAddress, () => updateAddress());
watch(mapInstance, () => updateAddress());

onMounted(() => {
  updateAddress();
  handleFind();
});
</script>

<template>
  <div class="text-gray-600">
    <span class="font-medium text-gray-800">Локация:</span>
    {{ currentAddressText }}
  </div>

  <div class="bg-gray-200 rounded-xl p-4 flex flex-col gap-4">
    <div class="flex flex-col gap-1.5">
      <label class="font-medium text-gray-700">Минимальный рейтинг</label>
      <InputNumber
        v-model="reqParams.min_rating"
        :min="0"
        :max="5"
        :step="0.5"
        placeholder="0 – 5"
        class="w-full"
      />
    </div>

    <div class="flex gap-3">
      <div class="flex flex-col gap-1.5 flex-1">
        <label class="font-medium text-gray-700">Начало работы</label>
        <TimePicker v-model="reqParams.work_starts" placeholder="С" />
      </div>
      <div class="flex flex-col gap-1.5 flex-1">
        <label class="font-medium text-gray-700">Конец работы</label>
        <TimePicker v-model="reqParams.work_ends" placeholder="До" />
      </div>
    </div>

    <div class="flex gap-3">
      <div
        class="flex-1 rounded-xl text-center p-2 bg-green-300 hover:bg-green-400"
        @click="handleFind"
      >
        Найти
      </div>
      <div
        class="flex-1 rounded-xl text-center p-2 bg-red-300 hover:bg-red-400"
        @click="handleClear"
      >
        Сбросить
      </div>
    </div>

    <div v-if="fetchingError" class="text-red-400 text-center font-medium">
      Ошибка, попробуйте ещё раз
    </div>
  </div>
</template>
