<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Student Directory</h1>
				<p class="text-xs text-slate-500 mt-0.5">Manage student enrollments, profiles, and academic records.</p>
			</div>
			<button
				@click="openRegisterModal"
				class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
			>
				<UserPlus class="w-4 h-4" />
				<span>Register Student</span>
			</button>
		</div>

		<!-- Filters & Search Bar -->
		<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex flex-wrap gap-3 items-center justify-between">
			<div class="flex-1 min-w-[260px] relative">
				<Search class="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
				<input
					v-model="searchQuery"
					@input="debounceSearch"
					type="text"
					placeholder="Search by name, admission # (FTS search)..."
					class="w-full pl-9 pr-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500"
				/>
			</div>

			<div class="flex items-center gap-2">
				<select
					v-model="statusFilter"
					@change="loadStudents"
					class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
				>
					<option value="">All Statuses</option>
					<option value="active">Active</option>
					<option value="inactive">Inactive</option>
					<option value="graduated">Graduated</option>
					<option value="withdrawn">Withdrawn</option>
				</select>

				<button
					@click="loadStudents"
					class="p-2 border border-slate-200 hover:bg-slate-50 rounded-lg text-slate-600 transition"
					title="Refresh List"
				>
					<RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
				</button>
			</div>
		</div>

		<!-- Students Table -->
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">Admission #</th>
						<th class="py-3 px-4">Student Name</th>
						<th class="py-3 px-4">Gender</th>
						<th class="py-3 px-4">Date of Birth</th>
						<th class="py-3 px-4">Status</th>
						<th class="py-3 px-4 text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-xs text-slate-700">
					<tr v-if="loading && students.length === 0">
						<td colspan="6" class="py-12 text-center text-slate-400">Loading student directory...</td>
					</tr>
					<tr v-else-if="students.length === 0">
						<td colspan="6" class="py-12 text-center text-slate-400">No students found matching your criteria.</td>
					</tr>
					<tr v-for="student in students" :key="student.ID" class="hover:bg-slate-50/70 transition">
						<td class="py-3 px-4 font-mono font-medium text-slate-900">{{ student.AdmissionNo }}</td>
						<td class="py-3 px-4 font-semibold text-slate-900">
							<router-link :to="`/students/${student.ID}`" class="hover:text-indigo-600 transition">
								{{ student.FirstName }} {{ student.MiddleName ? student.MiddleName + ' ' : '' }}{{ student.LastName }}
							</router-link>
						</td>
						<td class="py-3 px-4 capitalize text-slate-600">{{ student.Gender }}</td>
						<td class="py-3 px-4 text-slate-600">{{ formatDate(student.DateOfBirth) }}</td>
						<td class="py-3 px-4">
							<span
								class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold"
								:class="getStatusBadgeClass(student.Status)"
							>
								<span class="w-1.5 h-1.5 rounded-full" :class="getStatusDotClass(student.Status)"></span>
								{{ student.Status }}
							</span>
						</td>
						<td class="py-3 px-4 text-right space-x-2">
							<router-link
								:to="`/students/${student.ID}`"
								class="inline-flex items-center gap-1 text-[11px] font-medium text-indigo-600 hover:text-indigo-800 bg-indigo-50 hover:bg-indigo-100 px-2 py-1 rounded transition"
							>
								<Eye class="w-3.5 h-3.5" />
								<span>Profile</span>
							</router-link>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Register Student Modal -->
		<div v-if="showRegisterModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-xl w-full p-6 border border-slate-200 overflow-y-auto max-h-[90vh]">
				<h2 class="text-base font-bold text-slate-900 mb-1">Register New Student</h2>
				<p class="text-xs text-slate-500 mb-4">Enroll a student directly into an active academic year and class.</p>

				<form @submit.prevent="submitRegister" class="space-y-4">
					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">First Name *</label>
							<input v-model="form.firstName" type="text" required class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" />
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Last Name *</label>
							<input v-model="form.lastName" type="text" required class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" />
						</div>
					</div>

					<div class="grid grid-cols-3 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Gender *</label>
							<select v-model="form.gender" required class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20">
								<option value="male">Male</option>
								<option value="female">Female</option>
								<option value="other">Other</option>
							</select>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Date of Birth *</label>
							<input v-model="form.dobString" type="date" required class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" />
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Admission Date *</label>
							<input v-model="form.admissionDateString" type="date" required class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" />
						</div>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Academic Year *</label>
							<select v-model="form.academicYearId" required class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20">
								<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
							</select>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Class *</label>
							<select v-model="form.classId" @change="loadSectionsForClass" required class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20">
								<option v-for="c in classes" :key="c.ID" :value="c.ID">{{ c.Name }}</option>
							</select>
						</div>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Section</label>
							<select v-model="form.sectionId" class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20">
								<option value="">No Section</option>
								<option v-for="s in sections" :key="s.ID" :value="s.ID">{{ s.Name }}</option>
							</select>
						</div>
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Roll Number</label>
							<input v-model="form.rollNo" type="text" placeholder="e.g. 01" class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" />
						</div>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">National ID / Birth Cert #</label>
						<input v-model="form.nationalId" type="text" class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" />
					</div>

					<div v-if="errorMessage" class="text-xs text-rose-600 bg-rose-50 border border-rose-200 p-2.5 rounded-lg">
						{{ errorMessage }}
					</div>

					<div class="flex items-center justify-end gap-2 pt-2 border-t border-slate-100">
						<button
							type="button"
							@click="showRegisterModal = false"
							class="px-3 py-1.5 border border-slate-300 text-slate-700 rounded-lg text-xs hover:bg-slate-50 transition cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="submitting"
							class="px-4 py-1.5 bg-indigo-600 text-white rounded-lg text-xs font-medium hover:bg-indigo-700 transition disabled:opacity-50 cursor-pointer"
						>
							{{ submitting ? 'Registering...' : 'Register Student' }}
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
import type { StudentDTO, AcademicYearDTO, ClassDTO, SectionDTO } from '@/types/common'
import { UserPlus, Search, RefreshCw, Eye } from 'lucide-vue-next'

const services = getServices()

const students = ref<StudentDTO[]>([])
const academicYears = ref<AcademicYearDTO[]>([])
const classes = ref<ClassDTO[]>([])
const sections = ref<SectionDTO[]>([])

const searchQuery = ref('')
const statusFilter = ref('')
const loading = ref(false)
const showRegisterModal = ref(false)
const submitting = ref(false)
const errorMessage = ref('')

let searchTimer: any = null

const form = ref({
	firstName: '',
	lastName: '',
	gender: 'male' as const,
	dobString: '2018-01-01',
	admissionDateString: new Date().toISOString().split('T')[0],
	academicYearId: '',
	classId: '',
	sectionId: '',
	rollNo: '',
	nationalId: ''
})

async function loadStudents() {
	loading.value = true
	try {
		const res = await services.students.listStudents({
			status: statusFilter.value,
			searchQuery: searchQuery.value,
			limit: 50
		})
		if (res.success && res.data) {
			students.value = res.data
		}
	} finally {
		loading.value = false
	}
}

function debounceSearch() {
	clearTimeout(searchTimer)
	searchTimer = setTimeout(() => {
		loadStudents()
	}, 300)
}

async function loadMeta() {
	const [yearsRes, classesRes] = await Promise.all([
		services.academics.listAcademicYears(),
		services.academics.listClasses()
	])
	if (yearsRes.success && yearsRes.data) {
		academicYears.value = yearsRes.data
		const active = academicYears.value.find(y => y.Status === 'active')
		if (active) form.value.academicYearId = active.ID
		else if (academicYears.value.length > 0) form.value.academicYearId = academicYears.value[0].ID
	}
	if (classesRes.success && classesRes.data) {
		classes.value = classesRes.data
		if (classes.value.length > 0) {
			form.value.classId = classes.value[0].ID
			loadSectionsForClass()
		}
	}
}

async function loadSectionsForClass() {
	if (!form.value.academicYearId) return
	const res = await services.academics.listSectionsByYear(form.value.academicYearId)
	if (res.success && res.data) {
		sections.value = res.data.filter(s => s.ClassID === form.value.classId)
	}
}

function openRegisterModal() {
	errorMessage.value = ''
	showRegisterModal.value = true
}

function parseCivil(dateStr: string) {
	const [y, m, d] = dateStr.split('-').map(Number)
	return { Year: y, Month: m, Day: d }
}

async function submitRegister() {
	submitting.value = true
	errorMessage.value = ''
	try {
		const res = await services.students.registerStudent({
			firstName: form.value.firstName,
			lastName: form.value.lastName,
			gender: form.value.gender,
			dateOfBirth: parseCivil(form.value.dobString),
			admissionDate: parseCivil(form.value.admissionDateString),
			academicYearId: form.value.academicYearId,
			classId: form.value.classId,
			sectionId: form.value.sectionId || undefined,
			rollNo: form.value.rollNo || undefined,
			nationalId: form.value.nationalId
		})
		if (res.success) {
			showRegisterModal.value = false
			await loadStudents()
		} else {
			errorMessage.value = res.error?.message || 'Failed to register student'
		}
	} catch (err: any) {
		errorMessage.value = err.message || 'Error occurred'
	} finally {
		submitting.value = false
	}
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

function getStatusDotClass(status: string) {
	switch (status) {
		case 'active':
			return 'bg-emerald-500'
		case 'graduated':
			return 'bg-blue-500'
		case 'withdrawn':
			return 'bg-rose-500'
		default:
			return 'bg-slate-400'
	}
}

onMounted(() => {
	loadStudents()
	loadMeta()
})
</script>

