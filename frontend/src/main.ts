import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from '@/router'
import App from '@/App.vue'
import '@/assets/main.css'

/**
 * Poll briefly for Wails WebView2 bridge to inject window.go.
 * In desktop execution this resolves in ~20-50ms before app mount.
 */
async function waitForWailsRuntime(timeoutMs = 1500): Promise<boolean> {
	if (typeof window === 'undefined') return false
	const start = Date.now()
	while (Date.now() - start < timeoutMs) {
		const go = (window as any).go
		if (go && (go.wails || go.main || Object.keys(go).length > 0)) {
			return true
		}
		await new Promise(resolve => setTimeout(resolve, 15))
	}
	return false
}

async function bootstrap() {
	await waitForWailsRuntime(1500)

	const app = createApp(App)
	const pinia = createPinia()

	app.use(pinia)
	app.use(router)

	app.mount('#app')
}

bootstrap()

