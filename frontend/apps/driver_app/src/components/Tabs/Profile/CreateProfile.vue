<script setup lang="ts">
import { ref, nextTick, inject } from "vue";
import Cropper from "cropperjs";
import { TimePicker } from "ui";
import { MapsLocationPicker } from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import { useDriverProfile } from "../../../stores/driverStore";
import { ACTIVE_TAB_KEY } from "../../../injectionKeys";

const emit = defineEmits<{ (e: "created"): void }>();

const { createProfile, uploadProfileImage } = useDriverProfile();

const name = ref("");
const workStarts = ref("");
const workEnds = ref("");
const loading = ref(false);
const error = ref<string | null>(null);

const pickedLocation = ref<GeoPoint | null>(null);
const pickedAddress = ref("");

const showCropper = ref(false);
const cropImageSrc = ref("");
const cropperRef = ref<HTMLImageElement | null>(null);
let cropperInstance: Cropper | null = null;
const uploadLoading = ref(false);
const profileImageBase64 = ref<string | null>(null);

function onFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = (e) => {
    cropImageSrc.value = e.target?.result as string;
    showCropper.value = true;
    nextTick(() => initCropper());
  };
  reader.readAsDataURL(file);
  input.value = "";
}

function initCropper() {
  if (cropperInstance) {
    cropperInstance.destroy();
  }
  if (!cropperRef.value) return;

  cropperInstance = new Cropper(cropperRef.value);

  const selection = cropperInstance.getCropperSelection();
  if (selection) {
    selection.aspectRatio = 1;
    selection.initialCoverage = 1;
  }
}

async function onCrop() {
  if (!cropperInstance) return;
  uploadLoading.value = true;
  try {
    const cropperCanvas = cropperInstance.getCropperCanvas();
    const canvas = await cropperCanvas?.$toCanvas({ width: 300, height: 300 });
    if (!canvas) {
      throw new Error("Не удалось загрузить");
    }
    profileImageBase64.value = canvas.toDataURL("image/jpeg", 0.85);
    showCropper.value = false;
    cropperInstance.destroy();
    cropperInstance = null;
  } catch (err) {
    console.error("Crop failed:", err);
  } finally {
    uploadLoading.value = false;
  }
}

function onCropCancel() {
  showCropper.value = false;
  if (cropperInstance) {
    cropperInstance.destroy();
    cropperInstance = null;
  }
}

function removeImage() {
  profileImageBase64.value = null;
}

function onLocationPicked(point: GeoPoint, address: string) {
  pickedLocation.value = point;
  pickedAddress.value = address;
  activeTab.value = "profileTab";
}

async function onSubmit() {
  error.value = null;
  if (!pickedLocation.value) {
    error.value = "Выберете локацию";
    return;
  }
  loading.value = true;
  try {
    await createProfile(
      name.value.trim(),
      pickedLocation.value.lat,
      pickedLocation.value.lon,
      workStarts.value || undefined,
      workEnds.value || undefined,
    );

    if (profileImageBase64.value) {
      await uploadProfileImage(profileImageBase64.value);
    }

    emit("created");
  } catch (err) {
    error.value =
      err instanceof Error
        ? err.message
        : "Не удалось создать профиль водителя";
  } finally {
    loading.value = false;
  }
}

const activeTab = inject(ACTIVE_TAB_KEY)!;

function handleOnMapsPick() {
  activeTab.value = "mapsTab";
}
</script>

<template>
  <div
    class="flex flex-col items-center justify-center h-full gap-4 p-4 overflow-y-auto"
  >
    <h2 class="text-lg font-medium">Создать профиль водителя</h2>

    <div class="flex flex-col gap-3 w-full max-w-140">
      <label class="cursor-pointer relative group self-center">
        <div
          v-if="profileImageBase64"
          class="w-20 h-20 rounded-full overflow-hidden ring-2 ring-gray-300"
        >
          <img
            :src="profileImageBase64"
            alt="Avatar"
            class="w-full h-full object-cover"
          />
        </div>
        <div
          v-else
          class="w-20 h-20 rounded-full bg-gray-200 flex items-center justify-center ring-2 ring-gray-300"
        >
          <span class="text-2xl text-gray-400">+</span>
        </div>
        <div
          class="absolute inset-0 rounded-full bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <span class="text-white text-xs font-medium">{{
            profileImageBase64 ? "Изменить" : "Добавить"
          }}</span>
        </div>
        <input
          type="file"
          accept="image/*"
          class="hidden"
          @change="onFileSelect"
        />
      </label>

      <button
        v-if="profileImageBase64"
        @click="removeImage"
        class="text-red-400 hover:text-red-500 transition-colors self-center"
      >
        Убрать фото
      </button>

      <input
        v-model="name"
        type="text"
        placeholder="Введите имя"
        class="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
      />

      <div class="flex gap-2">
        <TimePicker v-model="workStarts" placeholder="Начало работы" />
        <TimePicker v-model="workEnds" placeholder="Конец работы" />
      </div>

      <MapsLocationPicker
        @click="handleOnMapsPick()"
        @pick="onLocationPicked"
      />

      <div v-if="pickedAddress" class="flex items-center gap-1 text-gray-500">
        <span class="text-green-500">●</span>
        {{ pickedAddress }}
      </div>

      <p v-if="error" class="text-red-500">{{ error }}</p>

      <button
        @click="onSubmit"
        :disabled="loading || !name.trim()"
        class="w-full px-4 py-2 bg-green-300 hover:bg-green-400 rounded-lg font-medium"
      >
        {{ loading ? "Создаём..." : "Создать" }}
      </button>
    </div>
  </div>

  <Teleport to="body">
    <div
      v-if="showCropper"
      class="fixed inset-0 z-50 bg-black/70 flex items-center justify-center"
    >
      <div class="bg-white rounded-lg p-4 w-[90vw] max-w-md">
        <h3 class="text-lg font-medium mb-4">Обрезать картинку</h3>
        <div class="w-full h-64 overflow-hidden">
          <img
            ref="cropperRef"
            :src="cropImageSrc"
            alt="Обрезать"
            class="max-w-full"
          />
        </div>
        <div class="flex justify-end gap-3 mt-4">
          <button
            @click="onCropCancel"
            class="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
          >
            Отменить
          </button>
          <button
            @click="onCrop"
            :disabled="uploadLoading"
            class="w-full px-4 py-2 bg-green-300 hover:bg-green-400 rounded-lg font-medium"
          >
            {{ uploadLoading ? "Сохранение..." : "Сохранить" }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
