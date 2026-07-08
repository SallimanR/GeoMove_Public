<script setup lang="ts">
import type { SearchResult } from "../../../geo/src/api/geocoding.ts";
import { computed } from "vue";

const props = defineProps<{
  result: SearchResult;
}>();

const location = computed(() => {
  const parts: string[] = [];
  const { city, state } = props.result.properties;
  if (city) parts.push(city);
  if (state && state !== city) parts.push(state);
  return parts.join(", ");
});
</script>

<template>
  <div class="flex flex-col gap-0.5">
    <div
      v-if="result.properties.name"
      class="font-semibold text-base text-gray-800"
    >
      {{ result.properties.name }}
    </div>

    <div class="text-sm text-gray-600">
      <template v-if="result.properties?.street">
        {{ result.properties.street }}
        <span v-if="result.properties?.housenumber">
          {{ result.properties.housenumber }}
        </span>
      </template>
      <span v-else-if="result.properties?.housenumber">
        {{ result.properties.housenumber }}
      </span>
    </div>

    <div v-if="location" class="text-xs text-gray-400">
      {{ location }}
    </div>
  </div>
</template>
