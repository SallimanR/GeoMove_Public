<script setup lang="ts">
import { ref } from "vue";
import MapsSearchInput from "./MapsSearchInput.vue";
import MapsRouteInput from "./MapsRouteInput.vue";
import MapsGPSLocation from "./MapsGPSLocation.vue";

type State = "from" | "to";

const mode = ref<State>("from");
</script>
<template>
  <div
    class="absolute inset-0 flex flex-col justify-between pointer-events-none"
  >
    <div class="pointer-events-auto m-4">
      <!-- <MapsSearchInput v-if="mode === 'from'" /> -->
      <MapsRouteInput v-if="mode === 'from'" />
    </div>

    <!-- <div class="flex gap-4 pointer-events-auto"> -->
    <!--   <button @click="mode = 'none'">None</button> -->
    <!--   <button @click="mode = 'search'">Search</button> -->
    <!--   <button @click="mode = 'route'">Route</button> -->
    <!--   <MapsGPSLocation /> -->
    <!-- </div> -->
  </div>
</template>
