<script setup lang="ts">
import { ref, watch, inject, nextTick } from "vue";
import Cropper from "cropperjs";
import { TimePicker } from "ui";
import InputMask from "primevue/inputmask";
import InputNumber from "primevue/inputnumber";
import { MapsLocationPicker } from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import type { Driver } from "driver";
import { useDriverProfile, resolveImageUrl } from "../../../stores/driverStore";
import { ACTIVE_TAB_KEY } from "../../../injectionKeys";

const props = defineProps<{
  driver: Driver;
}>();

const emit = defineEmits<{
  back: [];
}>();

const { updateProfile, uploadCarPhoto } = useDriverProfile();

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

const editCarPhotoMain = ref<string | null>(null);
const editCarPhotos = ref<string[]>([]);
const newCarPhotosBase64 = ref<string[]>([]);
const carPhotoUploading = ref(false);

const showCropper = ref(false);
const cropImageSrc = ref("");
const cropperRef = ref<HTMLImageElement | null>(null);
let cropperInstance: Cropper | null = null;
const uploadLoading = ref(false);
const newCarPhotoInput = ref<HTMLInputElement | null>(null);
let cropCallback: ((base64: string) => void) | null = null;

const activeTab = inject(ACTIVE_TAB_KEY)!;

function populateFromDriver(d: Driver) {
  editName.value = d.name;
  editPhone.value = d.phone ?? "";
  editMaxCarWeightKg.value = d.max_car_weight_kg ?? null;
  editMaxCarLengthMeters.value = d.max_car_length_meters ?? null;
  editWorkStarts.value = d.work_starts ? (d.work_starts.split("T")[1]?.slice(0, 5) ?? "") : "";
  editWorkEnds.value = d.work_ends ? (d.work_ends.split("T")[1]?.slice(0, 5) ?? "") : "";
  editLocation.value = { lat: d.lat, lon: d.lon };
  editAddress.value = d.address || `${d.lat.toFixed(6)}, ${d.lon.toFixed(6)}`;
  if (d.car_photo_main) editCarPhotoMain.value = d.car_photo_main;
  if (d.car_photos) editCarPhotos.value = [...d.car_photos];
}
watch(
  () => props.driver,
  (d) => populateFromDriver(d),
  { immediate: true },
);

function initCropper() {
  if (cropperInstance) cropperInstance.destroy();
  if (!cropperRef.value) return;
  cropperInstance = new Cropper(cropperRef.value);
  const selection = cropperInstance.getCropperSelection();
  if (selection) {
    selection.aspectRatio = 1;
    selection.initialCoverage = 1;
  }
}

function openCropper(file: File, cb: (base64: string) => void) {
  cropCallback = cb;
  const reader = new FileReader();
  reader.onload = (e) => {
    cropImageSrc.value = e.target?.result as string;
    showCropper.value = true;
    nextTick(() => initCropper());
  };
  reader.readAsDataURL(file);
}

async function onCrop() {
  if (!cropperInstance || !cropCallback) return;
  uploadLoading.value = true;
  try {
    const cropperCanvas = cropperInstance.getCropperCanvas();
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
    cropCallback(scaled.toDataURL("image/jpeg", 0.85));
    cropCallback = null;
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
  cropCallback = null;
  showCropper.value = false;
  if (cropperInstance) {
    cropperInstance.destroy();
    cropperInstance = null;
  }
}

function onCarMainFile(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  openCropper(file, (base64) => {
    editCarPhotoMain.value = base64;
  });
  input.value = "";
}

function removeCarMainPhoto() {
  editCarPhotoMain.value = null;
}

function addCarPhoto() {
  newCarPhotoInput.value?.click();
}

function onNewCarPhotoFile(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  openCropper(file, (base64) => {
    newCarPhotosBase64.value.push(base64);
  });
  input.value = "";
}

function removeExistingCarPhoto(index: number) {
  editCarPhotos.value.splice(index, 1);
}

function removeNewCarPhoto(index: number) {
  newCarPhotosBase64.value.splice(index, 1);
}

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
    let mainUrl = editCarPhotoMain.value;
    if (mainUrl?.startsWith("data:")) {
      carPhotoUploading.value = true;
      mainUrl = await uploadCarPhoto(mainUrl);
    }

    const finalPhotos: string[] = [...editCarPhotos.value];
    if (newCarPhotosBase64.value.length) {
      carPhotoUploading.value = true;
      for (const base64 of newCarPhotosBase64.value) {
        finalPhotos.push(await uploadCarPhoto(base64));
      }
    }
    carPhotoUploading.value = false;

    await updateProfile({
      name: editName.value.trim(),
      lat: editLocation.value.lat,
      lon: editLocation.value.lon,
      phone: editPhone.value.trim() || undefined,
      work_starts: editWorkStarts.value || undefined,
      work_ends: editWorkEnds.value || undefined,
      max_car_weight_kg: editMaxCarWeightKg.value ?? undefined,
      max_car_length_meters: editMaxCarLengthMeters.value ?? undefined,
      car_photo_main: mainUrl || undefined,
      car_photos: finalPhotos.length ? finalPhotos : undefined,
    });
    emit("back");
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Не удалось обновить профиль";
  } finally {
    loading.value = false;
    carPhotoUploading.value = false;
  }
}
</script>

<template>
  <div class="flex flex-col gap-4 overflow-y-auto p-4">
    <div class="flex items-center justify-between">
      <h2 class="text-lg font-medium">Редактирование профиля</h2>
      <button
        @click="emit('back')"
        class="rounded-lg bg-gray-200 px-3 py-1 text-sm hover:bg-gray-300"
      >
        ← Назад
      </button>
    </div>

    <div class="flex flex-col gap-3">
      <input
        v-model="editName"
        type="text"
        placeholder="Имя"
        class="w-full rounded-lg border border-gray-300 px-3 py-2 text-sm focus:ring-2 focus:ring-blue-400 focus:outline-none"
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

      <MapsLocationPicker @click="handleOnMapsPick()" @pick="onLocationPicked" />

      <div v-if="editAddress" class="flex items-center gap-1 text-gray-500">
        <span class="text-green-500">●</span>
        {{ editAddress }}
      </div>

      <div class="flex flex-col items-center gap-2">
        <span class="text-gray-500">Фото автомобиля</span>

        <label class="group relative cursor-pointer self-center">
          <div
            v-if="editCarPhotoMain"
            class="h-28 w-40 overflow-hidden rounded-lg ring-2 ring-gray-300"
          >
            <img
              :src="
                editCarPhotoMain.startsWith('data:')
                  ? editCarPhotoMain
                  : resolveImageUrl(editCarPhotoMain)
              "
              class="h-full w-full object-cover"
            />
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
              editCarPhotoMain ? "Изменить" : "Добавить"
            }}</span>
          </div>
          <input type="file" accept="image/*" class="hidden" @change="onCarMainFile" />
        </label>

        <button
          v-if="editCarPhotoMain"
          @click="removeCarMainPhoto"
          class="text-red-400 transition-colors hover:text-red-500"
        >
          Убрать
        </button>
      </div>

      <div class="flex flex-col items-center gap-2">
        <span class="text-gray-500">Дополнительные фото</span>

        <div class="flex flex-wrap justify-center gap-2">
          <div v-for="(url, idx) in editCarPhotos" :key="'existing-' + idx" class="group relative">
            <div class="h-20 w-20 overflow-hidden rounded-lg ring-2 ring-gray-300">
              <img :src="resolveImageUrl(url)" class="h-full w-full object-cover" />
            </div>
            <button
              @click="removeExistingCarPhoto(idx)"
              class="absolute -top-1.5 -right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-red-400 hover:bg-red-500"
            >
              <span class="text-xs leading-none text-white">×</span>
            </button>
          </div>

          <div
            v-for="(base64, idx) in newCarPhotosBase64"
            :key="'new-' + idx"
            class="group relative"
          >
            <div class="h-20 w-20 overflow-hidden rounded-lg ring-2 ring-gray-300">
              <img :src="base64" class="h-full w-full object-cover" />
            </div>
            <button
              @click="removeNewCarPhoto(idx)"
              class="absolute -top-1.5 -right-1.5 flex h-5 w-5 items-center justify-center rounded-full bg-red-400 hover:bg-red-500"
            >
              <span class="text-xs leading-none text-white">×</span>
            </button>
          </div>

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

      <p v-if="error" class="text-red-500">{{ error }}</p>

      <button
        @click="onSubmit"
        :disabled="loading || !editName.trim() || carPhotoUploading"
        class="w-full rounded-lg bg-green-300 px-4 py-2 font-medium hover:bg-green-400 disabled:opacity-50"
      >
        {{ carPhotoUploading ? "Загрузка фото..." : loading ? "Сохранение..." : "Сохранить" }}
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
