<script setup lang="ts">
withDefaults(
  defineProps<{
    count: number;
    color?: string;
    onZoomIn?: () => void;
  }>(),
  { color: "bg-blue-500" },
);
</script>

<template>
  <div class="flex cursor-pointer flex-col items-center" @click="onZoomIn?.()">
    <div
      class="flex h-12 w-12 items-center justify-center rounded-full text-base font-bold text-white shadow-lg"
      :class="color"
    >
      {{ count }}
    </div>
  </div>
</template>
