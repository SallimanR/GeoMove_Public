<script setup lang="ts">
import { $coords } from "../stores/mapsStore";

import { useMapsActions } from "../composables/useMapsActions.ts";

const { flyTo } = useMapsActions();

const gpsOptions: PositionOptions = {
  enableHighAccuracy: true,
  timeout: 5000,
  maximumAge: 0,
};

function onSuccess(pos: GeolocationPosition) {
  $coords.set({
    center: { lat: pos.coords.latitude, lon: pos.coords.longitude },
  });
  flyTo(pos.coords.latitude, pos.coords.longitude);
}

function onError(err: GeolocationPositionError) {
  console.log("failed to get gps position: ", err);
}

function handleFindMe(): void {
  navigator.geolocation.getCurrentPosition(onSuccess, onError, gpsOptions);
}
</script>

<template>
  <button @click="handleFindMe()">Find me</button>
</template>
