import { atom } from "nanostores";
import type { MovingPath } from "maps/composables/useMovingIconLayer.ts";

export const $movingDrivers = atom<MovingPath[]>([])
export const $lastFetchTime = atom<number>(Date.now())

export function updateMovingDrivers(newData: MovingPath[]): void {
	$movingDrivers.set(newData)
	$lastFetchTime.set(Date.now())
}
