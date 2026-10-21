<template>
	<div class="h-screen w-screen overflow-hidden flex flex-col print:h-auto print:w-auto print:overflow-visible">
		<!-- Browser Preview Warning Banner (when opened in external browser like Chrome instead of Wails window) -->
		<div
			v-if="isBrowserMode && showBrowserBanner"
			class="bg-amber-500 text-slate-900 text-xs px-4 py-1.5 flex items-center justify-between shrink-0 shadow-xs z-50 font-medium no-print"
		>
			<div class="flex items-center gap-2">
				<AlertTriangle class="w-4 h-4 shrink-0" />
				<span>
					<strong>Browser Preview:</strong> Running in an external browser tab. Use the native Wails Desktop Application window to interact with the offline SQLite database.
				</span>
			</div>
			<button @click="showBrowserBanner = false" class="text-slate-800 hover:text-black text-xs font-bold p-1 rounded hover:bg-amber-400">
				<X class="w-3.5 h-3.5" />
			</button>
		</div>

		<div class="flex-1 overflow-hidden print:overflow-visible">
			<router-view />
		</div>
		<ToastContainer class="no-print" />
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { AlertTriangle, X } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const isBrowserMode = ref(false)
const showBrowserBanner = ref(true)
const authStore = useAuthStore()

onMounted(async () => {
	// If window.go is not present after mount, user is running in external browser outside Wails
	if (typeof window !== 'undefined' && !(window as any).go) {
		isBrowserMode.value = true
	}
	if (authStore.isAuthenticated) {
		await authStore.refreshSchool()
	}
})
</script>

