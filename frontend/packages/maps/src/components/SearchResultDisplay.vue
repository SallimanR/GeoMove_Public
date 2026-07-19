<script setup lang="ts">
import { displayDistance, haversineDistance } from "@geomove/geo";
import type { SearchResult } from "../../../geo/src/api/geocoding.ts";
import { computed } from "vue";
import { useStore } from "@nanostores/vue";
import { $coords } from "../stores/mapsStore.ts";

const props = defineProps<{
  result: SearchResult;
}>();

const coords = useStore($coords);

const distanceText = computed(() => {
  const center = coords.value?.center;
  if (!center) return "";
  return displayDistance(
    haversineDistance(
      [
        props.result.geometry.coordinates[1],
        props.result.geometry.coordinates[0],
      ],
      [center.lat, center.lon],
    ),
  );
});

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
    <div>{{ distanceText }}</div>
  </div>
</template>
