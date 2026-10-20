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
				class="p-2 border border-slate-200 hover:bg-slate-50 rounded-lg text-slate-600 transition cursor-pointer"
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
						<td colspan="6" class="py-10 text-center text-slate-400">No applications found. Click "New Application" to add one.</td>
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

		<!-- In-App Confirmation / Action Choice Modal -->
		<div v-if="showNewAppChoiceModal" class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-150">
			<div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 border border-slate-200">
				<div class="flex items-start justify-between mb-4">
					<div class="flex items-center gap-3">
						<div class="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center shrink-0">
							<UserPlus class="w-5 h-5" />
						</div>
						<div>
							<h3 class="text-base font-bold text-slate-900">New Admission Application</h3>
							<p class="text-xs text-slate-500 mt-0.5">Choose how you would like to proceed:</p>
						</div>
					</div>
					<button @click="showNewAppChoiceModal = false" class="text-slate-400 hover:text-slate-600 p-1 cursor-pointer">
						<X class="w-4 h-4" />
					</button>
				</div>

				<div class="space-y-3 mb-5">
					<!-- Option 1: Direct Entry -->
					<div
						@click="startDirectSubmission"
						class="p-4 rounded-xl border border-slate-200 hover:border-indigo-400 hover:bg-indigo-50/40 transition cursor-pointer group"
					>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2.5">
								<FileEdit class="w-4 h-4 text-indigo-600" />
								<span class="text-xs font-bold text-slate-900 group-hover:text-indigo-600 transition">Direct Application Form</span>
							</div>
							<span class="text-[10px] bg-indigo-50 text-indigo-700 px-2 py-0.5 rounded-full font-medium">Quick Entry</span>
						</div>
						<p class="text-[11px] text-slate-500 mt-1 pl-6.5">Enter applicant details, guardian contact, and target class directly in the modal.</p>
					</div>

					<!-- Option 2: Form Builder -->
					<div
						@click="goToFormBuilder"
						class="p-4 rounded-xl border border-slate-200 hover:border-slate-400 hover:bg-slate-50/80 transition cursor-pointer group"
					>
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-2.5">
								<Layers class="w-4 h-4 text-slate-600" />
								<span class="text-xs font-bold text-slate-900 group-hover:text-slate-900 transition">Admission Form Builder</span>
							</div>
							<span class="text-[10px] bg-slate-100 text-slate-600 px-2 py-0.5 rounded-full font-medium">Custom Fields</span>
						</div>
						<p class="text-[11px] text-slate-500 mt-1 pl-6.5">Customize questionnaire fields, sections, validation rules, and publish schemas.</p>
					</div>
				</div>

				<div class="flex justify-end pt-3 border-t border-slate-100">
					<button
						type="button"
						@click="showNewAppChoiceModal = false"
						class="px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition cursor-pointer"
					>
						Cancel
					</button>
				</div>
			</div>
		</div>

		<!-- Direct Application Submission Modal -->
		<div v-if="showDirectSubmitModal" class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-150">
			<div class="bg-white rounded-2xl shadow-2xl max-w-lg w-full p-6 border border-slate-200 max-h-[90vh] overflow-y-auto">
				<div class="flex items-start justify-between mb-4 border-b border-slate-100 pb-3">
					<div>
						<h3 class="text-base font-bold text-slate-900">Direct Admission Application</h3>
						<p class="text-xs text-slate-500 mt-0.5">Submit an applicant for admissions review and enrollment.</p>
					</div>
					<button @click="showDirectSubmitModal = false" class="text-slate-400 hover:text-slate-600 p-1 cursor-pointer">
						<X class="w-4 h-4" />
					</button>
				</div>

				<form @submit.prevent="submitDirectApplication" class="space-y-4 text-xs">
					<!-- Academic Year & Applied Class -->
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block font-medium text-slate-700 mb-1">Academic Year *</label>
							<select v-model="newAppForm.academicYearId" required class="w-full border border-slate-300 rounded-lg p-2 bg-white text-xs">
								<option value="" disabled>Select Year</option>
								<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
							</select>
						</div>
						<div>
							<label class="block font-medium text-slate-700 mb-1">Applied Class *</label>
							<select v-model="newAppForm.appliedClassId" required class="w-full border border-slate-300 rounded-lg p-2 bg-white text-xs">
								<option value="" disabled>Select Class</option>
								<option v-for="c in classes" :key="c.ID" :value="c.ID">{{ c.Name }} ({{ c.Code }})</option>
							</select>
						</div>
					</div>

					<!-- Applicant Names -->
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block font-medium text-slate-700 mb-1">Applicant First Name *</label>
							<input v-model="newAppForm.applicantFirstName" required placeholder="e.g. Zain" class="w-full border border-slate-300 rounded-lg p-2 text-xs" />
						</div>
						<div>
							<label class="block font-medium text-slate-700 mb-1">Applicant Last Name *</label>
							<input v-model="newAppForm.applicantLastName" required placeholder="e.g. Abbas" class="w-full border border-slate-300 rounded-lg p-2 text-xs" />
						</div>
					</div>

					<!-- DOB & Gender -->
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block font-medium text-slate-700 mb-1">Date of Birth</label>
							<input v-model="newAppForm.applicantDob" type="date" class="w-full border border-slate-300 rounded-lg p-2 text-xs" />
						</div>
						<div>
							<label class="block font-medium text-slate-700 mb-1">Gender *</label>
							<select v-model="newAppForm.applicantGender" required class="w-full border border-slate-300 rounded-lg p-2 bg-white text-xs">
								<option value="male">Male</option>
								<option value="female">Female</option>
								<option value="other">Other</option>
							</select>
						</div>
					</div>

					<!-- Guardian Info -->
					<div class="border-t border-slate-100 pt-3">
						<h4 class="text-[11px] font-bold uppercase tracking-wider text-slate-400 mb-2">Guardian Information</h4>
						<div class="grid grid-cols-2 gap-3">
							<div>
								<label class="block font-medium text-slate-700 mb-1">Guardian Name *</label>
								<input v-model="newAppForm.guardianName" required placeholder="e.g. Tariq Abbas" class="w-full border border-slate-300 rounded-lg p-2 text-xs" />
							</div>
							<div>
								<label class="block font-medium text-slate-700 mb-1">Guardian Phone *</label>
								<input v-model="newAppForm.guardianPhone" required placeholder="e.g. +92-300-1234567" class="w-full border border-slate-300 rounded-lg p-2 text-xs" />
							</div>
						</div>
						<div class="mt-3">
							<label class="block font-medium text-slate-700 mb-1">Guardian Email</label>
							<input v-model="newAppForm.guardianEmail" type="email" placeholder="e.g. guardian@example.com" class="w-full border border-slate-300 rounded-lg p-2 text-xs" />
						</div>
					</div>

					<!-- Notes -->
					<div>
						<label class="block font-medium text-slate-700 mb-1">Notes / Previous School</label>
						<textarea v-model="newAppForm.notes" rows="2" placeholder="Optional background or previous schooling..." class="w-full border border-slate-300 rounded-lg p-2 text-xs"></textarea>
					</div>

					<!-- Error Alert -->
					<div v-if="submitError" class="p-2.5 rounded-lg bg-rose-50 text-rose-700 text-xs border border-rose-200">
						{{ submitError }}
					</div>

					<!-- Footer buttons -->
					<div class="flex items-center justify-end gap-2.5 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showDirectSubmitModal = false"
							class="px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition cursor-pointer"
							:disabled="isSubmitting"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
							:disabled="isSubmitting"
						>
							<span v-if="isSubmitting" class="animate-spin mr-1">⏳</span>
							<span>Submit Application</span>
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- Review Modal -->
		<div v-if="showReviewModal && selectedApp" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-5 border border-slate-200">
				<h3 class="text-sm font-bold text-slate-900 mb-2">Review Application: {{ selectedApp.ApplicationNo }}</h3>
				<form @submit.prevent="submitReview" class="space-y-3 text-xs">
					<div>
						<label class="block font-medium text-slate-700 mb-1">Decision / Status</label>
						<select v-model="reviewStatus" class="w-full border border-slate-300 rounded-lg p-2 bg-white">
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
						<button type="button" @click="showReviewModal = false" class="px-3 py-1.5 border border-slate-300 rounded-lg cursor-pointer">Cancel</button>
						<button type="submit" class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg font-medium cursor-pointer">Save Decision</button>
					</div>
				</form>
			</div>
		</div>

		<!-- Enroll Student Modal -->
		<div v-if="showEnrollModal && selectedApp" class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4 z-50 animate-in fade-in duration-150">
			<div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 border border-slate-200 transform transition-all animate-in zoom-in-95 duration-150 relative">
				<div class="flex items-start justify-between mb-4">
					<div>
						<h3 class="text-base font-bold text-slate-900">Enroll Student</h3>
						<p class="text-xs text-slate-500 mt-0.5">Convert applicant into an active enrolled student.</p>
					</div>
					<button
						type="button"
						@click="showEnrollModal = false"
						class="text-slate-400 hover:text-slate-600 rounded-lg p-1 transition cursor-pointer"
						title="Close"
					>
						<X class="w-4 h-4" />
					</button>
				</div>

				<!-- Applicant Summary Pill -->
				<div class="mb-4 p-3 bg-indigo-50/70 border border-indigo-100 rounded-xl flex items-center justify-between text-xs">
					<div>
						<p class="font-bold text-slate-900">
							{{ selectedApp.ApplicantFirstName }} {{ selectedApp.ApplicantLastName }}
						</p>
						<p class="text-[11px] text-slate-500">
							Application: <span class="font-mono font-semibold text-indigo-700">{{ selectedApp.ApplicationNo }}</span>
						</p>
					</div>
					<div class="text-right">
						<span class="text-[10px] uppercase font-bold text-slate-400 block tracking-wider">Applied Class</span>
						<span class="font-semibold text-slate-800">{{ selectedApp.ClassName || getClassName(selectedApp.AppliedClassID) }}</span>
					</div>
				</div>

				<form @submit.prevent="submitEnroll" class="space-y-3.5 text-xs">
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block font-medium text-slate-700 mb-1">Admission Date *</label>
							<input
								v-model="enrollAdmissionDate"
								type="date"
								required
								class="w-full border border-slate-300 rounded-lg p-2 text-xs bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
						<div>
							<label class="block font-medium text-slate-700 mb-1">Admission / Student ID</label>
							<input
								v-model="enrollAdmissionNo"
								type="text"
								:placeholder="defaultAdmissionNo"
								class="w-full border border-slate-300 rounded-lg p-2 text-xs font-mono bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
					</div>

					<!-- Class Selector -->
					<div>
						<div class="flex items-center justify-between mb-1">
							<label class="block font-medium text-slate-700">Enrolling Class *</label>
							<span class="text-[10px] text-slate-400 font-medium">{{ classes.length }} classes available</span>
						</div>
						<select
							v-model="enrollClassId"
							@change="onEnrollClassChange"
							required
							class="w-full border border-slate-300 rounded-lg p-2 text-xs bg-white text-slate-800 font-medium focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="" disabled>Select Class...</option>
							<option v-for="c in classes" :key="c.ID" :value="c.ID">
								{{ c.Name }}
							</option>
						</select>
					</div>

					<!-- Section Selector -->
					<div>
						<div class="flex items-center justify-between mb-1">
							<label class="block font-medium text-slate-700">Class Section</label>
							<span v-if="availableSections.length > 0" class="text-[10px] text-emerald-600 font-medium">
								{{ availableSections.length }} section(s) available
							</span>
							<span v-else class="text-[10px] text-amber-600 font-medium">
								No sections configured
							</span>
						</div>
						<select
							v-model="enrollSectionId"
							class="w-full border border-slate-300 rounded-lg p-2 text-xs bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="">No Section / General</option>
							<option v-for="s in availableSections" :key="s.ID" :value="s.ID">
								{{ s.Name }}
							</option>
						</select>
					</div>

					<!-- Roll Number -->
					<div>
						<label class="block font-medium text-slate-700 mb-1">Roll Number (Optional)</label>
						<input
							v-model="enrollRollNo"
							placeholder="e.g. 01"
							class="w-full border border-slate-300 rounded-lg p-2 text-xs font-mono bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
					</div>

					<!-- Action buttons -->
					<div class="flex items-center justify-end gap-2.5 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showEnrollModal = false"
							class="px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition cursor-pointer"
							:disabled="enrollingStudent"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-emerald-600 hover:bg-emerald-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
							:disabled="enrollingStudent || !enrollClassId"
						>
							<span v-if="enrollingStudent" class="animate-spin mr-1">⏳</span>
							<span>{{ enrollingStudent ? 'Enrolling...' : 'Confirm Enrollment' }}</span>
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import type { ApplicationDTO, SectionDTO, AcademicYearDTO, ClassDTO } from '@/types/common'
import { Plus, RefreshCw, Layers, UserPlus, FileEdit, X } from 'lucide-vue-next'

const router = useRouter()
const services = getServices()
const toast = useToast()

const applications = ref<ApplicationDTO[]>([])
const statusFilter = ref('')
const loading = ref(false)

// Modals
const showNewAppChoiceModal = ref(false)
const showDirectSubmitModal = ref(false)
const showReviewModal = ref(false)
const showEnrollModal = ref(false)

const selectedApp = ref<ApplicationDTO | null>(null)
const reviewStatus = ref('under_review')
const reviewNotes = ref('')

const enrollAdmissionDate = ref(new Date().toISOString().split('T')[0])
const enrollClassId = ref('')
const enrollSectionId = ref('')
const enrollRollNo = ref('')
const enrollAdmissionNo = ref('')
const defaultAdmissionNo = ref('')
const enrollingStudent = ref(false)
const allSections = ref<SectionDTO[]>([])
const availableSections = ref<SectionDTO[]>([])

function getClassName(classId?: string) {
	if (!classId) return ''
	const c = classes.value.find(cl => cl.ID === classId)
	return c ? c.Name : classId
}

function onEnrollClassChange() {
	enrollSectionId.value = ''
	filterSectionsForClass()
}

function filterSectionsForClass() {
	if (!enrollClassId.value) {
		availableSections.value = []
		return
	}
	availableSections.value = allSections.value.filter(s => s.ClassID === enrollClassId.value)
	if (availableSections.value.length === 1) {
		enrollSectionId.value = availableSections.value[0].ID
	}
}

// Direct Application Form Data
const academicYears = ref<AcademicYearDTO[]>([])
const classes = ref<ClassDTO[]>([])
const isSubmitting = ref(false)
const submitError = ref('')

const newAppForm = reactive({
	academicYearId: '',
	appliedClassId: '',
	applicantFirstName: '',
	applicantLastName: '',
	applicantDob: '',
	applicantGender: 'male',
	guardianName: '',
	guardianPhone: '',
	guardianEmail: '',
	notes: ''
})

async function loadApplications() {
	loading.value = true
	try {
		const res = await services.admissions.listApplications(statusFilter.value, '', 50)
		if (res.success && res.data) {
			applications.value = res.data
		}
	} catch (e: any) {
		toast.error(e.message || 'Failed to load applications')
	} finally {
		loading.value = false
	}
}

async function loadFormMetadata() {
	try {
		const [yearsRes, classesRes] = await Promise.all([
			services.academics.listAcademicYears(),
			services.academics.listClasses()
		])
		if (yearsRes.success && yearsRes.data) {
			academicYears.value = yearsRes.data
			const activeYear = yearsRes.data.find(y => y.Status === 'active') || yearsRes.data[0]
			if (activeYear && !newAppForm.academicYearId) {
				newAppForm.academicYearId = activeYear.ID
			}
		}
		if (classesRes.success && classesRes.data) {
			classes.value = classesRes.data
			if (classesRes.data.length > 0 && !newAppForm.appliedClassId) {
				newAppForm.appliedClassId = classesRes.data[0].ID
			}
		}
	} catch (e: any) {
		console.error('Failed to load academic years/classes:', e)
	}
}

function openNewAppModal() {
	showNewAppChoiceModal.value = true
}

function startDirectSubmission() {
	showNewAppChoiceModal.value = false
	submitError.value = ''
	loadFormMetadata()
	showDirectSubmitModal.value = true
}

function goToFormBuilder() {
	showNewAppChoiceModal.value = false
	router.push('/admissions/builder')
}

async function submitDirectApplication() {
	submitError.value = ''
	if (!newAppForm.academicYearId || !newAppForm.appliedClassId) {
		submitError.value = 'Academic Year and Applied Class are required.'
		return
	}
	if (!newAppForm.applicantFirstName.trim() || !newAppForm.applicantLastName.trim()) {
		submitError.value = 'Applicant first and last names are required.'
		return
	}
	if (!newAppForm.guardianName.trim() || !newAppForm.guardianPhone.trim()) {
		submitError.value = 'Guardian name and contact phone are required.'
		return
	}

	isSubmitting.value = true
	try {
		let dob: { Year: number; Month: number; Day: number } | undefined
		if (newAppForm.applicantDob) {
			const [y, m, d] = newAppForm.applicantDob.split('-').map(Number)
			if (y && m && d) {
				dob = { Year: y, Month: m, Day: d }
			}
		}

		const res = await services.admissions.submitApplication({
			academicYearId: newAppForm.academicYearId,
			appliedClassId: newAppForm.appliedClassId,
			applicantFirstName: newAppForm.applicantFirstName.trim(),
			applicantLastName: newAppForm.applicantLastName.trim(),
			applicantDob: dob,
			applicantGender: newAppForm.applicantGender,
			guardianName: newAppForm.guardianName.trim(),
			guardianPhoneNormalized: newAppForm.guardianPhone.trim(),
			guardianEmail: newAppForm.guardianEmail.trim(),
			response: newAppForm.notes ? { notes: newAppForm.notes } : {}
		})

		if (res.success) {
			toast.success('Admission application submitted successfully!')
			showDirectSubmitModal.value = false
			// Reset form fields
			newAppForm.applicantFirstName = ''
			newAppForm.applicantLastName = ''
			newAppForm.applicantDob = ''
			newAppForm.guardianName = ''
			newAppForm.guardianPhone = ''
			newAppForm.guardianEmail = ''
			newAppForm.notes = ''
			await loadApplications()
		} else {
			submitError.value = res.error?.message || 'Failed to submit application'
		}
	} catch (e: any) {
		submitError.value = e.message || 'An error occurred while submitting'
	} finally {
		isSubmitting.value = false
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
	try {
		const res = await services.admissions.updateStatus({
			applicationId: selectedApp.value.ID,
			newStatus: reviewStatus.value,
			reason: reviewNotes.value
		})
		if (res.success) {
			toast.success(`Application updated to "${reviewStatus.value.replace('_', ' ')}"`)
			showReviewModal.value = false
			await loadApplications()
		} else {
			toast.error(res.error?.message || 'Failed to save review decision')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred')
	}
}

async function openEnrollModal(app: ApplicationDTO) {
	selectedApp.value = app
	if (classes.value.length === 0) {
		await loadFormMetadata()
	}

	enrollClassId.value = app.AppliedClassID || (classes.value.length > 0 ? classes.value[0].ID : '')
	enrollSectionId.value = ''
	enrollRollNo.value = ''
	defaultAdmissionNo.value = app.ApplicationNo ? app.ApplicationNo.replace('APP-', 'ADM-') : 'ADM-2026-0001'
	enrollAdmissionNo.value = defaultAdmissionNo.value
	enrollAdmissionDate.value = new Date().toISOString().split('T')[0]

	try {
		const secRes = await services.academics.listSectionsByYear(app.AcademicYearID)
		if (secRes.success && secRes.data) {
			allSections.value = secRes.data
		} else {
			allSections.value = []
		}
	} catch (e) {
		allSections.value = []
	}

	filterSectionsForClass()
	showEnrollModal.value = true
}

async function submitEnroll() {
	if (!selectedApp.value) return
	if (!enrollClassId.value) {
		toast.error('Please select an Enrolling Class')
		return
	}
	enrollingStudent.value = true
	try {
		const [y, m, d] = enrollAdmissionDate.value.split('-').map(Number)
		const res = await services.admissions.enrollApplication({
			applicationId: selectedApp.value.ID,
			classId: enrollClassId.value,
			sectionId: enrollSectionId.value || undefined,
			rollNo: enrollRollNo.value.trim() || undefined,
			admissionNo: enrollAdmissionNo.value.trim() || defaultAdmissionNo.value,
			admissionDate: { Year: y, Month: m, Day: d }
		})
		if (res.success && res.data) {
			toast.success('Student enrolled successfully!')
			showEnrollModal.value = false
			router.push(`/students/${res.data}`)
		} else {
			toast.error(res.error?.message || 'Failed to enroll student')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred during enrollment')
	} finally {
		enrollingStudent.value = false
	}
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
	loadFormMetadata()
})
</script>
