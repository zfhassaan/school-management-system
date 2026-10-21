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

				<!-- School Logo Uploader -->
				<div class="p-4 bg-slate-50 border border-slate-200 rounded-xl flex items-center justify-between gap-4">
					<div class="flex items-center gap-4">
						<div class="w-16 h-16 rounded-xl border border-slate-300 bg-white flex items-center justify-center font-bold text-lg text-slate-400 overflow-hidden shrink-0 shadow-2xs">
							<img v-if="logoPreview || authStore.schoolLogo" :src="logoPreview || authStore.schoolLogo!" alt="School Logo Preview" class="w-full h-full object-contain p-1" />
							<div v-else class="text-center p-1">
								<ImageIcon class="w-6 h-6 mx-auto text-slate-300" />
								<span class="text-[9px] uppercase font-bold text-slate-400 block mt-0.5">No Logo</span>
							</div>
						</div>
						<div>
							<h3 class="text-xs font-bold text-slate-800">Official School Crest / Logo</h3>
							<p class="text-[11px] text-slate-500 mt-0.5 max-w-sm">
								Used on official admission forms, printable bills, collection receipts, and the application navigation.
							</p>
						</div>
					</div>

					<div class="flex items-center gap-2 shrink-0">
						<label class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium cursor-pointer transition shadow-2xs flex items-center gap-1.5">
							<Upload class="w-3.5 h-3.5" />
							<span>{{ logoPreview || authStore.schoolLogo ? 'Change Logo' : 'Upload Logo' }}</span>
							<input type="file" accept="image/png,image/jpeg,image/svg+xml,image/webp" class="hidden" @change="onLogoSelected" />
						</label>
						<button
							v-if="logoPreview || authStore.schoolLogo"
							type="button"
							@click="removeLogo"
							class="px-2.5 py-1.5 border border-slate-300 hover:bg-slate-100 text-rose-600 rounded-lg text-xs font-medium transition cursor-pointer"
						>
							Remove
						</button>
					</div>
				</div>

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

		<!-- Section 3: Database Security, Encryption & Backup Vault (Tight AES-256-GCM) -->
		<div class="bg-white p-6 rounded-2xl border border-slate-200 shadow-xs space-y-5">
			<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-b border-slate-200 pb-4">
				<div class="flex items-center gap-3">
					<div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-700 flex items-center justify-center border border-emerald-200 shrink-0">
						<Lock class="w-5 h-5" />
					</div>
					<div>
						<div class="flex items-center gap-2">
							<h2 class="text-sm font-bold text-slate-900">Database Encryption &amp; Backup Vault</h2>
							<span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800 border border-emerald-300 flex items-center gap-1">
								<ShieldCheck class="w-3 h-3 text-emerald-600" /> AES-256-GCM Encrypted
							</span>
						</div>
						<p class="text-xs text-slate-500 mt-0.5">
							The database on disk is tightly encrypted at rest. Unencrypted data can only be downloaded via authorized application backup.
						</p>
					</div>
				</div>

				<button
					type="button"
					@click="loadDbStatus"
					class="p-1.5 text-slate-500 hover:text-slate-800 hover:bg-slate-100 rounded-lg text-xs transition cursor-pointer self-start sm:self-center"
					title="Refresh Database Vault Status"
				>
					<RefreshCw class="w-4 h-4" />
				</button>
			</div>

			<!-- Vault Metadata Grid -->
			<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 text-xs">
				<div class="p-3 bg-slate-50 border border-slate-200 rounded-xl">
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block">Encryption Standard</span>
					<span class="font-bold text-slate-900 text-xs mt-0.5 block">AES-256-GCM</span>
					<span class="text-[10px] text-slate-500">Authenticated Galois/Counter</span>
				</div>
				<div class="p-3 bg-slate-50 border border-slate-200 rounded-xl">
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block">Key Derivation</span>
					<span class="font-bold text-slate-900 text-xs mt-0.5 block">PBKDF2-HMAC-SHA256</span>
					<span class="text-[10px] text-slate-500">100,000 Key Stretching Rounds</span>
				</div>
				<div class="p-3 bg-slate-50 border border-slate-200 rounded-xl">
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block">Encrypted Vault Size</span>
					<span class="font-bold text-slate-900 text-xs mt-0.5 block">{{ formatBytes(dbStatus?.vaultSizeBytes) }}</span>
					<span class="text-[10px] text-emerald-700 font-medium">Secured on Local Disk</span>
				</div>
				<div class="p-3 bg-slate-50 border border-slate-200 rounded-xl">
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block">Installed Storage Location</span>
					<span class="font-mono text-[10px] font-medium text-slate-700 truncate block mt-0.5" :title="dbStatus?.vaultPath">
						{{ dbStatus?.vaultPath ? dbStatus.vaultPath.split(/[\\/]/).slice(-3).join('/') : 'Installed Folder /data/school.db.enc' }}
					</span>
					<span class="text-[10px] text-slate-400">Same Installed Directory</span>
				</div>
			</div>

			<!-- Action Buttons for Backup and Decrypted Export -->
			<div class="pt-2 flex flex-wrap items-center gap-3">
				<!-- Export Decrypted Database -->
				<button
					type="button"
					@click="showExportDecryptedModal = true"
					class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl text-xs font-semibold shadow-xs flex items-center gap-2 transition cursor-pointer"
				>
					<Download class="w-4 h-4" />
					<span>Download Decrypted Database (.sqlite)</span>
				</button>

				<!-- Export Encrypted Vault File -->
				<button
					type="button"
					@click="exportEncryptedVault"
					:disabled="isExportingEncrypted"
					class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-800 rounded-xl text-xs font-semibold border border-slate-300 flex items-center gap-2 transition cursor-pointer"
				>
					<Lock class="w-4 h-4 text-slate-600" />
					<span>{{ isExportingEncrypted ? 'Exporting Vault...' : 'Download Encrypted Vault (.enc)' }}</span>
				</button>

				<!-- Restore Database -->
				<button
					type="button"
					@click="showRestoreModal = true"
					class="px-4 py-2 bg-amber-50 hover:bg-amber-100 text-amber-800 rounded-xl text-xs font-semibold border border-amber-300 flex items-center gap-2 transition cursor-pointer"
				>
					<Database class="w-4 h-4 text-amber-600" />
					<span>Restore Database from Backup</span>
				</button>
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

		<!-- Modal: Download Decrypted Database -->
		<div v-if="showExportDecryptedModal" class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-150">
			<div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 border border-slate-200 space-y-4">
				<div class="flex items-start justify-between">
					<div class="flex items-center gap-3">
						<div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-700 flex items-center justify-center border border-indigo-200 shrink-0">
							<Download class="w-5 h-5" />
						</div>
						<div>
							<h3 class="text-sm font-bold text-slate-900">Download Decrypted Database</h3>
							<p class="text-xs text-slate-500">Export clean, decrypted SQLite copy for backup.</p>
						</div>
					</div>
					<button @click="showExportDecryptedModal = false" class="text-slate-400 hover:text-slate-600 p-1">
						<X class="w-4 h-4" />
					</button>
				</div>

				<div class="p-3 bg-amber-50 border border-amber-200 rounded-xl flex items-start gap-2.5 text-xs text-amber-800">
					<AlertTriangle class="w-4 h-4 text-amber-600 shrink-0 mt-0.5" />
					<p>
						The exported file will be an <strong>unencrypted SQLite database (.sqlite)</strong>. Store it on a secure encrypted drive or protected offline storage.
					</p>
				</div>

				<form @submit.prevent="handleExportDecrypted" class="space-y-3 text-xs">
					<div>
						<label class="block font-semibold text-slate-700 mb-1">Confirm Administrator Password</label>
						<input
							v-model="adminExportPassword"
							type="password"
							required
							placeholder="Enter your login password to authorize export"
							class="w-full px-3 py-2 border rounded-lg text-xs bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>

					<div class="flex items-center justify-end gap-2.5 pt-2">
						<button
							type="button"
							@click="showExportDecryptedModal = false"
							class="px-3.5 py-2 border border-slate-300 text-slate-700 font-medium rounded-lg hover:bg-slate-50 transition cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="isExportingDecrypted"
							class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-semibold rounded-lg shadow-xs flex items-center gap-1.5 transition cursor-pointer"
						>
							<span v-if="isExportingDecrypted" class="animate-spin mr-1">⏳</span>
							<span>{{ isExportingDecrypted ? 'Decrypting & Exporting...' : 'Authorize & Download' }}</span>
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- Modal: Restore Database -->
		<div v-if="showRestoreModal" class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-150">
			<div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 border border-slate-200 space-y-4">
				<div class="flex items-start justify-between">
					<div class="flex items-center gap-3">
						<div class="w-10 h-10 rounded-xl bg-rose-50 text-rose-700 flex items-center justify-center border border-rose-200 shrink-0">
							<Database class="w-5 h-5" />
						</div>
						<div>
							<h3 class="text-sm font-bold text-slate-900">Restore Database from Backup</h3>
							<p class="text-xs text-slate-500">Restore system data from a .sqlite or .enc file.</p>
						</div>
					</div>
					<button @click="showRestoreModal = false" class="text-slate-400 hover:text-slate-600 p-1">
						<X class="w-4 h-4" />
					</button>
				</div>

				<div class="p-3 bg-rose-50 border border-rose-200 rounded-xl flex items-start gap-2.5 text-xs text-rose-800">
					<AlertTriangle class="w-4 h-4 text-rose-600 shrink-0 mt-0.5" />
					<p>
						<strong>Warning:</strong> Restoring will overwrite existing operational data. The system will automatically validate schema and re-encrypt the restored data.
					</p>
				</div>

				<form @submit.prevent="handleRestore" class="space-y-3 text-xs">
					<div>
						<label class="block font-semibold text-slate-700 mb-1">Backup File Path (.sqlite or .enc)</label>
						<input
							v-model="restoreBackupPath"
							type="text"
							required
							placeholder="C:\path\to\School_Database_Backup.sqlite"
							class="w-full px-3 py-2 border rounded-lg text-xs font-mono bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>

					<div>
						<label class="block font-semibold text-slate-700 mb-1">Confirm Administrator Password</label>
						<input
							v-model="adminRestorePassword"
							type="password"
							required
							placeholder="Enter your login password to authorize restore"
							class="w-full px-3 py-2 border rounded-lg text-xs bg-slate-50 focus:bg-white focus:ring-2 focus:ring-indigo-500 outline-hidden border-slate-300"
						/>
					</div>

					<div class="flex items-center justify-end gap-2.5 pt-2">
						<button
							type="button"
							@click="showRestoreModal = false"
							class="px-3.5 py-2 border border-slate-300 text-slate-700 font-medium rounded-lg hover:bg-slate-50 transition cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="isRestoring || !restoreBackupPath"
							class="px-4 py-2 bg-rose-600 hover:bg-rose-700 disabled:opacity-50 text-white font-semibold rounded-lg shadow-xs flex items-center gap-1.5 transition cursor-pointer"
						>
							<span v-if="isRestoring" class="animate-spin mr-1">⏳</span>
							<span>{{ isRestoring ? 'Restoring & Re-encrypting...' : 'Confirm & Restore' }}</span>
						</button>
					</div>
				</form>
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
import {
	Building2,
	Coins,
	Save,
	Github,
	Globe,
	ExternalLink,
	Upload,
	Image as ImageIcon,
	ShieldCheck,
	Lock,
	Database,
	Download,
	RefreshCw,
	AlertTriangle,
	X
} from 'lucide-vue-next'
import type { DatabaseStatusDTO } from '@/services/contracts'

const services = getServices()
const toast = useToast()
const authStore = useAuthStore()
const isSaving = ref(false)
const logoPreview = ref<string | null>(null)
const logoChanged = ref(false)

// Database Vault & Backup State
const dbStatus = ref<DatabaseStatusDTO | null>(null)
const showExportDecryptedModal = ref(false)
const showRestoreModal = ref(false)
const adminExportPassword = ref('')
const adminRestorePassword = ref('')
const restoreBackupPath = ref('')
const isExportingDecrypted = ref(false)
const isExportingEncrypted = ref(false)
const isRestoring = ref(false)

function onLogoSelected(e: Event) {
	const file = (e.target as HTMLInputElement).files?.[0]
	if (!file) return
	if (file.size > 2 * 1024 * 1024) {
		toast.error('Logo image must be under 2MB')
		return
	}
	const reader = new FileReader()
	reader.onload = (event) => {
		logoPreview.value = event.target?.result as string
		logoChanged.value = true
	}
	reader.readAsDataURL(file)
}

function removeLogo() {
	logoPreview.value = null
	logoChanged.value = true
	authStore.setSchoolLogo(null)
}

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
		await loadDbStatus()
	} catch (e: any) {
		toast.error(e.message || 'Failed to load school settings')
	}
})

function formatBytes(bytes?: number) {
	if (!bytes || bytes === 0) return '0 B'
	const k = 1024
	const sizes = ['B', 'KB', 'MB', 'GB']
	const i = Math.floor(Math.log(bytes) / Math.log(k))
	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`
}

async function loadDbStatus() {
	try {
		const res = await services.backup.getDatabaseStatus()
		if (res.success && res.data) {
			dbStatus.value = res.data
		}
	} catch (e) {
		console.error('Failed to load database status', e)
	}
}

async function handleExportDecrypted() {
	if (!adminExportPassword.value) {
		toast.error('Please enter administrator password to authorize export')
		return
	}
	isExportingDecrypted.value = true
	try {
		const res = await services.backup.exportDecryptedDatabase({
			adminPassword: adminExportPassword.value
		})
		if (res.success && res.data) {
			showExportDecryptedModal.value = false
			adminExportPassword.value = ''
			toast.success(`Decrypted SQLite database exported to: ${res.data}`)
			await loadDbStatus()
		} else {
			toast.error(res.error?.message || 'Authorization failed or export error')
		}
	} catch (e: any) {
		toast.error(e.message || 'Export failed')
	} finally {
		isExportingDecrypted.value = false
	}
}

async function exportEncryptedVault() {
	isExportingEncrypted.value = true
	try {
		const res = await services.backup.exportEncryptedDatabase({})
		if (res.success && res.data) {
			toast.success(`Encrypted database vault (.enc) saved to: ${res.data}`)
			await loadDbStatus()
		} else {
			toast.error(res.error?.message || 'Failed to export encrypted vault')
		}
	} catch (e: any) {
		toast.error(e.message || 'Export error')
	} finally {
		isExportingEncrypted.value = false
	}
}

async function handleRestore() {
	if (!restoreBackupPath.value) {
		toast.error('Please specify a valid backup file path')
		return
	}
	if (!adminRestorePassword.value) {
		toast.error('Please enter administrator password to confirm restore')
		return
	}
	isRestoring.value = true
	try {
		const res = await services.backup.restoreDatabase({
			adminPassword: adminRestorePassword.value,
			backupFilePath: restoreBackupPath.value.trim()
		})
		if (res.success) {
			showRestoreModal.value = false
			adminRestorePassword.value = ''
			restoreBackupPath.value = ''
			toast.success('Database restored and re-encrypted successfully! Reloading system...')
			await loadDbStatus()
			await authStore.refreshSchool()
		} else {
			toast.error(res.error?.message || 'Database restore failed')
		}
	} catch (e: any) {
		toast.error(e.message || 'Restore error')
	} finally {
		isRestoring.value = false
	}
}

async function saveAll() {
	isSaving.value = true
	try {
		const [pRes, sRes] = await Promise.all([
			services.school.updateSchoolProfile(profile),
			services.school.updateSettings(settings)
		])
		if (pRes.success && sRes.success) {
			if (logoChanged.value) {
				authStore.setSchoolLogo(logoPreview.value)
				logoChanged.value = false
			}
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
			toast.success('School profile, logo, and operational settings saved successfully!')
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

