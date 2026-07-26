<script setup lang="ts">
defineProps<{
  count: number;
  onZoomIn?: () => void;
}>();
</script>

<template>
  <div class="flex cursor-pointer flex-col items-center" @click="onZoomIn?.()">
    <div
      class="flex h-12 w-12 items-center justify-center rounded-full bg-blue-500 text-base font-bold text-white shadow-lg"
    >
      {{ count }}
    </div>
  </div>
</template>
