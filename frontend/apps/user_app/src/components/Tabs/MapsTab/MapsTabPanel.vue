<script setup lang="ts">
import { Maps, MapsOverlayControls } from "@geomove/maps";
import OrderDropdown from "./OrderDropdown.vue";

const styleApi = import.meta.env.VITE_STYLE_API;

// TODO:
// useMovingDriverLayer(undefined, 55.751, 37.618, MovingDriverPopups, false);
</script>

<template>
  <div class="relative flex flex-col h-full">
    <Maps :styleApi="styleApi" />
    <MapsOverlayControls :hideRouteInput="true" />
    <OrderDropdown />
  </div>
</template>
