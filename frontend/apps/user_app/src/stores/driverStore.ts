import type { Driver } from "driver/types/driver.ts";
import { atom } from "nanostores"

export const $selectedDriver = atom<Driver | null>(null);
export const $driverDropdownOpen = atom<boolean>(false);

const STATIC_FILES_URL_BASE = import.meta.env.VITE_STATIC_FILES_URL_BASE || "";

export function resolveImageUrl(path: string | undefined | null): string | null {
  if (!path) return null;
  if (path.startsWith("http")) return path;
  return `${STATIC_FILES_URL_BASE}${path}`;
}
