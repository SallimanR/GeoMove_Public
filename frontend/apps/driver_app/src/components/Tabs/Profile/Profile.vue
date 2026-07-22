<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useStore } from "@nanostores/vue";
import { usePushNotifications } from "notifications";

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
import DisplayProfile from "./DisplayProfile.vue";
import EditProfile from "./EditProfile.vue";
import { useDriverProfile } from "../../../stores/driverStore";

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

const { isSupported, init: initPush } = usePushNotifications();

const editing = ref(false);

onMounted(async () => {
  await checkAuth();
  if (isAuthenticated.value) {
    await fetchProfile();
    if (isSupported.value) initPush();
  }
});

watch(isAuthenticated, (auth) => {
  if (auth && isSupported.value) initPush();
});

function onLoginSuccess(u: { id: number; email: string | null }) {
  setUser({ id: u.id, email: u.email, phone: null, profile_image: null });
  fetchProfile();
  if (isSupported.value) initPush();
}

async function onLogout() {
  await logout();
  clearUser();
}

function onProfileUpdated() {
  fetchProfile();
  editing.value = false;
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

  <div v-else-if="driver">
    <DisplayProfile
      v-if="!editing"
      :driver="driver"
      @edit="editing = true"
      @logout="onLogout"
    />
    <EditProfile
      v-else
      :driver="driver"
      @back="onProfileUpdated"
    />
  </div>
</template>
