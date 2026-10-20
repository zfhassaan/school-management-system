<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Academic Structure</h1>
				<p class="text-xs text-slate-500 mt-0.5">Configure academic years, terms, classes, sections, and campuses.</p>
			</div>
			<div class="flex gap-2">
				<button
					v-if="currentTab === 'years'"
					@click="showYearModal = true"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>Add Academic Year</span>
				</button>
				<button
					v-if="currentTab === 'classes'"
					@click="showClassModal = true"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>Add Class</span>
				</button>
			</div>
		</div>

		<!-- Tabs -->
		<div class="flex border-b border-slate-200 space-x-6 text-xs font-medium">
			<button
				v-for="tab in ['years', 'classes', 'campuses']"
				:key="tab"
				@click="currentTab = tab"
				class="pb-3 capitalize cursor-pointer transition border-b-2"
				:class="currentTab === tab ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'"
			>
				{{ tab === 'years' ? 'Academic Years' : tab === 'classes' ? 'Classes & Sections' : 'Campuses' }}
			</button>
		</div>

		<!-- Tab: Academic Years -->
		<div v-if="currentTab === 'years'" class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">Year Name</th>
						<th class="py-3 px-4">Code</th>
						<th class="py-3 px-4">Start Date</th>
						<th class="py-3 px-4">End Date</th>
						<th class="py-3 px-4">Status</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-slate-700">
					<tr v-for="y in academicYears" :key="y.ID" class="hover:bg-slate-50/50">
						<td class="py-3 px-4 font-semibold text-slate-900">{{ y.Name }}</td>
						<td class="py-3 px-4 font-mono">{{ y.Code }}</td>
						<td class="py-3 px-4">{{ formatDate(y.StartsOn) }}</td>
						<td class="py-3 px-4">{{ formatDate(y.EndsOn) }}</td>
						<td class="py-3 px-4">
							<span
								class="px-2 py-0.5 rounded-full text-[10px] font-semibold capitalize"
								:class="y.Status === 'active' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-slate-100 text-slate-600 border border-slate-200'"
							>
								{{ y.Status }}
							</span>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Tab: Classes & Sections -->
		<div v-if="currentTab === 'classes'" class="space-y-4">
			<div v-for="cls in classes" :key="cls.ID" class="bg-white rounded-xl border border-slate-200 p-4 shadow-xs">
				<div class="flex items-center justify-between pb-3 border-b border-slate-100">
					<div>
						<span class="text-xs font-bold text-slate-900">{{ cls.Name }}</span>
						<span class="ml-2 font-mono text-[11px] text-slate-400">Code: {{ cls.Code }} • Level: {{ cls.LevelOrder }}</span>
					</div>
					<button
						@click="openAddSectionModal(cls.ID)"
						class="text-[11px] text-indigo-600 hover:text-indigo-800 font-medium px-2 py-1 bg-indigo-50 rounded cursor-pointer"
					>
						+ Add Section
					</button>
				</div>

				<div class="pt-3">
					<span class="text-[11px] font-bold text-slate-400 uppercase tracking-wider block mb-2">Sections</span>
					<div class="flex flex-wrap gap-2">
						<span
							v-for="sec in getSectionsForClass(cls.ID)"
							:key="sec.ID"
							class="px-3 py-1 bg-slate-50 border border-slate-200 rounded-lg text-xs font-medium text-slate-700 flex items-center gap-2"
						>
							<span>{{ sec.Name }}</span>
							<span class="text-[10px] text-slate-400">Cap: {{ sec.Capacity || '—' }}</span>
						</span>
						<span v-if="getSectionsForClass(cls.ID).length === 0" class="text-xs text-slate-400 italic">No sections created yet</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Tab: Campuses -->
		<div v-if="currentTab === 'campuses'" class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">Campus Name</th>
						<th class="py-3 px-4">Code</th>
						<th class="py-3 px-4">Address</th>
						<th class="py-3 px-4">Status</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-slate-700">
					<tr v-for="c in campuses" :key="c.ID" class="hover:bg-slate-50/50">
						<td class="py-3 px-4 font-semibold text-slate-900">{{ c.Name }}</td>
						<td class="py-3 px-4 font-mono">{{ c.Code }}</td>
						<td class="py-3 px-4">{{ c.Address || '—' }}</td>
						<td class="py-3 px-4">
							<span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
								Active
							</span>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Create Year Modal -->
		<div v-if="showYearModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-5 border border-slate-200">
				<h3 class="text-sm font-bold text-slate-900 mb-3">Add Academic Year</h3>
				<form @submit.prevent="submitCreateYear" class="space-y-3 text-xs">
					<div>
						<label class="block font-medium text-slate-700 mb-1">Name *</label>
						<input v-model="yearForm.name" placeholder="Academic Year 2026-2027" required class="w-full border border-slate-300 rounded-lg p-2" />
					</div>
					<div>
						<label class="block font-medium text-slate-700 mb-1">Code *</label>
						<input v-model="yearForm.code" placeholder="2026-27" required class="w-full border border-slate-300 rounded-lg p-2 font-mono" />
					</div>
					<div class="grid grid-cols-2 gap-2">
						<div>
							<label class="block font-medium text-slate-700 mb-1">Starts On *</label>
							<input v-model="yearForm.startsOn" type="date" required class="w-full border border-slate-300 rounded-lg p-2" />
						</div>
						<div>
							<label class="block font-medium text-slate-700 mb-1">Ends On *</label>
							<input v-model="yearForm.endsOn" type="date" required class="w-full border border-slate-300 rounded-lg p-2" />
						</div>
					</div>
					<div class="flex justify-end gap-2 pt-2">
						<button type="button" @click="showYearModal = false" class="px-3 py-1.5 border border-slate-300 rounded-lg">Cancel</button>
						<button type="submit" class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg font-medium">Save Year</button>
					</div>
				</form>
			</div>
		</div>

		<!-- Create Class Modal -->
		<div v-if="showClassModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-5 border border-slate-200">
				<h3 class="text-sm font-bold text-slate-900 mb-3">Add Class</h3>
				<form @submit.prevent="submitCreateClass" class="space-y-3 text-xs">
					<div>
						<label class="block font-medium text-slate-700 mb-1">Class Name *</label>
						<input v-model="classForm.name" placeholder="Grade 1" required class="w-full border border-slate-300 rounded-lg p-2" />
					</div>
					<div class="grid grid-cols-2 gap-2">
						<div>
							<label class="block font-medium text-slate-700 mb-1">Code *</label>
							<input v-model="classForm.code" placeholder="G1" required class="w-full border border-slate-300 rounded-lg p-2 font-mono" />
						</div>
						<div>
							<label class="block font-medium text-slate-700 mb-1">Level Order *</label>
							<input v-model.number="classForm.levelOrder" type="number" required class="w-full border border-slate-300 rounded-lg p-2" />
						</div>
					</div>
					<div class="flex justify-end gap-2 pt-2">
						<button type="button" @click="showClassModal = false" class="px-3 py-1.5 border border-slate-300 rounded-lg">Cancel</button>
						<button type="submit" class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg font-medium">Save Class</button>
					</div>
				</form>
			</div>
		</div>

		<!-- Create Section Modal -->
		<div v-if="showSectionModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl max-w-sm w-full p-5 border border-slate-200">
				<h3 class="text-sm font-bold text-slate-900 mb-3">Add Section</h3>
				<form @submit.prevent="submitCreateSection" class="space-y-3 text-xs">
					<div>
						<label class="block font-medium text-slate-700 mb-1">Section Name *</label>
						<input v-model="sectionForm.name" placeholder="Section A" required class="w-full border border-slate-300 rounded-lg p-2" />
					</div>
					<div>
						<label class="block font-medium text-slate-700 mb-1">Capacity</label>
						<input v-model.number="sectionForm.capacity" type="number" placeholder="40" class="w-full border border-slate-300 rounded-lg p-2" />
					</div>
					<div class="flex justify-end gap-2 pt-2">
						<button type="button" @click="showSectionModal = false" class="px-3 py-1.5 border border-slate-300 rounded-lg">Cancel</button>
						<button type="submit" class="px-3 py-1.5 bg-indigo-600 text-white rounded-lg font-medium">Save Section</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import type { AcademicYearDTO, ClassDTO, SectionDTO, CampusDTO } from '@/types/common'
import { Plus } from 'lucide-vue-next'

const services = getServices()

const currentTab = ref('years')
const academicYears = ref<AcademicYearDTO[]>([])
const classes = ref<ClassDTO[]>([])
const allSections = ref<SectionDTO[]>([])
const campuses = ref<CampusDTO[]>([])

const showYearModal = ref(false)
const showClassModal = ref(false)
const showSectionModal = ref(false)
const targetClassId = ref('')

const yearForm = ref({
	name: '',
	code: '',
	startsOn: '2026-09-01',
	endsOn: '2027-06-30'
})

const classForm = ref({
	name: '',
	code: '',
	levelOrder: 1
})

const sectionForm = ref({
	name: '',
	capacity: 40
})

async function loadData() {
	const [yearsRes, classesRes, campusesRes] = await Promise.all([
		services.academics.listAcademicYears(),
		services.academics.listClasses(),
		services.academics.listCampuses()
	])
	if (yearsRes.success && yearsRes.data) academicYears.value = yearsRes.data
	if (classesRes.success && classesRes.data) classes.value = classesRes.data
	if (campusesRes.success && campusesRes.data) campuses.value = campusesRes.data

	if (academicYears.value.length > 0) {
		const active = academicYears.value.find(y => y.Status === 'active') || academicYears.value[0]
		const secRes = await services.academics.listSectionsByYear(active.ID)
		if (secRes.success && secRes.data) allSections.value = secRes.data
	}
}

function getSectionsForClass(classId: string) {
	return allSections.value.filter(s => s.ClassID === classId)
}

function parseCivil(str: string) {
	const [y, m, d] = str.split('-').map(Number)
	return { Year: y, Month: m, Day: d }
}

async function submitCreateYear() {
	const res = await services.academics.createAcademicYear({
		name: yearForm.value.name,
		code: yearForm.value.code,
		startsOn: parseCivil(yearForm.value.startsOn),
		endsOn: parseCivil(yearForm.value.endsOn),
		status: 'active'
	})
	if (res.success) {
		showYearModal.value = false
		await loadData()
	}
}

async function submitCreateClass() {
	const res = await services.academics.createClass({
		name: classForm.value.name,
		code: classForm.value.code,
		levelOrder: classForm.value.levelOrder
	})
	if (res.success) {
		showClassModal.value = false
		await loadData()
	}
}

function openAddSectionModal(classId: string) {
	targetClassId.value = classId
	sectionForm.value.name = ''
	sectionForm.value.capacity = 40
	showSectionModal.value = true
}

async function submitCreateSection() {
	const activeYear = academicYears.value.find(y => y.Status === 'active') || academicYears.value[0]
	if (!activeYear) return
	const res = await services.academics.createSection({
		classId: targetClassId.value,
		academicYearId: activeYear.ID,
		name: sectionForm.value.name,
		capacity: sectionForm.value.capacity
	})
	if (res.success) {
		showSectionModal.value = false
		await loadData()
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

onMounted(() => {
	loadData()
})
</script>

