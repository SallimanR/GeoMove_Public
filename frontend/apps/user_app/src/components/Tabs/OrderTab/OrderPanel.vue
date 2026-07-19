<script setup lang="ts">
import { computed, inject, onMounted, ref } from "vue";
import { useStore } from "@nanostores/vue";
import { ACTIVE_TAB_KEY } from "src/injectionKeys";

import { orderClient } from "order/api/client.ts";
import { $orders, $activeTabSetter, setOrders } from "order/store/orderStore.ts";

import CreateOrder from "./CreateOrder.vue";
import DisplayOrder from "./DisplayOrder.vue";

const activeTab = inject(ACTIVE_TAB_KEY)!;
$activeTabSetter.set((tab: string) => {
  activeTab.value = tab;
});

const orders = useStore($orders);
const loaded = ref(false);

const terminalStatuses = new Set(["completed", "cancelled"]);

const hasActiveOrder = computed(() =>
  orders.value.some((o) => !terminalStatuses.has(o.status)),
);

async function loadMyOrders() {
  try {
    const { data } = await orderClient.GET("/order/my", {
      params: { query: { role: "customer" } },
    });
    if (data?.orders) setOrders(data.orders);
  } catch {
    // TODO:
  } finally {
    loaded.value = true;
  }
}

function handleOrderCreated() {
  loadMyOrders();
}

onMounted(() => {
  loadMyOrders();
});
</script>

<template>
  <div class="overflow-y-auto h-full p-3 flex flex-col gap-3">
    <template v-if="!loaded">
      <div class="flex items-center justify-center py-8 text-gray-400">
        Загрузка...
      </div>
    </template>

    <template v-else-if="!hasActiveOrder">
      <CreateOrder @created="handleOrderCreated" />
    </template>

    <template v-else>
      <div class="rounded-xl bg-blue-50 p-3 text-center text-blue-600 text-sm">
        У вас уже есть активный заказ. Дождитесь его завершения.
      </div>
    </template>

    <DisplayOrder v-if="orders.length > 0" :orders="orders" />
  </div>
</template>
