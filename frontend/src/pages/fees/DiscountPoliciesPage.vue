<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Discount & Concession Management</h1>
				<p class="text-xs text-slate-500 mt-0.5">Configure merit, sibling, staff, and financial aid rules, and review pending approval workflows.</p>
			</div>
			<div class="flex gap-2">
				<button
					v-if="currentTab === 'policies'"
					@click="openAddPolicyModal"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>Create Policy</span>
				</button>
			</div>
		</div>

		<!-- Tabs -->
		<div class="flex border-b border-slate-200 space-x-6 text-xs font-medium">
			<button
				v-for="tab in ['policies', 'approvals', 'assignment']"
				:key="tab"
				@click="currentTab = tab"
				class="pb-3 capitalize cursor-pointer transition border-b-2 flex items-center gap-2"
				:class="currentTab === tab ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'"
			>
				<span>{{ tab === 'policies' ? 'Discount Policies' : tab === 'approvals' ? 'Pending Approvals' : 'Assign Concessions' }}</span>
				<span
					v-if="tab === 'approvals' && pendingDiscounts.length > 0"
					class="px-1.5 py-0.2 rounded-full text-[10px] bg-rose-100 text-rose-700 font-bold"
				>
					{{ pendingDiscounts.length }}
				</span>
			</button>
		</div>

		<!-- TAB 1: Policies Catalog -->
		<div v-if="currentTab === 'policies'" class="space-y-4">
			<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
				<table class="w-full text-left border-collapse text-xs">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
							<th class="py-3 px-4">Code</th>
							<th class="py-3 px-4">Policy Name</th>
							<th class="py-3 px-4">Type</th>
							<th class="py-3 px-4">Benefit Rate / Amount</th>
							<th class="py-3 px-4">Target Fee Head</th>
							<th class="py-3 px-4">Workflow</th>
							<th class="py-3 px-4">Status</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 text-slate-700">
						<tr v-for="p in policies" :key="p.ID" class="hover:bg-slate-50/50">
							<td class="py-3 px-4 font-mono font-medium text-slate-900">{{ p.Code }}</td>
							<td class="py-3 px-4">
								<div class="font-semibold text-slate-900">{{ p.Name }}</div>
								<div v-if="p.Description" class="text-[11px] text-slate-400 truncate max-w-xs">{{ p.Description }}</div>
							</td>
							<td class="py-3 px-4">
								<span
									class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
									:class="p.DiscountType === 'percentage' ? 'bg-purple-50 text-purple-700 border border-purple-200' : 'bg-emerald-50 text-emerald-700 border border-emerald-200'"
								>
									{{ p.DiscountType }}
								</span>
							</td>
							<td class="py-3 px-4 font-bold text-slate-900">
								<span v-if="p.DiscountType === 'percentage'">{{ formatBasisPoints(p.RateBasisPoints) }}</span>
								<span v-else>{{ formatMinor(p.AmountMinor) }}</span>
							</td>
							<td class="py-3 px-4 text-slate-600">
								{{ p.FeeHeadID ? getFeeHeadName(p.FeeHeadID) : 'All Applicable Fees' }}
							</td>
							<td class="py-3 px-4">
								<span
									class="px-2 py-0.5 rounded-full text-[10px] font-semibold"
									:class="p.RequiresApproval ? 'bg-amber-50 text-amber-700 border border-amber-200' : 'bg-blue-50 text-blue-700 border border-blue-200'"
								>
									{{ p.RequiresApproval ? 'Requires Approval' : 'Auto-approved' }}
								</span>
							</td>
							<td class="py-3 px-4">
								<span
									class="px-2 py-0.5 rounded-full text-[10px] font-semibold"
									:class="p.IsActive ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-slate-100 text-slate-500 border border-slate-200'"
								>
									{{ p.IsActive ? 'Active' : 'Inactive' }}
								</span>
							</td>
						</tr>
						<tr v-if="policies.length === 0">
							<td colspan="7" class="py-8 text-center text-slate-400 italic">No discount policies configured. Click "Create Policy" to add one.</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- TAB 2: Pending Approvals Queue -->
		<div v-if="currentTab === 'approvals'" class="space-y-4">
			<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
				<table class="w-full text-left border-collapse text-xs">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
							<th class="py-3 px-4">Student</th>
							<th class="py-3 px-4">Admission #</th>
							<th class="py-3 px-4">Policy Code & Name</th>
							<th class="py-3 px-4">Benefit</th>
							<th class="py-3 px-4">Effective Dates</th>
							<th class="py-3 px-4 text-right">Approval Decision</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 text-slate-700">
						<tr v-for="d in pendingDiscounts" :key="d.ID" class="hover:bg-slate-50/50">
							<td class="py-3 px-4 font-semibold text-slate-900">
								{{ d.StudentFirstName }} {{ d.StudentLastName }}
							</td>
							<td class="py-3 px-4 font-mono text-slate-600">{{ d.AdmissionNo }}</td>
							<td class="py-3 px-4">
								<div class="font-medium text-slate-800">{{ d.PolicyName || 'Policy ' + d.DiscountPolicyID }}</div>
								<div class="font-mono text-[10px] text-slate-400">{{ d.PolicyCode }}</div>
							</td>
							<td class="py-3 px-4 font-bold text-slate-900">
								<span v-if="d.DiscountType === 'percentage'">{{ formatBasisPoints(d.RateBasisPoints || 0) }}</span>
								<span v-else>{{ formatMinor(d.PolicyAmountMinor || 0) }}</span>
							</td>
							<td class="py-3 px-4 text-slate-600">
								From: {{ formatCivilDate(d.ValidFrom) }}
								<span v-if="d.ValidTo" class="block text-[10px] text-slate-400">To: {{ formatCivilDate(d.ValidTo) }}</span>
							</td>
							<td class="py-3 px-4 text-right space-x-2">
								<button
									@click="decideDiscount(d.ID, true)"
									class="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded font-medium text-xs cursor-pointer"
								>
									Approve
								</button>
								<button
									@click="openRejectModal(d.ID)"
									class="px-2.5 py-1 bg-rose-50 hover:bg-rose-100 text-rose-700 border border-rose-200 rounded font-medium text-xs cursor-pointer"
								>
									Reject
								</button>
							</td>
						</tr>
						<tr v-if="pendingDiscounts.length === 0">
							<td colspan="6" class="py-8 text-center text-slate-400 italic">No pending discount approval requests in queue.</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- TAB 3: Assign Concession to Student -->
		<div v-if="currentTab === 'assignment'" class="grid grid-cols-12 gap-5">
			<!-- Left: Assignment Form -->
			<div class="col-span-5 bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
				<h3 class="font-bold text-sm text-slate-900">Grant Student Concession</h3>
				<form @submit.prevent="submitAssignForm" class="space-y-3">
					<div v-if="assignError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ assignError }}
					</div>
					<div v-if="assignSuccess" class="p-2.5 bg-emerald-50 text-emerald-700 border border-emerald-200 rounded-lg text-xs">
						{{ assignSuccess }}
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Select Student *</label>
						<select
							v-model="assignForm.studentId"
							@change="loadStudentDiscounts"
							required
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="" disabled>Select Student</option>
							<option v-for="s in studentsList" :key="s.ID" :value="s.ID">
								{{ s.FirstName }} {{ s.LastName }} ({{ s.AdmissionNo }})
							</option>
						</select>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Academic Year *</label>
						<select
							v-model="assignForm.academicYearId"
							@change="loadStudentDiscounts"
							required
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
						</select>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Discount Policy *</label>
						<select
							v-model="assignForm.discountPolicyId"
							required
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="" disabled>Select Policy</option>
							<option v-for="p in policies.filter(p => p.IsActive)" :key="p.ID" :value="p.ID">
								{{ p.Name }} ({{ p.Code }}) — {{ p.DiscountType === 'percentage' ? formatBasisPoints(p.RateBasisPoints) : formatMinor(p.AmountMinor) }}
							</option>
						</select>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Valid From *</label>
							<input
								v-model="assignForm.validFrom"
								type="date"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Valid To (Optional)</label>
							<input
								v-model="assignForm.validTo"
								type="date"
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
					</div>

					<div class="pt-2">
						<button
							type="submit"
							:disabled="assigning"
							class="w-full py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition cursor-pointer"
						>
							{{ assigning ? 'Applying...' : 'Apply Student Concession' }}
						</button>
					</div>
				</form>
			</div>

			<!-- Right: Student's Active & Past Concessions -->
			<div class="col-span-7 bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
				<h3 class="font-bold text-sm text-slate-900">Student's Concession Record</h3>
				<div v-if="assignForm.studentId" class="space-y-3">
					<table class="w-full text-left border-collapse text-xs">
						<thead>
							<tr class="bg-slate-50 border-b border-slate-200 text-[10px] font-bold uppercase tracking-wider text-slate-500">
								<th class="py-2.5 px-3">Policy</th>
								<th class="py-2.5 px-3">Benefit</th>
								<th class="py-2.5 px-3">Status</th>
								<th class="py-2.5 px-3">Validity</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-100">
							<tr v-for="sd in studentDiscounts" :key="sd.ID" class="hover:bg-slate-50/50">
								<td class="py-2.5 px-3 font-medium text-slate-800">
									{{ sd.PolicyName || 'Policy ' + sd.DiscountPolicyID }}
									<span class="font-mono text-[10px] text-slate-400 block">{{ sd.PolicyCode }}</span>
								</td>
								<td class="py-2.5 px-3 font-bold text-slate-900">
									<span v-if="sd.DiscountType === 'percentage'">{{ formatBasisPoints(sd.RateBasisPoints || 0) }}</span>
									<span v-else>{{ formatMinor(sd.PolicyAmountMinor || 0) }}</span>
								</td>
								<td class="py-2.5 px-3">
									<span
										class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
										:class="
											sd.Status === 'approved'
												? 'bg-emerald-50 text-emerald-700 border border-emerald-200'
												: sd.Status === 'pending'
												? 'bg-amber-50 text-amber-700 border border-amber-200'
												: 'bg-rose-50 text-rose-700 border border-rose-200'
										"
									>
										{{ sd.Status }}
									</span>
								</td>
								<td class="py-2.5 px-3 text-slate-600">
									{{ formatCivilDate(sd.ValidFrom) }}
									<span v-if="sd.ValidTo" class="text-slate-400"> → {{ formatCivilDate(sd.ValidTo) }}</span>
								</td>
							</tr>
							<tr v-if="studentDiscounts.length === 0">
								<td colspan="4" class="py-6 text-center text-slate-400 italic">No concession policies assigned to this student for this academic year.</td>
							</tr>
						</tbody>
					</table>
				</div>
				<div v-else class="text-xs text-slate-400 italic text-center py-10">
					Select a student on the left to inspect their concession history.
				</div>
			</div>
		</div>

		<!-- MODAL: Create Discount Policy -->
		<div v-if="showPolicyModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-md w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Create Discount Policy</h3>
					<button @click="showPolicyModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitPolicyForm" class="p-4 space-y-3">
					<div v-if="policyModalError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ policyModalError }}
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Policy Code *</label>
						<input
							v-model="policyForm.code"
							type="text"
							required
							placeholder="e.g. SIBLING_20, MERIT_50, STAFF_100"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Policy Name *</label>
						<input
							v-model="policyForm.name"
							type="text"
							required
							placeholder="e.g. Sibling Concession 20%"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Description</label>
						<textarea
							v-model="policyForm.description"
							rows="2"
							placeholder="Eligibility criteria..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Discount Type *</label>
							<select
								v-model="policyForm.discountType"
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							>
								<option value="percentage">Percentage (%)</option>
								<option value="fixed">Fixed Amount ($)</option>
							</select>
						</div>
						<div v-if="policyForm.discountType === 'percentage'">
							<label class="block text-xs font-medium text-slate-700 mb-1">Percentage (%) *</label>
							<input
								v-model.number="policyForm.ratePercent"
								type="number"
								min="0.01"
								max="100"
								step="0.01"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
						<div v-else>
							<label class="block text-xs font-medium text-slate-700 mb-1">Fixed Amount ($) *</label>
							<input
								v-model.number="policyForm.amountDecimal"
								type="number"
								min="0.01"
								step="0.01"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Target Fee Head (Optional)</label>
						<select
							v-model="policyForm.feeHeadId"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="">All Applicable Fees</option>
							<option v-for="h in feeHeads" :key="h.ID" :value="h.ID">{{ h.Name }} ({{ h.Code }})</option>
						</select>
					</div>

					<div class="pt-2 flex items-center gap-2">
						<input id="requiresApproval" v-model="policyForm.requiresApproval" type="checkbox" class="rounded text-indigo-600 focus:ring-indigo-500/20" />
						<label for="requiresApproval" class="text-xs text-slate-700 select-none">Requires Principal/Admin Approval Workflow</label>
					</div>

					<div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showPolicyModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="savingPolicy"
							class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							{{ savingPolicy ? 'Saving...' : 'Save Policy' }}
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- MODAL: Rejection Reason -->
		<div v-if="showRejectModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-sm w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Reject Concession Request</h3>
					<button @click="showRejectModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitReject" class="p-4 space-y-3">
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason for Rejection *</label>
						<textarea
							v-model="rejectionReason"
							required
							rows="3"
							placeholder="State the justification for denying this concession..."
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
import { ref, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import type {
	DiscountPolicyDTO,
	StudentDiscountDTO,
	FeeHeadDTO,
	AcademicYearDTO,
	StudentDTO
} from '@/types/common'
import { formatMinor, formatBasisPoints, formatCivilDate } from '@/utils/format'
import { Plus } from 'lucide-vue-next'

const services = getServices()

const currentTab = ref('policies')

const policies = ref<DiscountPolicyDTO[]>([])
const pendingDiscounts = ref<StudentDiscountDTO[]>([])
const feeHeads = ref<FeeHeadDTO[]>([])
const academicYears = ref<AcademicYearDTO[]>([])
const studentsList = ref<StudentDTO[]>([])
const studentDiscounts = ref<StudentDiscountDTO[]>([])

// Create Policy Modal
const showPolicyModal = ref(false)
const savingPolicy = ref(false)
const policyModalError = ref('')
const policyForm = ref({
	code: '',
	name: '',
	description: '',
	discountType: 'percentage',
	ratePercent: 20,
	amountDecimal: 0,
	feeHeadId: '',
	requiresApproval: true
})

// Reject Modal
const showRejectModal = ref(false)
const rejectingDiscountId = ref('')
const rejectionReason = ref('')

// Concession Assignment Form
const assigning = ref(false)
const assignError = ref('')
const assignSuccess = ref('')
const assignForm = ref({
	studentId: '',
	academicYearId: '',
	discountPolicyId: '',
	validFrom: new Date().toISOString().split('T')[0],
	validTo: ''
})

onMounted(async () => {
	await Promise.all([
		loadPolicies(),
		loadPendingDiscounts(),
		loadFeeHeads(),
		loadAcademicYears(),
		loadStudents()
	])
})

async function loadPolicies() {
	try {
		const res = await services.fees.listDiscountPolicies()
		if (res.success && res.data) {
			policies.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load discount policies', e)
	}
}

async function loadPendingDiscounts() {
	try {
		const res = await services.fees.listPendingDiscounts()
		if (res.success && res.data) {
			pendingDiscounts.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load pending discounts', e)
	}
}

async function loadFeeHeads() {
	try {
		const res = await services.fees.listFeeHeads()
		if (res.success && res.data) {
			feeHeads.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load fee heads', e)
	}
}

async function loadAcademicYears() {
	try {
		const res = await services.academics.listAcademicYears()
		if (res.success && res.data) {
			academicYears.value = res.data
			const active = res.data.find(y => y.Status === 'active')
			if (active) assignForm.value.academicYearId = active.ID
			else if (res.data.length > 0) assignForm.value.academicYearId = res.data[0].ID
		}
	} catch (e: any) {
		console.error('Failed to load academic years', e)
	}
}

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

async function loadStudentDiscounts() {
	if (!assignForm.value.studentId || !assignForm.value.academicYearId) return
	try {
		const res = await services.fees.listStudentDiscounts(
			assignForm.value.studentId,
			assignForm.value.academicYearId
		)
		if (res.success && res.data) {
			studentDiscounts.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load student discounts', e)
	}
}

function openAddPolicyModal() {
	policyModalError.value = ''
	policyForm.value = {
		code: '',
		name: '',
		description: '',
		discountType: 'percentage',
		ratePercent: 20,
		amountDecimal: 0,
		feeHeadId: '',
		requiresApproval: true
	}
	showPolicyModal.value = true
}

async function submitPolicyForm() {
	savingPolicy.value = true
	policyModalError.value = ''
	try {
		const basisPoints = policyForm.value.discountType === 'percentage'
			? Math.round(policyForm.value.ratePercent * 100)
			: 0
		const amountMinor = policyForm.value.discountType === 'fixed'
			? Math.round(policyForm.value.amountDecimal * 100)
			: 0

		const res = await services.fees.createDiscountPolicy({
			code: policyForm.value.code.toUpperCase(),
			name: policyForm.value.name,
			description: policyForm.value.description,
			discountType: policyForm.value.discountType,
			rateBasisPoints: basisPoints,
			amountMinor: amountMinor,
			feeHeadId: policyForm.value.feeHeadId || undefined,
			requiresApproval: policyForm.value.requiresApproval
		})
		if (!res.success) {
			policyModalError.value = res.error?.message || 'Failed to create discount policy'
			return
		}
		showPolicyModal.value = false
		await loadPolicies()
	} catch (e: any) {
		policyModalError.value = e.message || 'An error occurred'
	} finally {
		savingPolicy.value = false
	}
}

async function decideDiscount(discountId: string, approved: boolean, reason?: string) {
	try {
		const res = await services.fees.decideStudentDiscount({
			discountId: discountId,
			approved: approved,
			rejectionReason: reason
		})
		if (res.success) {
			await Promise.all([loadPendingDiscounts(), loadStudentDiscounts()])
		} else {
			alert(res.error?.message || 'Failed to process decision')
		}
	} catch (e: any) {
		alert(e.message || 'An error occurred')
	}
}

function openRejectModal(discountId: string) {
	rejectingDiscountId.value = discountId
	rejectionReason.value = ''
	showRejectModal.value = true
}

async function submitReject() {
	if (!rejectingDiscountId.value) return
	await decideDiscount(rejectingDiscountId.value, false, rejectionReason.value)
	showRejectModal.value = false
}

async function submitAssignForm() {
	assigning.value = true
	assignError.value = ''
	assignSuccess.value = ''
	try {
		const res = await services.fees.applyStudentDiscount({
			studentId: assignForm.value.studentId,
			discountPolicyId: assignForm.value.discountPolicyId,
			academicYearId: assignForm.value.academicYearId,
			validFrom: assignForm.value.validFrom,
			validTo: assignForm.value.validTo || undefined
		})
		if (!res.success) {
			assignError.value = res.error?.message || 'Failed to apply concession'
			return
		}
		assignSuccess.value = 'Concession application registered successfully.'
		await Promise.all([loadStudentDiscounts(), loadPendingDiscounts()])
	} catch (e: any) {
		assignError.value = e.message || 'An error occurred'
	} finally {
		assigning.value = false
	}
}

function getFeeHeadName(headId: string) {
	const h = feeHeads.value.find(hd => hd.ID === headId)
	return h ? h.Name : headId
}
</script>
