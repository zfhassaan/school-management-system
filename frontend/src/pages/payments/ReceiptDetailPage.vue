<template>
	<div class="max-w-4xl mx-auto space-y-6">
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
				<router-link
					v-if="linkedInvoiceId"
					:to="`/invoices/${linkedInvoiceId}`"
					class="py-2 px-3.5 bg-white border border-slate-200 hover:bg-slate-50 text-indigo-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<FileText class="w-4 h-4 text-indigo-500" />
					<span>View Invoice ({{ linkedInvoiceNo }})</span>
				</router-link>
				<button
					@click="printReceipt"
					class="py-2 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Printer class="w-4 h-4" />
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
		<div
			v-if="detail?.payment"
			class="receipt-paper bg-white rounded-2xl border border-slate-200 shadow-sm p-8 space-y-6 text-slate-800"
		>
			<!-- 1. Formal Institutional Header -->
			<div class="flex items-start justify-between border-b-2 border-slate-900 pb-5">
				<div class="flex items-start gap-4">
					<!-- School Monogram / Crest Badge -->
					<div class="w-16 h-16 rounded-xl bg-slate-900 text-white flex items-center justify-center font-black text-xl shrink-0 shadow-sm tracking-wider overflow-hidden">
						<img v-if="authStore.schoolLogo" :src="authStore.schoolLogo" alt="School Logo" class="w-full h-full object-contain p-1 bg-white" />
						<span v-else>{{ schoolInitials }}</span>
					</div>
					<div>
						<div class="flex items-center gap-2">
							<h1 class="text-2xl font-black text-slate-900 tracking-tight">{{ currentSchoolName }}</h1>
							<span v-if="currentSchoolCode" class="text-[10px] font-bold font-mono px-2 py-0.5 rounded bg-slate-100 text-slate-700 border border-slate-200">
								{{ currentSchoolCode }}
							</span>
						</div>
						<p v-if="currentSchoolLegalName" class="text-xs font-semibold text-slate-600 mt-0.5">
							{{ currentSchoolLegalName }}
						</p>
						<p class="text-xs text-slate-500 mt-0.5">
							{{ currentSchoolAddress }}
						</p>
						<p v-if="currentSchoolContact" class="text-[11px] text-slate-500 mt-0.5">
							{{ currentSchoolContact }}
						</p>
					</div>
				</div>

				<div class="text-right">
					<div class="inline-block text-[11px] font-black uppercase tracking-wider px-2.5 py-1 rounded bg-indigo-50 text-indigo-700 border border-indigo-200 mb-2">
						Fee Collection Voucher
					</div>
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-400">Receipt Number</div>
					<div class="text-xl font-black font-mono text-slate-900 mt-0.5">{{ detail.payment.ReceiptNo }}</div>
					<div class="mt-1.5 flex items-center justify-end gap-1.5">
						<span
							class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wide inline-flex items-center gap-1"
							:class="detail.payment.Status === 'posted' ? 'bg-emerald-50 text-emerald-700 border border-emerald-300' : 'bg-rose-50 text-rose-700 border border-rose-300'"
						>
							<CheckCircle2 v-if="detail.payment.Status === 'posted'" class="w-3 h-3 text-emerald-600" />
							<span>{{ detail.payment.Status === 'posted' ? 'PAID & POSTED' : detail.payment.Status }}</span>
						</span>
					</div>
				</div>
			</div>

			<!-- Reversal Banner if reversed -->
			<div v-if="detail.payment.Status === 'reversed'" class="p-3 bg-rose-50 border border-rose-200 rounded-lg text-xs text-rose-800">
				<div class="font-bold flex items-center gap-1.5 text-rose-700">
					<Ban class="w-4 h-4" />
					<span>TRANSACTION VOIDED / REVERSED</span>
				</div>
				<div class="mt-0.5 text-slate-700">This fee payment was formally reversed and credited back to the accounts ledger.</div>
				<div v-if="detail.reversal" class="text-[11px] text-rose-600 mt-1 font-medium">
					Reason: "{{ detail.reversal.Reason }}" • Reversed by {{ detail.reversal.ReversedByName || 'Admin' }} on {{ formatCivilDate(detail.reversal.ReversalDate) }}
				</div>
			</div>

			<!-- 2. Student, Academic & Payer Particulars (Formal Grid) -->
			<div class="grid grid-cols-2 md:grid-cols-4 gap-4 p-4 rounded-xl border border-slate-200 bg-slate-50/60 text-xs">
				<div>
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block mb-0.5">Student Name</span>
					<p class="font-bold text-slate-900 text-sm">
						{{ detail.payment.StudentFirstName }} {{ detail.payment.StudentLastName }}
					</p>
					<p class="text-[11px] font-mono font-medium text-slate-500 mt-0.5">
						ID: {{ detail.payment.AdmissionNo }}
					</p>
				</div>

				<div>
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block mb-0.5">Class / Grade</span>
					<p class="font-bold text-slate-900">
						{{ studentClassDisplay }}
					</p>
					<p class="text-[11px] text-slate-500 mt-0.5">
						Session: {{ studentSessionDisplay }}
					</p>
				</div>

				<div>
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block mb-0.5">Payer Particulars</span>
					<p class="font-semibold text-slate-900">
						{{ detail.payment.PayerName || 'Self / Guardian' }}
					</p>
					<p v-if="detail.payment.PayerPhone" class="text-[11px] text-slate-500 font-mono mt-0.5">
						Tel: {{ detail.payment.PayerPhone }}
					</p>
				</div>

				<div class="text-right md:text-left">
					<span class="text-[10px] font-bold uppercase tracking-wider text-slate-400 block mb-0.5">Payment Method & Date</span>
					<p class="font-bold text-slate-900 capitalize">
						{{ detail.payment.PaymentMethod.replace('_', ' ') }}
					</p>
					<p class="text-[11px] text-slate-600 mt-0.5">
						Date: {{ formatCivilDate(detail.payment.PaymentDate) }}
					</p>
					<p v-if="detail.payment.ReferenceNo" class="text-[10px] font-mono text-slate-400 mt-0.5">
						Ref: {{ detail.payment.ReferenceNo }}
					</p>
				</div>
			</div>

			<!-- 3. Settled Bill & Invoice Itemization -->
			<div class="space-y-3">
				<div class="flex items-center justify-between">
					<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">
						Settled Invoice Particulars & Fee Line Items
					</h3>
					<span v-if="detail.allocations.length > 0" class="text-[11px] font-semibold text-indigo-600">
						{{ detail.allocations.length }} Invoice Allocation(s)
					</span>
				</div>

				<!-- Table of Fee Allocations -->
				<div class="border border-slate-200 rounded-xl overflow-hidden shadow-2xs">
					<table class="w-full text-left text-xs border-collapse">
						<thead>
							<tr class="bg-slate-100 border-b border-slate-200 text-[10px] font-bold uppercase tracking-wider text-slate-600">
								<th class="py-2.5 px-3.5">Invoice #</th>
								<th class="py-2.5 px-3.5">Billing Period</th>
								<th class="py-2.5 px-3.5 text-right">Invoice Billed</th>
								<th class="py-2.5 px-3.5 text-right">Payment Allocated</th>
								<th class="py-2.5 px-3.5 text-right">Remaining Balance</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100 text-slate-700">
							<tr v-for="alloc in detail.allocations" :key="alloc.ID" class="hover:bg-slate-50/50">
								<td class="py-2.5 px-3.5 font-mono font-bold text-slate-900">
									{{ alloc.InvoiceNo }}
								</td>
								<td class="py-2.5 px-3.5 font-medium text-slate-700">
									{{ alloc.BillingPeriod }}
								</td>
								<td class="py-2.5 px-3.5 text-right font-medium text-slate-600">
									{{ getInvoiceTotalDisplay(alloc.InvoiceID) }}
								</td>
								<td class="py-2.5 px-3.5 text-right font-bold text-emerald-700">
									{{ formatMinor(alloc.AllocatedAmountMinor) }}
								</td>
								<td class="py-2.5 px-3.5 text-right font-mono font-semibold" :class="getInvoiceBalance(alloc.InvoiceID) === 0 ? 'text-emerald-600 font-bold' : 'text-slate-800'">
									{{ getInvoiceBalanceDisplay(alloc.InvoiceID) }}
								</td>
							</tr>

							<!-- Advance Student Credit / Deposit Line Item -->
							<tr v-if="detail.payment.UnallocatedAmountMinor > 0" class="hover:bg-slate-50/50 bg-indigo-50/20">
								<td class="py-2.5 px-3.5 font-mono font-bold text-slate-900">
									<span class="inline-flex items-center gap-1.5">
										<span class="px-1.5 py-0.5 rounded text-[10px] bg-indigo-100 text-indigo-700 font-semibold border border-indigo-200">ADVANCE</span>
										<span>ADV-CREDIT</span>
									</span>
								</td>
								<td class="py-2.5 px-3.5 font-medium text-slate-700">
									Student Advance Credit Ledger
								</td>
								<td class="py-2.5 px-3.5 text-right font-medium text-slate-600">
									{{ formatMinor(detail.payment.AmountMinor) }}
								</td>
								<td class="py-2.5 px-3.5 text-right font-bold text-indigo-700">
									{{ formatMinor(detail.payment.UnallocatedAmountMinor) }}
								</td>
								<td class="py-2.5 px-3.5 text-right font-mono font-semibold text-emerald-600 font-bold">
									0.00 (CREDITED)
								</td>
							</tr>

							<!-- Empty state only if zero allocations and zero unallocated -->
							<tr v-if="detail.allocations.length === 0 && detail.payment.UnallocatedAmountMinor === 0">
								<td colspan="5" class="py-6 text-center text-slate-500 italic">
									No fee allocations recorded for this voucher.
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- Advance Payment Note if advance amount exists -->
				<div v-if="detail.payment.UnallocatedAmountMinor > 0" class="p-3 bg-indigo-50/60 border border-indigo-100 rounded-xl text-xs text-indigo-900 flex items-center justify-between gap-3">
					<div class="flex items-center gap-2">
						<span class="w-2 h-2 rounded-full bg-indigo-500 shrink-0"></span>
						<span><strong>Advance Payment Deposit:</strong> Funds of {{ formatMinor(detail.payment.UnallocatedAmountMinor) }} credited directly to the student's credit advance ledger for subsequent tuition billings.</span>
					</div>
					<span class="font-mono font-bold text-indigo-700 shrink-0">{{ formatMinor(detail.payment.UnallocatedAmountMinor) }}</span>
				</div>

				<!-- Detailed Itemized Fee Breakdown if linked invoice items loaded -->
				<div v-if="linkedInvoiceDetail?.items && linkedInvoiceDetail.items.length > 0" class="mt-3 p-3.5 bg-slate-50 border border-slate-200 rounded-xl space-y-2">
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-500">
						Itemized Fee Head Breakdown ({{ linkedInvoiceDetail.invoice.InvoiceNo }})
					</div>
					<div class="grid grid-cols-2 md:grid-cols-3 gap-2 text-xs">
						<div
							v-for="item in linkedInvoiceDetail.items"
							:key="item.ID"
							class="flex items-center justify-between p-2 bg-white rounded-lg border border-slate-200"
						>
							<span class="font-medium text-slate-700 truncate pr-2">{{ item.FeeHeadName }}</span>
							<span class="font-bold text-slate-900 shrink-0">{{ formatMinor(item.NetAmountMinor) }}</span>
						</div>
					</div>
				</div>

				<!-- Standard Class Fee Head Schedule if pure advance payment deposit -->
				<div v-else-if="classFeeItems.length > 0" class="mt-3 p-3.5 bg-slate-50 border border-slate-200 rounded-xl space-y-2">
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-500">
						Standard Applicable Fee Heads ({{ studentClassDisplay }})
					</div>
					<div class="grid grid-cols-2 md:grid-cols-3 gap-2 text-xs">
						<div
							v-for="(item, idx) in classFeeItems"
							:key="item.ID || item.id || idx"
							class="flex items-center justify-between p-2 bg-white rounded-lg border border-slate-200"
						>
							<span class="font-medium text-slate-700 truncate pr-2">{{ item.FeeHeadName || item.feeHeadName }}</span>
							<span class="font-bold text-slate-900 shrink-0">{{ formatMinor(item.AmountMinor ?? item.amountMinor ?? 0) }}</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 4. Financial Totals & Balance Summary Box -->
			<div class="flex flex-col md:flex-row justify-between items-start gap-4 pt-2">
				<div class="text-xs text-slate-500 space-y-1 max-w-sm">
					<p class="font-bold text-slate-700 uppercase text-[10px] tracking-wider">Payment Verification & Mode</p>
					<p>Received by Cashier: <span class="font-semibold text-slate-800">{{ detail.payment.RecordedByName || 'Accounts Staff' }}</span></p>
					<p v-if="detail.payment.Notes" class="italic text-slate-600">Notes: {{ detail.payment.Notes }}</p>
					<p class="text-[11px] text-slate-400">All fees are processed under {{ currentSchoolName }} official billing regulations.</p>
				</div>

				<div class="w-full md:w-80 border border-slate-200 rounded-xl p-3.5 bg-slate-50/80 space-y-2 text-xs">
					<div class="flex justify-between py-1 text-slate-600 border-b border-slate-200">
						<span class="font-medium">Total Amount Received:</span>
						<span class="font-black text-slate-900 text-sm">{{ formatMinor(detail.payment.AmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-1 text-slate-600 border-b border-slate-200">
						<span>Allocated to Tuition / Bills:</span>
						<span class="font-bold text-emerald-700">{{ formatMinor(detail.payment.AllocatedAmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-1 text-slate-600 border-b border-slate-200">
						<span>Advance Student Credit:</span>
						<span :class="detail.payment.UnallocatedAmountMinor > 0 ? 'text-indigo-700 font-bold' : 'text-slate-500'">
							{{ formatMinor(detail.payment.UnallocatedAmountMinor) }}
						</span>
					</div>
					<div class="flex justify-between py-1.5 bg-white px-2.5 rounded-lg border border-slate-200 text-slate-900 font-bold">
						<span>Account Balance Due:</span>
						<span :class="remainingBalanceOverall === 0 ? 'text-emerald-700 font-bold' : 'text-rose-700'">
							{{ remainingBalanceOverall === 0 ? 'PAID IN FULL (0.00)' : formatMinor(remainingBalanceOverall) }}
						</span>
					</div>
				</div>
			</div>

			<!-- 5. Formal Certification & Dual Signatures -->
			<div class="border-t-2 border-slate-200 pt-8 mt-6 space-y-6">
				<div class="flex justify-between items-end text-xs text-slate-600">
					<!-- Depositor Signature -->
					<div class="text-center w-56">
						<div class="border-b border-slate-400 pb-1 mb-1 font-medium text-slate-700">
							{{ detail.payment.PayerName || detail.payment.StudentFirstName + ' ' + detail.payment.StudentLastName }}
						</div>
						<div class="text-[10px] text-slate-400 uppercase tracking-wider font-bold">Depositor / Payer Signature</div>
					</div>

					<!-- Official School Seal & Signatory -->
					<div class="text-center w-64">
						<div class="border-b border-slate-400 pb-1 mb-1 font-mono font-semibold text-slate-800">
							{{ detail.payment.RecordedByName || 'Accounts & Bursar Office' }}
						</div>
						<div class="text-[10px] text-slate-500 uppercase tracking-wider font-bold">
							Authorized Bursar Signature & Seal
						</div>
					</div>
				</div>

				<div class="pt-4 border-t border-slate-100 flex items-center justify-between text-[10px] text-slate-400">
					<span>Official Computerized Receipt • Issued by {{ currentSchoolName }}</span>
					<span>Document Ref: {{ detail.payment.ID.slice(0, 8) }} • {{ printedDateTime }}</span>
				</div>
			</div>
		</div>

		<div v-else-if="loading" class="p-16 text-center text-xs text-slate-400">
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
							placeholder="State clear reason for reversing this payment transaction..."
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
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import { useAuthStore } from '@/stores/auth'
import type { PaymentDetailDTO, EnrollmentDTO, InvoiceDetailDTO } from '@/types/common'
import { formatMinor, formatCivilDate } from '@/utils/format'
import { ArrowLeft, Printer, Ban, FileText, CheckCircle2 } from 'lucide-vue-next'

const route = useRoute()
const services = getServices()
const toast = useToast()
const authStore = useAuthStore()

const detail = ref<PaymentDetailDTO | null>(null)
const schoolProfile = ref<any>(null)
const studentEnrollment = ref<EnrollmentDTO | null>(null)
const linkedInvoiceDetail = ref<InvoiceDetailDTO | null>(null)
const classFeeItems = ref<any[]>([])
const loading = ref(false)

const showReverseModal = ref(false)
const reverseReason = ref('')

const currentSchoolName = computed(() => {
	const p = schoolProfile.value
	const a = authStore.school
	return p?.name || p?.Name || a?.name || 'School Management System'
})

const currentSchoolLegalName = computed(() => {
	const p = schoolProfile.value
	const a = authStore.school
	return p?.legalName || p?.LegalName || a?.legalName || ''
})

const currentSchoolCode = computed(() => {
	const p = schoolProfile.value
	const a = authStore.school
	return p?.code || p?.Code || a?.code || ''
})

const schoolInitials = computed(() => {
	const code = currentSchoolCode.value
	if (code && code !== 'OFFLINE' && code.length >= 2) {
		return code.length <= 3 ? code.toUpperCase() : code.substring(0, 2).toUpperCase()
	}
	const name = currentSchoolName.value
	const words = name.trim().split(/\s+/)
	if (words.length >= 2) {
		return (words[0][0] + words[1][0]).toUpperCase()
	}
	return name.substring(0, 2).toUpperCase()
})

const currentSchoolAddress = computed(() => {
	const p = schoolProfile.value
	const a = authStore.school
	const addr = p?.address || p?.Address || a?.address
	if (addr) {
		const parts = [addr.line1 || addr.Line1, addr.city || addr.City, addr.country || addr.Country].filter(Boolean)
		if (parts.length > 0) return parts.join(', ')
	}
	return 'Academic Campus & Financial Office'
})

const currentSchoolContact = computed(() => {
	const p = schoolProfile.value
	const a = authStore.school
	const phone = p?.phone || p?.Phone || a?.phone
	const email = p?.email || p?.Email || a?.email
	const parts = []
	if (phone) parts.push(`Tel: ${phone}`)
	if (email) parts.push(`Email: ${email}`)
	return parts.join(' • ')
})

const studentClassDisplay = computed(() => {
	if (studentEnrollment.value?.ClassName) {
		const sec = studentEnrollment.value.SectionName ? ` (${studentEnrollment.value.SectionName})` : ''
		return `${studentEnrollment.value.ClassName}${sec}`
	}
	return 'Enrolled Student'
})

const studentSessionDisplay = computed(() => {
	return studentEnrollment.value?.YearCode || 'Current Academic Session'
})

const linkedInvoiceId = computed(() => {
	if (detail.value?.allocations && detail.value.allocations.length > 0) {
		return detail.value.allocations[0].InvoiceID
	}
	return ''
})

const linkedInvoiceNo = computed(() => {
	if (detail.value?.allocations && detail.value.allocations.length > 0) {
		return detail.value.allocations[0].InvoiceNo || 'Invoice'
	}
	return ''
})

const remainingBalanceOverall = computed(() => {
	if (linkedInvoiceDetail.value?.invoice) {
		return linkedInvoiceDetail.value.invoice.BalanceMinor
	}
	return 0
})

const printedDateTime = computed(() => {
	const now = new Date()
	return now.toLocaleDateString(undefined, {
		year: 'numeric',
		month: 'short',
		day: 'numeric'
	})
})

function getInvoiceTotalDisplay(invoiceId: string): string {
	if (linkedInvoiceDetail.value?.invoice && linkedInvoiceDetail.value.invoice.ID === invoiceId) {
		return formatMinor(linkedInvoiceDetail.value.invoice.TotalAmountMinor)
	}
	return '—'
}

function getInvoiceBalance(invoiceId: string): number {
	if (linkedInvoiceDetail.value?.invoice && linkedInvoiceDetail.value.invoice.ID === invoiceId) {
		return linkedInvoiceDetail.value.invoice.BalanceMinor
	}
	return 0
}

function getInvoiceBalanceDisplay(invoiceId: string): string {
	if (linkedInvoiceDetail.value?.invoice && linkedInvoiceDetail.value.invoice.ID === invoiceId) {
		const bal = linkedInvoiceDetail.value.invoice.BalanceMinor
		return bal === 0 ? '0.00 (PAID)' : formatMinor(bal)
	}
	return '—'
}

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

			// 1. Fetch linked invoice details if allocated
			if (detail.value.allocations && detail.value.allocations.length > 0) {
				const primaryInvoiceId = detail.value.allocations[0].InvoiceID
				if (primaryInvoiceId) {
					try {
						const invRes = await services.invoices.getInvoice(primaryInvoiceId)
						if (invRes.success && invRes.data) {
							linkedInvoiceDetail.value = invRes.data
						}
					} catch (e) {
						console.error('Failed to load linked invoice details', e)
					}
				}
			}

			// 2. Fetch student enrollment details (class & section)
			if (detail.value.payment?.StudentID) {
				try {
					const enrRes = await services.enrollments.getActive(detail.value.payment.StudentID)
					if (enrRes.success && enrRes.data) {
						studentEnrollment.value = enrRes.data
					}
				} catch (e) {
					console.error('Failed to load student enrollment details', e)
				}
			}

			// 3. If no linked invoice items, fetch applicable class fee structure
			if (!linkedInvoiceDetail.value?.items?.length && studentEnrollment.value?.AcademicYearID && studentEnrollment.value?.ClassID) {
				try {
					const structRes = await services.fees.listFeeStructures(studentEnrollment.value.AcademicYearID)
					if (structRes.success && structRes.data) {
						const classStruct = structRes.data.find((s: any) => s.ClassID === studentEnrollment.value?.ClassID)
						if (classStruct?.ID) {
							const versRes = await services.fees.listFeeStructureVersions(classStruct.ID)
							if (versRes.success && versRes.data?.length) {
								const publishedVer = versRes.data.find((v: any) => (v.Status || v.status) === 'published') || versRes.data[0]
								if (publishedVer?.ID) {
									const verRes = await services.fees.getFeeStructureVersionDetail(publishedVer.ID)
									if (verRes.success && verRes.data?.items) {
										classFeeItems.value = verRes.data.items
									}
								}
							}
						}
					}
				} catch (e) {
					console.error('Failed to load class fee schedule', e)
				}
			}
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
	const originalTitle = document.title
	document.title = ''
	window.print()
	window.addEventListener('afterprint', () => {
		document.title = originalTitle
	}, { once: true })
	setTimeout(() => {
		document.title = originalTitle
	}, 1500)
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
			toast.success('Payment successfully reversed and marked as void.')
			await loadPayment()
		} else {
			toast.error(res.error?.message || 'Failed to reverse payment')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred while reversing payment')
	}
}
</script>

<style scoped>
@media print {
	@page {
		margin: 0 !important;
		size: auto;
	}

	.no-print {
		display: none !important;
	}

	.receipt-paper {
		border: none !important;
		box-shadow: none !important;
		padding: 12mm 14mm !important;
		margin: 0 auto !important;
		max-width: 100% !important;
		width: 100% !important;
		box-sizing: border-box !important;
	}
}
</style>
