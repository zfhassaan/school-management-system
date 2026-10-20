<template>
	<div class="h-screen w-screen overflow-hidden flex flex-col">
		<!-- Browser Preview Warning Banner (when opened in external browser like Chrome instead of Wails window) -->
		<div
			v-if="isBrowserMode && showBrowserBanner"
			class="bg-amber-500 text-slate-900 text-xs px-4 py-1.5 flex items-center justify-between shrink-0 shadow-xs z-50 font-medium"
		>
			<div class="flex items-center gap-2">
				<AlertTriangle class="w-4 h-4 shrink-0" />
				<span>
					<strong>Browser Preview:</strong> Running in an external browser tab. Use the native Wails Desktop Application window to interact with the offline SQLite database.
				</span>
			</div>
			<button @click="showBrowserBanner = false" class="text-slate-800 hover:text-black text-xs font-bold px-1.5 py-0.5 rounded hover:bg-amber-400">
				✕
			</button>
		</div>

		<div class="flex-1 overflow-hidden">
			<router-view />
		</div>
		<ToastContainer />
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { AlertTriangle } from 'lucide-vue-next'

const isBrowserMode = ref(false)
const showBrowserBanner = ref(true)

onMounted(() => {
	// If window.go is not present after mount, user is running in external browser outside Wails
	if (typeof window !== 'undefined' && !(window as any).go) {
		isBrowserMode.value = true
	}
})
</script>

