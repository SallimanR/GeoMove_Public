<script setup lang="ts">
import { Maps, MapsOverlayControls } from "@geomove/maps";
import { useMovingDriverLayer, MovingDriverProviders } from "geolocation";
import MovingDriverPopups from "./MovingDriverPopup.vue";

const styleApi = import.meta.env.VITE_STYLE_API;

// const provider = import.meta.env.VITE_DEV
//   ? undefined
//   : MovingDriverProviders.http();

useMovingDriverLayer(undefined, 55.751, 37.618, MovingDriverPopups, false);
</script>

<template>
  <div class="relative flex flex-col h-full">
    <Maps :styleApi="styleApi" />
    <MapsOverlayControls />
  </div>
</template>
