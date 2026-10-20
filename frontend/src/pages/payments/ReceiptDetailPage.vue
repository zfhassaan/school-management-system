<template>
	<div class="max-w-3xl mx-auto space-y-6">
		<!-- Top Action Bar (hidden on print) -->
		<div class="flex items-center justify-between no-print">
			<router-link
				to="/payments"
				class="inline-flex items-center gap-1.5 text-xs font-medium text-slate-500 hover:text-slate-800 transition"
			>
				<ArrowLeft class="w-4 h-4" />
				<span>Back to Payments</span>
			</router-link>

			<div class="flex items-center gap-2">
				<button
					@click="printReceipt"
					class="py-2 px-3.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Printer class="w-4 h-4 text-slate-500" />
					<span>Print Receipt</span>
				</button>
				<button
					v-if="detail?.payment && detail.payment.Status !== 'reversed'"
					@click="showReverseModal = true"
					class="py-2 px-3.5 bg-rose-50 border border-rose-200 hover:bg-rose-100 text-rose-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Ban class="w-4 h-4" />
					<span>Reverse Payment</span>
				</button>
			</div>
		</div>

		<!-- Receipt Paper Container -->
		<div v-if="detail?.payment" class="bg-white rounded-xl border border-slate-200 shadow-sm p-8 space-y-6 print:border-none print:shadow-none print:p-0">
			<!-- Header -->
			<div class="flex items-start justify-between border-b border-slate-100 pb-5">
				<div>
					<div class="text-[10px] font-bold uppercase tracking-wider text-indigo-600 mb-0.5">Fee Collection Receipt</div>
					<h1 class="text-xl font-black text-slate-900">{{ schoolProfile?.name || 'School Management System' }}</h1>
					<p class="text-xs text-slate-500">{{ schoolProfile?.address?.line1 || 'Academic Bursar & Cashier Office' }}</p>
					<p v-if="schoolProfile?.phone" class="text-xs text-slate-500">Tel: {{ schoolProfile.phone }}</p>
				</div>

				<div class="text-right">
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-400">Receipt Number</div>
					<div class="text-lg font-black font-mono text-slate-900">{{ detail.payment.ReceiptNo }}</div>
					<div class="mt-1">
						<span
							class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wide inline-block"
							:class="detail.payment.Status === 'posted' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-rose-50 text-rose-700 border border-rose-200'"
						>
							{{ detail.payment.Status }}
						</span>
					</div>
				</div>
			</div>

			<!-- Reversal Banner if reversed -->
			<div v-if="detail.payment.Status === 'reversed'" class="p-3 bg-rose-50 border border-rose-200 rounded-lg text-xs text-rose-800">
				<div class="font-bold">PAYMENT REVERSED:</div>
				<div class="mt-0.5">This transaction has been officially voided and reversed.</div>
				<div v-if="detail.reversal" class="text-[11px] text-rose-600 mt-1">
					Reversal Reason: "{{ detail.reversal.Reason }}" • Reversed by {{ detail.reversal.ReversedByName || 'Admin' }} on {{ formatCivilDate(detail.reversal.ReversalDate) }}
				</div>
			</div>

			<!-- Student & Payer Info Grid -->
			<div class="grid grid-cols-2 gap-4 bg-slate-50/70 p-4 rounded-xl border border-slate-100 text-xs">
				<div class="space-y-1">
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-400">Student Particulars</div>
					<div class="font-bold text-slate-900">{{ detail.payment.StudentFirstName }} {{ detail.payment.StudentLastName }}</div>
					<div class="text-slate-600 font-mono text-[11px]">Admission #: {{ detail.payment.AdmissionNo }}</div>
					<div v-if="detail.payment.PayerName" class="text-slate-600 pt-1">
						Payer: <span class="font-medium text-slate-800">{{ detail.payment.PayerName }}</span>
						<span v-if="detail.payment.PayerPhone" class="text-slate-400"> ({{ detail.payment.PayerPhone }})</span>
					</div>
				</div>

				<div class="space-y-1 text-right">
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-400">Payment Details</div>
					<div class="text-slate-700">Date: <span class="font-medium text-slate-900">{{ formatCivilDate(detail.payment.PaymentDate) }}</span></div>
					<div class="text-slate-700">Method: <span class="font-medium text-slate-900 capitalize">{{ detail.payment.PaymentMethod.replace('_', ' ') }}</span></div>
					<div v-if="detail.payment.ReferenceNo" class="text-slate-600 font-mono text-[11px]">Ref: {{ detail.payment.ReferenceNo }}</div>
					<div class="text-[10px] text-slate-400 pt-1">Cashier: {{ detail.payment.RecordedByName || 'Staff' }}</div>
				</div>
			</div>

			<!-- Bill Allocations Table -->
			<div class="space-y-2">
				<div class="text-xs font-bold uppercase tracking-wider text-slate-500">Bill Allocation Breakdown</div>
				<div class="border border-slate-200 rounded-lg overflow-hidden">
					<table class="w-full text-left text-xs border-collapse">
						<thead>
							<tr class="bg-slate-50 border-b border-slate-200 text-[10px] font-bold uppercase tracking-wider text-slate-500">
								<th class="py-2.5 px-3">Invoice #</th>
								<th class="py-2.5 px-3">Billing Period</th>
								<th class="py-2.5 px-3 text-right">Allocated Amount</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100 text-slate-700">
							<tr v-for="alloc in detail.allocations" :key="alloc.ID" class="hover:bg-slate-50/30">
								<td class="py-2.5 px-3 font-mono font-bold text-slate-900">{{ alloc.InvoiceNo }}</td>
								<td class="py-2.5 px-3 font-medium text-slate-700">{{ alloc.BillingPeriod }}</td>
								<td class="py-2.5 px-3 text-right font-bold text-emerald-700">
									{{ formatMinor(alloc.AllocatedAmountMinor) }}
								</td>
							</tr>
							<tr v-if="detail.allocations.length === 0">
								<td colspan="3" class="py-4 text-center text-slate-400 italic">
									Unallocated payment / advance deposit.
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>

			<!-- Financial Totals -->
			<div class="flex justify-end pt-2">
				<div class="w-64 space-y-1.5 text-xs">
					<div class="flex justify-between py-1 text-slate-600 border-b border-slate-100">
						<span>Total Tendered:</span>
						<span class="font-bold text-slate-900 text-sm">{{ formatMinor(detail.payment.AmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-1 text-slate-600 border-b border-slate-100">
						<span>Allocated to Bills:</span>
						<span class="font-medium text-emerald-700">{{ formatMinor(detail.payment.AllocatedAmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-1.5 text-slate-900 font-bold bg-slate-50 px-2 rounded-lg">
						<span>Excess / Credit Balance:</span>
						<span :class="detail.payment.UnallocatedAmountMinor > 0 ? 'text-blue-700 font-bold' : 'text-slate-500'">
							{{ formatMinor(detail.payment.UnallocatedAmountMinor) }}
						</span>
					</div>
				</div>
			</div>

			<!-- Footer / Signature -->
			<div class="border-t border-slate-100 pt-8 mt-6">
				<div class="flex justify-between items-end text-xs text-slate-500">
					<div>
						<p class="text-[10px] text-slate-400">Thank you for your payment.</p>
						<p class="text-[10px] text-slate-300 mt-1">Official computerized receipt issued by School Management System.</p>
					</div>
					<div class="text-center w-48">
						<div class="border-b border-slate-300 pb-1 mb-1 font-mono text-[11px] text-slate-600">
							{{ detail.payment.RecordedByName || 'Authorized Signatory' }}
						</div>
						<div class="text-[10px] text-slate-400 uppercase tracking-wider font-semibold">Authorized Signature</div>
					</div>
				</div>
			</div>
		</div>

		<div v-else-if="loading" class="p-12 text-center text-xs text-slate-400">
			Loading payment receipt...
		</div>

		<!-- MODAL: Reverse Payment -->
		<div v-if="showReverseModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50 no-print">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-sm w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Reverse Payment</h3>
					<button @click="showReverseModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitReverse" class="p-4 space-y-3">
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason for Reversal *</label>
						<textarea
							v-model="reverseReason"
							required
							rows="3"
							placeholder="State justification..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
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
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getServices } from '@/services/factory'
import type { PaymentDetailDTO, SchoolProfileDTO } from '@/types/common'
import { formatMinor, formatCivilDate } from '@/utils/format'
import { ArrowLeft, Printer, Ban } from 'lucide-vue-next'

const route = useRoute()
const services = getServices()

const detail = ref<PaymentDetailDTO | null>(null)
const schoolProfile = ref<SchoolProfileDTO | null>(null)
const loading = ref(false)

const showReverseModal = ref(false)
const reverseReason = ref('')

onMounted(async () => {
	await Promise.all([loadPayment(), loadSchoolProfile()])
})

async function loadPayment() {
	const id = route.params.id as string
	if (!id) return
	loading.value = true
	try {
		const res = await services.payments.getPayment(id)
		if (res.success && res.data) {
			detail.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load payment', e)
	} finally {
		loading.value = false
	}
}

async function loadSchoolProfile() {
	try {
		const res = await services.school.getSchool()
		if (res.success && res.data) {
			schoolProfile.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load school profile', e)
	}
}

function printReceipt() {
	window.print()
}

async function submitReverse() {
	if (!detail.value?.payment) return
	try {
		const res = await services.payments.reversePayment({
			paymentId: detail.value.payment.ID,
			reason: reverseReason.value
		})
		if (res.success) {
			showReverseModal.value = false
			await loadPayment()
		} else {
			alert(res.error?.message || 'Failed to reverse payment')
		}
	} catch (e: any) {
		alert(e.message || 'An error occurred')
	}
}
</script>

<style scoped>
@media print {
	.no-print {
		display: none !important;
	}
}
</style>

