import { atom } from "nanostores";
// import { type RealtimeDriver } from "src/types/realtimeDriver.ts";
import { type DriverRealtime } from "src/api/proto/location_update.ts";

export const $realtimeDrivers = atom<DriverRealtime[]>([])
export const $lastFetchTime = atom<number>(Date.now())

export function updateRealtimeDrivers(newData: DriverRealtime[]): void {
	$realtimeDrivers.set(newData)
	$lastFetchTime.set(Date.now())
}
