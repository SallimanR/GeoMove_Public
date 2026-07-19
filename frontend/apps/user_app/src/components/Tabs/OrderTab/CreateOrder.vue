<script setup lang="ts">
const emit = defineEmits<{
  created: [];
}>();

import { computed, ref } from "vue";
import { useStore } from "@nanostores/vue";

import { addressToText, displayDistance, SearchResult } from "@geomove/geo";
import {
  $endAddress,
  $endPoint,
  $routePath,
  $startAddress,
  $startPoint,
} from "@geomove/maps";

import { orderClient } from "order/api/client.ts";
import { addOrder } from "order/store/orderStore.ts";

import InputNumber from "primevue/inputnumber";
import Button from "primevue/button";

const startAddress = useStore($startAddress);
const endAddress = useStore($endAddress);
const startPoint = useStore($startPoint);
const endPoint = useStore($endPoint);
const routePath = useStore($routePath);

const startAddressText = computed(() => {
  if (!startAddress.value) return "Не выбрано";
  return addressToText(startAddress.value as SearchResult);
});
const endAddressText = computed(() => {
  if (!endAddress.value) return "Не выбрано";
  return addressToText(endAddress.value as SearchResult);
});

const distance = computed(() => {
  if (!routePath.value?.paths?.[0]?.distance) return null;
  return routePath.value.paths[0].distance;
});

const howManyWheelsBlocked = ref(0);
const isSubmitting = ref(false);
const submitError = ref<string | null>(null);
const submitSuccess = ref<string | null>(null);

const canPublish = computed(
  () => startPoint.value && endPoint.value && isSubmitting.value === false,
);

async function handlePublishOrder() {
  if (!startPoint.value || !endPoint.value) return;

  isSubmitting.value = true;
  submitError.value = null;
  submitSuccess.value = null;

  try {
    const { data, error: apiError } = await orderClient.POST("/order", {
      body: {
        from_lat: startPoint.value.lat,
        from_lon: startPoint.value.lon,
        from_address: startAddressText.value,
        to_lat: endPoint.value.lat,
        to_lon: endPoint.value.lon,
        to_address: endAddressText.value,
        how_many_wheels_blocked: howManyWheelsBlocked.value,
        total_distance_meters: distance.value
          ? Math.round(distance.value)
          : null,
      },
    });

    if (apiError) {
      submitError.value = apiError.error ?? "Ошибка при создании заказа";
    } else if (data) {
      addOrder(data);
      submitSuccess.value = "Заказ создан";
      emit("created");
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
    <div class="flex flex-col items-center gap-2 rounded-xl bg-gray-100 p-3">
      <p class="font-medium text-gray-700">Точка отправки:</p>
      <p class="rounded-xl p-1.5 text-center bg-white w-full">
        {{ startAddressText }}
      </p>
    </div>

    <div class="flex flex-col items-center gap-2 rounded-xl bg-gray-100 p-3">
      <p class="font-medium text-gray-700">Точка прибытия:</p>
      <p class="rounded-xl p-1.5 text-center bg-white w-full">
        {{ endAddressText }}
      </p>
    </div>

    <div
      v-if="distance"
      class="flex flex-col items-center gap-1 rounded-xl bg-gray-100 p-3"
    >
      <p class="font-medium text-gray-700">Расстояние:</p>
      <p class="text-lg font-semibold">{{ displayDistance(distance) }}</p>
    </div>

    <div class="flex flex-col gap-2 rounded-xl bg-gray-100 p-3">
      <label class="font-medium text-gray-700"
        >Количество заблокированных колёс</label
      >
      <InputNumber
        v-model="howManyWheelsBlocked"
        :min="1"
        :max="18"
        class="w-full"
      />
    </div>

    <div
      v-if="submitError"
      class="rounded-xl bg-red-100 p-3 text-center text-red-600 font-medium"
    >
      {{ submitError }}
    </div>

    <div
      v-if="submitSuccess"
      class="rounded-xl bg-green-100 p-3 text-center text-green-600 font-medium"
    >
      {{ submitSuccess }}
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
