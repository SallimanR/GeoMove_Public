<script setup lang="ts">
import type { Driver } from "driver/types/driver.ts";
import { $selectedDriver } from "src/stores/driverStore.ts";

const props = defineProps<Driver>();

function openCardPopup() {
  $selectedDriver.set(props);
}
</script>

<template>
  <div
    class="flex flex-col items-center bg-gray-100 rounded-2xl p-2"
    @click="openCardPopup"
  >
    <div v-if="props.rating" class="rating">
      <span>&#9733;</span> {{ props.rating }}
    </div>
    <div class="rounded-2xl p-2 bg-gray-200">{{ props.name }}</div>
    <img
      :src="'tow_image.jpg'"
      :alt="`Driver ${props.user_id}`"
      class="driver-avatar"
    />
  </div>
</template>

<style scoped>
.rating {
  font-size: 18px;
  font-weight: bold;
  color: #f5a623;
  margin-bottom: 4px;
}

.driver-avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px;
  margin: 4px 0;
}

.driver-id {
  font-size: 14px;
  color: #555;
  margin-top: 4px;
}
</style>
