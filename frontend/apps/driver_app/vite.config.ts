import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig({
	server: {
		host: '0.0.0.0',
		port: 4324,
		fs: {
			allow: ['../..'],
		},
		watch: {
			ignored: ['!**/packages/**'],
		},
	},

	resolve: {
		alias: {
			src: path.resolve('./src'),
		},
	},

	optimizeDeps: {
		exclude: ['@geomove/maps', '@geomove/geo'],
	},

	plugins: [vue(), tailwindcss()],

	build: {
		rolldownOptions: {
			output: {
				manualChunks(id) {
					if (id.includes('node_modules/maplibre-gl')) return 'maplibre'
				},
			},
		},
	},
})
