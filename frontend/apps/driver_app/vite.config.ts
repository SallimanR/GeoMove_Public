import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import vueDevTools from "vite-plugin-vue-devtools";
import tailwindcss from "@tailwindcss/vite";
import path from "path";
import fs from "fs";

function copyNotificationWorker() {
	const src = path.resolve("../../packages/notifications/public/notification-worker.js");
	const dest = path.resolve("public/notification-worker.js");
	fs.mkdirSync(path.dirname(dest), { recursive: true });
	fs.copyFileSync(src, dest);
}

export default defineConfig({
	server: {
		host: "0.0.0.0",
		port: 4324,
		fs: {
			allow: ["../.."],
		},
		watch: {
			ignored: ["!**/packages/**"],
		},
	},

	resolve: {
		alias: {
			src: path.resolve("./src"),
		},
	},

	optimizeDeps: {
		exclude: ["@geomove/maps", "@geomove/geo"],
	},

	plugins: [
		vue(),
		vueDevTools(),
		tailwindcss(),
		{
			name: "copy-notification-worker",
			configureServer() {
				copyNotificationWorker();
			},
			buildStart() {
				copyNotificationWorker();
			},
		},
	],

	build: {
		rolldownOptions: {
			output: {
				manualChunks(id) {
					if (id.includes("node_modules/maplibre-gl")) return "maplibre";
				},
			},
		},
	},
});
