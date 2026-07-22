<script setup lang="ts">
import { ref, watch, inject } from "vue";
import { TimePicker } from "ui";
import InputMask from "primevue/inputmask";
import InputNumber from "primevue/inputnumber";
import { MapsLocationPicker } from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import type { Driver } from "driver";
import { useDriverProfile } from "../../../stores/driverStore";
import { ACTIVE_TAB_KEY } from "../../../injectionKeys";

const props = defineProps<{
  driver: Driver;
}>();

const emit = defineEmits<{
  back: [];
}>();

const { updateProfile } = useDriverProfile();

const editName = ref("");
const editPhone = ref("");
const editMaxCarWeightKg = ref<number | null>(null);
const editMaxCarLengthMeters = ref<number | null>(null);
const editWorkStarts = ref("");
const editWorkEnds = ref("");
const editLocation = ref<GeoPoint>({ lat: 0, lon: 0 });
const editAddress = ref("");
const loading = ref(false);
const error = ref<string | null>(null);

const activeTab = inject(ACTIVE_TAB_KEY)!;

function populateFromDriver(d: Driver) {
  editName.value = d.name;
  editPhone.value = d.phone ?? "";
  editMaxCarWeightKg.value = d.max_car_weight_kg ?? null;
  editMaxCarLengthMeters.value = d.max_car_length_meters ?? null;
  editWorkStarts.value = d.work_starts ? d.work_starts.split("T")[1]?.slice(0, 5) ?? "" : "";
  editWorkEnds.value = d.work_ends ? d.work_ends.split("T")[1]?.slice(0, 5) ?? "" : "";
  editLocation.value = { lat: d.lat, lon: d.lon };
  editAddress.value = d.address || `${d.lat.toFixed(6)}, ${d.lon.toFixed(6)}`;
}
watch(() => props.driver, (d) => populateFromDriver(d), { immediate: true });

function onLocationPicked(point: GeoPoint, address: string) {
  editLocation.value = point;
  editAddress.value = address;
  activeTab.value = "profileTab";
}

function handleOnMapsPick() {
  activeTab.value = "mapsTab";
}

async function onSubmit() {
  error.value = null;
  loading.value = true;
  try {
    await updateProfile({
      name: editName.value.trim(),
      lat: editLocation.value.lat,
      lon: editLocation.value.lon,
      phone: editPhone.value.trim() || undefined,
      work_starts: editWorkStarts.value || undefined,
      work_ends: editWorkEnds.value || undefined,
      max_car_weight_kg: editMaxCarWeightKg.value ?? undefined,
      max_car_length_meters: editMaxCarLengthMeters.value ?? undefined,
    });
    emit("back");
  } catch (err) {
    error.value =
      err instanceof Error ? err.message : "Не удалось обновить профиль";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col gap-4 p-4 overflow-y-auto">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium">Редактирование профиля</h2>
      <button
        @click="emit('back')"
        class="px-3 py-1 bg-gray-200 hover:bg-gray-300 rounded-lg text-sm"
      >
        ← Назад
      </button>
    </div>

    <div class="flex flex-col gap-3">
      <input
        v-model="editName"
        type="text"
        placeholder="Имя"
        class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
      />

      <InputMask
        v-model="editPhone"
        mask="+7 (999) 999-99-99"
        placeholder="+7 (___) ___-__-__"
        class="w-full"
      />

      <div class="flex gap-2">
        <InputNumber
          v-model="editMaxCarWeightKg"
          placeholder="Макс. вес авто (кг)"
          :min="1"
          class="flex-1"
        />
        <InputNumber
          v-model="editMaxCarLengthMeters"
          placeholder="Макс. длина авто (м)"
          :min="0.1"
          :minFractionDigits="1"
          :maxFractionDigits="2"
          class="flex-1"
        />
      </div>

      <div class="flex gap-2">
        <TimePicker v-model="editWorkStarts" placeholder="Начало работы" />
        <TimePicker v-model="editWorkEnds" placeholder="Конец работы" />
      </div>

      <MapsLocationPicker
        @click="handleOnMapsPick()"
        @pick="onLocationPicked"
      />

      <div v-if="editAddress" class="flex items-center gap-1 text-gray-500">
        <span class="text-green-500">●</span>
        {{ editAddress }}
      </div>

      <p v-if="error" class="text-red-500">{{ error }}</p>

      <button
        @click="onSubmit"
        :disabled="loading || !editName.trim()"
        class="w-full px-4 py-2 bg-green-300 hover:bg-green-400 rounded-lg font-medium"
      >
        {{ loading ? "Сохранение..." : "Сохранить" }}
      </button>
    </div>
  </div>
</template>
