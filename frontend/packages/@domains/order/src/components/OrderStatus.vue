<template>
  <Timeline :value="events">
    <template #opposite="slotProps">
      <div class="text-xs leading-4 text-surface-500 dark:text-surface-400">
        {{ slotProps.item.date }}
      </div>
    </template>
    <template #content="slotProps">
      <div class="text-sm leading-4">
        {{ slotProps.item.status }}
      </div>
    </template>
  </Timeline>
</template>

<script setup>
import { ref } from "vue";

const events = ref([
  { status: "Ordered", date: "15/10/2026 10:30" },
  { status: "Processing", date: "15/10/2026 14:00" },
  { status: "Shipped", date: "15/10/2026 16:15" },
  { status: "Delivered", date: "16/10/2026 10:00" },
]);
</script>
