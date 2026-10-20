<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Admission Applications</h1>
				<p class="text-xs text-slate-500 mt-0.5">Track applicant submissions, conduct reviews, and enroll new students.</p>
			</div>
			<div class="flex items-center gap-2">
				<router-link
					to="/admissions/builder"
					class="py-2 px-3.5 border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5"
				>
					<Layers class="w-4 h-4 text-slate-500" />
					<span>Form Builder</span>
				</router-link>
				<button
					@click="openNewAppModal"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>New Application</span>
				</button>
			</div>
		</div>

		<!-- Filter Bar -->
		<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex flex-wrap gap-3 items-center justify-between">
			<div class="flex items-center gap-3">
				<select
					v-model="statusFilter"
					@change="loadApplications"
					class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700"
				>
					<option value="">All Statuses</option>
					<option value="submitted">Submitted</option>
					<option value="under_review">Under Review</option>
					<option value="approved">Approved</option>
					<option value="rejected">Rejected</option>
					<option value="enrolled">Enrolled</option>
				</select>
			</div>
			<button
				@click="loadApplications"
				class="p-2 border border-slate-200 hover:bg-slate-50 rounded-lg text-slate-600 transition"
			>
				<RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
			</button>
		</div>

		<!-- Applications Table -->
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">Application #</th>
						<th class="py-3 px-4">Applicant Name</th>
						<th class="py-3 px-4">Applied Class</th>
						<th class="py-3 px-4">Guardian Contact</th>
						<th class="py-3 px-4">Status</th>
						<th class="py-3 px-4 text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-slate-700">
					<tr v-if="loading && applications.length === 0">
						<td colspan="6" class="py-10 text-center text-slate-400">Loading applications...</td>
					</tr>
					<tr v-else-if="applications.length === 0">
						<td colspan="6" class="py-10 text-center text-slate-400">No applications found.</td>
					</tr>
					<tr v-for="app in applications" :key="app.ID" class="hover:bg-slate-50/50">
						<td class="py-3 px-4 font-mono font-medium text-slate-900">{{ app.ApplicationNo }}</td>
						<td class="py-3 px-4 font-semibold text-slate-900">{{ app.ApplicantFirstName }} {{ app.ApplicantLastName }}</td>
						<td class="py-3 px-4 font-medium">{{ app.ClassName || app.AppliedClassID }}</td>
						<td class="py-3 px-4 font-mono text-slate-600">{{ app.GuardianPhoneNormalized || '—' }}</td>
						<td class="py-3 px-4">
							<span
								class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
								:class="getStatusBadgeClass(app.Status)"
							>
								{{ app.Status.replace('_', ' ') }}
							</span>
						</td>
						<td class="py-3 px-4 text-right space-x-2">
							<button
								v-if="app.Status === 'approved'"
								@click="openEnrollModal(app)"
								class="px-2 py-1 bg-emerald-600 hover:bg-emerald-700 text-white rounded font-medium text-[11px] cursor-pointer"
							>
								Enroll Student
							</button>
							<button
								v-else-if="app.Status !== 'enrolled'"
								@click="openReviewModal(app)"
								class="px-2 py-1 bg-slate-100 hover:bg-slate-200 text-slate-700 rounded font-medium text-[11px] cursor-pointer"
							>
								Review
							</button>
							<span v-else class="text-emerald-600 font-medium text-[11px]">Enrolled</span>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Review Modal -->
		<div v-if="showReviewModal && selectedApp" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-5 border border-slate-200">
				<h3 class="text-sm font-bold text-slate-900 mb-2">Review Application: {{ selectedApp.ApplicationNo }}</h3>
				<form @submit.prevent="submitReview" class="space-y-3 text-xs">
					<div>
						<label class="block font-medium text-slate-700 mb-1">Decision / Status</label>
						<select v-model="reviewStatus" class="w-full border border-slate-300 rounded-lg p-2">
							<option value="under_review">Under Review</option>
							<option value="approved">Approved</option>
							<option value="rejected">Rejected</option>
						</select>
					</div>
					<div>
						<label class="block font-medium text-slate-700 mb-1">Notes / Reason</label>
						<textarea v-model="reviewNotes" rows="3" required class="w-full border border-slate-300 rounded-lg p-2" placeholder="Review comments..."></textarea>
					</div>
					<div class="flex justify-end gap-2 pt-2">
						<button type="button" @click="showReviewModal = false" class="px-3 py-1.5 border border-slate-300 rounded-lg">Cancel</button>
						<button type="submit" class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg font-medium">Save Decision</button>
					</div>
				</form>
			</div>
		</div>

		<!-- Enroll Student Modal -->
		<div v-if="showEnrollModal && selectedApp" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-5 border border-slate-200">
				<h3 class="text-sm font-bold text-slate-900 mb-1">Enroll Student</h3>
				<p class="text-xs text-slate-500 mb-3">Converts {{ selectedApp.ApplicantFirstName }} {{ selectedApp.ApplicantLastName }} into an active student.</p>

				<form @submit.prevent="submitEnroll" class="space-y-3 text-xs">
					<div>
						<label class="block font-medium text-slate-700 mb-1">Admission Date *</label>
						<input v-model="enrollAdmissionDate" type="date" required class="w-full border border-slate-300 rounded-lg p-2" />
					</div>
					<div>
						<label class="block font-medium text-slate-700 mb-1">Section</label>
						<select v-model="enrollSectionId" class="w-full border border-slate-300 rounded-lg p-2">
							<option value="">No Section</option>
							<option v-for="s in availableSections" :key="s.ID" :value="s.ID">{{ s.Name }}</option>
						</select>
					</div>
					<div>
						<label class="block font-medium text-slate-700 mb-1">Roll Number</label>
						<input v-model="enrollRollNo" placeholder="e.g. 01" class="w-full border border-slate-300 rounded-lg p-2" />
					</div>
					<div class="flex justify-end gap-2 pt-2">
						<button type="button" @click="showEnrollModal = false" class="px-3 py-1.5 border border-slate-300 rounded-lg">Cancel</button>
						<button type="submit" class="px-3 py-1.5 bg-emerald-600 text-white rounded-lg font-medium">Confirm Enrollment</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getServices } from '@/services/factory'
import type { ApplicationDTO, SectionDTO } from '@/types/common'
import { Plus, RefreshCw, Layers } from 'lucide-vue-next'

const router = useRouter()
const services = getServices()

const applications = ref<ApplicationDTO[]>([])
const statusFilter = ref('')
const loading = ref(false)

const showReviewModal = ref(false)
const showEnrollModal = ref(false)
const selectedApp = ref<ApplicationDTO | null>(null)

const reviewStatus = ref('under_review')
const reviewNotes = ref('')

const enrollAdmissionDate = ref(new Date().toISOString().split('T')[0])
const enrollSectionId = ref('')
const enrollRollNo = ref('')
const availableSections = ref<SectionDTO[]>([])

async function loadApplications() {
	loading.value = true
	try {
		const res = await services.admissions.listApplications(statusFilter.value, '', 50)
		if (res.success && res.data) {
			applications.value = res.data
		}
	} finally {
		loading.value = false
	}
}

function openReviewModal(app: ApplicationDTO) {
	selectedApp.value = app
	reviewStatus.value = app.Status === 'submitted' ? 'under_review' : app.Status
	reviewNotes.value = ''
	showReviewModal.value = true
}

async function submitReview() {
	if (!selectedApp.value) return
	const res = await services.admissions.updateStatus({
		applicationId: selectedApp.value.ID,
		newStatus: reviewStatus.value,
		reason: reviewNotes.value
	})
	if (res.success) {
		showReviewModal.value = false
		await loadApplications()
	}
}

async function openEnrollModal(app: ApplicationDTO) {
	selectedApp.value = app
	enrollSectionId.value = ''
	enrollRollNo.value = ''
	// Load sections for academic year
	const secRes = await services.academics.listSectionsByYear(app.AcademicYearID)
	if (secRes.success && secRes.data) {
		availableSections.value = secRes.data.filter(s => s.ClassID === app.AppliedClassID)
	}
	showEnrollModal.value = true
}

async function submitEnroll() {
	if (!selectedApp.value) return
	const [y, m, d] = enrollAdmissionDate.value.split('-').map(Number)
	const res = await services.admissions.enrollApplication({
		applicationId: selectedApp.value.ID,
		sectionId: enrollSectionId.value || undefined,
		rollNo: enrollRollNo.value || undefined,
		admissionDate: { Year: y, Month: m, Day: d }
	})
	if (res.success && res.data) {
		showEnrollModal.value = false
		router.push(`/students/${res.data}`)
	}
}

function openNewAppModal() {
	alert('Direct submission modal or redirect to admission form')
}

function getStatusBadgeClass(status: string) {
	switch (status) {
		case 'submitted':
			return 'bg-blue-50 text-blue-700 border border-blue-200'
		case 'under_review':
			return 'bg-amber-50 text-amber-700 border border-amber-200'
		case 'approved':
			return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
		case 'rejected':
			return 'bg-rose-50 text-rose-700 border border-rose-200'
		case 'enrolled':
			return 'bg-purple-50 text-purple-700 border border-purple-200'
		default:
			return 'bg-slate-100 text-slate-600 border border-slate-200'
	}
}

onMounted(() => {
	loadApplications()
})
</script>

