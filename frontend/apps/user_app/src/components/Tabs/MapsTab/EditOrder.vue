<script setup lang="ts">
import { ref, watch } from "vue";
import { updateOrderInStore } from "order/store/orderStore.ts";
import { orderClient } from "order/api/client.ts";
import type { Order } from "order";
import { $startPoint, $endPoint, $startAddress, $endAddress, MapsLocationPicker } from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import OrderFormFields from "./OrderFormFields.vue";
import Button from "primevue/button";
import InputText from "primevue/inputtext";

const props = defineProps<{
  order: Order;
}>();

const emit = defineEmits<{
  back: [];
  pickStart: [];
  pickDone: [];
  cancelPick: [];
}>();

const isSubmitting = ref(false);
const submitError = ref<string | null>(null);

const editFromLat = ref(0);
const editFromLon = ref(0);
const editFromAddress = ref("");
const editToLat = ref(0);
const editToLon = ref(0);
const editToAddress = ref("");
const editWheels = ref(0);
const editCarType = ref("Другое");
const editCarName = ref("");
const editCarWeightKg = ref(0);
const editCarLengthMeters = ref(0);
const editCarPhotoUrl = ref("");
const editCustomerMessage = ref("");

function populateFromOrder(o: Order | null) {
  if (!o) return;
  editFromLat.value = o.from_lat;
  editFromLon.value = o.from_lon;
  editFromAddress.value = o.from_address ?? "";
  editToLat.value = o.to_lat;
  editToLon.value = o.to_lon;
  editToAddress.value = o.to_address ?? "";
  editWheels.value = o.how_many_wheels_blocked;
  editCarType.value = (o as any).car_type ?? "Другое";
  editCarName.value = (o as any).car_name ?? "";
  editCarWeightKg.value = (o as any).car_weight_kg ?? 0;
  editCarLengthMeters.value = (o as any).car_length_meters ?? 0;
  editCarPhotoUrl.value = (o as any).car_photo_url ?? "";
  editCustomerMessage.value = (o as any).customer_message ?? "";
}
watch(() => props.order, (o) => populateFromOrder(o), { immediate: true });

function onFromPicked(point: GeoPoint, address: string) {
  editFromLat.value = point.lat;
  editFromLon.value = point.lon;
  editFromAddress.value = address;
  $startPoint.set(point);
  $startAddress.set({
    properties: { name: address },
    geometry: { type: "Point", coordinates: [point.lon, point.lat] },
  });
  emit("pickDone");
}

function onToPicked(point: GeoPoint, address: string) {
  editToLat.value = point.lat;
  editToLon.value = point.lon;
  editToAddress.value = address;
  $endPoint.set(point);
  $endAddress.set({
    properties: { name: address },
    geometry: { type: "Point", coordinates: [point.lon, point.lat] },
  });
  emit("pickDone");
}

async function handleSave() {
  if (!props.order) return;
  isSubmitting.value = true;
  submitError.value = null;

  try {
    const o = props.order;
    const { data, error } = await orderClient.PUT("/order/{order_id}", {
      params: { path: { order_id: o.id } },
      body: {
        from_lat: editFromLat.value,
        from_lon: editFromLon.value,
        from_address: editFromAddress.value || undefined,
        to_lat: editToLat.value,
        to_lon: editToLon.value,
        to_address: editToAddress.value || undefined,
        how_many_wheels_blocked: editWheels.value,
        car_weight_kg: editCarWeightKg.value,
        car_length_meters: editCarLengthMeters.value,
        car_type: editCarType.value,
        car_name: editCarName.value,
        car_photo_url: editCarPhotoUrl.value || null,
        customer_message: editCustomerMessage.value || null,
      },
    });

    if (error) {
      submitError.value = error?.error ?? "Ошибка при сохранении";
    } else if (data) {
      updateOrderInStore(data);
      emit("back");
    }
  } catch {
    submitError.value = "Не удалось сохранить изменения";
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col gap-3">
    <div class="flex items-center justify-between">
      <h3 class="font-semibold">Редактирование заказа #{{ order.id }}</h3>
      <Button
        @click="emit('back')"
        class="!bg-gray-200 !border-gray-200 !text-gray-700 !text-sm !px-3 !py-1"
      >
        Назад
      </Button>
    </div>

    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-2">
        <InputText
          :modelValue="editFromAddress || `${editFromLat}, ${editFromLon}`"
          placeholder="Откуда"
          class="flex-1"
          readonly
        />
        <MapsLocationPicker @pick="onFromPicked" @click="emit('pickStart')" @cancel="emit('cancelPick')" />
      </div>

      <div class="flex items-center gap-2">
        <InputText
          :modelValue="editToAddress || `${editToLat}, ${editToLon}`"
          placeholder="Куда"
          class="flex-1"
          readonly
        />
        <MapsLocationPicker @pick="onToPicked" @click="emit('pickStart')" @cancel="emit('cancelPick')" />
      </div>
    </div>

    <OrderFormFields
      :wheels="editWheels"
      :carType="editCarType"
      :carName="editCarName"
      :carWeightKg="editCarWeightKg"
      :carLengthMeters="editCarLengthMeters"
      :carPhotoUrl="editCarPhotoUrl"
      :customerMessage="editCustomerMessage"
      @update:wheels="editWheels = $event"
      @update:carType="editCarType = $event"
      @update:carName="editCarName = $event"
      @update:carWeightKg="editCarWeightKg = $event"
      @update:carLengthMeters="editCarLengthMeters = $event"
      @update:carPhotoUrl="editCarPhotoUrl = $event"
      @update:customerMessage="editCustomerMessage = $event"
    />

    <div
      v-if="submitError"
      class="rounded-xl bg-red-100 p-3 text-center text-red-600"
    >
      {{ submitError }}
    </div>

    <Button
      :loading="isSubmitting"
      @click="handleSave"
      class="w-full !bg-green-500 !border-green-500"
    >
      {{ isSubmitting ? "Сохранение..." : "Сохранить изменения" }}
    </Button>
  </div>
</template>
