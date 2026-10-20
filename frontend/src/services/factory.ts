import type { AppServices } from '@/services/contracts'
import { createWailsServices } from '@/services/wails'
import { createHttpServices } from '@/services/http'

let servicesInstance: AppServices | null = null

export function getServices(): AppServices {
	if (!servicesInstance) {
		const isDesktop = typeof window !== 'undefined' && !!window.go?.wails
		if (isDesktop) {
			servicesInstance = createWailsServices()
		} else {
			// Browser or dev server mode
			servicesInstance = createWailsServices() // Still returns Wails with fallback, or Http
		}
	}
	return servicesInstance
}

export function setServicesMode(mode: 'wails' | 'http'): AppServices {
	if (mode === 'http') {
		servicesInstance = createHttpServices()
	} else {
		servicesInstance = createWailsServices()
	}
	return servicesInstance
}

