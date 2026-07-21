<script setup lang="ts">
import { ref } from "vue";
import InputNumber from "primevue/inputnumber";
import InputText from "primevue/inputtext";
import Textarea from "primevue/textarea";
import Select from "primevue/select";

const carTypes = [
  { label: "Легковой", value: "Легковой" },
  { label: "Внедорожник", value: "Внедорожник" },
  { label: "Микроавтобус", value: "Микроавтобус" },
  { label: "Грузовик", value: "Грузовик" },
  { label: "Мотоцикл", value: "Мотоцикл" },
  { label: "Спецтехника", value: "Спецтехника" },
  { label: "Электромобиль", value: "Электромобиль" },
  { label: "Другое", value: "Другое" },
];

const props = withDefaults(
  defineProps<{
    wheels: number;
    carType: string;
    carName: string;
    carWeightKg: number;
    carLengthMeters: number;
    carPhotoUrl?: string;
    customerMessage: string;
    mode?: "create" | "edit";
  }>(),
  { mode: "create", carPhotoUrl: "" },
);

const fileInput = ref<HTMLInputElement | null>(null);

const emit = defineEmits<{
  "update:wheels": [value: number];
  "update:carType": [value: string];
  "update:carName": [value: string];
  "update:carWeightKg": [value: number];
  "update:carLengthMeters": [value: number];
  "update:carPhotoUrl": [value: string];
  "update:customerMessage": [value: string];
}>();

function handleFileSelect(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = () => {
    emit("update:carPhotoUrl", reader.result as string);
  };
  reader.readAsDataURL(file);
}
</script>

<template>
  <div class="flex flex-col gap-2">
    <label>Количество заблокированных колёс *</label>
    <InputNumber
      :modelValue="wheels"
      @update:modelValue="(v) => emit('update:wheels', v ?? 0)"
      :min="1"
      :max="18"
    />
  </div>

  <div class="flex flex-col gap-2">
    <label>Тип автомобиля *</label>
    <Select
      :modelValue="carType"
      @update:modelValue="(v) => emit('update:carType', v)"
      :options="carTypes"
      optionLabel="label"
      optionValue="value"
    />
  </div>

  <div class="flex flex-col gap-2">
    <label>Название автомобиля *</label>
    <InputText
      :modelValue="carName"
      @update:modelValue="(v) => emit('update:carName', v ?? '')"
      placeholder="Например: Toyota Camry"
    />
  </div>

  <div class="flex gap-2">
    <div class="flex flex-col gap-2 flex-1">
      <label>Вес авто (кг) *</label>
      <InputNumber
        :modelValue="carWeightKg"
        @update:modelValue="(v) => emit('update:carWeightKg', v ?? 0)"
        :min="1"
      />
    </div>
    <div class="flex flex-col gap-2 flex-1">
      <label>Длина авто (м) *</label>
      <InputNumber
        :modelValue="carLengthMeters"
        @update:modelValue="(v) => emit('update:carLengthMeters', v ?? 0)"
        :min="0.1"
        :minFractionDigits="1"
        :maxFractionDigits="2"
      />
    </div>
  </div>

  <div class="flex flex-col gap-2">
    <label>Фото автомобиля</label>
    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      hidden
      @change="handleFileSelect"
    />
    <button
      type="button"
      class="w-full rounded-xl p-2 bg-green-400 hover:bg-green-500 text-white transition-colors"
      @click="fileInput?.click()"
    >
      {{ carPhotoUrl ? "Фото выбрано" : "Выбрать фото" }}
    </button>
    <img
      v-if="carPhotoUrl"
      :src="carPhotoUrl"
      class="w-full max-h-48 object-contain rounded-lg"
    />
  </div>

  <div class="flex flex-col gap-2">
    <label>Сообщение для водителя</label>
    <Textarea
      :modelValue="customerMessage"
      @update:modelValue="(v) => emit('update:customerMessage', v ?? '')"
      placeholder="Например: автомобиль в кювете, нужен кран..."
      rows="2"
    />
  </div>
</template>
