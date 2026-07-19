<script setup lang="ts">
import { BottomPanel } from "ui";
import DriverSearch from "./DriverSearch.vue";
import DriverFreelyAvailableSearch from "./DriverFreelyAvailableSearch.vue";
</script>

<template>
  <BottomPanel openPlaceholder="Поиск">
    <DriverSearch />
    <DriverFreelyAvailableSearch />
  </BottomPanel>
</template>
