interface ImportMetaEnv {
	readonly VITE_VK_APP_ID: number
	readonly VITE_VK_REDIRECT_URL: string
	readonly VITE_GOOGLE_CLIENT_ID: string
	readonly VITE_STYLE_API: string
	readonly VITE_WS_HOST: string
	readonly VITE_STATIC_FILES_URL_BASE: string
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
