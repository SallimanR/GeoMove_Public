// Astro expects "PUBLIC" prefix excplicitly defined
interface ImportMetaEnv {
	readonly PUBLIC_MAP_STYLE_API: string
	readonly PUBLIC_MAP_TILES_AIP: string
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
