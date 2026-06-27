<script setup lang="ts">
import { ref } from "vue";

import { type SearchResult, type SearchResults, getMapSearch } from "geo";
import { $coords } from "../stores/mapsStore.ts";

const props = defineProps<{
  placeholder: string;
  onSearchResultClick: (lat: number, lon: number) => void;
}>();

const request = ref("");
const isTypingRequest = ref(false);
const results = ref<SearchResults>();

// TODO: emits
// const emits = defineEmits(["update:modelValue", "selectResult"]);

async function handleSearch() {
  isTypingRequest.value = true;
  results.value = await getMapSearch(
    request.value,
    $coords.value.center.lat,
    $coords.value.center.lon,
  );
  console.log("map search API response:", results.value);
}

async function handleSearchResultClick(searchResult: SearchResult) {
  isTypingRequest.value = false;
  const result = searchResult.properties;
  request.value = `${result.name || ""}
	${result.city || ""}
	${result.street || ""}
	${result.housenumber || ""}
  `;
  props.onSearchResultClick(
    searchResult.geometry.coordinates[1],
    searchResult.geometry.coordinates[0],
  );
}
</script>
<template>
  <input
    v-model="request"
    @input="handleSearch()"
    :placeholder="placeholder"
    class="rounded-2xl bg-gray-200 p-3"
  />
  <div
    @click="handleSearchResultClick(result)"
    v-show="request && isTypingRequest"
    v-for="result in results?.features"
    class="flex bg-gray-200 p-2 select-none"
  >
    <div>{{ result.properties?.city }}</div>
    <div>{{ result.properties?.name }}</div>
    <div>{{ result.properties?.street }}</div>
    <div>{{ result.properties?.housenumber }}</div>
  </div>
</template>
