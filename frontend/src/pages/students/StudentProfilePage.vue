<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<!-- Back Button & Header -->
		<div>
			<router-link to="/students" class="inline-flex items-center gap-1 text-xs text-slate-500 hover:text-slate-800 transition mb-2">
				<ArrowLeft class="w-3.5 h-3.5" />
				<span>Back to Directory</span>
			</router-link>

			<div v-if="student" class="bg-white p-6 rounded-xl border border-slate-200 shadow-xs flex flex-wrap items-center justify-between gap-4">
				<div class="flex items-center gap-4">
					<div class="w-14 h-14 rounded-full bg-indigo-50 border border-indigo-200 flex items-center justify-center text-indigo-700 font-bold text-lg">
						{{ student.FirstName[0] }}{{ student.LastName[0] }}
					</div>
					<div>
						<div class="flex items-center gap-2">
							<h1 class="text-xl font-bold text-slate-900">{{ student.FirstName }} {{ student.MiddleName ? student.MiddleName + ' ' : '' }}{{ student.LastName }}</h1>
							<span
								class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold capitalize"
								:class="getStatusBadgeClass(student.Status)"
							>
								{{ student.Status }}
							</span>
						</div>
						<div class="flex flex-wrap items-center gap-4 text-xs text-slate-500 mt-1 font-mono">
							<span>ADM: {{ student.AdmissionNo }}</span>
							<span v-if="activeEnrollment">• Class: {{ activeEnrollment.ClassName || 'Unassigned' }} {{ activeEnrollment.SectionName ? '(' + activeEnrollment.SectionName + ')' : '' }}</span>
							<span>• Gender: <span class="capitalize">{{ student.Gender }}</span></span>
						</div>
					</div>
				</div>

				<div class="flex items-center gap-2">
					<button
						@click="openChangeStatusModal"
						class="px-3 py-1.5 border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition cursor-pointer"
					>
						Change Status
					</button>
				</div>
			</div>
		</div>

		<!-- Navigation Tabs -->
		<div class="flex border-b border-slate-200 space-x-6 text-xs font-medium">
			<button
				v-for="tab in tabs"
				:key="tab.id"
				@click="currentTab = tab.id"
				class="pb-3 cursor-pointer transition border-b-2"
				:class="currentTab === tab.id ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'"
			>
				{{ tab.label }}
			</button>
		</div>

		<!-- Tab 1: Overview -->
		<div v-if="currentTab === 'overview' && student" class="grid grid-cols-2 gap-6">
			<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
				<h3 class="text-xs font-bold uppercase tracking-wider text-slate-400">Personal Information</h3>
				<div class="grid grid-cols-2 gap-4 text-xs">
					<div>
						<span class="text-slate-500 block">Full Name</span>
						<span class="font-medium text-slate-800">{{ student.FirstName }} {{ student.LastName }}</span>
					</div>
					<div>
						<span class="text-slate-500 block">Date of Birth</span>
						<span class="font-medium text-slate-800">{{ formatDate(student.DateOfBirth) }}</span>
					</div>
					<div>
						<span class="text-slate-500 block">Gender</span>
						<span class="font-medium text-slate-800 capitalize">{{ student.Gender }}</span>
					</div>
					<div>
						<span class="text-slate-500 block">National ID / Birth Cert</span>
						<span class="font-medium text-slate-800">{{ student.NationalID || '—' }}</span>
					</div>
					<div>
						<span class="text-slate-500 block">Admission Date</span>
						<span class="font-medium text-slate-800">{{ formatDate(student.AdmissionDate) }}</span>
					</div>
				</div>
			</div>

			<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
				<h3 class="text-xs font-bold uppercase tracking-wider text-slate-400">Notes & Remarks</h3>
				<p class="text-xs text-slate-600 leading-relaxed">{{ student.Notes || 'No notes on record.' }}</p>
			</div>
		</div>

		<!-- Tab 2: Academic History -->
		<div v-if="currentTab === 'academics'" class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<div class="p-4 border-b border-slate-100 flex items-center justify-between">
				<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">Enrollment History</h3>
			</div>
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-2.5 px-4">Academic Year</th>
						<th class="py-2.5 px-4">Class</th>
						<th class="py-2.5 px-4">Section</th>
						<th class="py-2.5 px-4">Roll No</th>
						<th class="py-2.5 px-4">Status</th>
						<th class="py-2.5 px-4">Enrolled On</th>
						<th class="py-2.5 px-4">Ended On</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-slate-700">
					<tr v-if="enrollmentHistory.length === 0">
						<td colspan="7" class="py-8 text-center text-slate-400">No enrollment history records found.</td>
					</tr>
					<tr v-for="enr in enrollmentHistory" :key="enr.ID" class="hover:bg-slate-50/50">
						<td class="py-2.5 px-4 font-medium">{{ enr.YearCode || enr.AcademicYearID }}</td>
						<td class="py-2.5 px-4 font-semibold text-slate-900">{{ enr.ClassName || enr.ClassID }}</td>
						<td class="py-2.5 px-4">{{ enr.SectionName || '—' }}</td>
						<td class="py-2.5 px-4 font-mono">{{ enr.RollNo || '—' }}</td>
						<td class="py-2.5 px-4">
							<span
								class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
								:class="enr.Status === 'active' ? 'bg-emerald-50 text-emerald-700' : 'bg-slate-100 text-slate-600'"
							>
								{{ enr.Status }}
							</span>
						</td>
						<td class="py-2.5 px-4">{{ formatDate(enr.EnrolledOn) }}</td>
						<td class="py-2.5 px-4 text-slate-400">{{ enr.EndedOn ? formatDate(enr.EndedOn) : 'Current' }}</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Tab 3: Guardians & Siblings -->
		<div v-if="currentTab === 'guardians'" class="space-y-6">
			<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs">
				<div class="flex items-center justify-between mb-4">
					<div>
						<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">Associated Guardians</h3>
						<p class="text-[11px] text-slate-500">Parents, legal guardians, and authorized pickup contacts.</p>
					</div>
					<button
						@click="openAddGuardianModal"
						class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
					>
						Add Guardian
					</button>
				</div>

				<div v-if="guardiansList.length === 0" class="py-8 text-center text-slate-400 text-xs">
					No guardians linked to this student yet.
				</div>

				<div v-else class="grid grid-cols-2 gap-4">
					<div v-for="g in guardiansList" :key="g.GuardianID" class="p-4 rounded-lg border border-slate-200 bg-slate-50/50 space-y-2">
						<div class="flex items-center justify-between">
							<div class="font-bold text-xs text-slate-900">{{ g.FirstName }} {{ g.LastName }}</div>
							<span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-indigo-50 text-indigo-700 capitalize border border-indigo-100">
								{{ g.Relationship }}
							</span>
						</div>
						<div class="text-xs text-slate-600 space-y-0.5 font-mono">
							<div>Phone: {{ g.PhonePrimary || '—' }}</div>
							<div>Email: {{ g.Email || '—' }}</div>
						</div>
						<div class="flex flex-wrap gap-1.5 pt-2 border-t border-slate-200/60 text-[10px]">
							<span v-if="g.IsPrimary" class="px-1.5 py-0.5 bg-emerald-100 text-emerald-800 rounded font-medium">Primary Contact</span>
							<span v-if="g.IsEmergencyContact" class="px-1.5 py-0.5 bg-rose-100 text-rose-800 rounded font-medium">Emergency Contact</span>
							<span v-if="g.CanPickup" class="px-1.5 py-0.5 bg-blue-100 text-blue-800 rounded font-medium">Authorized Pickup</span>
						</div>
					</div>
				</div>
			</div>

			<!-- Siblings Detection Card -->
			<div class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs">
				<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700 mb-1">Detected Siblings</h3>
				<p class="text-[11px] text-slate-500 mb-4">Automatically identified via shared guardian telephone records.</p>

				<div v-if="siblings.length === 0" class="text-xs text-slate-400 italic">
					No siblings detected with matching guardian contact details.
				</div>
				<div v-else class="grid grid-cols-3 gap-3">
					<div v-for="sib in siblings" :key="sib.ID" class="p-3 border border-slate-200 rounded-lg bg-slate-50/50">
						<div class="font-semibold text-xs text-slate-900">
							<router-link :to="`/students/${sib.ID}`" class="hover:text-indigo-600">
								{{ sib.FirstName }} {{ sib.LastName }}
							</router-link>
						</div>
						<div class="text-[11px] text-slate-500 font-mono mt-0.5">ADM: {{ sib.AdmissionNo }}</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Tab 4: Documents -->
		<div v-if="currentTab === 'documents'" class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">Student Documents</h3>
					<p class="text-[11px] text-slate-500">Birth certificates, medical records, and past academic certificates.</p>
				</div>
				<label class="px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer flex items-center gap-1.5">
					<Upload class="w-3.5 h-3.5" />
					<span>Upload Document</span>
					<input type="file" @change="handleFileUpload" class="hidden" />
				</label>
			</div>

			<div v-if="documents.length === 0" class="py-8 text-center text-slate-400 text-xs">
				No documents uploaded for this student yet.
			</div>

			<div v-else class="divide-y divide-slate-100 text-xs">
				<div v-for="doc in documents" :key="doc.ID" class="py-3 flex items-center justify-between">
					<div class="flex items-center gap-3">
						<FileText class="w-5 h-5 text-indigo-500" />
						<div>
							<div class="font-semibold text-slate-900">{{ doc.FileName }}</div>
							<div class="text-[11px] text-slate-400 font-mono">{{ (doc.SizeBytes / 1024).toFixed(1) }} KB • Type: {{ doc.DocumentTypeCode }}</div>
						</div>
					</div>
					<button
						@click="downloadDocument(doc.ID, doc.FileName)"
						class="text-indigo-600 hover:text-indigo-800 font-medium text-xs px-2.5 py-1 bg-indigo-50 rounded"
					>
						Download
					</button>
				</div>
			</div>
		</div>

		<!-- Change Status Modal -->
		<div v-if="showStatusModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-5 border border-slate-200">
				<h3 class="text-sm font-bold text-slate-900 mb-3">Change Student Status</h3>
				<form @submit.prevent="submitStatusChange" class="space-y-3">
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">New Status</label>
						<select v-model="newStatus" class="w-full text-xs border border-slate-300 rounded-lg p-2">
							<option value="active">Active</option>
							<option value="inactive">Inactive</option>
							<option value="graduated">Graduated</option>
							<option value="withdrawn">Withdrawn</option>
						</select>
					</div>
					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Reason</label>
						<textarea v-model="statusReason" required rows="2" class="w-full text-xs border border-slate-300 rounded-lg p-2" placeholder="Reason for change..."></textarea>
					</div>
					<div class="flex justify-end gap-2 pt-2">
						<button type="button" @click="showStatusModal = false" class="px-3 py-1.5 border border-slate-300 text-xs rounded-lg">Cancel</button>
						<button type="submit" class="px-3 py-1.5 bg-indigo-600 text-white text-xs rounded-lg font-medium">Update Status</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getServices } from '@/services/factory'
import type {
	StudentDTO,
	EnrollmentDTO,
	StudentGuardianLinkDTO,
	DocumentDTO
} from '@/types/common'
import { ArrowLeft, Upload, FileText } from 'lucide-vue-next'

const route = useRoute()
const services = getServices()

const currentTab = ref('overview')
const tabs = [
	{ id: 'overview', label: 'Overview' },
	{ id: 'academics', label: 'Academic History' },
	{ id: 'guardians', label: 'Guardians & Siblings' },
	{ id: 'documents', label: 'Documents' }
]

const student = ref<StudentDTO | null>(null)
const activeEnrollment = ref<EnrollmentDTO | null>(null)
const enrollmentHistory = ref<EnrollmentDTO[]>([])
const guardiansList = ref<StudentGuardianLinkDTO[]>([])
const siblings = ref<StudentDTO[]>([])
const documents = ref<DocumentDTO[]>([])

const showStatusModal = ref(false)
const newStatus = ref('active')
const statusReason = ref('')

async function loadProfile() {
	const studentId = route.params.id as string
	if (!studentId) return

	const [studentRes, activeRes, historyRes, guardiansRes, siblingsRes, docsRes] = await Promise.all([
		services.students.getStudent(studentId),
		services.enrollments.getActive(studentId),
		services.enrollments.listHistory(studentId),
		services.guardians.listForStudent(studentId),
		services.guardians.findSiblings(studentId),
		services.documents.listForStudent(studentId)
	])

	if (studentRes.success && studentRes.data) student.value = studentRes.data
	if (activeRes.success && activeRes.data) activeEnrollment.value = activeRes.data
	if (historyRes.success && historyRes.data) enrollmentHistory.value = historyRes.data
	if (guardiansRes.success && guardiansRes.data) guardiansList.value = guardiansRes.data
	if (siblingsRes.success && siblingsRes.data) siblings.value = siblingsRes.data
	if (docsRes.success && docsRes.data) documents.value = docsRes.data
}

function openChangeStatusModal() {
	if (student.value) newStatus.value = student.value.Status
	statusReason.value = ''
	showStatusModal.value = true
}

async function submitStatusChange() {
	if (!student.value) return
	const today = new Date().toISOString().split('T')[0].split('-').map(Number)
	const res = await services.students.changeStatus({
		studentId: student.value.ID,
		newStatus: newStatus.value,
		effectiveOn: { Year: today[0], Month: today[1], Day: today[2] },
		reason: statusReason.value
	})
	if (res.success) {
		showStatusModal.value = false
		await loadProfile()
	}
}

async function handleFileUpload(e: Event) {
	const input = e.target as HTMLInputElement
	if (!input.files || input.files.length === 0 || !student.value) return
	const file = input.files[0]
	const reader = new FileReader()
	reader.onload = async () => {
		const base64 = (reader.result as string).split(',')[1]
		await services.documents.uploadDocument({
			studentId: student.value?.ID,
			documentTypeCode: 'GENERAL',
			fileName: file.name,
			mimeType: file.type || 'application/octet-stream',
			dataBase64: base64
		})
		await loadProfile()
	}
	reader.readAsDataURL(file)
}

async function downloadDocument(id: string, fileName: string) {
	const res = await services.documents.getContent(id)
	if (res.success && res.data) {
		const byteCharacters = atob(res.data.dataBase64)
		const byteNumbers = new Array(byteCharacters.length)
		for (let i = 0; i < byteCharacters.length; i++) {
			byteNumbers[i] = byteCharacters.charCodeAt(i)
		}
		const byteArray = new Uint8Array(byteNumbers)
		const blob = new Blob([byteArray])
		const link = document.createElement('a')
		link.href = window.URL.createObjectURL(blob)
		link.download = fileName
		link.click()
	}
}

function openAddGuardianModal() {
	// Simple prompt or placeholder modal for linking
	alert('Open guardian form modal')
}

function formatDate(d: any) {
	if (!d) return '—'
	if (typeof d === 'string') return d
	if (d.Year && d.Month && d.Day) {
		return `${d.Year}-${String(d.Month).padStart(2, '0')}-${String(d.Day).padStart(2, '0')}`
	}
	return '—'
}

function getStatusBadgeClass(status: string) {
	switch (status) {
		case 'active':
			return 'bg-emerald-50 text-emerald-700 border border-emerald-200'
		case 'graduated':
			return 'bg-blue-50 text-blue-700 border border-blue-200'
		case 'withdrawn':
			return 'bg-rose-50 text-rose-700 border border-rose-200'
		default:
			return 'bg-slate-100 text-slate-600 border border-slate-200'
	}
}

watch(() => route.params.id, () => {
	loadProfile()
})

onMounted(() => {
	loadProfile()
})
</script>

