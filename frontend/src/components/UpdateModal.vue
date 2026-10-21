<template>
	<div v-if="isOpen" class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-150">
		<div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-6 border border-slate-200 space-y-4">
			<!-- Modal Header -->
			<div class="flex items-start justify-between">
				<div class="flex items-center gap-3">
					<div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-700 flex items-center justify-center border border-indigo-200 shrink-0">
						<Sparkles class="w-5 h-5" />
					</div>
					<div>
						<div class="flex items-center gap-2.5">
							<h3 class="text-sm font-bold text-slate-900">Application Updates</h3>
							<button
								@click="checkForUpdates"
								:disabled="isChecking || isApplying"
								class="px-2.5 py-1 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-semibold shadow-2xs flex items-center gap-1.5 transition cursor-pointer disabled:opacity-50"
							>
								<RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isChecking }" />
								<span>{{ isChecking ? 'Checking...' : 'Check' }}</span>
							</button>
						</div>
						<p class="text-xs text-slate-500 mt-0.5">Check for the newest version and patches.</p>
					</div>
				</div>
				<button @click="$emit('close')" class="text-slate-400 hover:text-slate-600 p-1 cursor-pointer">
					<X class="w-4 h-4" />
				</button>
			</div>

			<!-- Status: Checking in progress -->
			<div v-if="isChecking" class="py-8 text-center space-y-2">
				<RefreshCw class="w-8 h-8 mx-auto text-indigo-600 animate-spin" />
				<p class="text-xs font-semibold text-slate-700">Checking for latest release...</p>
			</div>

			<!-- Status: Idle / Ready to check -->
			<div v-else-if="!isChecking && !checkResult && !errorMessage" class="p-6 bg-slate-50 border border-slate-200 rounded-xl text-center space-y-1.5">
				<p class="text-xs font-semibold text-slate-700">Ready to check for updates</p>
				<p class="text-[11px] text-slate-500">
					Click the <strong>Check</strong> button above to verify if a newer release or patch is available.
				</p>
			</div>

			<!-- Status: Up to date -->
			<div v-else-if="checkResult && !checkResult.hasUpdate" class="p-4 bg-emerald-50 border border-emerald-200 rounded-xl text-center space-y-1.5">
				<CheckCircle2 class="w-8 h-8 mx-auto text-emerald-600" />
				<h4 class="text-sm font-bold text-emerald-900">You are on the latest version!</h4>
				<p class="text-xs text-emerald-700">
					Current version <strong class="font-mono">v{{ checkResult.currentVersion }}</strong> matches the latest release.
				</p>
			</div>

			<!-- Status: New Update Available -->
			<div v-else-if="checkResult && checkResult.hasUpdate" class="space-y-3">
				<div class="p-3.5 bg-indigo-50 border border-indigo-200 rounded-xl flex items-start gap-3">
					<ArrowUpCircle class="w-6 h-6 text-indigo-600 shrink-0 mt-0.5" />
					<div class="space-y-0.5">
						<div class="flex items-center gap-2">
							<span class="text-xs font-bold text-indigo-950">New Version Available:</span>
							<span class="px-2 py-0.5 rounded-full text-xs font-black bg-indigo-600 text-white font-mono">
								v{{ checkResult.latestVersion }}
							</span>
						</div>
						<p class="text-xs font-medium text-indigo-800">{{ checkResult.releaseTitle || 'Latest Release Patch' }}</p>
						<p class="text-[10px] text-slate-500">
							Published on {{ formatDate(checkResult.publishedAt) }} • File size: {{ formatBytes(checkResult.assetSizeBytes) }}
						</p>
					</div>
				</div>

				<!-- Release Notes Box -->
				<div v-if="checkResult.releaseNotes" class="space-y-1">
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-500">Changelog &amp; Release Notes</span>
					<div class="max-h-36 overflow-y-auto p-3 bg-slate-50 border border-slate-200 rounded-lg text-xs text-slate-700 font-mono whitespace-pre-wrap leading-relaxed">
						{{ checkResult.releaseNotes }}
					</div>
				</div>

				<!-- Update Action Block -->
				<div class="p-3 bg-slate-100 border border-slate-200 rounded-xl space-y-2">
					<div class="flex items-center justify-between text-xs">
						<span class="font-bold text-slate-800">Download Patch Asset:</span>
						<span class="font-mono text-[11px] text-slate-600 truncate max-w-xs">{{ checkResult.downloadFileName || 'Installer Package' }}</span>
					</div>

					<button
						v-if="!isApplying"
						@click="applyUpdate"
						:disabled="!checkResult.downloadURL"
						class="w-full py-2.5 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 text-white font-bold rounded-lg text-xs transition shadow-xs flex items-center justify-center gap-2 cursor-pointer"
					>
						<Download class="w-4 h-4" />
						<span>Download &amp; Update Application Now</span>
					</button>

					<div v-else class="py-2 text-center space-y-2">
						<div class="flex items-center justify-center gap-2 text-xs font-bold text-indigo-700">
							<RefreshCw class="w-4 h-4 animate-spin" />
							<span>{{ updateProgressStatus }}</span>
						</div>
						<div class="w-full bg-slate-200 rounded-full h-2 overflow-hidden">
							<div class="bg-indigo-600 h-2 rounded-full transition-all duration-300 animate-pulse w-3/4"></div>
						</div>
						<p class="text-[10px] text-slate-500">The app will automatically restart once the patch is applied.</p>
					</div>
				</div>
			</div>

			<!-- Error message if any -->
			<div v-if="errorMessage" class="p-3 bg-rose-50 border border-rose-200 rounded-xl flex items-start gap-2.5 text-xs text-rose-800">
				<AlertTriangle class="w-4 h-4 text-rose-600 shrink-0 mt-0.5" />
				<div>
					<p class="font-bold">Check Failed</p>
					<p class="text-[11px] mt-0.5 text-rose-700">{{ errorMessage }}</p>
				</div>
			</div>

			<!-- Footer -->
			<div class="flex items-center justify-between pt-2 border-t border-slate-100 text-xs text-slate-400">
				<span>Current Installed: <strong class="text-slate-600 font-mono">v{{ currentAppVersion }}</strong></span>
				<button
					@click="$emit('close')"
					:disabled="isApplying"
					class="px-3.5 py-1.5 border border-slate-300 hover:bg-slate-50 text-slate-700 rounded-lg text-xs font-medium transition cursor-pointer"
				>
					Close
				</button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import type { UpdateInfoDTO } from '@/services/contracts'
import { Sparkles, X, RefreshCw, CheckCircle2, ArrowUpCircle, Download, AlertTriangle } from 'lucide-vue-next'

const props = defineProps<{
	isOpen: boolean
}>()

defineEmits<{
	(e: 'close'): void
}>()

const services = getServices()
const toast = useToast()

const repoInput = ref(localStorage.getItem('sms_update_repo') || 'zfhassaan/school-management-system')
const currentAppVersion = ref('1.0.0')
const isChecking = ref(false)
const isApplying = ref(false)
const updateProgressStatus = ref('Downloading update package...')
const checkResult = ref<UpdateInfoDTO | null>(null)
const errorMessage = ref<string | null>(null)

function formatBytes(bytes?: number) {
	if (!bytes || bytes === 0) return '0 B'
	const k = 1024
	const sizes = ['B', 'KB', 'MB', 'GB']
	const i = Math.floor(Math.log(bytes) / Math.log(k))
	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

function formatDate(iso?: string) {
	if (!iso) return 'Recent'
	try {
		return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
	} catch {
		return iso
	}
}

async function loadVersion() {
	try {
		const res = await services.updates.getVersion()
		if (res.success && res.data) {
			currentAppVersion.value = res.data
		}
	} catch (e) {
		console.warn('Failed to fetch version', e)
	}
}

async function checkForUpdates() {
	if (!repoInput.value.trim()) {
		errorMessage.value = 'Please enter a valid GitHub repository (owner/repo)'
		return
	}

	localStorage.setItem('sms_update_repo', repoInput.value.trim())
	isChecking.value = true
	errorMessage.value = null
	checkResult.value = null

	try {
		const res = await services.updates.checkForUpdates(repoInput.value.trim())
		if (res.success && res.data) {
			checkResult.value = res.data
			if (res.data.hasUpdate) {
				toast.info(`New version v${res.data.latestVersion} found on GitHub!`)
			} else {
				toast.success('Your application is completely up to date.')
			}
		} else {
			errorMessage.value = res.error?.message || 'Could not find any releases on the specified GitHub repository.'
		}
	} catch (e: any) {
		errorMessage.value = e.message || 'Network error connecting to GitHub'
	} finally {
		isChecking.value = false
	}
}

async function applyUpdate() {
	if (!checkResult.value?.downloadURL) return
	isApplying.value = true
	updateProgressStatus.value = 'Downloading update package from GitHub...'

	try {
		const res = await services.updates.downloadAndApplyUpdate({
			downloadURL: checkResult.value.downloadURL,
			downloadFileName: checkResult.value.downloadFileName || 'update.exe'
		})

		if (res.success) {
			updateProgressStatus.value = 'Patch downloaded! Restarting application...'
			toast.success('Update downloaded! Restarting application in 2 seconds...')
		} else {
			toast.error(res.error?.message || 'Failed to apply update patch')
			isApplying.value = false
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred while updating')
		isApplying.value = false
	}
}

watch(
	() => props.isOpen,
	(open) => {
		if (open) {
			loadVersion()
			checkForUpdates()
		}
	}
)

onMounted(() => {
	loadVersion()
	if (props.isOpen) {
		checkForUpdates()
	}
})
</script>
