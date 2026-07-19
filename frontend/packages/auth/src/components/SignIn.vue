<script setup lang="ts">
import GoogleLogin from "./Providers/GoogleLogin.vue";

const emit = defineEmits<{
  (e: "login-success", user: { id: number; email: string | null }): void;
}>();
</script>

<template>
  <GoogleLogin @login-success="(user) => emit('login-success', user)" />
</template>
