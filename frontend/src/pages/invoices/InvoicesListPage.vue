<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Student Invoices & Billing</h1>
				<p class="text-xs text-slate-500 mt-0.5">Generate batch tuition invoices, monitor fee schedules, and track payment balances.</p>
			</div>
			<button
				@click="openBatchModal"
				class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
			>
				<Plus class="w-4 h-4" />
				<span>Generate Batch Invoices</span>
			</button>
		</div>

		<!-- Filters & Search Bar -->
		<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex flex-wrap gap-3 items-center justify-between">
			<div class="flex items-center gap-3 flex-wrap">
				<div>
					<label class="text-[11px] font-semibold text-slate-600 mr-2">Year:</label>
					<select
						v-model="yearFilter"
						@change="loadInvoices"
						class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					>
						<option value="">All Years</option>
						<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }}</option>
					</select>
				</div>

				<div>
					<label class="text-[11px] font-semibold text-slate-600 mr-2">Status:</label>
					<select
						v-model="statusFilter"
						@change="loadInvoices"
						class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					>
						<option value="">All Statuses</option>
						<option value="draft">Draft</option>
						<option value="issued">Issued</option>
						<option value="partially_paid">Partially Paid</option>
						<option value="paid">Paid</option>
						<option value="overdue">Overdue</option>
						<option value="void">Void</option>
					</select>
				</div>
			</div>

			<div class="flex items-center gap-2">
				<button
					@click="loadInvoices"
					class="p-2 border border-slate-200 hover:bg-slate-50 rounded-lg text-slate-600 transition cursor-pointer"
					title="Refresh List"
				>
					<RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
				</button>
			</div>
		</div>

		<!-- Invoices Table -->
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">Invoice #</th>
						<th class="py-3 px-4">Student</th>
						<th class="py-3 px-4">Period</th>
						<th class="py-3 px-4">Issue / Due Date</th>
						<th class="py-3 px-4 text-right">Total Amount</th>
						<th class="py-3 px-4 text-right">Balance Due</th>
						<th class="py-3 px-4 text-center">Status</th>
						<th class="py-3 px-4 text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-slate-700">
					<tr v-for="inv in invoices" :key="inv.ID" class="hover:bg-slate-50/50">
						<td class="py-3 px-4 font-mono font-bold text-slate-900">
							<router-link :to="`/invoices/${inv.ID}`" class="hover:text-indigo-600 hover:underline">
								{{ inv.InvoiceNo }}
							</router-link>
						</td>
						<td class="py-3 px-4">
							<div class="font-semibold text-slate-900">{{ inv.StudentFirstName }} {{ inv.StudentLastName }}</div>
							<div class="font-mono text-[10px] text-slate-400">Adm: {{ inv.AdmissionNo }}</div>
						</td>
						<td class="py-3 px-4 font-medium text-slate-600">{{ inv.BillingPeriod }}</td>
						<td class="py-3 px-4 text-slate-600">
							<div>{{ formatCivilDate(inv.IssueDate) }}</div>
							<div class="text-[10px] text-slate-400">Due: {{ formatCivilDate(inv.DueDate) }}</div>
						</td>
						<td class="py-3 px-4 text-right font-medium text-slate-800">
							{{ formatMinor(inv.TotalAmountMinor) }}
							<div v-if="inv.DiscountAmountMinor > 0" class="text-[10px] text-emerald-600">
								-{{ formatMinor(inv.DiscountAmountMinor) }} disc
							</div>
						</td>
						<td class="py-3 px-4 text-right font-bold" :class="inv.BalanceMinor > 0 && inv.Status !== 'void' ? 'text-rose-600' : 'text-slate-700'">
							{{ formatMinor(inv.BalanceMinor) }}
						</td>
						<td class="py-3 px-4 text-center">
							<span
								class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
								:class="getStatusBadgeClass(inv.Status)"
							>
								{{ inv.Status.replace('_', ' ') }}
							</span>
						</td>
						<td class="py-3 px-4 text-right space-x-2">
							<router-link
								:to="`/invoices/${inv.ID}`"
								class="text-indigo-600 hover:text-indigo-800 font-medium text-xs"
							>
								View
							</router-link>
							<button
								v-if="inv.Status !== 'void'"
								@click="openVoidModal(inv.ID)"
								class="text-rose-500 hover:text-rose-700 font-medium text-xs cursor-pointer"
							>
								Void
							</button>
						</td>
					</tr>
					<tr v-if="invoices.length === 0">
						<td colspan="8" class="py-10 text-center text-slate-400 italic">
							No invoices match the specified criteria. Click "Generate Batch Invoices" to issue bills for enrolled students.
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- MODAL: Batch Generate Invoices -->
		<div v-if="showBatchModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-md w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Generate Invoices (Batch by Class)</h3>
					<button @click="showBatchModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitBatchGenerate" class="p-4 space-y-3">
					<div v-if="batchError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ batchError }}
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Academic Year *</label>
						<select
							v-model="batchForm.academicYearId"
							required
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
						</select>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Target Class *</label>
						<select
							v-model="batchForm.classId"
							required
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="" disabled>Select Class</option>
							<option v-for="c in classes" :key="c.ID" :value="c.ID">{{ c.Name }} ({{ c.Code }})</option>
						</select>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Billing Period Code *</label>
						<input
							v-model="batchForm.billingPeriod"
							type="text"
							required
							placeholder="e.g. 2026-09 or 2026-TERM1"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
						<p class="text-[10px] text-slate-400 mt-0.5">Enforces idempotency: students will not be double-billed for the same period.</p>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Issue Date *</label>
							<input
								v-model="batchForm.issueDate"
								type="date"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Due Date *</label>
							<input
								v-model="batchForm.dueDate"
								type="date"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
					</div>

					<div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showBatchModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="generatingBatch"
							class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							{{ generatingBatch ? 'Generating Invoices...' : 'Generate Batch' }}
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- MODAL: Void Invoice -->
		<div v-if="showVoidModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-sm w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Void Invoice</h3>
					<button @click="showVoidModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitVoidInvoice" class="p-4 space-y-3">
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason for Voiding *</label>
						<textarea
							v-model="voidReason"
							required
							rows="3"
							placeholder="State justification (e.g. incorrect class structure applied, duplicate bill)..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>
					<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
						<button
							type="button"
							@click="showVoidModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-1.5 bg-rose-600 hover:bg-rose-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							Confirm Void
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import type { InvoiceDTO, AcademicYearDTO, ClassDTO } from '@/types/common'
import { formatMinor, formatCivilDate } from '@/utils/format'
import { Plus, RefreshCw } from 'lucide-vue-next'

const services = getServices()

const invoices = ref<InvoiceDTO[]>([])
const academicYears = ref<AcademicYearDTO[]>([])
const classes = ref<ClassDTO[]>([])
const loading = ref(false)

const yearFilter = ref('')
const statusFilter = ref('')

// Batch Modal
const showBatchModal = ref(false)
const generatingBatch = ref(false)
const batchError = ref('')
const batchForm = ref({
	academicYearId: '',
	classId: '',
	billingPeriod: '',
	issueDate: new Date().toISOString().split('T')[0],
	dueDate: ''
})

// Void Modal
const showVoidModal = ref(false)
const voidingInvoiceId = ref('')
const voidReason = ref('')

onMounted(async () => {
	await Promise.all([loadAcademicYears(), loadClasses()])
	await loadInvoices()
})

async function loadAcademicYears() {
	try {
		const res = await services.academics.listAcademicYears()
		if (res.success && res.data) {
			academicYears.value = res.data
			const active = res.data.find(y => y.Status === 'active')
			if (active) {
				yearFilter.value = active.ID
				batchForm.value.academicYearId = active.ID
			} else if (res.data.length > 0) {
				yearFilter.value = res.data[0].ID
				batchForm.value.academicYearId = res.data[0].ID
			}
		}
	} catch (e: any) {
		console.error('Failed to load academic years', e)
	}
}

async function loadClasses() {
	try {
		const res = await services.academics.listClasses()
		if (res.success && res.data) {
			classes.value = res.data
			if (classes.value.length > 0) {
				batchForm.value.classId = classes.value[0].ID
			}
		}
	} catch (e: any) {
		console.error('Failed to load classes', e)
	}
}

async function loadInvoices() {
	loading.value = true
	try {
		const res = await services.invoices.listInvoices({
			academicYearId: yearFilter.value || undefined,
			status: statusFilter.value || undefined,
			limit: 100
		})
		if (res.success && res.data) {
			invoices.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load invoices', e)
	} finally {
		loading.value = false
	}
}

function openBatchModal() {
	batchError.value = ''
	const today = new Date()
	const dueDate = new Date()
	dueDate.setDate(today.getDate() + 15)

	const month = String(today.getMonth() + 1).padStart(2, '0')
	batchForm.value.billingPeriod = `${today.getFullYear()}-${month}`
	batchForm.value.issueDate = today.toISOString().split('T')[0]
	batchForm.value.dueDate = dueDate.toISOString().split('T')[0]
	showBatchModal.value = true
}

async function submitBatchGenerate() {
	generatingBatch.value = true
	batchError.value = ''
	try {
		const res = await services.invoices.generateInvoicesBatch({
			academicYearId: batchForm.value.academicYearId,
			classId: batchForm.value.classId,
			billingPeriod: batchForm.value.billingPeriod,
			issueDate: batchForm.value.issueDate,
			dueDate: batchForm.value.dueDate
		})
		if (!res.success) {
			batchError.value = res.error?.message || 'Failed to generate batch invoices'
			return
		}
		const count = res.data ? res.data.length : 0
		alert(`Successfully generated ${count} invoices for the selected class period.`)
		showBatchModal.value = false
		await loadInvoices()
	} catch (e: any) {
		batchError.value = e.message || 'An error occurred'
	} finally {
		generatingBatch.value = false
	}
}

function openVoidModal(invoiceId: string) {
	voidingInvoiceId.value = invoiceId
	voidReason.value = ''
	showVoidModal.value = true
}

async function submitVoidInvoice() {
	if (!voidingInvoiceId.value) return
	try {
		const res = await services.invoices.voidInvoice({
			invoiceId: voidingInvoiceId.value,
			reason: voidReason.value
		})
		if (res.success) {
			showVoidModal.value = false
			await loadInvoices()
		} else {
			alert(res.error?.message || 'Failed to void invoice')
		}
	} catch (e: any) {
		alert(e.message || 'An error occurred')
	}
}

function getStatusBadgeClass(status: string) {
	switch (status) {
		case 'draft':
			return 'bg-amber-50 text-amber-700 border border-amber-200'
		case 'issued':
			return 'bg-blue-50 text-blue-700 border border-blue-200'
		case 'partially_paid':
			return 'bg-yellow-50 text-yellow-700 border border-yellow-200'
		case 'paid':
			return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
		case 'overdue':
			return 'bg-rose-50 text-rose-700 border border-rose-200'
		case 'void':
			return 'bg-slate-100 text-slate-500 border border-slate-200'
		default:
			return 'bg-slate-100 text-slate-600'
	}
}
</script>

