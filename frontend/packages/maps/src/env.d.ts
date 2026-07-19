declare module '*.css'
declare module 'maplibre-gl/dist/maplibre-gl.css'

interface ImportMetaEnv {
	readonly VITE_MAP_STYLE_API: string
	readonly VITE_MAP_TILES_API: string
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
