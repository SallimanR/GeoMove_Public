<script setup lang="ts">
import { useStore } from "@nanostores/vue";
import { $user } from "auth";
import { resolveImageUrl } from "../../../stores/driverStore";
import { usePushNotifications } from "notifications";
import type { Driver } from "driver";

function formatTime(iso: string | undefined | null): string {
  if (!iso) return "—";
  return iso.split("T")[1]?.slice(0, 5) ?? "—";
}

defineProps<{
  driver: Driver;
}>();

const emit = defineEmits<{
  edit: [];
  logout: [];
}>();

const user = useStore($user);
const { isSupported, isSubscribed } = usePushNotifications();
</script>

<template>
  <div class="flex flex-col items-center gap-4 p-4">
    <div class="h-24 w-24 overflow-hidden rounded-full ring-2 ring-gray-300">
      <img
        v-if="driver.profile_image"
        :src="resolveImageUrl(driver.profile_image)"
        alt="Profile"
        class="h-full w-full object-cover"
      />
      <div v-else class="flex h-full w-full items-center justify-center bg-gray-200">
        <span class="text-3xl text-gray-500">{{
          driver.name?.charAt(0).toUpperCase() || "?"
        }}</span>
      </div>
    </div>

    <p class="text-lg font-medium">{{ driver.name }}</p>
    <p class="text-gray-500">{{ user?.email }}</p>

    <p v-if="driver.phone" class="text-gray-600">{{ driver.phone }}</p>

    <div v-if="driver.rating" class="flex items-center gap-1">
      <span class="text-yellow-500">★ {{ driver.rating.toFixed(1) }}</span>
    </div>

    <div class="flex w-full max-w-80 flex-col items-center gap-2 rounded-xl bg-gray-200 p-4">
      <div v-if="driver.work_starts || driver.work_ends" class="flex gap-2">
        <span class="text-gray-500">Часы работы:</span>
        <span> {{ formatTime(driver.work_starts) }} – {{ formatTime(driver.work_ends) }} </span>
      </div>

      <div v-if="driver.address" class="flex gap-2">
        <span class="text-gray-500">Адрес:</span>
        <span>{{ driver.address }}</span>
      </div>

      <div v-if="driver.max_car_weight_kg" class="flex gap-2">
        <span class="text-gray-500">Макс. вес авто:</span>
        <span>{{ driver.max_car_weight_kg }} кг</span>
      </div>

      <div v-if="driver.max_car_length_meters" class="flex gap-2">
        <span class="text-gray-500">Макс. длина авто:</span>
        <span>{{ driver.max_car_length_meters }} м</span>
      </div>

      <div v-if="driver.car_photo_main" class="flex flex-col items-center gap-1">
        <span class="text-gray-500">Фото авто:</span>
        <img
          :src="resolveImageUrl(driver.car_photo_main)"
          class="max-h-48 max-w-full rounded-lg object-contain"
        />
      </div>
      <div v-if="driver.car_photos?.length" class="flex flex-col items-center gap-1">
        <span class="text-gray-500">Доп. фото:</span>
        <div class="flex flex-wrap justify-center gap-2">
          <img
            v-for="(url, idx) in driver.car_photos"
            :key="idx"
            :src="resolveImageUrl(url)"
            class="h-24 w-24 rounded-lg object-cover"
          />
        </div>
      </div>
    </div>

    <p v-if="isSupported && isSubscribed" class="text-sm text-green-600">Уведомления включены</p>

    <div class="flex w-full max-w-60 flex-col gap-2">
      <button
        @click="emit('edit')"
        class="w-full rounded-lg bg-blue-400 px-4 py-2 hover:bg-blue-500"
      >
        Редактировать
      </button>

      <button
        @click="emit('logout')"
        class="w-full rounded-lg bg-red-400 px-4 py-2 hover:bg-red-500"
      >
        Выйти из аккаунта
      </button>
    </div>
  </div>
</template>
