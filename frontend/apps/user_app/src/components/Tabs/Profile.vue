<script setup lang="ts">
import { onMounted, inject, watch } from "vue";
import { useStore } from "@nanostores/vue";
import { useRouter } from "vue-router";
import SingIn from "auth/components/SignIn.vue";
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
import { usePushNotifications } from "notifications";
import { ACTIVE_TAB_KEY } from "src/injectionKeys";

const router = useRouter();
const user = useStore($user);
const isAuthenticated = useStore($isAuthenticated);
const loading = useStore($loading);
const error = useStore($error);
const activeTab = inject(ACTIVE_TAB_KEY);

const { isSupported, isSubscribed, error: notifyError, init } = usePushNotifications();

onMounted(() => {
  checkAuth();
});

watch(isAuthenticated, (auth) => {
  if (auth && isSupported.value) {
    init();
  }
});

function onLoginSuccess(u: { id: number; email: string | null }) {
  setUser({ id: u.id, email: u.email, phone: null, profile_image: null });
  if (activeTab) {
    activeTab.value = "mapsTab";
  }
  if (isSupported.value) {
    init();
  }
}

async function onLogout() {
  await logout();
  clearUser();
}
</script>

<template>
  <div v-if="loading" class="flex items-center justify-center h-full">
    <p class="text-gray-500">Загружаем профиль...</p>
  </div>

  <div v-else-if="error" class="flex items-center justify-center h-full">
    <p class="text-red-500">{{ error }}</p>
  </div>

  <div
    v-else-if="isAuthenticated && user"
    class="flex flex-col items-center justify-center h-full gap-4 p-4"
  >
    <div
      v-if="user.profile_image"
      class="w-24 h-24 rounded-full overflow-hidden"
    >
      <img
        :src="user.profile_image"
        alt="Profile"
        class="w-full h-full object-cover"
      />
    </div>
    <div
      v-else
      class="w-24 h-24 rounded-full bg-gray-200 flex items-center justify-center"
    >
      <span class="text-3xl text-gray-500">{{
        user.email?.charAt(0).toUpperCase() || "?"
      }}</span>
    </div>
    <p class="text-lg font-medium">{{ user.email }}</p>
    <p v-if="notifyError" class="text-sm text-amber-600">{{ notifyError }}</p>
    <p v-if="isSubscribed" class="text-xs text-green-600">Уведомления включены</p>
    <button
      @click="router.push('/triphistory')"
      class="px-6 py-2.5 bg-blue-500 text-white rounded-lg shadow-sm hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-offset-2 transition-all duration-150 text-sm font-medium"
    >
      История поездок
    </button>
    <button
      @click="onLogout"
      class="px-6 py-2.5 bg-red-400 text-white rounded-lg shadow-sm hover:bg-red-500 focus:outline-none focus:ring-2 focus:ring-red-400 focus:ring-offset-2 transition-all duration-150 text-sm font-medium"
    >
      Выйти из аккаунта
    </button>
  </div>

  <div v-else class="flex items-center justify-center h-full">
    <SingIn @login-success="onLoginSuccess" />
  </div>
</template>
