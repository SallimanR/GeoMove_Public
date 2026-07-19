import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
	history: createWebHistory(),
	routes: [
		{
			path: '/',
			name: 'home',
			component: () => import('./components/Tabs/TabBar.vue'),
		},
		{
			path: '/triphistory',
			name: 'tripHistory',
			component: () => import('./views/TripHistory.vue'),
		},
	],
})
