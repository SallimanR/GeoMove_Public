import { addPopupToMap } from "@geomove/maps";
import { $driverStore } from "driver/store/driverStore.ts";
import DriverPopup from "src/components/Tabs/MapsTab/DriverPopup.vue";

export function displayDriverPopups() {
	for (const driver of $driverStore.get()) {
		addPopupToMap(
			driver.lat,
			driver.lon,
			DriverPopup,
			driver,
			"drivers",
		);
	}
}
