import { atom } from "nanostores";
import type { Driver } from "../types/driver.ts";
import type { FreelyAvailableDriver } from "../types/freely_available.ts";

export const $driverStore = atom<Driver[]>([])
export const $freelyAvailableDriverStore = atom<FreelyAvailableDriver[]>([])
