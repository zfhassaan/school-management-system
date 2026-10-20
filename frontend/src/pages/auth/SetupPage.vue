<template>
	<div class="min-h-screen w-screen flex items-center justify-center bg-slate-100 p-4">
		<div class="max-w-2xl w-full bg-white rounded-2xl shadow-xl border border-slate-200/80 p-8">
			<!-- Header -->
			<div class="text-center mb-8">
				<div class="w-12 h-12 bg-indigo-600 text-white rounded-xl mx-auto flex items-center justify-center font-bold text-lg mb-3 shadow-md shadow-indigo-200">
					SMS
				</div>
				<h1 class="text-xl font-bold text-slate-900">First-Run School Setup</h1>
				<p class="text-xs text-slate-500 mt-1">Initialize your school profile and create the first Super Administrator account.</p>
			</div>

			<!-- Step 1: Form -->
			<div v-if="!recoveryCodes">
				<form @submit.prevent="handleSetup" class="space-y-6">
					<!-- School Section -->
					<div class="space-y-3">
						<h2 class="text-xs font-bold uppercase tracking-wider text-indigo-700 flex items-center gap-1.5 border-b pb-1">
							<Building2 class="w-4 h-4" /> 1. School Information
						</h2>
						<div class="grid grid-cols-2 gap-4">
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">School Name *</label>
								<input
									v-model="form.schoolName"
									type="text"
									required
									placeholder="e.g. Al-Huda Grammar School"
									class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
								/>
							</div>
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">School Code *</label>
								<input
									v-model="form.schoolCode"
									type="text"
									required
									placeholder="e.g. ALHUDA"
									class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300 uppercase"
								/>
							</div>
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">Timezone *</label>
								<select
									v-model="form.timezone"
									class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
								>
									<option value="Asia/Karachi">Asia/Karachi (PKT +05:00)</option>
									<option value="UTC">UTC (+00:00)</option>
									<option value="Asia/Dubai">Asia/Dubai (GST +04:00)</option>
									<option value="Asia/Riyadh">Asia/Riyadh (AST +03:00)</option>
									<option value="Asia/Dhaka">Asia/Dhaka (BST +06:00)</option>
									<option value="Asia/Kolkata">Asia/Kolkata (IST +05:30)</option>
								</select>
							</div>
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">Currency Code *</label>
								<div class="flex gap-2">
									<input
										v-model="form.currencyCode"
										type="text"
										required
										maxlength="3"
										placeholder="PKR"
										class="w-2/3 px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300 uppercase"
									/>
									<select
										v-model="form.currencyMinorUnits"
										class="w-1/3 px-2 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
										title="Minor decimal places"
									>
										<option :value="0">0 dec</option>
										<option :value="2">2 dec</option>
									</select>
								</div>
							</div>
						</div>
					</div>

					<!-- Super Admin Section -->
					<div class="space-y-3 pt-2">
						<h2 class="text-xs font-bold uppercase tracking-wider text-indigo-700 flex items-center gap-1.5 border-b pb-1">
							<ShieldCheck class="w-4 h-4" /> 2. Super Administrator Account
						</h2>
						<div class="grid grid-cols-2 gap-4">
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">Username *</label>
								<input
									v-model="form.adminUsername"
									type="text"
									required
									placeholder="admin"
									class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
								/>
							</div>
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">Display Name *</label>
								<input
									v-model="form.adminDisplayName"
									type="text"
									required
									placeholder="Headmaster / Admin"
									class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
								/>
							</div>
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">Email Address</label>
								<input
									v-model="form.adminEmail"
									type="email"
									placeholder="admin@school.internal"
									class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
								/>
							</div>
							<div>
								<label class="block text-xs font-semibold text-slate-700 mb-1">Password * (Min 10 chars)</label>
								<input
									v-model="form.adminPassword"
									type="password"
									required
									minlength="10"
									placeholder="••••••••••••"
									class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
								/>
							</div>
						</div>
					</div>

					<div v-if="errorMessage" class="p-3 rounded-lg bg-rose-50 border border-rose-200 text-rose-700 text-xs">
						{{ errorMessage }}
					</div>

					<button
						type="submit"
						:disabled="isLoading"
						class="w-full py-3 px-4 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-semibold rounded-xl text-sm transition shadow-sm flex items-center justify-center gap-2"
					>
						<Loader2 v-if="isLoading" class="w-4 h-4 animate-spin" />
						<span>{{ isLoading ? 'Initializing System...' : 'Complete First-Run Setup' }}</span>
					</button>
				</form>
			</div>

			<!-- Step 2: Recovery Codes Presentation -->
			<div v-else class="space-y-6">
				<div class="p-4 bg-amber-50 border border-amber-200 rounded-xl">
					<div class="flex items-start gap-3">
						<KeyRound class="w-5 h-5 text-amber-700 shrink-0 mt-0.5" />
						<div>
							<h3 class="text-sm font-bold text-amber-900">Save Your Account Recovery Codes</h3>
							<p class="text-xs text-amber-800 mt-1">
								These 10 one-time codes can restore access if you ever lose your Super Admin password.
								<strong>They will not be displayed again.</strong> Please store them in a secure physical or password-vault location.
							</p>
						</div>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-2 bg-slate-50 p-4 rounded-xl border border-slate-200 font-mono text-xs text-slate-800">
					<div v-for="(code, idx) in recoveryCodes" :key="idx" class="p-2 bg-white rounded border border-slate-200 text-center">
						<span class="text-slate-400 mr-2">{{ idx + 1 }}.</span>{{ code }}
					</div>
				</div>

				<div class="flex gap-3">
					<button
						type="button"
						@click="copyCodes"
						class="flex-1 py-2.5 px-4 bg-slate-100 hover:bg-slate-200 text-slate-800 font-medium rounded-xl text-sm transition flex items-center justify-center gap-2"
					>
						<Copy class="w-4 h-4" />
						<span>{{ isCopied ? 'Copied to Clipboard!' : 'Copy All Codes' }}</span>
					</button>

					<button
						type="button"
						@click="proceedToLogin"
						class="flex-1 py-2.5 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-xl text-sm transition flex items-center justify-center gap-2 shadow-sm"
					>
						<Check class="w-4 h-4" />
						<span>I Have Saved These Codes</span>
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Building2, ShieldCheck, KeyRound, Copy, Check, Loader2 } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
	schoolName: 'Horizon International Academy',
	schoolCode: 'HORIZON',
	timezone: 'Asia/Karachi',
	currencyCode: 'PKR',
	currencyMinorUnits: 2,
	adminUsername: 'admin',
	adminDisplayName: 'Headmaster',
	adminEmail: 'admin@horizon.edu',
	adminPassword: ''
})

const isLoading = ref(false)
const errorMessage = ref('')
const recoveryCodes = ref<string[] | null>(null)
const isCopied = ref(false)

async function handleSetup() {
	errorMessage.value = ''
	isLoading.value = true
	try {
		const res = await authStore.firstRunSetup(form)
		if (res.success && res.data) {
			recoveryCodes.value = res.data.recoveryCodes
		} else {
			errorMessage.value = res.error || 'Setup failed'
		}
	} catch (e: any) {
		errorMessage.value = e.message || 'An error occurred'
	} finally {
		isLoading.value = false
	}
}

function copyCodes() {
	if (!recoveryCodes.value) return
	const text = recoveryCodes.value.map((c, i) => `${i + 1}. ${c}`).join('\n')
	navigator.clipboard.writeText(text)
	isCopied.value = true
	setTimeout(() => {
		isCopied.value = false
	}, 2000)
}

function proceedToLogin() {
	router.push('/login')
}
</script>

