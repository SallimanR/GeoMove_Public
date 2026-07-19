import { createApp } from 'vue'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import GoogleSignInPlugin from 'vue3-google-signin'
import { configureGeo } from '@geomove/geo'
import { router } from './router'
import AppShell from './components/AppShell.vue'
import './styles/global.css'

configureGeo({
  routingApi: import.meta.env.VITE_ROUTING_API,
  geocodingApi: import.meta.env.VITE_GEO_SEARCH_API,
})

const app = createApp(AppShell)

app.use(router)
app.use(PrimeVue, {
	theme: {
		preset: Aura,
		options: {
			prefix: 'p',
			darkModeSelector: 'light',
			cssLayer: false,
		},
	},
})

app.use(GoogleSignInPlugin, {
	clientId: import.meta.env.VITE_GOOGLE_CLIENT_ID,
})

app.mount('#app')
