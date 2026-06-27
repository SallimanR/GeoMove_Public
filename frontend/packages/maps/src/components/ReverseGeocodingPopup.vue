<script setup lang="ts">
const props = defineProps<{
  street: string;
  housenumber: string;
  city?: string;
}>();
</script>

<template>
  <div class="p-2">
    <strong>Address</strong><br />
    <span>{{ props.street }} {{ props.housenumber }}</span>
    <span v-if="city">, {{ props.city }}</span>
  </div>
</template>
