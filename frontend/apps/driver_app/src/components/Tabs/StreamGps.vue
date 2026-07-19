<script setup lang="ts">
import { onUnmounted } from "vue";
import { streamGPS, GpsProviders } from "geolocation";
import { moscowMockPath } from "geolocation/ws-gps-data";

const WS_HOST = import.meta.env.VITE_WS_HOST ?? "localhost:8100";
const WS_URL = `${location.protocol === "https:" ? "wss" : "ws"}://${WS_HOST}/ws/tow_driver`;

const provider = import.meta.env.DEV
  ? GpsProviders.ringBuffer(moscowMockPath)
  : GpsProviders.browser();

const { isStreaming, currentPosition, error, start, stop } = streamGPS(
  WS_URL,
  provider,
);

onUnmounted(() => stop());

function toggleStreaming() {
  if (isStreaming.value) {
    stop();
  } else {
    start();
  }
}
</script>

<template>
  <div class="flex flex-col gap-4 p-4">
    <button
      @click="toggleStreaming"
      :class="[
        'w-full py-3 rounded-lg font-medium transition-colors',
        isStreaming
          ? 'bg-red-500 text-white hover:bg-red-600'
          : 'bg-green-500 text-white hover:bg-green-600',
      ]"
    >
      {{ isStreaming ? "Остановить трансляцию" : "Начать трансляцию GPS" }}
    </button>

    <div
      v-if="isStreaming"
      class="flex flex-col gap-2 bg-gray-50 rounded-lg p-3"
    >
      <div class="flex items-center gap-2">
        <span class="w-2.5 h-2.5 rounded-full bg-green-500 animate-pulse" />
        <span class="text-green-600 font-medium">Трансляция активна</span>
      </div>

      <div
        v-if="currentPosition"
        class="flex flex-col gap-1 text-sm text-gray-600"
      >
        <div class="flex justify-between">
          <span>Широта:</span>
          <span>{{ currentPosition.lat.toFixed(6) }}</span>
        </div>
        <div class="flex justify-between">
          <span>Долгота:</span>
          <span>{{ currentPosition.lng.toFixed(6) }}</span>
        </div>
        <div class="flex justify-between">
          <span>Точность GPS:</span>
          <span
            :class="
              currentPosition.accuracy > 50
                ? 'text-orange-500'
                : 'text-green-600'
            "
          >
            &plusmn;{{ currentPosition.accuracy.toFixed(1) }}м
          </span>
        </div>
        <div class="flex justify-between">
          <span>Обновлено:</span>
          <span v-if="currentPosition.timestamp"
            >~{{
              ((Date.now() - currentPosition.timestamp) / 1000).toFixed(1)
            }}с назад</span
          >
        </div>
      </div>
    </div>

    <p v-if="error" class="text-red-500 text-sm text-center">{{ error }}</p>
  </div>
</template>
