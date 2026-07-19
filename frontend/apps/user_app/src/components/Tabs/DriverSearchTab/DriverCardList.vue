<script setup lang="ts">
import { ref } from "vue";
import DriverCard from "./DriverCard.vue";
import type { Driver } from "driver/types/driver.ts";
import { $driverStore } from "driver/store/driverStore.ts";

const driversList = ref<Driver[]>([]);
$driverStore.subscribe((drivers) => {
  driversList.value = drivers;
});
</script>

<template>
  <h2 class="text-lg font-medium text-center">Эвакуаторы</h2>

  <div class="flex flex-col gap-3 p-4">
    <div v-for="driver in driversList" :key="driver.user_id">
      <DriverCard v-bind="driver" />
    </div>
    <div
      v-if="!driversList.length"
      class="text-center text-sm text-gray-400 py-8"
    >
      Нет доступных водителей
    </div>
  </div>
</template>
