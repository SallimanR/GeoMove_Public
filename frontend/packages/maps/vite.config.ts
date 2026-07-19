import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import tailwindcss from "@tailwindcss/vite"
import dts from "vite-plugin-dts"
import { resolve } from "path"

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    dts({
      tsconfigPath: "./tsconfig.build.json",
      include: ["src"],
      outDir: "dist",
      rollupTypes: false,
    }),
  ],
  build: {
    lib: {
      entry: resolve(__dirname, "src/build-entry.ts"),
      name: "GeomoveMaps",
      formats: ["es"],
      fileName: () => "index.js",
    },
    rollupOptions: {
      external: [
        "vue",
        "maplibre-gl",
        "@deck.gl/core",
        "@deck.gl/layers",
        "@deck.gl/mapbox",
        "deck.gl",
        "nanostores",
        "@nanostores/vue",
        "pmtiles",
        "@geomove/geo",
        "@capacitor/cli",
        "@capacitor/core",
        "@capacitor/geolocation",
      ],
      output: {
        globals: {
          vue: "Vue",
        },
      },
    },
    cssCodeSplit: false,
  },
})
