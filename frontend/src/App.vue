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

		<!-- Startup Auto-Update Prompt Modal (Only triggered if internet is connected and newer version exists on GitHub) -->
		<div
			v-if="showStartupUpdate && startupUpdateInfo"
			class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-200 no-print"
		>
			<div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-6 border border-slate-200 space-y-4">
				<!-- Header -->
				<div class="flex items-start justify-between">
					<div class="flex items-center gap-3">
						<div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-700 flex items-center justify-center border border-indigo-200 shrink-0">
							<Sparkles class="w-5 h-5 text-indigo-600" />
						</div>
						<div>
							<div class="flex items-center gap-2">
								<h3 class="text-sm font-bold text-slate-900">Application Update Available</h3>
								<span class="px-2 py-0.5 rounded-full text-xs font-black bg-indigo-600 text-white font-mono">
									v{{ startupUpdateInfo.latestVersion }}
								</span>
							</div>
							<p class="text-xs text-slate-500 mt-0.5">Current Installed: <strong class="font-mono">v{{ startupUpdateInfo.currentVersion }}</strong></p>
						</div>
					</div>
					<button
						v-if="!isApplyingStartupUpdate"
						@click="dismissStartupUpdate"
						class="text-slate-400 hover:text-slate-600 p-1 cursor-pointer"
						title="Dismiss and continue to app"
					>
						<X class="w-4 h-4" />
					</button>
				</div>

				<!-- Zero Data Loss Guarantee Card -->
				<div class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-xs text-emerald-900 space-y-1">
					<div class="flex items-center gap-1.5 font-bold text-emerald-950">
						<ShieldCheck class="w-4 h-4 text-emerald-600 shrink-0" />
						<span>Client Data Safe • Zero Information Loss</span>
					</div>
					<p class="text-[11px] text-emerald-800 leading-relaxed">
						All student enrollments, fee records, financial vouchers, and settings are preserved automatically during the update.
					</p>
				</div>

				<!-- Release Highlights if available -->
				<div v-if="startupUpdateInfo.releaseNotes" class="space-y-1">
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-500">What's New in v{{ startupUpdateInfo.latestVersion }}</span>
					<div class="max-h-28 overflow-y-auto p-2.5 bg-slate-50 border border-slate-200 rounded-lg text-xs text-slate-700 font-mono whitespace-pre-wrap leading-relaxed">
						{{ startupUpdateInfo.releaseNotes }}
					</div>
				</div>

				<!-- Download & Apply Progress -->
				<div v-if="isApplyingStartupUpdate" class="p-4 bg-indigo-50 border border-indigo-200 rounded-xl text-center space-y-2">
					<div class="flex items-center justify-center gap-2 text-xs font-bold text-indigo-700">
						<RefreshCw class="w-4 h-4 animate-spin" />
						<span>{{ startupUpdateStatus }}</span>
					</div>
					<div class="w-full bg-indigo-200 rounded-full h-2 overflow-hidden">
						<div class="bg-indigo-600 h-2 rounded-full transition-all duration-300 animate-pulse w-3/4"></div>
					</div>
					<p class="text-[10px] text-slate-500">The application will automatically restart once the update completes.</p>
				</div>

				<!-- Action Buttons -->
				<div v-else class="flex gap-2.5 pt-2">
					<button
						type="button"
						@click="dismissStartupUpdate"
						class="flex-1 py-2.5 px-3 border border-slate-300 hover:bg-slate-50 text-slate-700 font-semibold rounded-xl text-xs transition cursor-pointer"
					>
						Continue Without Updating
					</button>
					<button
						type="button"
						@click="applyStartupUpdate"
						:disabled="!startupUpdateInfo.downloadURL"
						class="flex-1 py-2.5 px-3 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white font-bold rounded-xl text-xs transition shadow-sm flex items-center justify-center gap-1.5 cursor-pointer"
					>
						<Download class="w-4 h-4" />
						<span>Update Now (v{{ startupUpdateInfo.latestVersion }})</span>
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { AlertTriangle, X, Sparkles, ShieldCheck, Download, RefreshCw } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { getServices } from '@/services/factory'
import type { UpdateInfoDTO } from '@/services/contracts'

const isBrowserMode = ref(false)
const showBrowserBanner = ref(true)
const authStore = useAuthStore()
const services = getServices()

// Startup Update Checking State
const showStartupUpdate = ref(false)
const isApplyingStartupUpdate = ref(false)
const startupUpdateStatus = ref('Downloading update package from GitHub...')
const startupUpdateInfo = ref<UpdateInfoDTO | null>(null)

function dismissStartupUpdate() {
	showStartupUpdate.value = false
}

async function checkStartupUpdate() {
	// 1. Strict offline check: if navigator reports offline, skip immediately
	if (typeof navigator !== 'undefined' && !navigator.onLine) {
		return
	}

	// 2. Wait briefly (700ms) for initial view to render and Wails bridge to settle
	await new Promise(r => setTimeout(r, 700))

	// 3. Query GitHub Releases API with strict 3.5s timeout
	try {
		const repo = localStorage.getItem('sms_update_repo') || ''
		const checkPromise = services.updates.checkForUpdates(repo)
		const timeoutPromise = new Promise<never>((_, reject) =>
			setTimeout(() => reject(new Error('Update check timeout')), 3500)
		)
		const res: any = await Promise.race([checkPromise, timeoutPromise])

		if (res && res.success && res.data && res.data.hasUpdate) {
			startupUpdateInfo.value = res.data
			showStartupUpdate.value = true
		}
	} catch (e) {
		// In case of internet drop, timeout, or DNS failure: silently ignore and let the user work
		console.debug('Startup update check skipped or timed out:', e)
	}
}

async function applyStartupUpdate() {
	if (!startupUpdateInfo.value?.downloadURL) return
	isApplyingStartupUpdate.value = true
	startupUpdateStatus.value = 'Downloading update package from GitHub...'

	try {
		const res = await services.updates.downloadAndApplyUpdate({
			downloadURL: startupUpdateInfo.value.downloadURL,
			downloadFileName: startupUpdateInfo.value.downloadFileName || 'SchoolManagement-windows-amd64.exe'
		})
		if (res.success) {
			startupUpdateStatus.value = 'Update applied! Restarting application in 2 seconds...'
		} else {
			startupUpdateStatus.value = 'Failed: ' + (res.error?.message || 'Update installation error')
			setTimeout(() => {
				isApplyingStartupUpdate.value = false
			}, 3000)
		}
	} catch (e: any) {
		startupUpdateStatus.value = 'Error: ' + (e.message || 'An error occurred')
		setTimeout(() => {
			isApplyingStartupUpdate.value = false
		}, 3000)
	}
}

onMounted(async () => {
	// If window.go is not present after mount, user is running in external browser outside Wails
	if (typeof window !== 'undefined' && !(window as any).go) {
		isBrowserMode.value = true
	}
	if (authStore.isAuthenticated) {
		await authStore.refreshSchool()
	}

	// Check for updates on startup (runs only if internet is available)
	checkStartupUpdate()
})
</script>
