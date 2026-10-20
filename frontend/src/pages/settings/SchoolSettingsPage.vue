<template>
	<div class="space-y-6 max-w-5xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">School Profile & System Settings</h1>
				<p class="text-xs text-slate-500 mt-0.5">Configure school identity, financial conventions, and numbering sequences.</p>
			</div>
			<button
				@click="saveAll"
				:disabled="isSaving"
				class="py-2 px-4 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-2"
			>
				<Save class="w-4 h-4" />
				<span>{{ isSaving ? 'Saving Changes...' : 'Save Configuration' }}</span>
			</button>
		</div>

		<!-- Two Column Sections -->
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<!-- Section 1: School Profile -->
			<div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs space-y-4">
				<h2 class="text-sm font-bold text-slate-900 flex items-center gap-2 border-b pb-2">
					<Building2 class="w-4 h-4 text-indigo-600" /> School Identity
				</h2>

				<div>
					<label class="block text-xs font-semibold text-slate-700 mb-1">School Name *</label>
					<input
						v-model="profile.name"
						type="text"
						class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
					/>
				</div>

				<div>
					<label class="block text-xs font-semibold text-slate-700 mb-1">Legal / Registered Entity Name</label>
					<input
						v-model="profile.legalName"
						type="text"
						class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
					/>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">School Code *</label>
						<input
							v-model="profile.code"
							type="text"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300 uppercase"
						/>
					</div>
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Timezone *</label>
						<select
							v-model="profile.timezone"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						>
							<option value="Asia/Karachi">Asia/Karachi</option>
							<option value="UTC">UTC</option>
							<option value="Asia/Dubai">Asia/Dubai</option>
							<option value="Asia/Riyadh">Asia/Riyadh</option>
						</select>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Contact Phone</label>
						<input
							v-model="profile.phone"
							type="text"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Official Email</label>
						<input
							v-model="profile.email"
							type="email"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>
				</div>

				<div>
					<label class="block text-xs font-semibold text-slate-700 mb-1">Website</label>
					<input
						v-model="profile.website"
						type="text"
						class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
					/>
				</div>
			</div>

			<!-- Section 2: Billing & Operational Settings -->
			<div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs space-y-4">
				<h2 class="text-sm font-bold text-slate-900 flex items-center gap-2 border-b pb-2">
					<Coins class="w-4 h-4 text-emerald-600" /> Operational & Billing Rules
				</h2>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Invoice Prefix *</label>
						<input
							v-model="settings.invoicePrefix"
							type="text"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300 uppercase"
						/>
					</div>
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Receipt Prefix *</label>
						<input
							v-model="settings.receiptPrefix"
							type="text"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300 uppercase"
						/>
					</div>
				</div>

				<div>
					<label class="block text-xs font-semibold text-slate-700 mb-1">Admission No Format</label>
					<input
						v-model="settings.admissionNoFormat"
						type="text"
						placeholder="{YYYY}-{SEQ:5}"
						class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300 font-mono text-xs"
					/>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Default Due Day (1-28)</label>
						<input
							v-model.number="settings.defaultDueDay"
							type="number"
							min="1"
							max="28"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Late Grace Days</label>
						<input
							v-model.number="settings.lateFeeGraceDays"
							type="number"
							min="0"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Rounding Unit (Paisa/Minor)</label>
						<input
							v-model.number="settings.roundingUnitMinor"
							type="number"
							min="1"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Rounding Mode</label>
						<select
							v-model="settings.roundingMode"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						>
							<option value="half_up">Half-Up (Standard)</option>
							<option value="half_even">Half-Even (Banker's)</option>
							<option value="down">Truncate / Down</option>
							<option value="up">Ceiling / Up</option>
						</select>
					</div>
				</div>

				<div>
					<label class="block text-xs font-semibold text-slate-700 mb-1">Discount Combination Policy</label>
					<select
						v-model="settings.discountCombination"
						class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
					>
						<option value="exclusive_overrides">Exclusive Overrides Stackables</option>
						<option value="best_for_student">Automatic Best for Student</option>
					</select>
				</div>

				<div class="pt-2">
					<label class="flex items-center gap-2 cursor-pointer">
						<input
							v-model="settings.requireSecondApprover"
							type="checkbox"
							class="rounded text-indigo-600 focus:ring-indigo-500 h-4 w-4 border-slate-300"
						/>
						<span class="text-xs font-semibold text-slate-800">Require Second Approver for Payment Reversals & Voids</span>
					</label>
				</div>
			</div>
		</div>

		<!-- Developer & System Architecture Showcase -->
		<div class="bg-gradient-to-r from-slate-900 via-indigo-950 to-slate-900 rounded-2xl p-6 text-white shadow-md border border-slate-800">
			<div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
				<div class="flex items-center gap-4">
					<div class="p-2 bg-slate-800/80 rounded-2xl border border-slate-700 shadow-inner">
						<DeveloperLogo sizeClass="w-12 h-12" />
					</div>
					<div>
						<div class="flex items-center gap-2">
							<h3 class="text-base font-bold text-white">Hassaan Ali</h3>
							<span class="text-[11px] font-mono font-medium px-2 py-0.5 rounded-full bg-indigo-500/20 text-indigo-300 border border-indigo-400/30">@zfhassaan</span>
						</div>
						<p class="text-xs text-indigo-300 font-medium">Software Architect & Full Stack Engineer</p>
						<p class="text-xs text-slate-400 mt-1 max-w-xl leading-relaxed">
							Engineered with Go, Wails v2, SQLite WAL mode, and Vue 3. Features zero floating-point arithmetic, strict sequential audit trails, and modular hexagonal clean architecture.
						</p>
					</div>
				</div>

				<div class="flex items-center gap-2.5 shrink-0 self-end md:self-center">
					<button
						type="button"
						@click="openExternalUrl('https://github.com/zfhassaan')"
						class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-white/10 hover:bg-white/20 text-white flex items-center gap-2 transition border border-white/10"
					>
						<Github class="w-4 h-4" />
						<span>GitHub</span>
						<ExternalLink class="w-3 h-3 text-slate-400" />
					</button>
					<button
						type="button"
						@click="openExternalUrl('https://zfhassaan.github.io')"
						class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-indigo-600 hover:bg-indigo-500 text-white flex items-center gap-2 transition shadow-xs"
					>
						<Globe class="w-4 h-4" />
						<span>Portfolio</span>
						<ExternalLink class="w-3 h-3 text-indigo-200" />
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import { useAuthStore } from '@/stores/auth'
import DeveloperLogo from '@/components/DeveloperLogo.vue'
import { openExternalUrl } from '@/utils/browser'
import { Building2, Coins, Save, Github, Globe, ExternalLink } from 'lucide-vue-next'

const services = getServices()
const toast = useToast()
const authStore = useAuthStore()
const isSaving = ref(false)

const profile = reactive({
	name: '',
	legalName: '',
	code: '',
	timezone: 'Asia/Karachi',
	phone: '',
	email: '',
	website: '',
	locale: 'en',
	address: {
		line1: '',
		city: '',
		country: ''
	}
})

const settings = reactive({
	schoolId: '',
	admissionNoFormat: '{YYYY}-{SEQ:5}',
	invoicePrefix: 'INV',
	receiptPrefix: 'RCP',
	roundingUnitMinor: 1,
	roundingMode: 'half_up',
	defaultDueDay: 10,
	lateFeeGraceDays: 0,
	discountCombination: 'exclusive_overrides',
	requireSecondApprover: false,
	dateFormat: 'dd/MM/yyyy'
})

onMounted(async () => {
	try {
		const [pRes, sRes] = await Promise.all([
			services.school.getSchool(),
			services.school.getSettings()
		])
		if (pRes.success && pRes.data) {
			const d = pRes.data as any
			profile.name = d.name || d.Name || ''
			profile.legalName = d.legalName || d.LegalName || ''
			profile.code = d.code || d.Code || ''
			profile.timezone = d.timezone || d.Timezone || 'Asia/Karachi'
			profile.phone = d.phone || d.Phone || ''
			profile.email = d.email || d.Email || ''
			profile.website = d.website || d.Website || ''
			profile.locale = d.locale || d.Locale || 'en'
			const addr = d.address || d.Address || {}
			profile.address = {
				line1: addr.line1 || addr.Line1 || '',
				city: addr.city || addr.City || '',
				country: addr.country || addr.Country || ''
			}
		}
		if (sRes.success && sRes.data) {
			Object.assign(settings, sRes.data)
		}
	} catch (e: any) {
		toast.error(e.message || 'Failed to load school settings')
	}
})

async function saveAll() {
	isSaving.value = true
	try {
		const [pRes, sRes] = await Promise.all([
			services.school.updateSchoolProfile(profile),
			services.school.updateSettings(settings)
		])
		if (pRes.success && sRes.success) {
			authStore.updateSchoolDetails({
				name: profile.name,
				legalName: profile.legalName,
				code: profile.code,
				timezone: profile.timezone,
				phone: profile.phone,
				email: profile.email,
				website: profile.website,
				address: profile.address
			})
			await authStore.refreshSchool()
			toast.success('School profile and operational settings saved successfully!')
		} else {
			toast.error(pRes.error?.message || sRes.error?.message || 'Failed to save configuration')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred during save')
	} finally {
		isSaving.value = false
	}
}
</script>

