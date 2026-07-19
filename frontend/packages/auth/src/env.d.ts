interface ImportMetaEnv {
	readonly VITE_AUTH_API_BASE: string
	readonly VITE_GOOGLE_CLIENT_ID: string
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
