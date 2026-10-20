<template>
	<div class="max-w-4xl mx-auto space-y-6">
		<!-- Navigation & Top Actions (hidden on print) -->
		<div class="flex items-center justify-between no-print">
			<router-link
				to="/invoices"
				class="inline-flex items-center gap-1.5 text-xs font-medium text-slate-500 hover:text-slate-800 transition"
			>
				<ArrowLeft class="w-4 h-4" />
				<span>Back to Invoices</span>
			</router-link>

			<div class="flex items-center gap-2">
				<button
					@click="printInvoice"
					class="py-2 px-3.5 bg-white border border-slate-200 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Printer class="w-4 h-4 text-slate-500" />
					<span>Print Statement</span>
				</button>
				<button
					v-if="detail?.invoice && detail.invoice.Status !== 'void'"
					@click="showVoidModal = true"
					class="py-2 px-3.5 bg-rose-50 border border-rose-200 hover:bg-rose-100 text-rose-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Ban class="w-4 h-4" />
					<span>Void Invoice</span>
				</button>
			</div>
		</div>

		<!-- Invoice Paper Container -->
		<div v-if="detail?.invoice" class="bg-white rounded-xl border border-slate-200 shadow-sm p-8 space-y-8 print:border-none print:shadow-none print:p-0">
			<!-- Header -->
			<div class="flex items-start justify-between border-b border-slate-100 pb-6">
				<div>
					<div class="text-xs font-bold uppercase tracking-wider text-indigo-600 mb-1">Official Student Bill</div>
					<h1 class="text-2xl font-black text-slate-900">{{ schoolProfile?.name || 'School Management System' }}</h1>
					<p class="text-xs text-slate-500 mt-0.5">{{ schoolProfile?.address?.line1 || 'Academic Campus & Financial Office' }}</p>
					<p v-if="schoolProfile?.phone" class="text-xs text-slate-500">Tel: {{ schoolProfile.phone }}</p>
				</div>

				<div class="text-right">
					<div class="text-xs font-bold uppercase tracking-wider text-slate-400">Invoice Number</div>
					<div class="text-xl font-black font-mono text-slate-900 mt-0.5">{{ detail.invoice.InvoiceNo }}</div>
					<div class="mt-2">
						<span
							class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wide inline-block"
							:class="getStatusBadgeClass(detail.invoice.Status)"
						>
							{{ detail.invoice.Status.replace('_', ' ') }}
						</span>
					</div>
				</div>
			</div>

			<!-- Void Alert Banner -->
			<div v-if="detail.invoice.Status === 'void'" class="p-3 bg-rose-50 border border-rose-200 rounded-lg text-xs text-rose-800">
				<span class="font-bold">VOIDED INVOICE:</span> This bill has been revoked and marked void.
				<span v-if="detail.invoice.VoidReason" class="italic ml-1">Reason: {{ detail.invoice.VoidReason }}</span>
			</div>

			<!-- Meta Details: Billed Student & Schedule Dates -->
			<div class="grid grid-cols-2 gap-6 bg-slate-50/70 p-4 rounded-xl border border-slate-100 text-xs">
				<div class="space-y-1">
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-400">Billed Student</div>
					<div class="font-bold text-slate-900 text-sm">
						{{ detail.invoice.StudentFirstName }} {{ detail.invoice.StudentLastName }}
					</div>
					<div class="text-slate-600 font-mono">Admission #: {{ detail.invoice.AdmissionNo }}</div>
				</div>

				<div class="space-y-1 text-right">
					<div class="text-[10px] font-bold uppercase tracking-wider text-slate-400">Billing Information</div>
					<div class="text-slate-700">Billing Period: <span class="font-bold font-mono text-slate-900">{{ detail.invoice.BillingPeriod }}</span></div>
					<div class="text-slate-600">Issue Date: <span class="font-medium text-slate-900">{{ formatCivilDate(detail.invoice.IssueDate) }}</span></div>
					<div class="text-slate-600">Due Date: <span class="font-bold text-rose-700">{{ formatCivilDate(detail.invoice.DueDate) }}</span></div>
				</div>
			</div>

			<!-- Line Items Table -->
			<div class="space-y-2">
				<div class="text-xs font-bold uppercase tracking-wider text-slate-500">Fee Item Breakdown</div>
				<div class="border border-slate-200 rounded-xl overflow-hidden">
					<table class="w-full text-left text-xs border-collapse">
						<thead>
							<tr class="bg-slate-50 border-b border-slate-200 text-[10px] font-bold uppercase tracking-wider text-slate-500">
								<th class="py-3 px-4">Fee Head / Description</th>
								<th class="py-3 px-4 text-right">Schedule Amount</th>
								<th class="py-3 px-4 text-right">Concession / Discount</th>
								<th class="py-3 px-4 text-right">Net Amount</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100 text-slate-700">
							<tr v-for="item in detail.items" :key="item.ID" class="hover:bg-slate-50/30">
								<td class="py-3 px-4">
									<div class="font-semibold text-slate-900">{{ item.FeeHeadName }}</div>
									<div v-if="item.Description" class="text-[11px] text-slate-400">{{ item.Description }}</div>
								</td>
								<td class="py-3 px-4 text-right font-medium text-slate-700">
									{{ formatMinor(item.OriginalAmountMinor) }}
								</td>
								<td class="py-3 px-4 text-right">
									<span v-if="item.DiscountAmountMinor > 0" class="text-emerald-600 font-medium">
										-{{ formatMinor(item.DiscountAmountMinor) }}
									</span>
									<span v-else class="text-slate-400">—</span>
								</td>
								<td class="py-3 px-4 text-right font-bold text-slate-900">
									{{ formatMinor(item.NetAmountMinor) }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>

			<!-- Financial Summary & Balance -->
			<div class="flex justify-end pt-2">
				<div class="w-72 space-y-2 text-xs">
					<div class="flex justify-between py-1 text-slate-600 border-b border-slate-100">
						<span>Subtotal:</span>
						<span class="font-medium text-slate-900">{{ formatMinor(detail.invoice.SubtotalAmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-1 text-emerald-600 border-b border-slate-100">
						<span>Total Concessions:</span>
						<span class="font-bold">-{{ formatMinor(detail.invoice.DiscountAmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-1 text-slate-900 font-bold border-b border-slate-200 text-sm">
						<span>Total Due:</span>
						<span class="text-indigo-700">{{ formatMinor(detail.invoice.TotalAmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-1 text-slate-600 border-b border-slate-100">
						<span>Paid to Date:</span>
						<span class="font-medium text-slate-900">{{ formatMinor(detail.invoice.PaidAmountMinor) }}</span>
					</div>
					<div class="flex justify-between py-2 text-slate-900 font-black text-sm bg-slate-50 px-2 rounded-lg">
						<span>Balance Remaining:</span>
						<span :class="detail.invoice.BalanceMinor > 0 && detail.invoice.Status !== 'void' ? 'text-rose-600' : 'text-slate-900'">
							{{ formatMinor(detail.invoice.BalanceMinor) }}
						</span>
					</div>
				</div>
			</div>

			<!-- Payment Instructions / Footer -->
			<div class="border-t border-slate-100 pt-6 text-[11px] text-slate-400 space-y-1">
				<p class="font-semibold text-slate-600">Payment Instructions</p>
				<p>Please present this invoice at the cashier's office or deposit directly into the institution's designated bank account. Ensure the Invoice Number ({{ detail.invoice.InvoiceNo }}) is quoted as the payment reference.</p>
				<p class="pt-4 text-[10px] text-slate-300">Generated automatically by School Management System. This is a computer-generated statement.</p>
			</div>
		</div>

		<div v-else-if="loading" class="p-12 text-center text-xs text-slate-400">
			Loading invoice details...
		</div>

		<!-- MODAL: Void Invoice -->
		<div v-if="showVoidModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50 no-print">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-sm w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Void Invoice</h3>
					<button @click="showVoidModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitVoid" class="p-4 space-y-3">
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason for Voiding *</label>
						<textarea
							v-model="voidReason"
							required
							rows="3"
							placeholder="State justification..."
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
import { useRoute } from 'vue-router'
import { getServices } from '@/services/factory'
import type { InvoiceDetailDTO, SchoolProfileDTO } from '@/types/common'
import { formatMinor, formatCivilDate } from '@/utils/format'
import { ArrowLeft, Printer, Ban } from 'lucide-vue-next'

const route = useRoute()
const services = getServices()

const detail = ref<InvoiceDetailDTO | null>(null)
const schoolProfile = ref<SchoolProfileDTO | null>(null)
const loading = ref(false)

const showVoidModal = ref(false)
const voidReason = ref('')

onMounted(async () => {
	await Promise.all([loadInvoice(), loadSchoolProfile()])
})

async function loadInvoice() {
	const id = route.params.id as string
	if (!id) return
	loading.value = true
	try {
		const res = await services.invoices.getInvoice(id)
		if (res.success && res.data) {
			detail.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load invoice', e)
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

function printInvoice() {
	window.print()
}

async function submitVoid() {
	if (!detail.value?.invoice) return
	try {
		const res = await services.invoices.voidInvoice({
			invoiceId: detail.value.invoice.ID,
			reason: voidReason.value
		})
		if (res.success) {
			showVoidModal.value = false
			await loadInvoice()
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

<style scoped>
@media print {
	.no-print {
		display: none !important;
	}
}
</style>
