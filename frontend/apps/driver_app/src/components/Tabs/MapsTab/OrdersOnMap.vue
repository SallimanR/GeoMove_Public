<script setup lang="ts">
import { watch, ref, onBeforeUnmount } from "vue";
import { $mapInstance, addPopupToMap, removePopupsByGroup, Marker } from "@geomove/maps";
import { $isAuthenticated } from "auth";
import { useDriverProfile } from "../../../stores/driverStore";
import { useOrders } from "../../../composables/useOrders";
import { useOrderNotifications } from "order";
import OrderCard from "../../OrderCard.vue";
import OrderMapPopup from "./OrderMapPopup.vue";
import OrderClusterPopup from "./OrderClusterPopup.vue";
import type { Order } from "order";

const PREVIEW_GROUP = "order-preview-popups";
const DETAIL_GROUP = "order-detail-popups";
const CLUSTER_GROUP = "order-cluster-popups";
const CLUSTER_THRESHOLD = 11;

const { exists: driverExists } = useDriverProfile();
const {
  allOrders,
  isMyOrder,
  isDeclined,
  acceptingId,
  cancellingId,
  fetchOrders,
  acceptOrder,
  cancelOrder,
} = useOrders();

useOrderNotifications(() => {
  console.log("[OrdersOnMap] notification received");
  void fetchOrders();
});

const detailOrderId = ref<number | null>(null);
const clustering = ref(true);
let locationMarker: Marker | null = null;
let zoomUnsub: (() => void) | undefined;

function handleShowOnMap(lat: number, lon: number) {
  const map = $mapInstance.get();
  if (!map) return;

  removePopupsByGroup(PREVIEW_GROUP);
  removePopupsByGroup(DETAIL_GROUP);
  removePopupsByGroup(CLUSTER_GROUP);
  detailOrderId.value = null;

  if (locationMarker) locationMarker.remove();
  locationMarker = new Marker({ color: "#3b82f6" }).setLngLat([lon, lat]).addTo(map);

  map.flyTo({ center: [lon, lat], zoom: 15 });
  renderPopups();
}

function closeDetail() {
  removePopupsByGroup(DETAIL_GROUP);
  detailOrderId.value = null;
}

function showDetailPopup(order: Order) {
  removePopupsByGroup(DETAIL_GROUP);

  addPopupToMap(
    order.from_lat,
    order.from_lon,
    OrderCard,
    {
      order,
      isMyOrder: isMyOrder(order),
      accepting: acceptingId.value === order.id,
      cancelling: cancellingId.value === order.id,
      showAcceptButton: !isMyOrder(order),
      showCancelButton:
        isMyOrder(order) && (order.status === "accepted" || order.status === "in_progress"),
      onClose: closeDetail,
      onAccept: (orderId: number) => {
        void acceptOrder(orderId).then(() => {
          removePopupsByGroup(PREVIEW_GROUP);
          removePopupsByGroup(DETAIL_GROUP);
          removePopupsByGroup(CLUSTER_GROUP);
          detailOrderId.value = null;
          renderPopups();
        });
      },
      onCancel: (orderId: number, reason: string) => {
        void cancelOrder(orderId, reason).then(() => {
          removePopupsByGroup(PREVIEW_GROUP);
          removePopupsByGroup(DETAIL_GROUP);
          removePopupsByGroup(CLUSTER_GROUP);
          detailOrderId.value = null;
          renderPopups();
        });
      },
      onShowOnMap: (lat: number, lon: number) => {
        handleShowOnMap(lat, lon);
      },
    },
    DETAIL_GROUP,
    { offset: [0, -20], closeOnClick: false },
  );
}

function zoomInCluster(lat: number, lon: number) {
  const map = $mapInstance.get();
  if (!map) return;
  map.flyTo({ center: [lon, lat], zoom: CLUSTER_THRESHOLD + 1 });
}

function renderClusters() {
  removePopupsByGroup(PREVIEW_GROUP);
  removePopupsByGroup(DETAIL_GROUP);
  removePopupsByGroup(CLUSTER_GROUP);
  detailOrderId.value = null;

  const map = $mapInstance.get();
  if (!map) return;

  const threshold = Math.max(40, 150 - (map.getZoom() - 5) * 22);

  const projected = allOrders.value.map((order) => {
    const p = map.project([order.from_lon, order.from_lat]);
    return { order, x: p.x, y: p.y };
  });

  const clusters: { lat: number; lon: number; count: number; orders: typeof projected }[] = [];

  for (const pt of projected) {
    let found = false;
    for (const cl of clusters) {
      const cx = cl.orders.reduce((s, o) => s + o.x, 0) / cl.orders.length;
      const cy = cl.orders.reduce((s, o) => s + o.y, 0) / cl.orders.length;
      const dx = pt.x - cx;
      const dy = pt.y - cy;
      if (Math.sqrt(dx * dx + dy * dy) < threshold) {
        cl.lat = (cl.lat * cl.count + pt.order.from_lat) / (cl.count + 1);
        cl.lon = (cl.lon * cl.count + pt.order.from_lon) / (cl.count + 1);
        cl.count++;
        cl.orders.push(pt);
        found = true;
        break;
      }
    }
    if (!found) {
      clusters.push({
        lat: pt.order.from_lat,
        lon: pt.order.from_lon,
        count: 1,
        orders: [pt],
      });
    }
  }

  for (const cl of clusters) {
    addPopupToMap(
      cl.lat,
      cl.lon,
      OrderClusterPopup,
      { count: cl.count, onZoomIn: () => zoomInCluster(cl.lat, cl.lon) },
      CLUSTER_GROUP,
      { offset: [0, -20], closeOnClick: false },
    );
  }
}

function renderSingles() {
  removePopupsByGroup(PREVIEW_GROUP);
  removePopupsByGroup(DETAIL_GROUP);
  removePopupsByGroup(CLUSTER_GROUP);
  detailOrderId.value = null;

  for (const order of allOrders.value) {
    addPopupToMap(
      order.from_lat,
      order.from_lon,
      OrderMapPopup,
      {
        order,
        isMyOrder: isMyOrder(order),
        isDeclined: isDeclined(order),
        onClick: () => {
          if (detailOrderId.value === order.id) {
            closeDetail();
            return;
          }
          detailOrderId.value = order.id;
          showDetailPopup(order);
        },
      },
      PREVIEW_GROUP,
      { offset: [0, -20], closeOnClick: false },
    );
  }
}

function renderPopups() {
  if (clustering.value) {
    renderClusters();
  } else {
    renderSingles();
  }
}

const unsubMap = $mapInstance.subscribe((map) => {
  if (zoomUnsub) zoomUnsub();
  if (map) {
    clustering.value = map.getZoom() < CLUSTER_THRESHOLD;
    zoomUnsub = () => {
      map.off("zoom", () => {});
    };
    map.on("zoom", () => {
      const z = map.getZoom();
      const wasClustering = clustering.value;
      clustering.value = z < CLUSTER_THRESHOLD;
      if (clustering.value !== wasClustering) {
        renderPopups();
      }
    });
    renderPopups();
  }
});

const unsubAuth = $isAuthenticated.subscribe((auth) => {
  if (auth && driverExists.value) {
    void fetchOrders();
  } else {
    removePopupsByGroup(PREVIEW_GROUP);
    removePopupsByGroup(DETAIL_GROUP);
    removePopupsByGroup(CLUSTER_GROUP);
    detailOrderId.value = null;
  }
});

watch(allOrders, renderPopups, { immediate: true });

watch(
  driverExists,
  (drv) => {
    if (drv && $isAuthenticated.get()) {
      void fetchOrders();
    } else {
      removePopupsByGroup(PREVIEW_GROUP);
      removePopupsByGroup(DETAIL_GROUP);
      removePopupsByGroup(CLUSTER_GROUP);
      detailOrderId.value = null;
    }
  },
  { immediate: true },
);

onBeforeUnmount(() => {
  if (locationMarker) locationMarker.remove();
  if (zoomUnsub) zoomUnsub();
  removePopupsByGroup(PREVIEW_GROUP);
  removePopupsByGroup(DETAIL_GROUP);
  removePopupsByGroup(CLUSTER_GROUP);
  unsubMap();
  unsubAuth();
});
</script>

<template>
  <div />
</template>
