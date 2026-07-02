<script setup lang="ts">
import { ref } from "vue";
import { type SearchResult, type SearchResults, getMapSearch } from "geo";
import { $coords } from "../stores/mapsStore.ts";

const props = defineProps<{
  modelValue: string;
  placeholder: string;
  onSearchResultClick: (lat: number, lon: number) => void;
}>();

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const isTypingRequest = ref(false);
const results = ref<SearchResults>();

async function handleSearch(query: string) {
  if (!query.trim()) {
    results.value = undefined;
    return;
  }
  isTypingRequest.value = true;
  results.value = await getMapSearch(
    query,
    $coords.value.center.lat,
    $coords.value.center.lon,
  );
}

// Emit the new value upwards, and perform search
function onInput(e: Event) {
  const value = (e.target as HTMLInputElement).value;
  emit("update:modelValue", value);
  handleSearch(value);
}

// When a result is clicked: emit a display name and call the parent callback
async function handleSearchResultClick(searchResult: SearchResult) {
  isTypingRequest.value = false;
  const { properties } = searchResult;

  const parts = [
    properties.name,
    properties.city,
    properties.street,
    properties.housenumber,
  ].filter(Boolean);
  const displayName =
    parts.join(", ") ||
    `${searchResult.geometry.coordinates[1].toFixed(4)}, ${searchResult.geometry.coordinates[0].toFixed(4)}`;

  // Update the input field with the address
  emit("update:modelValue", displayName);

  // Pass the coordinates to the parent (which updates the store)
  props.onSearchResultClick(
    searchResult.geometry.coordinates[1],
    searchResult.geometry.coordinates[0],
  );
}
</script>

<template>
  <div class="relative">
    <input
      :value="modelValue"
      @input="onInput"
      :placeholder="placeholder"
      class="rounded-2xl bg-gray-200 p-3 w-full"
    />
    <!-- Search results dropdown -->
    <div
      v-if="results && isTypingRequest && modelValue"
      class="absolute bg-white border rounded shadow-lg mt-1 max-h-60 overflow-y-auto w-full z-10"
    >
      <div
        v-for="result in results.features"
        @click="handleSearchResultClick(result)"
        class="p-2 hover:bg-gray-100 cursor-pointer"
      >
        {{ result.properties?.name }}, {{ result.properties?.city }},
        {{ result.properties?.street }}, {{ result.properties?.housenumber }}
      </div>
    </div>
  </div>
</template>
