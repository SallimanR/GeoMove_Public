<script setup lang="ts">
const props = defineProps<{ id: number }>();

function openDriverRealtimePopup() {
  console.log("realtime id", props.id);
}
</script>

<template>
  <div
    class="flex flex-col items-center realtime-driver-popup"
    @click="openDriverRealtimePopup"
  >
    <img
      :src="'tow_image.jpg'"
      :alt="`Водитель ${props.id}`"
      class="driver-avatar"
    />
  </div>
</template>

<style scoped>
.driver-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #3cb300;
  margin: 2px 0;
}
</style>
