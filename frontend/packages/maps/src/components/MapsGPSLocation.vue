<script setup lang="ts">
import { ref } from 'vue'
import { $coords } from '../stores/mapsStore'
import { useMapsActions } from '../composables/useMapsActions'

const { flyTo } = useMapsActions()

const loading = ref(false)
const error = ref<string | null>(null)

async function getPositionCapacitor(): Promise<{ lat: number; lng: number }> {
  try {
    const mod = await new Function('specifier', 'return import(specifier)')('@capacitor/geolocation')
    const geo: any = mod.Geolocation

    const perm = await geo.checkPermissions()
    if (perm.location !== 'granted') {
      const req = await geo.requestPermissions()
      if (req.location !== 'granted') {
        throw new Error('Доступ к геолокации отклонён')
      }
    }

    const pos = await geo.getCurrentPosition({
      enableHighAccuracy: true,
      timeout: 10000,
    })

    if (!pos?.coords) {
      throw new Error('Не удалось определить местоположение')
    }

    return { lat: pos.coords.latitude, lng: pos.coords.longitude }
  } catch {
    const pos = await new Promise<GeolocationPosition>((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(resolve, reject, {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 0,
      })
    })
    return { lat: pos.coords.latitude, lng: pos.coords.longitude }
  }
}

async function handleFindMe() {
  loading.value = true
  error.value = null

  try {
    let lat: number, lng: number

    let isNative = false
    try {
      const mod = await new Function('specifier', 'return import(specifier)')('@capacitor/core')
      isNative = mod.Capacitor.isNativePlatform()
    } catch {}

    if (isNative) {
      const coords = await getPositionCapacitor()
      lat = coords.lat
      lng = coords.lng
    } else {
      const pos = await new Promise<GeolocationPosition>((resolve, reject) => {
        navigator.geolocation.getCurrentPosition(resolve, reject, {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 0,
        })
      })
      lat = pos.coords.latitude
      lng = pos.coords.longitude
    }

    $coords.set({ center: { lat, lon: lng } })
    flyTo(lat, lng)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Ошибка геолокации'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <button
    @click="handleFindMe"
    :disabled="loading"
    type="button"
    class="p-2 rounded-lg bg-gray-200 hover:bg-gray-300 disabled:opacity-50 transition-colors"
    :title="error ?? 'Найти меня'"
  >
    <svg
      v-if="loading"
      class="animate-spin h-5 w-5"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
    </svg>
    <svg
      v-else
      class="h-5 w-5"
      :class="error ? 'text-red-500' : ''"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      stroke-width="1.5"
      stroke="currentColor"
    >
      <path stroke-linecap="round" stroke-linejoin="round" d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z" />
    </svg>
  </button>
</template>
