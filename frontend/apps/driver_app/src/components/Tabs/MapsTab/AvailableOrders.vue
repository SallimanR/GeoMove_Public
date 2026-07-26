<script setup lang="ts">
import { ref, onMounted, watch, inject, computed } from "vue";
import { useStore } from "@nanostores/vue";
import { $isAuthenticated } from "auth";
import { $mapInstance, Marker } from "@geomove/maps";
import { useDriverProfile } from "../../../stores/driverStore";
import { useOrderNotifications } from "order";
import { useOrders } from "../../../composables/useOrders";
import { CLOSE_BOTTOM_PANEL_KEY } from "../../../injectionKeys";
import OrderCard from "../../OrderCard.vue";
import Button from "primevue/button";

const isAuthenticated = useStore($isAuthenticated);
const { exists: driverExists } = useDriverProfile();

const {
  availableOrders,
  myOrders,
  declinedOrders,
  loading,
  error,
  acceptingId,
  decliningId,
  restoringId,
  cancellingId,
  isMyOrder,
  fetchOrders,
  acceptOrder,
  declineOrder,
  undoDeclineOrder,
  cancelOrder,
} = useOrders();

type Tab = "available" | "mine" | "declined";
const activeTab = ref<Tab>("available");

const tabs = computed(() => [
  { key: "available" as Tab, label: "Доступные", count: availableOrders.value.length },
  { key: "mine" as Tab, label: "Мои", count: myOrders.value.length },
  { key: "declined" as Tab, label: "Отклонённые", count: declinedOrders.value.length },
]);

const currentOrders = computed(() => {
  switch (activeTab.value) {
    case "available":
      return availableOrders.value;
    case "mine":
      return myOrders.value;
    case "declined":
      return declinedOrders.value;
  }
});

const closeBottomPanel = inject(CLOSE_BOTTOM_PANEL_KEY);

let locationMarker: Marker | null = null;

function handleShowOnMap(lat: number, lon: number) {
  const map = $mapInstance.get();
  if (!map) return;

  closeBottomPanel?.();

  if (locationMarker) locationMarker.remove();
  locationMarker = new Marker({ color: "#3b82f6" }).setLngLat([lon, lat]).addTo(map);

  map.flyTo({ center: [lon, lat], zoom: 15 });
}

useOrderNotifications(() => {
  console.log("[AvailableOrders] notification received");
  fetchOrders();
});

onMounted(() => {
  if (isAuthenticated.value && driverExists.value) fetchOrders();
});

watch([isAuthenticated, driverExists], ([auth, drv]) => {
  if (auth && drv) fetchOrders();
});
</script>

<template>
  <div class="flex h-full max-h-full flex-col gap-3 p-4">
    <h3 class="text-center font-semibold">Заказы</h3>

    <div class="flex gap-1 rounded-lg bg-gray-100 p-1">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        @click="activeTab = tab.key"
        class="flex-1 rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
        :class="
          activeTab === tab.key
            ? 'bg-white text-gray-900 shadow-sm'
            : 'text-gray-500 hover:text-gray-700'
        "
      >
        {{ tab.label }} ({{ tab.count }})
      </button>
    </div>

    <div class="flex-1 overflow-y-auto">
      <div v-if="!isAuthenticated" class="py-8 text-center text-gray-500">
        Войдите в аккаунт чтобы видеть заказы
      </div>

      <div v-else-if="!driverExists" class="py-8 text-center text-gray-500">
        Сначала создайте профиль водителя
      </div>

      <div v-else-if="loading" class="py-8 text-center text-gray-500">Загрузка...</div>

      <div v-else-if="error" class="py-8 text-center text-red-500">{{ error }}</div>

      <div v-else-if="currentOrders.length === 0" class="py-8 text-center text-gray-500">
        {{ activeTab === "declined" ? "Нет отклонённых заказов" : "Нет доступных заказов" }}
      </div>

      <div v-else class="flex flex-col gap-2">
        <OrderCard
          v-for="order in currentOrders"
          :key="order.id"
          :order="order"
          :isMyOrder="activeTab === 'mine' || isMyOrder(order)"
          :showAcceptButton="activeTab === 'available' && !isMyOrder(order)"
          :showRestoreButton="activeTab === 'declined'"
          :showCancelButton="
            activeTab === 'mine' && (order.status === 'accepted' || order.status === 'in_progress')
          "
          :accepting="acceptingId === order.id"
          :declining="decliningId === order.id"
          :restoring="restoringId === order.id"
          :cancelling="cancellingId === order.id"
          :onAccept="(id) => acceptOrder(id)"
          :onDecline="(id) => declineOrder(id)"
          :onRestore="(id) => undoDeclineOrder(id)"
          :onCancel="(id, reason) => cancelOrder(id, reason)"
          @showOnMap="handleShowOnMap"
        />
      </div>
    </div>

    <Button class="w-full !border-gray-200 !bg-gray-200 !text-gray-700" @click="fetchOrders">
      Обновить
    </Button>
  </div>
</template>
