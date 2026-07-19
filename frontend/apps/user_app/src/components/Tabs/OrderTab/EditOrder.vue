<script setup lang="ts">
import { onMounted, ref } from "vue";
import { $locationPicking, setPickCallback } from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import type { Order } from "order/types/order.ts";
import { orderClient } from "order/api/client.ts";
import {
  $orders,
  $pendingEdit,
  $activeTabSetter,
  setOrders,
  type LocationEntry,
} from "order/store/orderStore.ts";
import { useStore } from "@nanostores/vue";
import { displayDistance } from "@geomove/geo";

import Button from "primevue/button";
import Dialog from "primevue/dialog";
import InputNumber from "primevue/inputnumber";
import InputText from "primevue/inputtext";

const props = defineProps<{
  order: Order;
}>();

const orders = useStore($orders);
const pendingEdit = useStore($pendingEdit);

const editVisible = ref(false);
const editSubmitting = ref(false);
const editError = ref<string | null>(null);
const editWheels = ref(props.order.how_many_wheels_blocked);
const editPrice = ref(props.order.price_rubles ?? 0);
const editFrom = ref<LocationEntry>({
  lat: props.order.from_lat,
  lon: props.order.from_lon,
  address: "",
});
const editTo = ref<LocationEntry>({
  lat: props.order.to_lat,
  lon: props.order.to_lon,
  address: "",
});

const cancelVisible = ref(false);
const cancelReason = ref("");
const cancelSubmitting = ref(false);
const cancelError = ref<string | null>(null);

onMounted(() => {
  if (!pendingEdit.value || pendingEdit.value.orderId !== props.order.id) return;
  editWheels.value = pendingEdit.value.wheels;
  editPrice.value = pendingEdit.value.price;
  editFrom.value = pendingEdit.value.from;
  editTo.value = pendingEdit.value.to;
  editError.value = null;
  editVisible.value = true;
});

function openEdit() {
  editWheels.value = props.order.how_many_wheels_blocked;
  editPrice.value = props.order.price_rubles ?? 0;
  editFrom.value = {
    lat: props.order.from_lat,
    lon: props.order.from_lon,
    address: props.order.from_address || "",
  };
  editTo.value = {
    lat: props.order.to_lat,
    lon: props.order.to_lon,
    address: props.order.to_address || "",
  };
  editError.value = null;
  editVisible.value = true;
}

function startMapPick(pickType: "from" | "to") {
  setPickCallback((point: GeoPoint, address: string) => {
    const target = pickType === "from" ? editFrom.value : editTo.value;
    const updated = {
      orderId: props.order.id,
      wheels: editWheels.value,
      price: editPrice.value,
      from: pickType === "from"
        ? { lat: point.lat, lon: point.lon, address }
        : editFrom.value,
      to: pickType === "to"
        ? { lat: point.lat, lon: point.lon, address }
        : editTo.value,
    };
    $pendingEdit.set(updated);

    const activeTabFn = $activeTabSetter.get();
    if (activeTabFn) activeTabFn("orderTab");
  });

  $locationPicking.set(true);

  editVisible.value = false;

  const updated = {
    orderId: props.order.id,
    wheels: editWheels.value,
    price: editPrice.value,
    from: editFrom.value,
    to: editTo.value,
  };
  $pendingEdit.set(updated);

  const activeTabFn = $activeTabSetter.get();
  if (activeTabFn) activeTabFn("mapsTab");
}

async function saveEdit() {
  editSubmitting.value = true;
  editError.value = null;

  try {
    const { data, error } = await orderClient.PUT("/order/{order_id}", {
      params: { path: { order_id: props.order.id } },
      body: {
        from_lat: editFrom.value.lat,
        from_lon: editFrom.value.lon,
        from_address: editFrom.value.address,
        to_lat: editTo.value.lat,
        to_lon: editTo.value.lon,
        to_address: editTo.value.address,
        how_many_wheels_blocked: editWheels.value,
        total_distance_meters: props.order.total_distance_meters,
        price_rubles: editPrice.value || null,
      },
    });

    if (error) {
      editError.value = error?.error ?? "Не удалось сохранить";
    } else if (data) {
      const updated = orders.value.map((o) =>
        o.id === data.id ? data : o,
      );
      setOrders(updated);
      editVisible.value = false;
      $pendingEdit.set(null);
    }
  } catch {
    editError.value = "Не удалось сохранить";
  } finally {
    editSubmitting.value = false;
  }
}

function closeEdit() {
  editVisible.value = false;
  $pendingEdit.set(null);
}

function openCancel() {
  cancelReason.value = "";
  cancelError.value = null;
  cancelVisible.value = true;
}

async function confirmCancel() {
  cancelSubmitting.value = true;
  cancelError.value = null;

  try {
    const { error } = await orderClient.PATCH("/order/{order_id}/status", {
      params: { path: { order_id: props.order.id } },
      body: {
        status: "cancelled",
        cancellation_reason: cancelReason.value || undefined,
      },
    });

    if (error) {
      cancelError.value = error?.error ?? "Не удалось отменить заказ";
    } else {
      cancelVisible.value = false;
      const updated = orders.value.map((o) =>
        o.id === props.order.id
          ? { ...o, status: "cancelled" as const, cancellation_reason: cancelReason.value || null }
          : o,
      );
      setOrders(updated);
    }
  } catch {
    cancelError.value = "Не удалось отменить заказ";
  } finally {
    cancelSubmitting.value = false;
  }
}

const editable = (status: string) =>
  status === "forming" || status === "pending";

const cancellable = (status: string) =>
  status === "accepted" || status === "in_progress";
</script>

<template>
  <div class="flex justify-end gap-2 mt-1">
    <Button
      v-if="editable(order.status)"
      severity="secondary"
      size="small"
      @click="openEdit"
    >
      Редактировать
    </Button>

    <Button
      v-if="cancellable(order.status)"
      severity="danger"
      size="small"
      @click="openCancel"
    >
      Отменить заказ
    </Button>
  </div>

  <Dialog
    v-model:visible="editVisible"
    header="Редактирование заказа"
    :modal="true"
    :closable="!editSubmitting"
    class="w-80"
    @update:visible="(val) => !val && closeEdit()"
  >
    <div class="flex flex-col gap-3">
      <div class="flex flex-col gap-1">
        <label class="text-sm text-gray-600">Точка отправки</label>
        <div class="flex items-center gap-2">
          <span class="text-xs text-gray-500 truncate flex-1">
            {{ editFrom.address || `${editFrom.lat.toFixed(5)}, ${editFrom.lon.toFixed(5)}` }}
          </span>
          <button
            @click="startMapPick('from')"
            type="button"
            class="p-2 rounded-lg bg-gray-200 hover:bg-gray-300 text-sm shrink-0"
          >
            На карте
          </button>
        </div>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm text-gray-600">Точка прибытия</label>
        <div class="flex items-center gap-2">
          <span class="text-xs text-gray-500 truncate flex-1">
            {{ editTo.address || `${editTo.lat.toFixed(5)}, ${editTo.lon.toFixed(5)}` }}
          </span>
          <button
            @click="startMapPick('to')"
            type="button"
            class="p-2 rounded-lg bg-gray-200 hover:bg-gray-300 text-sm shrink-0"
          >
            На карте
          </button>
        </div>
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm text-gray-600">Количество заблокированных колёс</label>
        <InputNumber
          v-model="editWheels"
          :min="1"
          :max="18"
          :disabled="editSubmitting"
          class="w-full"
        />
      </div>

      <div class="flex flex-col gap-1">
        <label class="text-sm text-gray-600">Цена (₽)</label>
        <InputNumber
          v-model="editPrice"
          :min="0"
          :disabled="editSubmitting"
          class="w-full"
        />
      </div>

      <div class="text-xs text-gray-400">
        Расстояние: {{ displayDistance(props.order.total_distance_meters ?? 0) }}
      </div>

      <div
        v-if="editError"
        class="text-sm text-red-500"
      >
        {{ editError }}
      </div>

      <div class="flex gap-2 justify-end">
        <Button
          severity="secondary"
          size="small"
          :disabled="editSubmitting"
          @click="closeEdit"
        >
          Отмена
        </Button>
        <Button
          size="small"
          :loading="editSubmitting"
          @click="saveEdit"
        >
          Сохранить
        </Button>
      </div>
    </div>
  </Dialog>

  <Dialog
    v-model:visible="cancelVisible"
    header="Отмена заказа"
    :modal="true"
    :closable="!cancelSubmitting"
    class="w-80"
  >
    <div class="flex flex-col gap-3">
      <label class="text-sm text-gray-600">Укажите причину отмены:</label>
      <InputText
        v-model="cancelReason"
        :disabled="cancelSubmitting"
        placeholder="Например: передумал, ошибка в адресе..."
        class="w-full"
      />

      <div
        v-if="cancelError"
        class="text-sm text-red-500"
      >
        {{ cancelError }}
      </div>

      <div class="flex gap-2 justify-end">
        <Button
          severity="secondary"
          size="small"
          :disabled="cancelSubmitting"
          @click="cancelVisible = false"
        >
          Назад
        </Button>
        <Button
          severity="danger"
          size="small"
          :loading="cancelSubmitting"
          @click="confirmCancel"
        >
          Отменить
        </Button>
      </div>
    </div>
  </Dialog>
</template>
