interface ImportMetaEnv {
	readonly VITE_ROUTING_API: string
	readonly VITE_GEO_SEARCH_API: string
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
