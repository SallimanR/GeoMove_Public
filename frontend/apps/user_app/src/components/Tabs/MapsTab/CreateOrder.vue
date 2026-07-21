<script setup lang="ts">
import { computed, ref } from "vue";
import { useStore } from "@nanostores/vue";
import {
  $startPoint,
  $endPoint,
  $routePath,
  $startAddress,
  $endAddress,
  MapsLocationPicker,
} from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import { displayDistance, addressToText } from "@geomove/geo";
import { orderClient } from "order/api/client.ts";
import {
  addOrder,
  $createOrderForm,
  $orderRoute,
  type OrderRoute,
} from "order/store/orderStore.ts";
import { driverClient } from "driver/api/client.ts";
import { $driverStore } from "driver/store/driverStore.ts";
import { displayDriverPopups } from "src/maps/displayDrivers.ts";
import OrderFormFields from "./OrderFormFields.vue";
import Button from "primevue/button";
import InputText from "primevue/inputtext";

const emit = defineEmits<{
  pickStart: [];
  pickEnd: [];
  pickDone: [];
  cancelPick: [];
  close: [];
}>();

const form = useStore($createOrderForm);
const routePath = useStore($routePath);
const startPoint = useStore($startPoint);
const endPoint = useStore($endPoint);
const startAddress = useStore($startAddress);
const endAddress = useStore($endAddress);

const isSubmitting = ref(false);
const submitError = ref<string | null>(null);

const fromText = computed(() =>
  startAddress.value ? addressToText(startAddress.value) : "",
);
const toText = computed(() =>
  endAddress.value ? addressToText(endAddress.value) : "",
);

const distance = computed(() => {
  if (!routePath.value?.paths?.[0]?.distance) return null;
  return routePath.value.paths[0].distance;
});

const canPublish = computed(
  () =>
    startPoint.value &&
    endPoint.value &&
    form.value.carWeightKg &&
    form.value.carLengthMeters &&
    form.value.carName &&
    !isSubmitting.value,
);

function onFromPicked(point: GeoPoint, address: string) {
  $startPoint.set(point);
  $startAddress.set({
    properties: { name: address },
    geometry: { type: "Point", coordinates: [point.lon, point.lat] },
  });
  $orderRoute.set({
    ...$orderRoute.get(),
    fromLat: point.lat,
    fromLon: point.lon,
    fromText: address,
  } as OrderRoute);
  emit("pickDone");
}

function onToPicked(point: GeoPoint, address: string) {
  $endPoint.set(point);
  $endAddress.set({
    properties: { name: address },
    geometry: { type: "Point", coordinates: [point.lon, point.lat] },
  });
  $orderRoute.set({
    ...$orderRoute.get(),
    toLat: point.lat,
    toLon: point.lon,
    toText: address,
  } as OrderRoute);
  emit("pickDone");
}

function handleClearRoute() {
  $startPoint.set(null);
  $startAddress.set(null);
  $endPoint.set(null);
  $endAddress.set(null);
  $orderRoute.set(null);
}

async function handlePublishOrder() {
  if (!startPoint.value || !endPoint.value) return;

  isSubmitting.value = true;
  submitError.value = null;

  try {
    const { data, error } = await orderClient.POST("/order", {
      body: {
        from_lat: startPoint.value.lat,
        from_lon: startPoint.value.lon,
        from_address: fromText.value,
        to_lat: endPoint.value.lat,
        to_lon: endPoint.value.lon,
        to_address: toText.value,
        how_many_wheels_blocked: form.value.wheels,
        total_distance_meters: distance.value
          ? Math.round(distance.value)
          : null,
        car_weight_kg: form.value.carWeightKg,
        car_length_meters: form.value.carLengthMeters,
        car_type: form.value.carType,
        car_name: form.value.carName,
        car_photo_url: form.value.carPhotoUrl || null,
        customer_message: form.value.customerMessage || null,
      },
    });

    if (error) {
      submitError.value = error?.error ?? "Ошибка при создании заказа";
    } else if (data) {
      addOrder(data);
      emit("close");

      const start = $startPoint.get();
      const { data: driverData } = await driverClient.POST("/driver/filter", {
        body: { user_lat: start?.lat ?? 55, user_lon: start?.lon ?? 37 },
      });
      if (driverData?.drivers) {
        $driverStore.set(driverData.drivers);
        displayDriverPopups();
      }
    }
  } catch {
    submitError.value = "Не удалось отправить заказ";
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col gap-3">
    <h3 class="font-semibold text-center">Создать заказ</h3>

    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-2">
        <InputText
          :modelValue="fromText"
          placeholder="Откуда"
          class="flex-1"
          readonly
        />
        <MapsLocationPicker @pick="onFromPicked" @click="emit('pickStart')" @cancel="emit('cancelPick')" />
      </div>

      <div class="flex items-center gap-2">
        <InputText :modelValue="toText" placeholder="Куда" class="flex-1" readonly />
        <MapsLocationPicker @pick="onToPicked" @click="emit('pickStart')" @cancel="emit('cancelPick')" />
      </div>

      <div v-if="routePath" class="flex gap-2">
        <button
          @click="handleClearRoute"
          type="button"
          class="rounded-xl bg-red-300 p-2 text-sm flex-1"
        >
          Сбросить маршрут
        </button>
      </div>

      <div
        v-if="routePath && routePath.paths[0]?.distance"
        class="text-center"
      >
        Расстояние: {{ displayDistance(routePath.paths[0].distance) }}
      </div>
    </div>

    <OrderFormFields
      :wheels="form.wheels"
      :carType="form.carType"
      :carName="form.carName"
      :carWeightKg="form.carWeightKg"
      :carLengthMeters="form.carLengthMeters"
      :carPhotoUrl="form.carPhotoUrl"
      :customerMessage="form.customerMessage"
      @update:wheels="(v) => $createOrderForm.set({ ...$createOrderForm.get(), wheels: v })"
      @update:carType="(v) => $createOrderForm.set({ ...$createOrderForm.get(), carType: v })"
      @update:carName="(v) => $createOrderForm.set({ ...$createOrderForm.get(), carName: v })"
      @update:carWeightKg="(v) => $createOrderForm.set({ ...$createOrderForm.get(), carWeightKg: v })"
      @update:carLengthMeters="(v) => $createOrderForm.set({ ...$createOrderForm.get(), carLengthMeters: v })"
      @update:carPhotoUrl="(v) => $createOrderForm.set({ ...$createOrderForm.get(), carPhotoUrl: v })"
      @update:customerMessage="(v) => $createOrderForm.set({ ...$createOrderForm.get(), customerMessage: v })"
    />

    <div
      v-if="submitError"
      class="rounded-xl bg-red-100 p-3 text-center text-red-600"
    >
      {{ submitError }}
    </div>

    <Button
      :disabled="!canPublish"
      :loading="isSubmitting"
      @click="handlePublishOrder"
      class="w-full"
      :class="canPublish ? '!bg-green-500 !border-green-500' : ''"
    >
      {{ isSubmitting ? "Отправка..." : "Создать заказ" }}
    </Button>
  </div>
</template>
