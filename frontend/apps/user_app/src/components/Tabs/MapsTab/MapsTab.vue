<script setup lang="ts">
import { onMounted } from "vue";
import { Maps, MapsOverlayControls } from "@geomove/maps";
import OrderDropdown from "./OrderDropdown.vue";
import DriverDropdown from "./DriverDropdown.vue";
import DriverCardPopup from "./DriverCardPopup.vue";
import FreelyAvailableDriverCardPopup from "./FreelyAvailableDriverCardPopup.vue";
import {
  initDriverPopupZoom,
  displayDriverPopups,
  displayFreelyAvailableDriverPopups,
} from "../../../maps/displayDrivers";

const styleApi = import.meta.env.VITE_STYLE_API;

function redrawDrivers() {
  displayDriverPopups();
  displayFreelyAvailableDriverPopups();
}

onMounted(() => {
  initDriverPopupZoom(redrawDrivers);
});
</script>

<template>
  <div class="relative flex h-full flex-col">
    <Maps :styleApi="styleApi" />
    <MapsOverlayControls :hideRouteInput="true" />
    <OrderDropdown />
    <DriverDropdown />
    <DriverCardPopup />
    <FreelyAvailableDriverCardPopup />
  </div>
</template>
