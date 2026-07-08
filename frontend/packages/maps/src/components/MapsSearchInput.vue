<script setup lang="ts">
import { ref } from "vue";
import { useMapsActions } from "../composables/useMapsActions.ts";
import MapsSearchField from "./MapsSearchField.vue";

const { flyTo } = useMapsActions();

const search = ref("");

function onSearchResultClick(lat: number, lon: number) {
  flyTo(lat, lon);
}

function onClear() {
  search.value = "";
}
</script>
<template>
  <MapsSearchField
    v-model="search"
    :onSearchResultClick="onSearchResultClick"
    :onClear="onClear"
    placeholder="Поиск"
  />
</template>
