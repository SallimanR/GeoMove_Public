<script setup lang="ts">
import { useMapsActions } from "../composables/useMapsActions.ts";
import MapsSearchField from "./MapsSearchField.vue";

const { flyTo } = useMapsActions();

function onSearchResultClick(lat: number, lon: number) {
  flyTo(lat, lon);
}
</script>
<template>
  <MapsSearchField
    :onSearchResultClick="onSearchResultClick"
    placeholder="Search"
  />
</template>
