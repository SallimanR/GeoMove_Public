<script setup lang="ts">
import Select from "primevue/select";
import { buildTimeOptions } from "../utils/timePicker";

const model = defineModel<string>();

defineProps<{ placeholder?: string }>();

const timeOptions = buildTimeOptions();
</script>

<template>
  <Select
    v-model="model"
    :options="timeOptions"
    :placeholder="placeholder"
    class="w-full"
  />
</template>
