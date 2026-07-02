<script setup lang="ts">
import { $deckOverlay, $mapInstance } from "@stores/mapsStore";
import { RealtimeDriversAnimator } from "geolocation/animation";
import { type Map as MaplibreMaps } from "maplibre-gl";
import { onUnmounted, ref } from "vue";

const animator = ref<RealtimeDriversAnimator>();

const map = $mapInstance;
const deckOverlay = $deckOverlay;

function handleRealtimeDriversStart() {
  if (!map) return;
  const mapInstance = map.get() as MaplibreMaps | null;
  if (!mapInstance) return;
  if (deckOverlay.value) {
    animator.value = new RealtimeDriversAnimator(
      mapInstance,
      deckOverlay.value,
    );
    animator.value.start();
  }
}

function handleRealtimeDriversStop() {
  if (animator.value) {
    animator.value.stop();
  }
}

onUnmounted(() => {
  handleRealtimeDriversStop();
});
</script>
<template>
  <div class="flex flex-col">
    <button class="p-2" @click="handleRealtimeDriversStart()">
      Start realtime drivers
    </button>
    <button class="p-2" @click="handleRealtimeDriversStop()">
      Stop realtime drivers
    </button>
  </div>
</template>
