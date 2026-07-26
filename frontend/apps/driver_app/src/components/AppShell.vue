<script setup lang="ts">
import { onMounted } from "vue";
import { $isAuthenticated, checkAuth } from "auth";
import { useDriverProfile } from "../stores/driverStore";

const { fetchProfile } = useDriverProfile();

onMounted(async () => {
  await checkAuth();
  if ($isAuthenticated.get()) {
    await fetchProfile();
  }
});
</script>

<template>
  <router-view v-slot="{ Component }">
    <keep-alive>
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>
