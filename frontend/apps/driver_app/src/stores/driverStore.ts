import { ref, shallowRef } from "vue";
import { driverClient, freelyAvailableDriverClient } from "driver/api/client.ts";
import { authClient } from "auth/api/client.ts";
import type { Driver } from "driver/types/driver.ts";
import type {
	FreelyAvailableResponse,
	CreateFreelyAvailableRequest,
	UpdateFreelyAvailableRequest,
} from "driver/types/freelyAvailable.ts";

const STATIC_FILES_URL_BASE = import.meta.env.VITE_STATIC_FILES_URL_BASE || "https://localhost:8050";

export function resolveImageUrl(path: string | null): string | null {
	if (!path) return null;
	if (path.startsWith("http")) return path;
	return `${STATIC_FILES_URL_BASE}${path}`;
}

const driver = shallowRef<Driver | null>(null);
const exists = ref(false);
const loading = ref(false);
const error = ref<string | null>(null);

const freelyAvailable = shallowRef<FreelyAvailableResponse | null>(null);
const faExists = ref(false);
const faLoading = ref(false);
const faError = ref<string | null>(null);

export function useDriverProfile() {
	async function fetchProfile() {
		loading.value = true;
		error.value = null;
		try {
			const { data, response } = await driverClient.GET("/driver/profile");
			if (response.status === 404 || !data) {
				driver.value = null;
				exists.value = false;
				return;
			}
			driver.value = data;
			exists.value = true;
		} catch (err) {
			driver.value = null;
			exists.value = false;
			error.value = err instanceof Error ? err.message : "Не удалось загрузить профиль";
		} finally {
			loading.value = false;
		}
	}

	async function uploadProfileImage(imageBase64: string): Promise<string> {
		const { data, error: apiErr } = await authClient.POST("/profile/image", {
			body: { image: imageBase64 },
		});
		if (apiErr) throw new Error("Не удалось отправить картинку");
		return resolveImageUrl(data?.image_url) ?? data?.image_url ?? "";
	}

	async function createProfile(
		name: string,
		lat: number,
		lon: number,
		workStarts?: string,
		workEnds?: string,
	): Promise<void> {
		const { error: apiErr } = await driverClient.POST("/driver/profile", {
			body: {
				name,
				lat,
				lon,
				work_starts: workStarts || undefined,
				work_ends: workEnds || undefined,
			},
		});
		if (apiErr) {
			throw new Error("Не удалось создать профиль водителя");
		}
	}

	async function createFreelyAvailable(req: CreateFreelyAvailableRequest): Promise<void> {
		const { error: apiErr } = await freelyAvailableDriverClient.POST("/driver/freely-available", {
			body: req,
		});
		if (apiErr) {
			throw new Error("Не удалось создать");
		}
	}

	async function fetchFreelyAvailable(userId: number): Promise<void> {
		faLoading.value = true;
		faError.value = null;
		try {
			const { data, response } = await freelyAvailableDriverClient.GET(
				"/driver/{user_id}/freely-available",
				{ params: { path: { user_id: userId } } },
			);
			if (response.status === 404 || !data) {
				freelyAvailable.value = null;
				faExists.value = false;
				return;
			}
			freelyAvailable.value = data;
			faExists.value = true;
		} catch (err) {
			freelyAvailable.value = null;
			faExists.value = false;
			faError.value = err instanceof Error ? err.message : "Ошибка загрузки";
		} finally {
			faLoading.value = false;
		}
	}

	async function updateFreelyAvailable(req: UpdateFreelyAvailableRequest): Promise<void> {
		const { error: apiErr } = await freelyAvailableDriverClient.PUT("/driver/freely-available", {
			body: req,
		});
		if (apiErr) {
			throw new Error("Не удалось обновить");
		}
	}

	async function deleteFreelyAvailable(): Promise<void> {
		const { error: apiErr } = await freelyAvailableDriverClient.DELETE("/driver/freely-available");
		if (apiErr) {
			throw new Error("Не удалось удалить");
		}
	}

	return {
		driver,
		exists,
		loading,
		error,
		fetchProfile,
		uploadProfileImage,
		createProfile,
		freelyAvailable,
		faExists,
		faLoading,
		faError,
		fetchFreelyAvailable,
		createFreelyAvailable,
		updateFreelyAvailable,
		deleteFreelyAvailable,
	};
}
