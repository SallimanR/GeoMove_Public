<script setup lang="ts">
import MapsRouteInput from "./MapsRouteInput.vue";
import MapsSearchInput from "./MapsSearchInput.vue";
import MapsCurrentLocationBox from "./MapsCurrentLocationBox.vue";
import MapsGPSLocation from "./MapsGPSLocation.vue";

const props = withDefaults(
  defineProps<{
    hideRouteInput?: boolean;
  }>(),
  {
    hideRouteInput: false,
  },
);
</script>

<template>
  <div class="absolute inset-0 flex flex-col pointer-events-none">
    <div class="flex items-start justify-between gap-2 mt-4 mr-4 ml-4">
      <MapsSearchInput class="w-full max-w-200 pointer-events-auto" />
      <MapsGPSLocation class="shrink-0 pointer-events-auto" />
    </div>
    <MapsCurrentLocationBox />
    <div v-if="!props.hideRouteInput" class="pointer-events-auto">
      <MapsRouteInput />
    </div>
  </div>
</template>
