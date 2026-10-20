<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Financial Adjustments & Credit Ledgers</h1>
				<p class="text-xs text-slate-500 mt-0.5">Manage fee waivers, debt write-offs, and inspect student credit advance ledgers.</p>
			</div>
			<div class="flex gap-2">
				<router-link
					to="/payments"
					class="py-2 px-3.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Receipt class="w-4 h-4 text-slate-500" />
					<span>Back to Payments</span>
				</router-link>
				<button
					v-if="currentTab === 'adjustments'"
					@click="openRequestAdjustmentModal"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>Request Adjustment</span>
				</button>
			</div>
		</div>

		<!-- Tabs -->
		<div class="flex border-b border-slate-200 space-x-6 text-xs font-medium">
			<button
				v-for="tab in ['adjustments', 'credits']"
				:key="tab"
				@click="currentTab = tab"
				class="pb-3 capitalize cursor-pointer transition border-b-2 flex items-center gap-2"
				:class="currentTab === tab ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'"
			>
				<span>{{ tab === 'adjustments' ? 'Adjustments & Waivers' : 'Student Credit Balances' }}</span>
				<span
					v-if="tab === 'adjustments' && pendingCount > 0"
					class="px-1.5 py-0.2 rounded-full text-[10px] bg-rose-100 text-rose-700 font-bold"
				>
					{{ pendingCount }}
				</span>
			</button>
		</div>

		<!-- TAB 1: Adjustments -->
		<div v-if="currentTab === 'adjustments'" class="space-y-4">
			<!-- Status Filter Bar -->
			<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex items-center justify-between">
				<div class="flex items-center gap-3">
					<label class="text-[11px] font-semibold text-slate-600">Filter Status:</label>
					<select
						v-model="statusFilter"
						@change="loadAdjustments"
						class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					>
						<option value="">All Statuses</option>
						<option value="pending">Pending Approval</option>
						<option value="approved">Approved</option>
						<option value="rejected">Rejected</option>
					</select>
				</div>
				<button
					@click="loadAdjustments"
					class="p-2 border border-slate-200 hover:bg-slate-50 rounded-lg text-slate-600 transition cursor-pointer"
					title="Refresh"
				>
					<RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
				</button>
			</div>

			<!-- Adjustments Table -->
			<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
				<table class="w-full text-left border-collapse text-xs">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
							<th class="py-3 px-4">Student</th>
							<th class="py-3 px-4">Invoice #</th>
							<th class="py-3 px-4">Type</th>
							<th class="py-3 px-4 text-right">Adjustment Amount</th>
							<th class="py-3 px-4">Reason / Notes</th>
							<th class="py-3 px-4">Requested By</th>
							<th class="py-3 px-4 text-center">Status</th>
							<th class="py-3 px-4 text-right">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 text-slate-700">
						<tr v-for="adj in adjustments" :key="adj.ID" class="hover:bg-slate-50/50">
							<td class="py-3 px-4">
								<div class="font-semibold text-slate-900">{{ adj.StudentFirstName }} {{ adj.StudentLastName }}</div>
								<div class="font-mono text-[10px] text-slate-400">Adm: {{ adj.AdmissionNo }}</div>
							</td>
							<td class="py-3 px-4 font-mono font-medium text-slate-700">
								{{ adj.InvoiceNo || 'Account Balance' }}
							</td>
							<td class="py-3 px-4">
								<span
									class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
									:class="getAdjustmentBadgeClass(adj.AdjustmentType)"
								>
									{{ adj.AdjustmentType.replace('_', ' ') }}
								</span>
							</td>
							<td class="py-3 px-4 text-right font-bold text-slate-900">
								{{ formatMinor(adj.AmountMinor) }}
							</td>
							<td class="py-3 px-4">
								<div class="text-slate-800">{{ adj.Reason }}</div>
								<div v-if="adj.DecisionNotes" class="text-[10px] text-slate-400 mt-0.5">Decision: {{ adj.DecisionNotes }}</div>
							</td>
							<td class="py-3 px-4 text-slate-600">
								<div>{{ adj.RequestedByName || 'Staff' }}</div>
								<div class="text-[10px] text-slate-400">{{ formatCivilDate(adj.RequestedAt) }}</div>
							</td>
							<td class="py-3 px-4 text-center">
								<span
									class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
									:class="
										adj.Status === 'approved'
											? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
											: adj.Status === 'pending'
											? 'bg-amber-50 text-amber-700 border border-amber-200'
											: 'bg-rose-50 text-rose-700 border border-rose-200'
									"
								>
									{{ adj.Status }}
								</span>
							</td>
							<td class="py-3 px-4 text-right space-x-2">
								<template v-if="adj.Status === 'pending'">
									<button
										@click="decideAdjustment(adj.ID, true)"
										class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded font-medium text-xs cursor-pointer"
									>
										Approve
									</button>
									<button
										@click="openRejectModal(adj.ID)"
										class="px-2.5 py-1 bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200 rounded font-medium text-xs cursor-pointer"
									>
										Reject
									</button>
								</template>
								<span v-else class="text-slate-400 text-[11px]">—</span>
							</td>
						</tr>
						<tr v-if="adjustments.length === 0">
							<td colspan="8" class="py-8 text-center text-slate-400 italic">No financial adjustment requests found.</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- TAB 2: Student Credit Balances -->
		<div v-if="currentTab === 'credits'" class="grid grid-cols-12 gap-5">
			<!-- Left: Student Selector -->
			<div class="col-span-4 bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
				<h3 class="font-bold text-sm text-slate-900">Select Student</h3>
				<div>
					<StudentSearchSelect
						v-model="selectedCreditStudentId"
						:students="studentsList"
						label="Student"
						placeholder="Search student by name or ID..."
						@select="loadCreditBalance"
					/>
				</div>

				<div v-if="creditInfo" class="pt-4 border-t border-slate-100 space-y-2">
					<div class="text-[11px] font-bold uppercase tracking-wider text-slate-400">Available Advance Credit</div>
					<div class="text-2xl font-black text-indigo-700">
						{{ formatMinor(creditInfo.balanceMinor) }}
					</div>
					<p class="text-[11px] text-slate-500">
						This credit balance is automatically applied towards future tuition invoices or can be refunded upon withdrawal.
					</p>
				</div>
			</div>

			<!-- Right: Credit Balance Ledger -->
			<div class="col-span-8 bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
				<h3 class="font-bold text-sm text-slate-900">Credit Balance Ledger & History</h3>
				<div v-if="creditInfo" class="border border-slate-200 rounded-lg overflow-hidden">
					<table class="w-full text-left text-xs border-collapse">
						<thead>
							<tr class="bg-slate-50 border-b border-slate-200 text-[10px] font-bold uppercase tracking-wider text-slate-500">
								<th class="py-2.5 px-3">Date</th>
								<th class="py-2.5 px-3">Type</th>
								<th class="py-2.5 px-3 text-right">Amount</th>
								<th class="py-2.5 px-3 text-right">Balance After</th>
								<th class="py-2.5 px-3">Notes</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100 text-slate-700">
							<tr v-for="entry in creditInfo.ledger" :key="entry.ID" class="hover:bg-slate-50/50">
								<td class="py-2.5 px-3 text-slate-600">{{ formatCivilDate(entry.CreatedAt) }}</td>
								<td class="py-2.5 px-3">
									<span class="px-1.5 py-0.5 rounded text-[10px] font-semibold capitalize bg-slate-100 text-slate-700">
										{{ entry.ReferenceType.replace('_', ' ') }}
									</span>
								</td>
								<td
									class="py-2.5 px-3 text-right font-bold"
									:class="entry.AmountMinor >= 0 ? 'text-emerald-600' : 'text-rose-600'"
								>
									{{ entry.AmountMinor >= 0 ? '+' : '' }}{{ formatMinor(entry.AmountMinor) }}
								</td>
								<td class="py-2.5 px-3 text-right font-medium text-slate-900">
									{{ formatMinor(entry.BalanceAfterMinor) }}
								</td>
								<td class="py-2.5 px-3 text-slate-500 text-[11px]">{{ entry.Notes || '—' }}</td>
							</tr>
							<tr v-if="creditInfo.ledger.length === 0">
								<td colspan="5" class="py-6 text-center text-slate-400 italic">No credit transactions recorded.</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div v-else class="text-xs text-slate-400 italic text-center py-12">
					Select a student on the left to inspect their advance credit balance and ledger transactions.
				</div>
			</div>
		</div>

		<!-- MODAL: Request Financial Adjustment -->
		<div v-if="showAdjustmentModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-md w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Request Financial Adjustment</h3>
					<button @click="showAdjustmentModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitAdjustmentRequest" class="p-4 space-y-3">
					<div v-if="adjustmentError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ adjustmentError }}
					</div>

					<div>
						<StudentSearchSelect
							v-model="adjForm.studentId"
							:students="studentsList"
							label="Select Student"
							required
							placeholder="Search student by name or ID..."
							@select="loadStudentInvoices"
						/>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Target Invoice (Optional)</label>
						<select
							v-model="adjForm.invoiceId"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="">General Account Adjustment</option>
							<option v-for="inv in studentInvoices" :key="inv.ID" :value="inv.ID">
								{{ inv.InvoiceNo }} (Period: {{ inv.BillingPeriod }} • Balance: {{ formatMinor(inv.BalanceMinor) }})
							</option>
						</select>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Adjustment Type *</label>
							<select
								v-model="adjForm.adjustmentType"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							>
								<option value="waiver">Fee Waiver</option>
								<option value="write_off">Debt Write-Off</option>
								<option value="debit_correction">Debit Correction</option>
								<option value="credit_correction">Credit Correction</option>
							</select>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Amount (PKR) *</label>
							<input
								v-model.number="adjForm.amountDecimal"
								type="number"
								min="0.01"
								step="0.01"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 font-bold"
							/>
						</div>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason & Justification *</label>
						<textarea
							v-model="adjForm.reason"
							required
							rows="3"
							placeholder="State clear justification for administrative approval..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>

					<div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showAdjustmentModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="submittingAdj"
							class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							{{ submittingAdj ? 'Submitting...' : 'Submit Request' }}
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- MODAL: Rejection Notes -->
		<div v-if="showRejectModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-sm w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Reject Adjustment Request</h3>
					<button @click="showRejectModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitReject" class="p-4 space-y-3">
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason for Rejection *</label>
						<textarea
							v-model="rejectNotes"
							required
							rows="3"
							placeholder="State justification..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>
					<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
						<button
							type="button"
							@click="showRejectModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-1.5 bg-rose-600 hover:bg-rose-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							Confirm Reject
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import StudentSearchSelect from '@/components/StudentSearchSelect.vue'
import type {
	FinancialAdjustmentDTO,
	StudentDTO,
	InvoiceDTO,
	StudentCreditBalanceInfoDTO
} from '@/types/common'
import { formatMinor, formatCivilDate } from '@/utils/format'
import { Plus, RefreshCw, Receipt } from 'lucide-vue-next'

const services = getServices()
const toast = useToast()

const currentTab = ref('adjustments')
const adjustments = ref<FinancialAdjustmentDTO[]>([])
const studentsList = ref<StudentDTO[]>([])
const studentInvoices = ref<InvoiceDTO[]>([])
const loading = ref(false)
const statusFilter = ref('')

// Credits Tab
const selectedCreditStudentId = ref('')
const creditInfo = ref<StudentCreditBalanceInfoDTO | null>(null)

// Request Modal
const showAdjustmentModal = ref(false)
const submittingAdj = ref(false)
const adjustmentError = ref('')
const adjForm = ref({
	studentId: '',
	invoiceId: '',
	adjustmentType: 'waiver',
	amountDecimal: 50,
	reason: ''
})

// Reject Modal
const showRejectModal = ref(false)
const rejectingAdjId = ref('')
const rejectNotes = ref('')

const pendingCount = computed(() => {
	return adjustments.value.filter(a => a.Status === 'pending').length
})

onMounted(async () => {
	await Promise.all([loadStudents(), loadAdjustments()])
})

async function loadStudents() {
	try {
		const res = await services.students.listStudents({ limit: 100 })
		if (res.success && res.data) {
			studentsList.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load students', e)
	}
}

async function loadAdjustments() {
	loading.value = true
	try {
		const res = await services.payments.listFinancialAdjustments({
			status: statusFilter.value || undefined,
			limit: 100
		})
		if (res.success && res.data) {
			adjustments.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load adjustments', e)
	} finally {
		loading.value = false
	}
}

async function loadStudentInvoices() {
	if (!adjForm.value.studentId) return
	try {
		const res = await services.invoices.listInvoices({
			studentId: adjForm.value.studentId,
			limit: 50
		})
		if (res.success && res.data) {
			studentInvoices.value = res.data.filter(i => i.Status !== 'void')
		}
	} catch (e: any) {
		console.error('Failed to load student invoices', e)
	}
}

async function loadCreditBalance() {
	if (!selectedCreditStudentId.value) return
	try {
		const res = await services.payments.getStudentCreditBalance(selectedCreditStudentId.value)
		if (res.success && res.data) {
			creditInfo.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load credit balance', e)
	}
}

function openRequestAdjustmentModal() {
	adjustmentError.value = ''
	adjForm.value = {
		studentId: '',
		invoiceId: '',
		adjustmentType: 'waiver',
		amountDecimal: 50,
		reason: ''
	}
	studentInvoices.value = []
	showAdjustmentModal.value = true
}

async function submitAdjustmentRequest() {
	submittingAdj.value = true
	adjustmentError.value = ''
	try {
		const amountMinor = Math.round(adjForm.value.amountDecimal * 100)
		const res = await services.payments.requestFinancialAdjustment({
			studentId: adjForm.value.studentId,
			invoiceId: adjForm.value.invoiceId || undefined,
			adjustmentType: adjForm.value.adjustmentType,
			amountMinor: amountMinor,
			reason: adjForm.value.reason
		})
		if (!res.success) {
			adjustmentError.value = res.error?.message || 'Failed to submit adjustment request'
			return
		}
		showAdjustmentModal.value = false
		await loadAdjustments()
	} catch (e: any) {
		adjustmentError.value = e.message || 'An error occurred'
	} finally {
		submittingAdj.value = false
	}
}

async function decideAdjustment(id: string, approved: boolean, notes = '') {
	try {
		const res = await services.payments.decideFinancialAdjustment({
			adjustmentId: id,
			approved: approved,
			decisionNotes: notes
		})
		if (res.success) {
			toast.success(approved ? 'Adjustment successfully approved.' : 'Adjustment rejected.')
			await loadAdjustments()
		} else {
			toast.error(res.error?.message || 'Failed to decide adjustment')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred while processing adjustment')
	}
}

function openRejectModal(id: string) {
	rejectingAdjId.value = id
	rejectNotes.value = ''
	showRejectModal.value = true
}

async function submitReject() {
	if (!rejectingAdjId.value) return
	await decideAdjustment(rejectingAdjId.value, false, rejectNotes.value)
	showRejectModal.value = false
}

function getAdjustmentBadgeClass(type: string) {
	switch (type) {
		case 'waiver':
			return 'bg-purple-50 text-purple-700 border border-purple-200'
		case 'write_off':
			return 'bg-rose-50 text-rose-700 border border-rose-200'
		case 'credit_correction':
			return 'bg-blue-50 text-blue-700 border border-blue-200'
		case 'debit_correction':
			return 'bg-amber-50 text-amber-700 border border-amber-200'
		default:
			return 'bg-slate-100 text-slate-600'
	}
}
</script>

