<script setup lang="ts">
import { onMounted } from "vue";
import { useStore } from "@nanostores/vue";

import {
  $user,
  $isAuthenticated,
  $loading,
  $error,
  checkAuth,
  setUser,
  clearUser,
  logout,
} from "auth";
import SingIn from "auth/components/SignIn.vue";
import CreateProfile from "./CreateProfile.vue";
import { useDriverProfile, resolveImageUrl } from "../../../stores/driverStore";

const user = useStore($user);
const isAuthenticated = useStore($isAuthenticated);
const authLoading = useStore($loading);
const authError = useStore($error);
const {
  driver,
  exists,
  loading: driverLoading,
  error: driverError,
  fetchProfile,
} = useDriverProfile();

onMounted(async () => {
  await checkAuth();
  if (isAuthenticated.value) {
    await fetchProfile();
  }
});

function onLoginSuccess(u: { id: number; email: string | null }) {
  setUser({ id: u.id, email: u.email, phone: null, profile_image: null });
  fetchProfile();
}

async function onLogout() {
  await logout();
  clearUser();
}
</script>

<template>
  <div
    v-if="authLoading || driverLoading"
    class="flex items-center justify-center h-full"
  >
    <p class="text-gray-500">Загрузка...</p>
  </div>

  <div
    v-else-if="authError || driverError"
    class="flex items-center justify-center h-full"
  >
    <p class="text-red-500">{{ authError || driverError }}</p>
  </div>

  <div
    v-else-if="!isAuthenticated"
    class="flex items-center justify-center h-full"
  >
    <SingIn @login-success="onLoginSuccess" />
  </div>

  <div v-else-if="!exists">
    <CreateProfile @created="fetchProfile" />
  </div>

  <div
    v-else-if="driver"
    class="flex flex-col items-center justify-center h-full gap-4 p-4"
  >
    <label class="cursor-pointer relative group">
      <div
        v-if="driver.profile_image"
        class="w-24 h-24 rounded-full overflow-hidden ring-2 ring-gray-300"
      >
        <img
          :src="resolveImageUrl(driver.profile_image)"
          alt="Profile"
          class="w-full h-full object-cover"
        />
      </div>
      <div
        v-else
        class="w-24 h-24 rounded-full bg-gray-200 flex items-center justify-center ring-2 ring-gray-300"
      >
        <span class="text-3xl text-gray-500">{{
          driver.name?.charAt(0).toUpperCase() || "?"
        }}</span>
      </div>
    </label>

    <p class="text-lg font-medium">{{ driver.name }}</p>
    <p class="text-gray-500">{{ user?.email }}</p>

    <div v-if="driver.rating" class="flex items-center gap-1">
      <span class="text-yellow-500 text-sm"
        >★ {{ driver.rating.toFixed(1) }}</span
      >
    </div>

    <button
      @click="onLogout"
      class="px-6 py-2 bg-red-400 hover:bg-red-500 rounded-lg"
    >
      Выйти из аккаунта
    </button>
  </div>
</template>
