import { Popup, PopupOptions } from "maplibre-gl";
import { $mapInstance } from "../stores/mapsStore";
import { createApp, type Component, type App, type ComponentPublicInstance } from "vue";

export interface PopupEntry {
	popup: Popup;
	app?: App;
	instance?: ComponentPublicInstance;
	destroy: () => void;
}

const popupGroups = new Map<string, PopupEntry[]>();

export function addPopupToMap<T extends object = {}>(
	lat: number,
	lon: number,
	component: Component,
	props: T = {} as T,
	group: string,
	popupOptions: PopupOptions = { offset: 0, closeOnClick: false },
): PopupEntry | undefined {
	const map = $mapInstance.value;
	if (!map) return undefined;

	const popup = new Popup({
		closeButton: false,
		closeOnClick: false,
		anchor: "bottom",
		...popupOptions,
	});
	popup
		.setLngLat([lon, lat])
		.setHTML(`<div id='${group}-map-popup'></div>`)
		.addTo(map);

	const mountPoint = popup
		.getElement()
		.querySelector(`#${group}-map-popup`);
	if (!mountPoint) return undefined

	const app = createApp(component, props as Record<string, any>)
	const instance = app.mount(mountPoint);

	popup.on('close', () => {
		if (app) {
			app.unmount();
		}
	})

	const destroy = (): void => {
		popup.remove();
		app.unmount();
	};

	const entry: PopupEntry = { popup, app, instance, destroy };
	if (!popupGroups.has(group)) {
		popupGroups.set(group, []);
	}
	popupGroups.get(group)!.push(entry);

	return entry;
}

export function removePopupsByGroup(group: string): void {
	const entries = popupGroups.get(group);
	if (!entries) return;
	entries.forEach(({ destroy }) => { destroy(); });
	popupGroups.delete(group);
}

export function removeAllPopups(): void {
	for (const entries of popupGroups.values()) {
		entries.forEach(({ destroy }) => { destroy(); });
	}
	popupGroups.clear();
}
