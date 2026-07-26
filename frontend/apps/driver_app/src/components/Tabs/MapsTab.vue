<script setup lang="ts">
import { ref, computed, provide, inject } from "vue";
import { $mapInstance, Maps, MapsOverlayControls, useRouteDisplay } from "@geomove/maps";
import OrdersOnMap from "./MapsTab/OrdersOnMap.vue";
import AvailableOrders from "./MapsTab/AvailableOrders.vue";
import { useOrders } from "../../composables/useOrders";
import { ACTIVE_TAB_KEY, CLOSE_BOTTOM_PANEL_KEY } from "../../injectionKeys";
import { useDriverProfile } from "src/stores/driverStore";

const styleApi = import.meta.env.VITE_STYLE_API;
const { allOrders } = useOrders();

const open = ref(false);

const activeTab = inject(ACTIVE_TAB_KEY)!;
const { exists: driverExists } = useDriverProfile();

const availableCount = computed(() => allOrders.value.filter((o) => o.status === "pending").length);

$mapInstance.subscribe((map) => {
  if (!map) return;
  useRouteDisplay(map as maplibregl.MapLibreMap);
});

function openPanel() {
  open.value = true;
}

function closePanel() {
  open.value = false;
}

provide(CLOSE_BOTTOM_PANEL_KEY, closePanel);
</script>

<template>
  <div class="relative flex h-full flex-col">
    <Maps :styleApi="styleApi" />
    <OrdersOnMap />
    <MapsOverlayControls :hideRouteInput="true" />

    <div
      v-if="!open"
      class="pointer-events-auto absolute bottom-4 left-1/2 w-full max-w-[25rem] -translate-x-1/2"
    >
      <div
        v-if="driverExists"
        @click="openPanel"
        class="cursor-pointer rounded-xl bg-white p-3 text-center text-gray-500 shadow-lg transition hover:bg-gray-50"
      >
        Доступно заказов: {{ availableCount }}
      </div>
      <div
        v-else
        @click="activeTab = 'profileTab'"
        class="cursor-pointer rounded-xl bg-white p-3 text-center shadow-lg transition hover:bg-gray-50"
      >
        Чтобы смотреть заказы, необходимо
        <div class="rounded-xl bg-gray-100 py-2 font-medium text-green-500">
          Войти в профиль водителя
        </div>
      </div>
    </div>

    <div
      v-show="open"
      class="pointer-events-auto absolute inset-0 z-[100] flex flex-col"
      @click="open = false"
    >
      <div class="absolute inset-0 bg-black/30" />
      <div class="relative flex flex-1 items-center justify-center">
        <span class="text-lg font-medium text-white select-none">Нажмите чтобы закрыть</span>
      </div>

      <div
        class="relative mx-auto max-h-[80vh] w-full overflow-y-auto rounded-t-2xl bg-white p-4"
        @click.stop
      >
        <AvailableOrders />
      </div>
    </div>
  </div>
</template>
