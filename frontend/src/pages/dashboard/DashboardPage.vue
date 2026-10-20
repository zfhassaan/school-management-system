<template>
	<div class="space-y-6 max-w-7xl mx-auto">
		<!-- Welcome Banner -->
		<div class="bg-gradient-to-r from-indigo-800 via-indigo-900 to-slate-900 rounded-2xl p-6 text-white shadow-md relative overflow-hidden">
			<div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-4">
				<div>
					<div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-white/10 text-indigo-100 text-xs font-medium mb-2.5 backdrop-blur-xs">
						<span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
						<span>Offline-First Desktop System Active</span>
						<span v-if="activeYearName" class="text-indigo-200">• Session: {{ activeYearName }}</span>
					</div>
					<h1 class="text-2xl font-bold tracking-tight">Welcome, {{ authStore.user?.displayName || 'Administrator' }}</h1>
					<p class="text-xs text-indigo-200 mt-1 max-w-xl">
						Managing <span class="font-semibold text-white">{{ authStore.school?.name || 'School System' }}</span>
						<span v-if="authStore.school?.code"> (Code: {{ authStore.school.code }})</span>.
						All student records, tuition invoices, and receipts are safely encrypted and stored locally.
					</p>
				</div>

				<div class="flex items-center gap-2.5 shrink-0">
					<router-link
						to="/payments"
						class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-emerald-500 hover:bg-emerald-400 text-white flex items-center gap-2 transition shadow-xs cursor-pointer"
					>
						<CreditCard class="w-4 h-4" />
						<span>Record Payment</span>
					</router-link>
					<router-link
						to="/admissions"
						class="px-3.5 py-2 rounded-xl text-xs font-semibold bg-white/15 hover:bg-white/25 text-white flex items-center gap-2 transition border border-white/10 cursor-pointer"
					>
						<UserPlus class="w-4 h-4" />
						<span>New Admission</span>
					</router-link>
				</div>
			</div>
		</div>

		<!-- Key School Operations & Financial Metrics -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
			<!-- Total Students -->
			<router-link
				to="/students"
				class="bg-white p-5 rounded-2xl border border-slate-200/90 shadow-xs hover:border-indigo-300 hover:shadow-md transition-all group block"
			>
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-[11px] font-bold uppercase tracking-wider text-slate-500">Enrolled Students</span>
					<div class="w-9 h-9 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center group-hover:scale-110 transition-transform">
						<GraduationCap class="w-5 h-5" />
					</div>
				</div>
				<div class="text-2xl font-black text-slate-900 tracking-tight">
					<span v-if="loading" class="animate-pulse text-slate-300">...</span>
					<span v-else>{{ studentsCount }}</span>
				</div>
				<p class="text-[11px] text-slate-500 mt-1 flex items-center gap-1.5">
					<span class="inline-block w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
					<span>Active across {{ classesCount }} grade classes</span>
				</p>
			</router-link>

			<!-- Admissions Pipeline -->
			<router-link
				to="/admissions"
				class="bg-white p-5 rounded-2xl border border-slate-200/90 shadow-xs hover:border-indigo-300 hover:shadow-md transition-all group block"
			>
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-[11px] font-bold uppercase tracking-wider text-slate-500">Admissions Pipeline</span>
					<div class="w-9 h-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center group-hover:scale-110 transition-transform">
						<FileText class="w-5 h-5" />
					</div>
				</div>
				<div class="text-2xl font-black text-slate-900 tracking-tight">
					<span v-if="loading" class="animate-pulse text-slate-300">...</span>
					<span v-else>{{ applicationsCount }}</span>
					<span class="text-xs font-medium text-slate-500 ml-1.5">Applications</span>
				</div>
				<p class="text-[11px] text-amber-700 mt-1 font-medium">
					{{ submittedAppsCount }} pending review • {{ approvedAppsCount }} ready to enroll
				</p>
			</router-link>

			<!-- Total Collections (PKR) -->
			<router-link
				to="/payments"
				class="bg-white p-5 rounded-2xl border border-slate-200/90 shadow-xs hover:border-emerald-300 hover:shadow-md transition-all group block"
			>
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-[11px] font-bold uppercase tracking-wider text-slate-500">Fees Collected</span>
					<div class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center group-hover:scale-110 transition-transform">
						<Coins class="w-5 h-5" />
					</div>
				</div>
				<div class="text-2xl font-black text-emerald-700 tracking-tight">
					<span v-if="loading" class="animate-pulse text-slate-300">...</span>
					<span v-else>{{ formatMinor(totalCollectedMinor) }}</span>
				</div>
				<p class="text-[11px] text-slate-500 mt-1">
					{{ paymentsCount }} receipts issued this session
				</p>
			</router-link>

			<!-- Outstanding Receivables (PKR) -->
			<router-link
				to="/invoices"
				class="bg-white p-5 rounded-2xl border border-slate-200/90 shadow-xs hover:border-rose-300 hover:shadow-md transition-all group block"
			>
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-[11px] font-bold uppercase tracking-wider text-slate-500">Receivables Due</span>
					<div class="w-9 h-9 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center group-hover:scale-110 transition-transform">
						<Receipt class="w-5 h-5" />
					</div>
				</div>
				<div class="text-2xl font-black text-slate-900 tracking-tight">
					<span v-if="loading" class="animate-pulse text-slate-300">...</span>
					<span v-else>{{ formatMinor(totalOutstandingMinor) }}</span>
				</div>
				<p class="text-[11px] text-rose-600 mt-1 font-medium">
					From {{ invoicesCount }} student billing invoices
				</p>
			</router-link>
		</div>

		<!-- Activity & Pipeline Overview Grid -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Recent Receipts & Collections -->
			<div class="bg-white rounded-2xl border border-slate-200/90 shadow-xs p-5 space-y-4">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-2.5">
						<div class="w-8 h-8 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center">
							<Receipt class="w-4 h-4" />
						</div>
						<div>
							<h3 class="text-sm font-bold text-slate-900">Recent Fee Collections</h3>
							<p class="text-[11px] text-slate-500">Latest cashier receipts and bank deposits</p>
						</div>
					</div>
					<router-link
						to="/payments"
						class="text-xs font-semibold text-indigo-600 hover:text-indigo-800 flex items-center gap-1 transition"
					>
						<span>View All</span>
						<ArrowRight class="w-3.5 h-3.5" />
					</router-link>
				</div>

				<div class="overflow-x-auto">
					<table class="w-full text-left text-xs">
						<thead>
							<tr class="border-b border-slate-100 text-slate-400 font-semibold text-[10px] uppercase">
								<th class="pb-2">Receipt #</th>
								<th class="pb-2">Student</th>
								<th class="pb-2">Method</th>
								<th class="pb-2 text-right">Amount (PKR)</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100 text-slate-700">
							<tr v-if="recentPayments.length === 0">
								<td colspan="4" class="py-6 text-center text-slate-400 text-xs">No fee payments recorded yet.</td>
							</tr>
							<tr v-for="p in recentPayments" :key="p.ID" class="hover:bg-slate-50/60">
								<td class="py-2.5 font-mono font-medium text-slate-900">{{ p.ReceiptNo }}</td>
								<td class="py-2.5 font-semibold text-slate-900">
									{{ p.StudentFirstName ? `${p.StudentFirstName} ${p.StudentLastName || ''}` : p.PayerName || 'Student' }}
								</td>
								<td class="py-2.5">
									<span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-slate-100 text-slate-700 capitalize">
										{{ p.PaymentMethod.replace('_', ' ') }}
									</span>
								</td>
								<td class="py-2.5 text-right font-mono font-bold text-emerald-700">
									{{ formatMinor(p.AmountMinor) }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>

			<!-- Recent Admission Applications -->
			<div class="bg-white rounded-2xl border border-slate-200/90 shadow-xs p-5 space-y-4">
				<div class="flex items-center justify-between">
					<div class="flex items-center gap-2.5">
						<div class="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center">
							<UserPlus class="w-4 h-4" />
						</div>
						<div>
							<h3 class="text-sm font-bold text-slate-900">Recent Applications</h3>
							<p class="text-[11px] text-slate-500">New student applicants awaiting review</p>
						</div>
					</div>
					<router-link
						to="/admissions"
						class="text-xs font-semibold text-indigo-600 hover:text-indigo-800 flex items-center gap-1 transition"
					>
						<span>View All</span>
						<ArrowRight class="w-3.5 h-3.5" />
					</router-link>
				</div>

				<div class="overflow-x-auto">
					<table class="w-full text-left text-xs">
						<thead>
							<tr class="border-b border-slate-100 text-slate-400 font-semibold text-[10px] uppercase">
								<th class="pb-2">Application #</th>
								<th class="pb-2">Applicant</th>
								<th class="pb-2">Class</th>
								<th class="pb-2 text-right">Status</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100 text-slate-700">
							<tr v-if="recentApplications.length === 0">
								<td colspan="4" class="py-6 text-center text-slate-400 text-xs">No admissions submitted yet.</td>
							</tr>
							<tr v-for="app in recentApplications" :key="app.ID" class="hover:bg-slate-50/60">
								<td class="py-2.5 font-mono font-medium text-slate-900">{{ app.ApplicationNo }}</td>
								<td class="py-2.5 font-semibold text-slate-900">
									{{ app.ApplicantFirstName }} {{ app.ApplicantLastName }}
								</td>
								<td class="py-2.5 text-slate-600 font-medium">
									{{ app.ClassName || 'General' }}
								</td>
								<td class="py-2.5 text-right">
									<span
										class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
										:class="getAppBadgeClass(app.Status)"
									>
										{{ app.Status.replace('_', ' ') }}
									</span>
								</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
		</div>

		<!-- Quick Operational Workflows -->
		<div class="space-y-3">
			<h3 class="text-xs font-bold uppercase tracking-wider text-slate-500">Quick Administrative Workflows</h3>
			<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
				<router-link
					to="/admissions"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-3.5"
				>
					<div class="w-10 h-10 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<UserPlus class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-xs font-bold text-slate-900 group-hover:text-indigo-600 transition">Admissions & Intake</h4>
						<p class="text-[11px] text-slate-500 mt-0.5">Review candidate forms, conduct interviews, and enroll students.</p>
					</div>
				</router-link>

				<router-link
					to="/invoices"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-3.5"
				>
					<div class="w-10 h-10 rounded-xl bg-purple-50 text-purple-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<Receipt class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-xs font-bold text-slate-900 group-hover:text-indigo-600 transition">Monthly Billing & Invoices</h4>
						<p class="text-[11px] text-slate-500 mt-0.5">Issue class-wide fee billing periods with automated sibling discounts.</p>
					</div>
				</router-link>

				<router-link
					to="/payments"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-3.5"
				>
					<div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<CreditCard class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-xs font-bold text-slate-900 group-hover:text-indigo-600 transition">Fee Collections & Receipts</h4>
						<p class="text-[11px] text-slate-500 mt-0.5">Post cashier receipts, bank wire references, and FIFO allocations.</p>
					</div>
				</router-link>

				<router-link
					to="/students"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-3.5"
				>
					<div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<GraduationCap class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-xs font-bold text-slate-900 group-hover:text-indigo-600 transition">Students Directory</h4>
						<p class="text-[11px] text-slate-500 mt-0.5">Manage student master profiles, guardians, and academic history.</p>
					</div>
				</router-link>

				<router-link
					to="/fees/structures"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-3.5"
				>
					<div class="w-10 h-10 rounded-xl bg-sky-50 text-sky-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<DollarSign class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-xs font-bold text-slate-900 group-hover:text-indigo-600 transition">Fee Schedules & Concessions</h4>
						<p class="text-[11px] text-slate-500 mt-0.5">Define recurring tuition heads, discount policies, and fee structures.</p>
					</div>
				</router-link>

				<router-link
					to="/settings/school"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-3.5"
				>
					<div class="w-10 h-10 rounded-xl bg-slate-100 text-slate-700 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<Building2 class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-xs font-bold text-slate-900 group-hover:text-indigo-600 transition">School Settings</h4>
						<p class="text-[11px] text-slate-500 mt-0.5">Configure school legal name, code, timezones, and currency options.</p>
					</div>
				</router-link>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getServices } from '@/services/factory'
import { formatMinor } from '@/utils/format'
import type { PaymentDTO, ApplicationDTO, AcademicYearDTO, InvoiceDTO } from '@/types/common'
import {
	GraduationCap,
	FileText,
	Receipt,
	Coins,
	CreditCard,
	UserPlus,
	Building2,
	DollarSign,
	ArrowRight
} from 'lucide-vue-next'

const authStore = useAuthStore()
const services = getServices()

const loading = ref(true)

const studentsCount = ref(0)
const applicationsCount = ref(0)
const submittedAppsCount = ref(0)
const approvedAppsCount = ref(0)

const totalInvoicedMinor = ref(0)
const totalCollectedMinor = ref(0)
const totalOutstandingMinor = ref(0)
const invoicesCount = ref(0)
const paymentsCount = ref(0)

const activeYearName = ref('')
const classesCount = ref(0)

const recentPayments = ref<PaymentDTO[]>([])
const recentApplications = ref<ApplicationDTO[]>([])

function getAppBadgeClass(status: string) {
	switch (status) {
		case 'submitted':
			return 'bg-blue-50 text-blue-700 border border-blue-200'
		case 'under_review':
			return 'bg-amber-50 text-amber-700 border border-amber-200'
		case 'approved':
			return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
		case 'enrolled':
			return 'bg-purple-50 text-purple-700 border border-purple-200'
		default:
			return 'bg-slate-100 text-slate-600 border border-slate-200'
	}
}

async function loadDashboardData() {
	loading.value = true
	try {
		const [
			studentsRes,
			appsRes,
			invoicesRes,
			paymentsRes,
			yearsRes,
			classesRes
		] = await Promise.all([
			services.students.listStudents({ limit: 100 }),
			services.admissions.listApplications('', '', 10),
			services.invoices.listInvoices({ limit: 100 }),
			services.payments.listPayments({ limit: 10 }),
			services.academics.listAcademicYears(),
			services.academics.listClasses()
		])

		if (studentsRes.success && studentsRes.data) {
			studentsCount.value = studentsRes.data.length
		}

		if (appsRes.success && appsRes.data) {
			applicationsCount.value = appsRes.data.length
			submittedAppsCount.value = appsRes.data.filter((a: ApplicationDTO) => a.Status === 'submitted' || a.Status === 'under_review').length
			approvedAppsCount.value = appsRes.data.filter((a: ApplicationDTO) => a.Status === 'approved').length
			recentApplications.value = appsRes.data.slice(0, 5)
		}

		if (invoicesRes.success && invoicesRes.data) {
			invoicesCount.value = invoicesRes.data.length
			let invTotal = 0
			let outTotal = 0
			for (const inv of (invoicesRes.data as InvoiceDTO[])) {
				if (inv.Status !== 'void') {
					invTotal += inv.TotalAmountMinor || 0
					outTotal += inv.BalanceMinor || 0
				}
			}
			totalInvoicedMinor.value = invTotal
			totalOutstandingMinor.value = outTotal
		}

		if (paymentsRes.success && paymentsRes.data) {
			paymentsCount.value = paymentsRes.data.length
			let colTotal = 0
			for (const p of (paymentsRes.data as PaymentDTO[])) {
				if (p.Status === 'posted') {
					colTotal += p.AmountMinor || 0
				}
			}
			totalCollectedMinor.value = colTotal
			recentPayments.value = paymentsRes.data.slice(0, 5)
		}

		if (yearsRes.success && yearsRes.data) {
			const active = (yearsRes.data as AcademicYearDTO[]).find((y: AcademicYearDTO) => y.Status === 'active') || yearsRes.data[0]
			if (active) activeYearName.value = `${active.Name} (${active.Code})`
		}

		if (classesRes.success && classesRes.data) {
			classesCount.value = classesRes.data.length
		}
	} catch (e) {
		console.warn('Dashboard metrics fetch error:', e)
	} finally {
		loading.value = false
	}
}

onMounted(() => {
	loadDashboardData()
})
</script>
