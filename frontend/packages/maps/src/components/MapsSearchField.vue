<script setup lang="ts">
import { ref } from "vue";
import {
  type SearchResult,
  type SearchResultList,
  addressToText,
  getMapSearch,
} from "@geomove/geo";
import { $coords } from "../stores/mapsStore.ts";
import SearchResultDisplay from "./SearchResultDisplay.vue";

const props = defineProps<{
  modelValue: string;
  placeholder: string;
  onSearchResultClick: (lat: number, lon: number) => void;
  onClear: () => void;
}>();

const isFocused = ref(false);

const emit = defineEmits<{
  (e: "update:modelValue", value: string): void;
}>();

const isTypingRequest = ref(false);
const results = ref<SearchResultList>();

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

  const displayName = addressToText(searchResult);

  // Update the input field with the address
  emit("update:modelValue", displayName);

  props.onSearchResultClick(
    searchResult.geometry.coordinates[1],
    searchResult.geometry.coordinates[0],
  );
}

function clearInput() {
  props.onClear();

  emit("update:modelValue", "");
  results.value = undefined;
  isTypingRequest.value = false;
}
</script>

<template>
  <div
    class="rounded-xl"
    @focusin="isFocused = true"
    @focusout="isFocused = false"
  >
    <div class="flex flex-row items-center gap-2">
      <input
        :value="modelValue"
        @input="onInput"
        :placeholder="placeholder"
        class="rounded-xl bg-gray-200 p-3 w-full"
      />
      <button
        @click="clearInput"
        class="shrink w-8 h-8 rounded-full bg-gray-300 hover:bg-gray-400 flex items-center justify-center text-gray-700 font-bold"
      >
        ✕
      </button>
    </div>

    <!-- Search results dropdown -->
    <div
      v-if="results && isTypingRequest && modelValue"
      v-show="isFocused"
      class="bg-white rounded-xl max-h-[40vh] w-full z-10 overflow-y-auto"
    >
      <div
        v-for="result in results.features"
        @mousedown.prevent
        @click="handleSearchResultClick(result)"
        class="m-1 p-2 rounded-xl bg-gray-200 hover:bg-gray-300 cursor-pointer"
      >
        <SearchResultDisplay :result="result" />
      </div>
    </div>
  </div>
</template>
