<script setup lang="ts">
import { ref, watch } from "vue";

const props = withDefaults(
  defineProps<{ modelValue?: boolean; openPlaceholder?: string }>(),
  {
    modelValue: true,
    openPlaceholder: "открыть",
  },
);
const emit = defineEmits<{ (e: "update:modelValue", v: boolean): void }>();

const visible = ref(props.modelValue);
watch(
  () => props.modelValue,
  (v) => {
    visible.value = v;
  },
);

function close() {
  visible.value = false;
  emit("update:modelValue", false);
}

function open() {
  visible.value = true;
  emit("update:modelValue", true);
}
</script>

<template>
  <div
    v-show="visible"
    class="fixed bottom-4 left-4 right-4 rounded-2xl bg-gray-200 overflow-y-auto p-4 flex flex-col gap-4 z-10"
  >
    <button
      class="absolute top-3 right-3 w-8 h-8 rounded-full bg-gray-300 hover:bg-gray-400 flex items-center justify-center text-gray-600 font-bold transition z-10"
      @click="close"
    >
      &#10005;
    </button>
    <slot />
  </div>

  <button
    v-show="!visible"
    class="fixed bottom-4 right-4 p-2 hover:bg-gray-200 rounded-full bg-gray-600 shadow-lg flex items-center justify-center text-lg font-bold transition z-20"
    @click="open"
  >
    {{ props.openPlaceholder }}
  </button>
</template>
