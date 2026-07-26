<script setup lang="ts">
import { computed } from "vue";
import { useStore } from "@nanostores/vue";
import { $driverDropdownOpen } from "src/stores/driverStore.ts";
import {
  type FreelyAvailableDriver,
  type Driver,
  $driverStore,
  $freelyAvailableDriverStore,
} from "driver";

import DriverCard from "./DriverCard.vue";
import FreelyAvailableDriverCard from "./FreelyAvailableDriverCard.vue";

const driverStore = useStore($driverStore);
const faDriverStore = useStore($freelyAvailableDriverStore);
const open = useStore($driverDropdownOpen);

const totalDrivers = computed(() => driverStore.value.length + faDriverStore.value.length);

function openPanel() {
  $driverDropdownOpen.set(true);
}

function closePanel() {
  $driverDropdownOpen.set(false);
}
</script>

<template>
  <div
    v-if="totalDrivers > 0 && !open"
    class="pointer-events-auto absolute bottom-16 left-1/2 w-fit max-w-100 -translate-x-1/2"
  >
    <div
      @click="openPanel"
      class="mb-2 cursor-pointer rounded-xl bg-blue-500 px-4 py-2 text-center whitespace-nowrap text-white shadow-lg transition hover:bg-blue-600"
    >
      эвакуаторов рядом: {{ totalDrivers }}
    </div>
  </div>

  <div
    v-if="totalDrivers > 0 && open"
    class="pointer-events-auto absolute inset-0 z-100 flex flex-col"
    @click="closePanel"
  >
    <div class="absolute inset-0 bg-black/30" />
    <div class="relative flex flex-1 items-center justify-center">
      <span class="text-lg font-medium text-white select-none"> Нажмите чтобы закрыть </span>
    </div>

    <div
      class="relative mx-auto max-h-[80vh] w-full max-w-300 overflow-y-auto rounded-t-2xl bg-white p-4"
      @click.stop
    >
      <h3 class="mb-3 text-center font-semibold">Доступные эвакуаторы</h3>

      <div class="flex flex-col gap-3">
        <DriverCard
          v-for="driver in driverStore as Driver[]"
          :key="'d-' + driver.user_id"
          v-bind="driver"
        />
      </div>

      <div
        v-if="!driverStore.length && !faDriverStore.length"
        class="py-8 text-center text-sm text-gray-400"
      >
        Нет доступных водителей
      </div>

      <template v-if="faDriverStore.length > 0">
        <h3 class="mt-6 mb-3 text-center font-semibold">Свободные эвакуаторы</h3>

        <div class="flex flex-col gap-3">
          <FreelyAvailableDriverCard
            v-for="fa in faDriverStore as FreelyAvailableDriver[]"
            :key="'fa-' + fa.user_id"
            v-bind="fa"
          />
        </div>
      </template>
    </div>
  </div>
</template>
