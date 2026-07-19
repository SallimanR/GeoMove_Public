import type { Driver } from "driver/types/driver.ts";
import { atom } from "nanostores"

export const $selectedDriver = atom<Driver | null>(null);
