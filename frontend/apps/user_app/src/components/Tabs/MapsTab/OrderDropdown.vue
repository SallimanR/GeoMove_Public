<script setup lang="ts">
import { computed, ref, watch, onMounted, onUnmounted, toRaw } from "vue";
import { useStore } from "@nanostores/vue";
import {
  $startPoint,
  $endPoint,
  $startAddress,
  $endAddress,
  $mapInstance,
  useRouteDisplay,
} from "@geomove/maps";
import { addressToText, getReverseGeocoding } from "@geomove/geo";
import { Marker } from "maplibre-gl";
import { $orders, $orderRoute, setOrders, type OrderRoute } from "order/store/orderStore.ts";
import { orderClient } from "order/api/client.ts";
import type { Order } from "order";
import CreateOrder from "./CreateOrder.vue";
import DisplayOrder from "./DisplayOrder.vue";
import EditOrder from "./EditOrder.vue";

const orders = useStore($orders);
const open = ref(false);
const editing = ref(false);
const picking = ref(false);

onMounted(async () => {
  const { data } = await orderClient.GET("/order/my", {
    params: { query: { role: "customer" } },
  });
  if (data?.orders) setOrders(data.orders as Order[]);
});

const activeOrder = computed<Order | null>(() => {
  return orders.value.find(
    (o) => !["completed", "cancelled"].includes(o.status),
  ) ?? null;
});

const startMarker = new Marker({ draggable: true, color: "#40fc0c" });
const endMarker = new Marker({ draggable: true, color: "#fc5b55" });

function getMap() {
  const raw = toRaw($mapInstance.get());
  return raw as unknown as import("maplibre-gl").Map | null;
}

$startPoint.subscribe((p) => {
  const map = getMap();
  if (p && map) {
    startMarker.setLngLat([p.lon, p.lat]).addTo(map);
  } else {
    startMarker.remove();
  }
});

$endPoint.subscribe((p) => {
  const map = getMap();
  if (p && map) {
    endMarker.setLngLat([p.lon, p.lat]).addTo(map);
  } else {
    endMarker.remove();
  }
});

startMarker.on("dragend", () => {
  const map = getMap();
  if (!map) return;
  const pos = startMarker.getLngLat();
  $startPoint.set({ lat: pos.lat, lon: pos.lng });
  getReverseGeocoding(pos.lat, pos.lng).then((r) => {
    $startAddress.set(r);
    $orderRoute.set({
      ...$orderRoute.get(),
      fromLat: pos.lat,
      fromLon: pos.lng,
      fromText: addressToText(r),
    } as OrderRoute);
  });
});

endMarker.on("dragend", () => {
  const map = getMap();
  if (!map) return;
  const pos = endMarker.getLngLat();
  $endPoint.set({ lat: pos.lat, lon: pos.lng });
  getReverseGeocoding(pos.lat, pos.lng).then((r) => {
    $endAddress.set(r);
    $orderRoute.set({
      ...$orderRoute.get(),
      toLat: pos.lat,
      toLon: pos.lng,
      toText: addressToText(r),
    } as OrderRoute);
  });
});

function openPanel() {
  open.value = true;
}

function onPickStart() {
  open.value = false;
  picking.value = true;
}

function onPickDone() {
  open.value = true;
  picking.value = false;
}

function setPointsFromOrder(order: Order | null) {
  if (order) {
    $startPoint.set({ lat: order.from_lat, lon: order.from_lon });
    $endPoint.set({ lat: order.to_lat, lon: order.to_lon });
    if (order.from_address) {
      $startAddress.set({
        properties: { name: order.from_address },
        geometry: { type: "Point", coordinates: [order.from_lon, order.from_lat] },
      });
    }
    if (order.to_address) {
      $endAddress.set({
        properties: { name: order.to_address },
        geometry: { type: "Point", coordinates: [order.to_lon, order.to_lat] },
      });
    }
  } else {
    $startPoint.set(null);
    $endPoint.set(null);
    $startAddress.set(null);
    $endAddress.set(null);
  }
}

watch(activeOrder, (order) => {
  editing.value = false;
  setPointsFromOrder(order ?? null);
});

function onCancelPick() {
  open.value = true;
  picking.value = false;
}

const unsub = $mapInstance.subscribe((map) => {
  if (map) {
    const m = map as unknown as import("maplibre-gl").Map;
    if (!m.isStyleLoaded()) {
      m.once("style.load", () => useRouteDisplay(m));
    } else {
      useRouteDisplay(m);
    }
    const sp = $startPoint.get();
    if (sp) startMarker.setLngLat([sp.lon, sp.lat]).addTo(m);
    const ep = $endPoint.get();
    if (ep) endMarker.setLngLat([ep.lon, ep.lat]).addTo(m);

    if (!activeOrder.value) {
      const route = $orderRoute.get();
      if (route && !sp && !ep) {
        if (route.fromLat && route.fromLon) {
          $startPoint.set({ lat: route.fromLat, lon: route.fromLon });
          $startAddress.set({
            properties: { name: route.fromText },
            geometry: { type: "Point", coordinates: [route.fromLon, route.fromLat] },
          });
        }
        if (route.toLat && route.toLon) {
          $endPoint.set({ lat: route.toLat, lon: route.toLon });
          $endAddress.set({
            properties: { name: route.toText },
            geometry: { type: "Point", coordinates: [route.toLon, route.toLat] },
          });
        }
      }
    }
  }
});

onUnmounted(() => {
  unsub();
  startMarker.remove();
  endMarker.remove();
});
</script>

<template>
  <div
    v-if="!open && !picking"
    class="absolute bottom-4 left-1/2 -translate-x-1/2 w-full max-w-100 pointer-events-auto"
  >
    <div
      @click="openPanel"
      class="bg-white rounded-xl p-3 shadow-lg cursor-pointer text-gray-500 hover:bg-gray-50 transition text-center"
    >
      {{ activeOrder ? "Мой заказ" : "Создать заказ" }}
    </div>
  </div>

  <div v-show="open" class="absolute inset-0 z-100 pointer-events-auto flex flex-col" @click="open = false">
    <div class="absolute inset-0 bg-black/30" />
    <div class="relative flex-1 flex items-center justify-center">
      <span class="text-white text-lg font-medium select-none">
        Нажмите чтобы закрыть
      </span>
    </div>

    <div
      class="relative w-full max-w-300 rounded-t-2xl bg-white overflow-y-auto max-h-[80vh] p-4 mx-auto"
      @click.stop
    >
      <template v-if="activeOrder">
        <DisplayOrder v-if="!editing" :order="activeOrder" @edit="editing = true" />
        <EditOrder v-else :order="activeOrder" @back="editing = false" @pick-start="onPickStart" @pick-done="onPickDone" @cancel-pick="onCancelPick" />
      </template>
      <CreateOrder
        v-else
        @pick-start="onPickStart"
        @pick-done="onPickDone"
        @cancel-pick="onCancelPick"
        @close="open = false"
      />
    </div>
  </div>
</template>
