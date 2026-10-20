<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Student Class Promotions</h1>
				<p class="text-xs text-slate-500 mt-0.5">Bulk transition students between academic years with audit trail and historical tracking.</p>
			</div>
			<button
				@click="requestPromotion"
				:disabled="executing || selectedStudentIds.length === 0 || isSameAcademicYear"
				class="py-2 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
				:title="isSameAcademicYear ? 'Select a different Target Academic Year to promote students' : ''"
			>
				<ArrowRight class="w-4 h-4" />
				<span>{{ executing ? 'Promoting...' : `Promote Selected (${selectedStudentIds.length})` }}</span>
			</button>
		</div>

		<!-- Step 1: Configuration Card -->
		<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
			<div class="flex items-center justify-between">
				<h3 class="text-xs font-bold uppercase tracking-wider text-slate-700">1. Select Transition Path</h3>
				<span v-if="academicYears.length > 1" class="text-[11px] text-slate-400">
					{{ academicYears.length }} Academic Years Configured
				</span>
			</div>

			<!-- Warning if same academic year -->
			<div
				v-if="isSameAcademicYear"
				class="p-3 bg-amber-50 border border-amber-200 rounded-lg flex items-start gap-2.5 text-xs text-amber-800"
			>
				<AlertTriangle class="w-4 h-4 text-amber-600 shrink-0 mt-0.5" />
				<div>
					<p class="font-bold">Academic Year Notice</p>
					<p class="mt-0.5 text-amber-700">
						Source and Target Academic Years are currently both set to <strong>{{ fromYearName || 'the same year' }}</strong>.
						Class promotions transition students into a subsequent academic year.
						<template v-if="academicYears.length <= 1">
							Only 1 academic year exists; please create the upcoming academic year (e.g. 2027-2028) in <em>Settings &gt; Academic Years</em> before running promotions.
						</template>
						<template v-else>
							Please choose the upcoming year under <em>To Academic Year</em>.
						</template>
					</p>
				</div>
			</div>

			<div class="grid grid-cols-4 gap-3 text-xs">
				<div>
					<label class="block font-medium text-slate-700 mb-1">From Academic Year *</label>
					<select v-model="fromYearId" @change="onSourceChange" class="w-full border border-slate-300 rounded-lg p-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20">
						<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
					</select>
				</div>

				<div>
					<label class="block font-medium text-slate-700 mb-1">Source Class *</label>
					<select v-model="sourceClassId" @change="onSourceChange" class="w-full border border-slate-300 rounded-lg p-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20">
						<option v-for="c in classes" :key="c.ID" :value="c.ID">{{ c.Name }}</option>
					</select>
				</div>

				<div>
					<label class="block font-medium text-slate-700 mb-1">To Academic Year *</label>
					<select v-model="toYearId" @change="loadTargetSections" class="w-full border border-slate-300 rounded-lg p-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20" :class="{ 'border-amber-400 bg-amber-50/30': isSameAcademicYear }">
						<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
					</select>
				</div>

				<div>
					<label class="block font-medium text-slate-700 mb-1">Target Class *</label>
					<select v-model="targetClassId" @change="loadTargetSections" class="w-full border border-slate-300 rounded-lg p-2 bg-white text-slate-800 focus:outline-none focus:ring-2 focus:ring-indigo-500/20">
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
					<button @click="selectAll" class="text-indigo-600 hover:text-indigo-800 font-medium cursor-pointer">Select All</button>
					<span class="text-slate-300">|</span>
					<button @click="deselectAll" class="text-slate-500 hover:text-slate-700 font-medium cursor-pointer">Deselect All</button>
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
						<td colspan="5" class="py-8 text-center text-slate-400 italic">No active students found in this class.</td>
					</tr>
					<tr v-for="item in roster" :key="item.ID" class="hover:bg-slate-50/50">
						<td class="py-2.5 px-4">
							<input type="checkbox" :value="item.StudentID" v-model="selectedStudentIds" class="rounded border-slate-300 text-indigo-600 cursor-pointer" />
						</td>
						<td class="py-2.5 px-4 font-mono font-medium text-slate-700">
							{{ studentsMap[item.StudentID]?.AdmissionNo || item.RollNo || item.StudentID.slice(0, 8) }}
						</td>
						<td class="py-2.5 px-4 font-semibold text-slate-900">
							{{ studentsMap[item.StudentID] ? `${studentsMap[item.StudentID].FirstName} ${studentsMap[item.StudentID].LastName}` : (item.ClassName ? `${item.ClassName} Student` : 'Student') }}
						</td>
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

		<!-- In-App Modal: Confirm Promotion Execution -->
		<ConfirmModal
			:show="showConfirmModal"
			title="Confirm Student Class Promotion"
			confirm-text="Confirm Promotion"
			cancel-text="Cancel"
			variant="primary"
			:loading="executing"
			@confirm="submitPromotionBatch"
			@cancel="showConfirmModal = false"
		>
			<div class="space-y-3 bg-slate-50 p-3.5 rounded-xl border border-slate-200">
				<div class="grid grid-cols-2 gap-3 text-xs">
					<div>
						<span class="text-[10px] uppercase font-bold text-slate-400 block tracking-wider">From Path</span>
						<p class="font-bold text-slate-800">{{ fromYearName }}</p>
						<p class="text-indigo-600 font-medium">{{ sourceClassName }}</p>
					</div>
					<div>
						<span class="text-[10px] uppercase font-bold text-slate-400 block tracking-wider">Target Path</span>
						<p class="font-bold text-slate-800">{{ toYearName }}</p>
						<p class="text-indigo-600 font-medium">{{ targetClassName }}</p>
					</div>
				</div>

				<div class="pt-2.5 border-t border-slate-200/80 flex items-center justify-between text-xs text-slate-700">
					<span>Students to Transition:</span>
					<span class="font-bold text-indigo-700 bg-indigo-100/70 px-2.5 py-0.5 rounded-full text-xs">
						{{ selectedStudentIds.length }} Selected
					</span>
				</div>
			</div>
		</ConfirmModal>

		<!-- In-App Modal: Alert / Status Modal -->
		<ConfirmModal
			:show="statusModal.show"
			:title="statusModal.title"
			:message="statusModal.message"
			:variant="statusModal.variant"
			:confirm-text="statusModal.confirmText"
			:cancel-text="statusModal.cancelText"
			:hide-cancel="statusModal.hideCancel"
			@confirm="handleStatusConfirm"
			@cancel="statusModal.show = false"
		/>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getServices } from '@/services/factory'
import type { AcademicYearDTO, ClassDTO, SectionDTO, EnrollmentDTO, StudentDTO } from '@/types/common'
import { ArrowRight, AlertTriangle } from 'lucide-vue-next'
import ConfirmModal from '@/components/ConfirmModal.vue'

const router = useRouter()
const services = getServices()

const academicYears = ref<AcademicYearDTO[]>([])
const classes = ref<ClassDTO[]>([])
const targetSections = ref<SectionDTO[]>([])
const roster = ref<EnrollmentDTO[]>([])
const studentsMap = ref<Record<string, StudentDTO>>({})

const fromYearId = ref('')
const toYearId = ref('')
const sourceClassId = ref('')
const targetClassId = ref('')

const selectedStudentIds = ref<string[]>([])
const targetSectionMap = ref<Record<string, string>>({})
const actionMap = ref<Record<string, string>>({})

const loadingStudents = ref(false)
const executing = ref(false)
const showConfirmModal = ref(false)

interface StatusModalState {
	show: boolean
	title: string
	message: string
	variant: 'primary' | 'danger' | 'warning' | 'info' | 'success'
	confirmText: string
	cancelText: string
	hideCancel: boolean
	onConfirm?: () => void
}

const statusModal = ref<StatusModalState>({
	show: false,
	title: '',
	message: '',
	variant: 'info',
	confirmText: 'OK',
	cancelText: 'Cancel',
	hideCancel: true
})

function openStatusModal(options: Partial<StatusModalState>) {
	statusModal.value = {
		show: true,
		title: options.title || '',
		message: options.message || '',
		variant: options.variant || 'info',
		confirmText: options.confirmText || 'OK',
		cancelText: options.cancelText || 'Cancel',
		hideCancel: options.hideCancel ?? true,
		onConfirm: options.onConfirm
	}
}

function handleStatusConfirm() {
	const callback = statusModal.value.onConfirm
	statusModal.value.show = false
	if (callback) {
		callback()
	}
}

const fromYearName = computed(() => academicYears.value.find(y => y.ID === fromYearId.value)?.Name || '')
const toYearName = computed(() => academicYears.value.find(y => y.ID === toYearId.value)?.Name || '')
const sourceClassName = computed(() => classes.value.find(c => c.ID === sourceClassId.value)?.Name || '')
const targetClassName = computed(() => classes.value.find(c => c.ID === targetClassId.value)?.Name || '')
const isSameAcademicYear = computed(() => !!fromYearId.value && !!toYearId.value && fromYearId.value === toYearId.value)

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

		// If multiple academic years exist, pick the other one for target year by default
		if (academicYears.value.length > 1) {
			const otherYear = academicYears.value.find(y => y.ID !== fromYearId.value)
			toYearId.value = otherYear ? otherYear.ID : academicYears.value[1].ID
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

	await loadStudents()
}

async function loadStudents() {
	try {
		const res = await services.students.listStudents({ limit: 500 })
		if (res.success && res.data) {
			const map: Record<string, StudentDTO> = {}
			for (const s of res.data) {
				map[s.ID] = s
			}
			studentsMap.value = map
		}
	} catch (e) {
		console.error('Failed to load students profiles', e)
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

function requestPromotion() {
	if (selectedStudentIds.value.length === 0) {
		openStatusModal({
			title: 'No Students Selected',
			message: 'Please select at least one student from the roster to execute promotion.',
			variant: 'warning',
			confirmText: 'Dismiss',
			hideCancel: true
		})
		return
	}

	if (isSameAcademicYear.value) {
		openStatusModal({
			title: 'Different Academic Years Required',
			message: 'The source and target academic years cannot be the same. Class promotions transition students into a subsequent academic year.\n\nPlease select a different target academic year, or configure the upcoming academic year in Academic Settings.',
			variant: 'warning',
			confirmText: 'Understood',
			hideCancel: true
		})
		return
	}

	if (!targetClassId.value) {
		openStatusModal({
			title: 'Target Class Required',
			message: 'Please select a Target Class before running promotion.',
			variant: 'warning',
			confirmText: 'Dismiss',
			hideCancel: true
		})
		return
	}

	showConfirmModal.value = true
}

async function submitPromotionBatch() {
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

		showConfirmModal.value = false

		if (res.success) {
			const count = selectedStudentIds.value.length
			openStatusModal({
				title: 'Promotion Completed Successfully',
				message: `Successfully processed promotions for ${count} student(s) into ${targetClassName.value} (${toYearName.value}). Student enrollments have been updated.`,
				variant: 'success',
				confirmText: 'View Students Directory',
				cancelText: 'Stay on Promotions',
				hideCancel: false,
				onConfirm: () => {
					router.push('/students')
				}
			})
			await onSourceChange()
		} else {
			let errorMsg = res.error?.message || 'Failed to execute promotion'
			if (errorMsg.includes('from_academic_year_id <> to_academic_year_id') || errorMsg.includes('cannot be identical')) {
				errorMsg = 'Source and target academic years cannot be identical. Promotions require distinct academic years.'
			}
			openStatusModal({
				title: 'Promotion Execution Failed',
				message: errorMsg,
				variant: 'danger',
				confirmText: 'Dismiss',
				hideCancel: true
			})
		}
	} catch (e: any) {
		showConfirmModal.value = false
		openStatusModal({
			title: 'Promotion Execution Failed',
			message: e?.message || 'An unexpected error occurred while executing the promotion batch.',
			variant: 'danger',
			confirmText: 'Dismiss',
			hideCancel: true
		})
	} finally {
		executing.value = false
	}
}

onMounted(() => {
	loadMeta()
})
</script>
