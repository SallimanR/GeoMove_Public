<script setup lang="ts">
import { ref, onMounted, inject, watch } from "vue";
import { useStore } from "@nanostores/vue";
import {
  MapsLocationPicker,
  $locationPicking,
  setPickCallback,
} from "@geomove/maps";
import type { GeoPoint } from "@geomove/maps";
import { TimePicker } from "ui";
import DatePicker from "primevue/datepicker";
import SingIn from "auth/components/SignIn.vue";
import { $user, $isAuthenticated, $loading, checkAuth, setUser } from "auth";
import { useDriverProfile } from "../../stores/driverStore";
import { ACTIVE_TAB_KEY } from "../../injectionKeys";

const user = useStore($user);
const isAuthenticated = useStore($isAuthenticated);
const authLoading = useStore($loading);
const {
  driver,
  exists: driverExists,
  loading: driverLoading,
  fetchProfile,
  freelyAvailable,
  faExists,
  faLoading,
  createFreelyAvailable,
  fetchFreelyAvailable,
  updateFreelyAvailable,
  deleteFreelyAvailable,
} = useDriverProfile();

const activeTab = inject(ACTIVE_TAB_KEY)!;

const editMode = ref(false);
const loading = ref(false);
const error = ref<string | null>(null);

const fromDate = ref<Date | null>(null);
const fromTime = ref("");
const toDate = ref<Date | null>(null);
const toTime = ref("");
const enRouteOrder = ref(false);
const tariffPerKm = ref<number | undefined>(undefined);

interface LocationEntry {
  lat: number;
  lon: number;
  address: string;
}
const fromLocation = ref<LocationEntry | null>(null);
const toLocations = ref<LocationEntry[]>([]);
let nextLocationIndex = 0;
interface PendingPick {
  type: "from" | "to";
  index?: number;
}
const pendingPick = ref<PendingPick | null>(null);

onMounted(async () => {
  await checkAuth();
  if (isAuthenticated.value) {
    await fetchProfile();
  }
});

watch([isAuthenticated, driverExists], async () => {
  if (isAuthenticated.value && driverExists.value && user.value) {
    await fetchFreelyAvailable(user.value.id);
  }
});

function onLoginSuccess(u: { id: number; email: string | null }) {
  setUser({ id: u.id, email: u.email, phone: null, profile_image: null });
  fetchProfile();
}

function startCreate() {
  fromDate.value = null;
  fromTime.value = "";
  toDate.value = null;
  toTime.value = "";
  enRouteOrder.value = false;
  tariffPerKm.value = undefined;
  fromLocation.value = null;
  toLocations.value = [];
  editMode.value = true;
}

function startEdit() {
  if (!freelyAvailable.value) return;
  fromDate.value = isoToDate(freelyAvailable.value.from_date);
  fromTime.value = toTimeFromIso(freelyAvailable.value.from_date);
  toDate.value = isoToDate(freelyAvailable.value.to_date);
  toTime.value = toTimeFromIso(freelyAvailable.value.to_date);
  enRouteOrder.value = freelyAvailable.value.en_route_order ?? false;
  tariffPerKm.value = freelyAvailable.value.tariff_per_km ?? undefined;
  fromLocation.value = {
    lat: freelyAvailable.value.from_location.lat,
    lon: freelyAvailable.value.from_location.lon,
    address: freelyAvailable.value.from_location.address || "",
  };
  toLocations.value = (freelyAvailable.value.to_locations || []).map((l) => ({
    lat: l.lat,
    lon: l.lon,
    address: l.address || "",
  }));
  editMode.value = true;
}

function cancelEdit() {
  editMode.value = false;
  error.value = null;
}

function onLocationPicked(point: GeoPoint, address: string) {
  if (!pendingPick.value) return;
  const entry: LocationEntry = { lat: point.lat, lon: point.lon, address };
  if (pendingPick.value.type === "from") {
    fromLocation.value = entry;
  } else if (
    pendingPick.value.type === "to" &&
    pendingPick.value.index !== undefined
  ) {
    toLocations.value[pendingPick.value.index] = entry;
  }
  pendingPick.value = null;
  activeTab.value = "orderSearchTab";
}

function pickFromLocation() {
  pendingPick.value = { type: "from" };
  activeTab.value = "mapsTab";
}

function pickToLocation(index: number) {
  pendingPick.value = { type: "to", index };
  activeTab.value = "mapsTab";
}

function addToLocation() {
  const idx = toLocations.value.length;
  toLocations.value.push({ lat: 0, lon: 0, address: "" });
  pendingPick.value = { type: "to", index: idx };
  setPickCallback((point: GeoPoint, address: string) =>
    onLocationPicked(point, address),
  );
  $locationPicking.set(true);
}

function removeToLocation(index: number) {
  toLocations.value.splice(index, 1);
}

async function onSubmit() {
  error.value = null;
  if (!fromLocation.value) {
    error.value = "Выберите точку отправления";
    return;
  }
  loading.value = true;
  try {
    const body = {
      from_date: dateToIso(fromDate.value, fromTime.value),
      to_date: dateToIso(toDate.value, toTime.value),
      from_location: {
        lat: fromLocation.value.lat,
        lon: fromLocation.value.lon,
      },
      to_locations: toLocations.value.map((l) => ({ lat: l.lat, lon: l.lon })),
      en_route_order: enRouteOrder.value || undefined,
      tariff_per_km: tariffPerKm.value ?? undefined,
    };

    if (faExists.value) {
      await updateFreelyAvailable(body);
    } else {
      await createFreelyAvailable(body);
    }

    if (user.value) {
      await fetchFreelyAvailable(user.value.id);
    }
    editMode.value = false;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Ошибка сохранения";
  } finally {
    loading.value = false;
  }
}

async function onDelete() {
  loading.value = true;
  try {
    await deleteFreelyAvailable();
    freelyAvailable.value = null;
    faExists.value = false;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Ошибка удаления";
  } finally {
    loading.value = false;
  }
}

function isoToDate(iso: string): Date | null {
  const d = new Date(iso);
  return isNaN(d.getTime()) ? null : d;
}

function dateToIso(d: Date | null, time: string): string {
  if (!d) return "";
  const pad = (n: number) => String(n).padStart(2, "0");
  const datePart = `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
  return new Date(`${datePart}T${time || "00:00"}`).toISOString();
}

function toTimeFromIso(isoString: string): string {
  const d = new Date(isoString);
  if (isNaN(d.getTime())) return "";
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function formatDateTime(isoString: string): string {
  const d = new Date(isoString);
  if (isNaN(d.getTime())) return "—";
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${pad(d.getDate())}.${pad(d.getMonth() + 1)}.${d.getFullYear()}, ${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function formatTariff(val: number | null): string {
  if (val == null) return "—";
  return `${val} ₽/км`;
}
</script>

<template>
  <div
    v-if="authLoading || driverLoading || faLoading"
    class="flex items-center justify-center h-full p-4"
  >
    <p class="text-gray-500">Загрузка...</p>
  </div>

  <div
    v-else-if="!isAuthenticated"
    class="flex items-center justify-center h-full"
  >
    <SingIn @login-success="onLoginSuccess" />
  </div>

  <div
    v-else-if="!driverExists"
    class="flex items-center justify-center h-full p-4"
  >
    <p class="text-gray-500">Сначала создайте профиль водителя</p>
  </div>

  <!-- Edit/Create form -->
  <div
    v-else-if="editMode"
    class="flex flex-col h-full overflow-y-auto p-4 gap-3"
  >
    <h2 class="text-lg font-medium text-center">Свободный эвакуатор</h2>

    <div class="flex flex-col gap-2">
      <label class="text-gray-600">Дата и время начала</label>
      <div class="flex gap-2">
        <DatePicker v-model="fromDate" showIcon fluid :showOnFocus="false" class="flex-1" />
        <TimePicker v-model="fromTime" placeholder="Время" class="flex-1" />
      </div>
    </div>

    <div class="flex flex-col gap-2">
      <label class="text-gray-600">Дата и время окончания</label>
      <div class="flex gap-2">
        <DatePicker v-model="toDate" showIcon fluid :showOnFocus="false" class="flex-1" />
        <TimePicker v-model="toTime" placeholder="Время" class="flex-1" />
      </div>
    </div>

    <div class="flex flex-col gap-2">
      <label class="text-gray-600">Точка отправления:</label>
      <div class="flex flex-row">
        <MapsLocationPicker
          @click="pickFromLocation()"
          @pick="onLocationPicked"
        />
      </div>
      <span v-if="fromLocation" class="text-gray-500">
        <span class="text-green-500">●</span>
        {{
          fromLocation.address ||
          fromLocation.lat.toFixed(5) + ", " + fromLocation.lon.toFixed(5)
        }}
      </span>
    </div>

    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between">
        <label class="text-gray-600">Точки назначения:</label>
        <button
          @click="addToLocation"
          type="button"
          class="rounded-xl p-2 bg-green-300 hover:text-blue-600"
        >
          + Добавить
        </button>
      </div>
      <div
        v-for="(loc, idx) in toLocations"
        :key="idx"
        class="flex flex-col gap-1"
      >
        <div class="flex items-center gap-1">
          <span class="text-gray-400">Точка {{ idx + 1 }}:</span>
        </div>
        <span v-if="loc.lat !== 0 || loc.address" class="text-gray-500">
          {{ loc.address || loc.lat.toFixed(5) + ", " + loc.lon.toFixed(5) }}
        </span>

        <div class="flex flex-row gap-2">
          <MapsLocationPicker
            @click="pickToLocation(idx)"
            @pick="onLocationPicked"
          />
          <button
            @click="removeToLocation(idx)"
            type="button"
            class="rounded-xl p-2 bg-red-300 hover:bg-red-400"
          >
            Убрать
          </button>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-3">
      <label class="text-gray-600">Попутный заказ</label>
      <input v-model="enRouteOrder" type="checkbox" class="w-4 h-4" />
    </div>

    <div class="flex flex-col gap-2">
      <label class="text-gray-600">Тариф за км (₽)</label>
      <input
        v-model.number="tariffPerKm"
        type="number"
        step="0.01"
        placeholder="50"
        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
      />
    </div>

    <p v-if="error" class="text-red-400">{{ error }}</p>

    <div class="flex gap-2">
      <button
        @click="cancelEdit"
        type="button"
        class="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg text-gray-600 hover:bg-gray-50 transition-colors"
      >
        Отмена
      </button>
      <button
        @click="onSubmit"
        :disabled="loading || !fromDate || !toDate || !fromLocation"
        class="flex-1 px-4 py-2.5 bg-blue-500 text-white rounded-lg shadow-sm hover:bg-blue-600 disabled:opacity-50 transition-colors font-medium"
      >
        {{ loading ? "Сохранение..." : "Сохранить" }}
      </button>
    </div>
  </div>

  <!-- Create button when no entry -->
  <div
    v-else-if="!faExists && !editMode"
    class="flex flex-col items-center justify-center h-full p-4 gap-4"
  >
    <p class="text-gray-500">У вас нет активной заявки</p>
    <button
      @click="startCreate"
      class="px-6 py-2.5 bg-blue-500 text-white rounded-lg shadow-sm hover:bg-blue-600 transition-colors font-medium"
    >
      Свободный эвакуатор
    </button>
  </div>

  <!-- Display mode -->
  <div
    v-else-if="faExists && !editMode && freelyAvailable"
    class="flex flex-col h-full overflow-y-auto p-4 gap-3"
  >
    <h2 class="text-lg font-medium text-center">Свободный эвакуатор</h2>

    <div class="flex flex-col gap-2 bg-gray-50 rounded-lg p-3">
      <div class="flex justify-between">
        <span class="text-gray-500">Период:</span>
        <span
          >{{ formatDateTime(freelyAvailable.from_date) }} —
          {{ formatDateTime(freelyAvailable.to_date) }}</span
        >
      </div>
      <div class="flex justify-between">
        <span class="text-gray-500">Тариф:</span>
        <span>{{ formatTariff(freelyAvailable.tariff_per_km) }}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-gray-500">Попутный заказ:</span>
        <span>{{ freelyAvailable.en_route_order ? "Да" : "Нет" }}</span>
      </div>
    </div>

    <div class="flex flex-col gap-1">
      <p class="text-gray-500">Точка отправления:</p>
      <p class="">
        {{
          freelyAvailable.from_location.address ||
          freelyAvailable.from_location.lat.toFixed(5) +
            ", " +
            freelyAvailable.from_location.lon.toFixed(5)
        }}
      </p>
    </div>

    <div
      v-if="freelyAvailable.to_locations?.length"
      class="flex flex-col gap-1"
    >
      <p class="text-gray-500">Точки назначения:</p>
      <div
        v-for="(loc, idx) in freelyAvailable.to_locations"
        :key="idx"
        class="pl-2"
      >
        {{ idx + 1 }}.
        {{ loc.address || loc.lat.toFixed(5) + ", " + loc.lon.toFixed(5) }}
      </div>
    </div>

    <p v-if="error" class="text-red-500">{{ error }}</p>

    <div class="flex gap-2 mt-2">
      <button
        @click="onDelete"
        :disabled="loading"
        class="flex-1 px-4 py-2.5 border border-red-300 text-red-500 rounded-lg hover:bg-red-50 disabled:opacity-50 transition-colors"
      >
        {{ loading ? "Удаление..." : "Удалить" }}
      </button>
      <button
        @click="startEdit"
        class="flex-1 px-4 py-2.5 bg-blue-500 text-white rounded-lg shadow-sm hover:bg-blue-600 transition-colors font-medium"
      >
        Изменить
      </button>
    </div>
  </div>
</template>
