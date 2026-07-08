<script setup lang="ts">
import MapsRouteInput from "./MapsRouteInput.vue";
import MapsSearchInput from "./MapsSearchInput.vue";
</script>

<template>
  <div
    class="absolute inset-0 flex flex-col justify-between pointer-events-none"
  >
    <div class="pointer-events-auto m-4">
      <MapsSearchInput />
      <MapsRouteInput />
    </div>
  </div>
</template>
