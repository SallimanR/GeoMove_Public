<script setup lang="ts">
import { ref, nextTick, inject } from "vue";
import Cropper from "cropperjs";
import { TimePicker } from "ui";
import InputMask from "primevue/inputmask";
import InputNumber from "primevue/inputnumber";
import { MapsLocationPicker } from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import { useDriverProfile } from "../../../stores/driverStore";
import { ACTIVE_TAB_KEY } from "../../../injectionKeys";

const emit = defineEmits<{ (e: "created"): void }>();

const { createProfile, uploadProfileImage, uploadCarPhoto } = useDriverProfile();

const name = ref("");
const phone = ref("");
const maxCarWeightKg = ref<number | null>(null);
const maxCarLengthMeters = ref<number | null>(null);
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

const carPhotoMainBase64 = ref<string | null>(null);
const carPhotosBase64 = ref<string[]>([]);
const carPhotoUploading = ref(false);

type CropTarget =
  | { type: "profile" }
  | { type: "car_main" }
  | { type: "car_extra"; index: number }
  | { type: "car_new" };
const cropTarget = ref<CropTarget | null>(null);
const newCarPhotoInput = ref<HTMLInputElement | null>(null);

function addCarPhoto() {
  newCarPhotoInput.value?.click();
}

function onNewCarPhotoFile(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  cropTarget.value = { type: "car_new" };
  const reader = new FileReader();
  reader.onload = (e) => {
    cropImageSrc.value = e.target?.result as string;
    showCropper.value = true;
    nextTick(() => initCropper());
  };
  reader.readAsDataURL(file);
  input.value = "";
}

function onFileSelect(target: CropTarget, event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  cropTarget.value = target;
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
  if (!cropperInstance || !cropTarget.value) return;
  uploadLoading.value = true;
  try {
    const cropperCanvas = cropperInstance.getCropperCanvas();
    const target = cropTarget.value;
    let base64: string;
    if (target.type === "profile") {
      const canvas = await cropperCanvas?.$toCanvas({ width: 300, height: 300 });
      if (!canvas) throw new Error("Не удалось загрузить");
      base64 = canvas.toDataURL("image/jpeg", 0.85);
    } else {
      const srcCanvas = await cropperCanvas?.$toCanvas({ width: 1200 });
      if (!srcCanvas) throw new Error("Не удалось загрузить");
      const side = Math.min(srcCanvas.width, srcCanvas.height);
      const temp = document.createElement("canvas");
      temp.width = side;
      temp.height = side;
      const tctx = temp.getContext("2d")!;
      tctx.drawImage(srcCanvas, (srcCanvas.width - side) / 2, (srcCanvas.height - side) / 2, side, side, 0, 0, side, side);
      const scaled = document.createElement("canvas");
      scaled.width = 600;
      scaled.height = 600;
      const sctx = scaled.getContext("2d")!;
      sctx.drawImage(temp, 0, 0, 600, 600);
      base64 = scaled.toDataURL("image/jpeg", 0.85);
    }

    if (target.type === "profile") {
      profileImageBase64.value = base64;
    } else if (target.type === "car_main") {
      carPhotoMainBase64.value = base64;
    } else if (target.type === "car_extra") {
      carPhotosBase64.value[target.index] = base64;
    } else if (target.type === "car_new") {
      carPhotosBase64.value.push(base64);
    }

    cropTarget.value = null;
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
  cropTarget.value = null;
  showCropper.value = false;
  if (cropperInstance) {
    cropperInstance.destroy();
    cropperInstance = null;
  }
}

function removeImage() {
  profileImageBase64.value = null;
}

function removeCarMainPhoto() {
  carPhotoMainBase64.value = null;
}

function removeCarPhoto(index: number) {
  carPhotosBase64.value.splice(index, 1);
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
    let carPhotoMainUrl: string | undefined;
    const carPhotoUrls: string[] = [];

    if (carPhotoMainBase64.value) {
      carPhotoUploading.value = true;
      carPhotoMainUrl = await uploadCarPhoto(carPhotoMainBase64.value);
    }
    for (const photo of carPhotosBase64.value) {
      if (photo) {
        const url = await uploadCarPhoto(photo);
        carPhotoUrls.push(url);
      }
    }
    carPhotoUploading.value = false;

    await createProfile(
      name.value.trim(),
      pickedLocation.value.lat,
      pickedLocation.value.lon,
      workStarts.value || undefined,
      workEnds.value || undefined,
      phone.value.trim() || undefined,
      maxCarWeightKg.value ?? undefined,
      maxCarLengthMeters.value ?? undefined,
      carPhotoMainUrl,
      carPhotoUrls.length ? carPhotoUrls : undefined,
    );

    if (profileImageBase64.value) {
      await uploadProfileImage(profileImageBase64.value);
    }

    emit("created");
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Не удалось создать профиль водителя";
  } finally {
    loading.value = false;
    carPhotoUploading.value = false;
  }
}

const activeTab = inject(ACTIVE_TAB_KEY)!;

function handleOnMapsPick() {
  activeTab.value = "mapsTab";
}
</script>

<template>
  <div class="flex h-full flex-col items-center justify-center gap-4 overflow-y-auto p-4">
    <h2 class="text-lg font-medium">Создать профиль водителя</h2>

    <div class="flex w-full max-w-140 flex-col gap-3">
      <label class="group relative cursor-pointer self-center">
        <div
          v-if="profileImageBase64"
          class="h-20 w-20 overflow-hidden rounded-full ring-2 ring-gray-300"
        >
          <img :src="profileImageBase64" alt="Avatar" class="h-full w-full object-cover" />
        </div>
        <div
          v-else
          class="flex h-20 w-20 items-center justify-center rounded-full bg-gray-200 ring-2 ring-gray-300"
        >
          <span class="text-2xl text-gray-400">+</span>
        </div>
        <div
          class="absolute inset-0 flex items-center justify-center rounded-full bg-black/40 opacity-0 transition-opacity group-hover:opacity-100"
        >
          <span class="text-xs font-medium text-white">{{
            profileImageBase64 ? "Изменить" : "Добавить"
          }}</span>
        </div>
        <input type="file" accept="image/*" class="hidden" @change="onFileSelect" />
      </label>

      <button
        v-if="profileImageBase64"
        @click="removeImage"
        class="self-center text-red-400 transition-colors hover:text-red-500"
      >
        Убрать фото
      </button>

      <div class="flex flex-col items-center gap-2">
        <span class="text-gray-500">Фото автомобиля</span>

        <label class="group relative cursor-pointer self-center">
          <div
            v-if="carPhotoMainBase64"
            class="h-28 w-40 overflow-hidden rounded-lg ring-2 ring-gray-300"
          >
            <img :src="carPhotoMainBase64" class="h-full w-full object-cover" />
          </div>
          <div
            v-else
            class="flex h-28 w-40 items-center justify-center rounded-lg bg-gray-200 ring-2 ring-gray-300"
          >
            <span class="text-2xl text-gray-400">+</span>
          </div>
          <div
            class="absolute inset-0 flex items-center justify-center rounded-lg bg-black/40 opacity-0 transition-opacity group-hover:opacity-100 active:opacity-100"
          >
            <span class="text-xs font-medium text-white">{{
              carPhotoMainBase64 ? "Изменить" : "Добавить"
            }}</span>
          </div>
          <input
            type="file"
            accept="image/*"
            class="hidden"
            @change="onFileSelect({ type: 'car_main' }, $event)"
          />
        </label>

        <button
          v-if="carPhotoMainBase64"
          @click="removeCarMainPhoto"
          class="text-red-400 transition-colors hover:text-red-500"
        >
          Убрать
        </button>
      </div>

      <div class="flex flex-col items-center gap-2">
        <span class="text-gray-500">Дополнительные фото</span>

        <div class="flex flex-wrap justify-center gap-2">
          <label
            v-for="(photo, idx) in carPhotosBase64"
            :key="idx"
            class="group relative cursor-pointer"
          >
            <div class="h-20 w-20 overflow-hidden rounded-lg ring-2 ring-gray-300">
              <img :src="photo" class="h-full w-full object-cover" />
            </div>

            <div
              class="absolute inset-0 flex items-center justify-center rounded-lg bg-black/40 opacity-0 transition-opacity group-hover:opacity-100 active:opacity-100"
            >
              <span class="text-xs font-medium text-white">Изменить</span>
            </div>

            <button
              @click.stop="removeCarPhoto(idx)"
              class="absolute -top-1.5 -right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-red-400 hover:bg-red-500"
            >
              <span class="text-xs leading-none text-white">×</span>
            </button>

            <input
              type="file"
              accept="image/*"
              class="hidden"
              @change="onFileSelect({ type: 'car_extra', index: idx }, $event)"
            />
          </label>

          <button
            @click="addCarPhoto"
            class="flex h-20 w-20 items-center justify-center rounded-lg bg-gray-200 ring-2 ring-gray-300 transition-colors hover:bg-gray-300"
          >
            <span class="text-xl text-gray-500">+</span>
          </button>
          <input
            ref="newCarPhotoInput"
            type="file"
            accept="image/*"
            class="hidden"
            @change="onNewCarPhotoFile"
          />
        </div>
      </div>

      <input
        v-model="name"
        type="text"
        placeholder="Введите имя"
        class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-400 focus:outline-none"
      />

      <InputMask
        v-model="phone"
        mask="+7 (999) 999-99-99"
        placeholder="+7 (___) ___-__-__"
        class="w-full"
      />

      <div class="flex gap-2">
        <InputNumber
          v-model="maxCarWeightKg"
          placeholder="Макс. вес авто (кг)"
          :min="1"
          class="flex-1"
        />
        <InputNumber
          v-model="maxCarLengthMeters"
          placeholder="Макс. длина авто (м)"
          :min="0.1"
          :minFractionDigits="1"
          :maxFractionDigits="2"
          class="flex-1"
        />
      </div>

      <div class="flex gap-2">
        <TimePicker v-model="workStarts" placeholder="Начало работы" />
        <TimePicker v-model="workEnds" placeholder="Конец работы" />
      </div>

      <MapsLocationPicker @click="handleOnMapsPick()" @pick="onLocationPicked" />

      <div v-if="pickedAddress" class="flex items-center gap-1 text-gray-500">
        <span class="text-green-500">●</span>
        {{ pickedAddress }}
      </div>

      <p v-if="error" class="text-red-500">{{ error }}</p>

      <button
        @click="onSubmit"
        :disabled="loading || !name.trim() || carPhotoUploading"
        class="w-full rounded-lg bg-green-300 px-4 py-2 font-medium hover:bg-green-400 disabled:opacity-50"
      >
        {{ carPhotoUploading ? "Загрузка фото..." : loading ? "Создаём..." : "Создать" }}
      </button>
    </div>
  </div>

  <Teleport to="body">
    <div v-if="showCropper" class="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
      <div class="w-[90vw] max-w-md rounded-lg bg-white p-4">
        <h3 class="mb-4 text-lg font-medium">Обрезать картинку</h3>
        <div class="h-64 w-full overflow-hidden">
          <img ref="cropperRef" :src="cropImageSrc" alt="Обрезать" class="max-w-full" />
        </div>
        <div class="mt-4 flex justify-end gap-3">
          <button
            @click="onCropCancel"
            class="rounded-lg px-4 py-2 text-sm text-gray-600 transition-colors hover:bg-gray-100"
          >
            Отменить
          </button>
          <button
            @click="onCrop"
            :disabled="uploadLoading"
            class="w-full rounded-lg bg-green-300 px-4 py-2 font-medium hover:bg-green-400"
          >
            {{ uploadLoading ? "Сохранение..." : "Сохранить" }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
