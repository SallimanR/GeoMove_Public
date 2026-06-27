import { defineConfig, memoryCache } from "astro/config";
import { visualizer } from "rollup-plugin-visualizer";

import vue from "@astrojs/vue";
import node from "@astrojs/node";

import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  output: "server",
  adapter: node({
    mode: "standalone",
  }),

  integrations: [
    vue({
      appEntrypoint: "./src/pages/_app.ts",
      devtools: false,
    }),
  ],

  cache: {
    provider: memoryCache(),
  },

  vite: {
    plugins: [tailwindcss(), visualizer()],
  },
});
