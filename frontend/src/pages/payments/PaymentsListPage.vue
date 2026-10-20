<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Payments & Receipts</h1>
				<p class="text-xs text-slate-500 mt-0.5">Record student fee collections, generate sequential receipts, and manage payment allocations.</p>
			</div>
			<div class="flex gap-2">
				<router-link
					to="/payments/adjustments"
					class="py-2 px-3.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Scale class="w-4 h-4 text-slate-500" />
					<span>Financial Adjustments</span>
				</router-link>
				<button
					@click="openRecordPaymentModal"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>Record Payment</span>
				</button>
			</div>
		</div>

		<!-- Filters & Search Bar -->
		<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex flex-wrap gap-3 items-center justify-between">
			<div class="flex items-center gap-3 flex-wrap flex-1">
				<div class="relative min-w-[220px]">
					<Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
					<input
						v-model="searchFilter"
						type="text"
						placeholder="Search by student, ID, or receipt #..."
						class="w-full pl-8.5 pr-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 bg-white"
					/>
				</div>

				<div>
					<label class="text-[11px] font-semibold text-slate-600 mr-2">Method:</label>
					<select
						v-model="methodFilter"
						@change="loadPayments"
						class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					>
						<option value="">All Methods</option>
						<option value="cash">Cash</option>
						<option value="bank_transfer">Bank Transfer</option>
						<option value="cheque">Cheque</option>
						<option value="pos_card">POS / Card</option>
						<option value="mobile_money">Mobile Money</option>
					</select>
				</div>

				<div>
					<label class="text-[11px] font-semibold text-slate-600 mr-2">Status:</label>
					<select
						v-model="statusFilter"
						@change="loadPayments"
						class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					>
						<option value="">All Statuses</option>
						<option value="posted">Posted</option>
						<option value="reversed">Reversed</option>
					</select>
				</div>
			</div>

			<div class="flex items-center gap-2">
				<button
					@click="loadPayments"
					class="p-2 border border-slate-200 hover:bg-slate-50 rounded-lg text-slate-600 transition cursor-pointer"
					title="Refresh List"
				>
					<RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
				</button>
			</div>
		</div>

		<!-- Payments Table -->
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">Receipt #</th>
						<th class="py-3 px-4">Student</th>
						<th class="py-3 px-4">Date</th>
						<th class="py-3 px-4">Method & Ref</th>
						<th class="py-3 px-4 text-right">Amount Received</th>
						<th class="py-3 px-4 text-right">Allocated</th>
						<th class="py-3 px-4 text-right">Credit (Overpaid)</th>
						<th class="py-3 px-4 text-center">Status</th>
						<th class="py-3 px-4 text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-slate-700">
					<tr v-if="filteredPayments.length === 0">
						<td colspan="9" class="py-8 text-center text-slate-400">
							{{ paymentsList.length === 0 ? 'No payments recorded yet.' : 'No payments matching search criteria.' }}
						</td>
					</tr>
					<tr v-for="p in filteredPayments" :key="p.ID" class="hover:bg-slate-50/50">
						<td class="py-3 px-4 font-mono font-bold text-slate-900">
							<router-link :to="`/payments/${p.ID}`" class="hover:text-indigo-600 hover:underline">
								{{ p.ReceiptNo }}
							</router-link>
						</td>
						<td class="py-3 px-4">
							<div class="font-semibold text-slate-900">{{ p.StudentFirstName }} {{ p.StudentLastName }}</div>
							<div class="font-mono text-[10px] text-slate-400">Adm: {{ p.AdmissionNo }}</div>
						</td>
						<td class="py-3 px-4 text-slate-600">
							{{ formatCivilDate(p.PaymentDate) }}
						</td>
						<td class="py-3 px-4">
							<span class="capitalize font-medium text-slate-800">{{ p.PaymentMethod.replace('_', ' ') }}</span>
							<div v-if="p.ReferenceNo" class="font-mono text-[10px] text-slate-400">Ref: {{ p.ReferenceNo }}</div>
						</td>
						<td class="py-3 px-4 text-right font-bold text-slate-900">
							{{ formatMinor(p.AmountMinor) }}
						</td>
						<td class="py-3 px-4 text-right text-emerald-600 font-medium">
							{{ formatMinor(p.AllocatedAmountMinor) }}
						</td>
						<td class="py-3 px-4 text-right">
							<span v-if="p.UnallocatedAmountMinor > 0" class="text-blue-600 font-bold">
								+{{ formatMinor(p.UnallocatedAmountMinor) }}
							</span>
							<span v-else class="text-slate-400">—</span>
						</td>
						<td class="py-3 px-4 text-center">
							<span
								class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
								:class="p.Status === 'posted' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-rose-50 text-rose-700 border border-rose-200'"
							>
								{{ p.Status }}
							</span>
						</td>
						<td class="py-3 px-4 text-right space-x-2">
							<router-link
								:to="`/payments/${p.ID}`"
								class="text-indigo-600 hover:text-indigo-800 font-medium text-xs"
							>
								Receipt
							</router-link>
							<button
								v-if="p.Status !== 'reversed'"
								@click="openReverseModal(p.ID, p.ReceiptNo)"
								class="text-rose-500 hover:text-rose-700 font-medium text-xs cursor-pointer"
							>
								Reverse
							</button>
						</td>
					</tr>
					<tr v-if="paymentsList.length === 0">
						<td colspan="9" class="py-10 text-center text-slate-400 italic">
							No payments recorded. Click "Record Payment" to post a collection.
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- MODAL: Record Payment -->
		<div v-if="showPaymentModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-md w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Record Student Payment</h3>
					<button @click="showPaymentModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitRecordPayment" class="p-4 space-y-3">
					<div v-if="paymentError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ paymentError }}
					</div>

					<StudentSearchSelect
						v-model="paymentForm.studentId"
						:students="studentsList"
						label="Select Student"
						required
						placeholder="Type student name or ID (e.g. Diana, ADM-2026-0004)..."
						@select="onStudentSelected"
					/>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Amount (PKR) *</label>
							<input
								v-model.number="paymentForm.amountDecimal"
								type="number"
								min="1"
								step="1"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 font-bold"
							/>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Payment Date *</label>
							<input
								v-model="paymentForm.paymentDate"
								type="date"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Payment Method *</label>
							<select
								v-model="paymentForm.paymentMethod"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							>
								<option value="cash">Cash</option>
								<option value="bank_transfer">Bank Transfer</option>
								<option value="cheque">Cheque</option>
								<option value="pos_card">POS / Card</option>
								<option value="mobile_money">Mobile Money</option>
							</select>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Reference # (Optional)</label>
							<input
								v-model="paymentForm.referenceNo"
								type="text"
								placeholder="e.g. Cheque / Txn ID"
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Payer Name</label>
							<input
								v-model="paymentForm.payerName"
								type="text"
								placeholder="Parent / Guardian"
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Payer Phone</label>
							<input
								v-model="paymentForm.payerPhone"
								type="text"
								placeholder="Contact number"
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Notes</label>
						<textarea
							v-model="paymentForm.notes"
							rows="2"
							placeholder="Optional memo..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>

					<div class="p-2.5 bg-blue-50 border border-blue-100 rounded-lg text-[11px] text-blue-700">
						💡 Payments are automatically allocated against the student's oldest unpaid invoices (FIFO). Any excess overpayment is preserved as credit for future billing.
					</div>

					<div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showPaymentModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="savingPayment"
							class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							{{ savingPayment ? 'Posting Payment...' : 'Post Payment & Print' }}
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- MODAL: Reverse Payment -->
		<div v-if="showReverseModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-sm w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Reverse Payment ({{ reverseReceiptNo }})</h3>
					<button @click="showReverseModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitReversePayment" class="p-4 space-y-3">
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason for Reversal *</label>
						<textarea
							v-model="reverseReason"
							required
							rows="3"
							placeholder="State justification (e.g. cheque bounced, incorrect amount entered)..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>
					<div class="p-2 bg-amber-50 border border-amber-200 rounded text-[10px] text-amber-800">
						⚠️ Reversing this payment will restore the outstanding balances on all associated student invoices.
					</div>
					<div class="flex justify-end gap-2 pt-2 border-t border-slate-100">
						<button
							type="button"
							@click="showReverseModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-1.5 bg-rose-600 hover:bg-rose-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							Confirm Reversal
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import StudentSearchSelect from '@/components/StudentSearchSelect.vue'
import type { PaymentDTO, StudentDTO } from '@/types/common'
import { formatMinor, formatCivilDate } from '@/utils/format'
import { Plus, RefreshCw, Scale, Search } from 'lucide-vue-next'

const router = useRouter()
const services = getServices()
const toast = useToast()

const paymentsList = ref<PaymentDTO[]>([])
const studentsList = ref<StudentDTO[]>([])
const loading = ref(false)

const methodFilter = ref('')
const statusFilter = ref('')
const searchFilter = ref('')

const filteredPayments = computed(() => {
	const q = searchFilter.value.trim().toLowerCase()
	if (!q) return paymentsList.value
	return paymentsList.value.filter(p => {
		const name = (p.StudentFirstName ? `${p.StudentFirstName} ${p.StudentLastName || ''}` : p.PayerName || '').toLowerCase()
		const rcp = (p.ReceiptNo || '').toLowerCase()
		const ref = (p.ReferenceNo || '').toLowerCase()
		const adm = (p.AdmissionNo || '').toLowerCase()
		return name.includes(q) || rcp.includes(q) || ref.includes(q) || adm.includes(q)
	})
})

function onStudentSelected(s: StudentDTO | null) {
	if (s && !paymentForm.value.payerName) {
		paymentForm.value.payerName = `${s.FirstName} ${s.LastName}'s Guardian`
	}
}

// Payment Modal
const showPaymentModal = ref(false)
const savingPayment = ref(false)
const paymentError = ref('')
const paymentForm = ref({
	studentId: '',
	paymentDate: new Date().toISOString().split('T')[0],
	paymentMethod: 'cash',
	amountDecimal: 100,
	referenceNo: '',
	payerName: '',
	payerPhone: '',
	notes: ''
})

// Reverse Modal
const showReverseModal = ref(false)
const reversePaymentId = ref('')
const reverseReceiptNo = ref('')
const reverseReason = ref('')

onMounted(async () => {
	await Promise.all([loadStudents(), loadPayments()])
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

async function loadPayments() {
	loading.value = true
	try {
		const res = await services.payments.listPayments({
			status: statusFilter.value || undefined,
			paymentMethod: methodFilter.value || undefined,
			limit: 100
		})
		if (res.success && res.data) {
			paymentsList.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load payments', e)
	} finally {
		loading.value = false
	}
}

function openRecordPaymentModal() {
	paymentError.value = ''
	paymentForm.value = {
		studentId: '',
		paymentDate: new Date().toISOString().split('T')[0],
		paymentMethod: 'cash',
		amountDecimal: 100,
		referenceNo: '',
		payerName: '',
		payerPhone: '',
		notes: ''
	}
	showPaymentModal.value = true
}

async function submitRecordPayment() {
	savingPayment.value = true
	paymentError.value = ''
	try {
		const amountMinor = Math.round(paymentForm.value.amountDecimal * 100)
		const res = await services.payments.recordPayment({
			studentId: paymentForm.value.studentId,
			paymentDate: paymentForm.value.paymentDate,
			paymentMethod: paymentForm.value.paymentMethod,
			amountMinor: amountMinor,
			referenceNo: paymentForm.value.referenceNo || undefined,
			payerName: paymentForm.value.payerName || undefined,
			payerPhone: paymentForm.value.payerPhone || undefined,
			notes: paymentForm.value.notes || undefined
		})
		if (!res.success) {
			paymentError.value = res.error?.message || 'Failed to record payment'
			return
		}

		showPaymentModal.value = false
		await loadPayments()
		if (res.data?.paymentId) {
			router.push(`/payments/${res.data.paymentId}`)
		}
	} catch (e: any) {
		paymentError.value = e.message || 'An error occurred'
	} finally {
		savingPayment.value = false
	}
}

function openReverseModal(paymentId: string, receiptNo: string) {
	reversePaymentId.value = paymentId
	reverseReceiptNo.value = receiptNo
	reverseReason.value = ''
	showReverseModal.value = true
}

async function submitReversePayment() {
	if (!reversePaymentId.value) return
	try {
		const res = await services.payments.reversePayment({
			paymentId: reversePaymentId.value,
			reason: reverseReason.value
		})
		if (res.success) {
			showReverseModal.value = false
			toast.success('Payment successfully reversed and ledger updated.')
			await loadPayments()
		} else {
			toast.error(res.error?.message || 'Failed to reverse payment')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred while reversing payment')
	}
}
</script>

