<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="flex flex-col h-full">
    <div class="flex items-center gap-3 px-4 py-3 border-b border-gray-200">
      <button
        @click="router.back()"
        class="flex items-center justify-center w-10 h-10 rounded-full hover:bg-gray-100 transition-colors"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <polyline points="15 18 9 12 15 6" />
        </svg>
      </button>
      <h2 class="text-lg font-semibold">История поездок</h2>
    </div>

    <div class="flex-1 flex items-center justify-center">
      <p class="text-gray-400 text-sm">Нет завершённых поездок</p>
    </div>
  </div>
</template>
