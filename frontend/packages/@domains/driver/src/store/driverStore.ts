import { atom } from "nanostores";
import type { Driver } from "../types/driver.ts";
import type { FreelyAvailableDriver } from "driver/types/freelyAvailable.js";

export const $driverStore = atom<Driver[]>([])
export const $freelyAvailableDriverStore = atom<FreelyAvailableDriver[]>([])
