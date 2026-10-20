<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Student Class Promotions</h1>
				<p class="text-xs text-slate-500 mt-0.5">Bulk transition students between academic years with audit trail and historical tracking.</p>
			</div>
			<button
				@click="executePromotion"
				:disabled="executing || selectedStudentIds.length === 0"
				class="py-2 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
			>
				<ArrowRight class="w-4 h-4" />
				<span>{{ executing ? 'Promoting...' : `Promote Selected (${selectedStudentIds.length})` }}</span>
			</button>
		</div>

		<!-- Step 1: Configuration Card -->
		<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
			<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">1. Select Transition Path</h3>

			<div class="grid grid-cols-4 gap-3 text-xs">
				<div>
					<label class="block font-medium text-slate-700 mb-1">From Academic Year *</label>
					<select v-model="fromYearId" @change="onSourceChange" class="w-full border border-slate-300 rounded-lg p-2 bg-white">
						<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
					</select>
				</div>

				<div>
					<label class="block font-medium text-slate-700 mb-1">Source Class *</label>
					<select v-model="sourceClassId" @change="onSourceChange" class="w-full border border-slate-300 rounded-lg p-2 bg-white">
						<option v-for="c in classes" :key="c.ID" :value="c.ID">{{ c.Name }}</option>
					</select>
				</div>

				<div>
					<label class="block font-medium text-slate-700 mb-1">To Academic Year *</label>
					<select v-model="toYearId" class="w-full border border-slate-300 rounded-lg p-2 bg-white">
						<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
					</select>
				</div>

				<div>
					<label class="block font-medium text-slate-700 mb-1">Target Class *</label>
					<select v-model="targetClassId" @change="loadTargetSections" class="w-full border border-slate-300 rounded-lg p-2 bg-white">
						<option v-for="c in classes" :key="c.ID" :value="c.ID">{{ c.Name }}</option>
					</select>
				</div>
			</div>
		</div>

		<!-- Step 2: Student Selection & Promotion Roster -->
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<div class="p-4 border-b border-slate-100 flex items-center justify-between">
				<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">2. Select Students to Promote</h3>
				<div class="flex items-center gap-2 text-xs">
					<button @click="selectAll" class="text-indigo-600 hover:text-indigo-800 font-medium">Select All</button>
					<span class="text-slate-300">|</span>
					<button @click="deselectAll" class="text-slate-500 hover:text-slate-700 font-medium">Deselect All</button>
				</div>
			</div>

			<table class="w-full text-left border-collapse text-xs">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-2.5 px-4 w-10">Select</th>
						<th class="py-2.5 px-4">Admission #</th>
						<th class="py-2.5 px-4">Student Name</th>
						<th class="py-2.5 px-4">Target Section</th>
						<th class="py-2.5 px-4">Action</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-slate-700">
					<tr v-if="loadingStudents">
						<td colspan="5" class="py-8 text-center text-slate-400">Loading class roster...</td>
					</tr>
					<tr v-else-if="roster.length === 0">
						<td colspan="5" class="py-8 text-center text-slate-400">No active students found in this class.</td>
					</tr>
					<tr v-for="item in roster" :key="item.ID" class="hover:bg-slate-50/50">
						<td class="py-2.5 px-4">
							<input type="checkbox" :value="item.StudentID" v-model="selectedStudentIds" class="rounded border-slate-300 text-indigo-600" />
						</td>
						<td class="py-2.5 px-4 font-mono font-medium">{{ item.RollNo || item.ID.slice(0, 8) }}</td>
						<td class="py-2.5 px-4 font-semibold text-slate-900">{{ item.ClassName }} Student ({{ item.StudentID.slice(0, 8) }})</td>
						<td class="py-2.5 px-4">
							<select v-model="targetSectionMap[item.StudentID]" class="border border-slate-300 rounded px-2 py-1 bg-white text-xs">
								<option value="">No Section</option>
								<option v-for="sec in targetSections" :key="sec.ID" :value="sec.ID">{{ sec.Name }}</option>
							</select>
						</td>
						<td class="py-2.5 px-4">
							<select v-model="actionMap[item.StudentID]" class="border border-slate-300 rounded px-2 py-1 bg-white text-xs">
								<option value="promoted">Promote</option>
								<option value="retained">Retain (Repeat)</option>
								<option value="graduated">Graduated</option>
							</select>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getServices } from '@/services/factory'
import type { AcademicYearDTO, ClassDTO, SectionDTO, EnrollmentDTO } from '@/types/common'
import { ArrowRight } from 'lucide-vue-next'

const router = useRouter()
const services = getServices()

const academicYears = ref<AcademicYearDTO[]>([])
const classes = ref<ClassDTO[]>([])
const targetSections = ref<SectionDTO[]>([])
const roster = ref<EnrollmentDTO[]>([])

const fromYearId = ref('')
const toYearId = ref('')
const sourceClassId = ref('')
const targetClassId = ref('')

const selectedStudentIds = ref<string[]>([])
const targetSectionMap = ref<Record<string, string>>({})
const actionMap = ref<Record<string, string>>({})

const loadingStudents = ref(false)
const executing = ref(false)

async function loadMeta() {
	const [yearsRes, classesRes] = await Promise.all([
		services.academics.listAcademicYears(),
		services.academics.listClasses()
	])
	if (yearsRes.success && yearsRes.data) {
		academicYears.value = yearsRes.data
		const active = academicYears.value.find(y => y.Status === 'active')
		if (active) fromYearId.value = active.ID
		else if (academicYears.value.length > 0) fromYearId.value = academicYears.value[0].ID

		if (academicYears.value.length > 1) {
			toYearId.value = academicYears.value[1].ID
		} else {
			toYearId.value = fromYearId.value
		}
	}
	if (classesRes.success && classesRes.data) {
		classes.value = classesRes.data
		if (classes.value.length > 0) {
			sourceClassId.value = classes.value[0].ID
			targetClassId.value = classes.value.length > 1 ? classes.value[1].ID : classes.value[0].ID
			await onSourceChange()
			await loadTargetSections()
		}
	}
}

async function onSourceChange() {
	if (!fromYearId.value) return
	loadingStudents.value = true
	try {
		// Load active enrollments for source class
		const secRes = await services.academics.listSectionsByYear(fromYearId.value)
		if (secRes.success && secRes.data) {
			const matchingSections = secRes.data.filter(s => s.ClassID === sourceClassId.value)
			const allEnrs: EnrollmentDTO[] = []
			for (const sec of matchingSections) {
				const rRes = await services.enrollments.listSectionRoster(sec.ID)
				if (rRes.success && rRes.data) allEnrs.push(...rRes.data)
			}
			roster.value = allEnrs
			selectedStudentIds.value = allEnrs.map(e => e.StudentID)
			allEnrs.forEach(e => {
				actionMap.value[e.StudentID] = 'promoted'
			})
		}
	} finally {
		loadingStudents.value = false
	}
}

async function loadTargetSections() {
	if (!toYearId.value) return
	const secRes = await services.academics.listSectionsByYear(toYearId.value)
	if (secRes.success && secRes.data) {
		targetSections.value = secRes.data.filter(s => s.ClassID === targetClassId.value)
	}
}

function selectAll() {
	selectedStudentIds.value = roster.value.map(e => e.StudentID)
}

function deselectAll() {
	selectedStudentIds.value = []
}

async function executePromotion() {
	if (selectedStudentIds.value.length === 0) return
	executing.value = true
	try {
		const today = new Date().toISOString().split('T')[0].split('-').map(Number)
		const items = selectedStudentIds.value.map(id => ({
			studentId: id,
			targetClassId: targetClassId.value,
			targetSectionId: targetSectionMap.value[id] || undefined,
			action: actionMap.value[id] || 'promoted'
		}))

		const res = await services.enrollments.executePromotionBatch({
			fromAcademicYearId: fromYearId.value,
			toAcademicYearId: toYearId.value,
			effectiveDate: { Year: today[0], Month: today[1], Day: today[2] },
			notes: 'Bulk promotion executed from UI',
			items
		})
		if (res.success) {
			alert(`Successfully promoted ${selectedStudentIds.value.length} students!`)
			router.push('/students')
		} else {
			alert(res.error?.message || 'Failed to execute promotion')
		}
	} finally {
		executing.value = false
	}
}

onMounted(() => {
	loadMeta()
})
</script>
