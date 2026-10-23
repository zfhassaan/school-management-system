<template>
	<div class="space-y-6 max-w-5xl mx-auto">
		<!-- Header -->
		<div>
			<router-link to="/admissions" class="inline-flex items-center gap-1 text-xs text-slate-500 hover:text-slate-800 transition mb-2">
				<ArrowLeft class="w-3.5 h-3.5" />
				<span>Back to Applications</span>
			</router-link>

			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-xl font-bold text-slate-900">Admission Form Builder</h1>
					<p class="text-xs text-slate-500 mt-0.5">Design dynamic admission form fields, validation rules, and publish versioned schemas.</p>
				</div>
				<div class="flex items-center gap-2">
					<button
						@click="loadInstitutionalTemplate"
						type="button"
						class="px-3 py-1.5 border border-indigo-200 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 font-medium rounded-lg text-xs transition cursor-pointer flex items-center gap-1.5 shadow-2xs"
						title="Load complete standard institutional form template (4 sections, 28 fields)"
					>
						<RotateCcw class="w-3.5 h-3.5" />
						<span>Load Standard Template</span>
					</button>
					<button
						@click="saveDraftVersion"
						:disabled="saving"
						class="px-3.5 py-1.5 border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition cursor-pointer"
					>
						Save New Draft
					</button>
					<button
						@click="publishCurrentVersion"
						:disabled="publishing || currentVersion?.Status === 'published'"
						class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
					>
						<CheckCircle2 class="w-4 h-4" />
						<span>{{ currentVersion?.Status === 'published' ? 'Published' : 'Publish Version' }}</span>
					</button>
				</div>
			</div>
		</div>

		<!-- Definition and Version Selector -->
		<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex items-center justify-between gap-4">
			<div class="flex items-center gap-4 text-xs">
				<span class="font-medium text-slate-700">Form Definition:</span>
				<select v-model="selectedFormDefId" @change="onFormDefChange" class="border border-slate-300 rounded-lg py-1 px-3 bg-white font-medium">
					<option v-for="d in definitions" :key="d.ID || d.id" :value="d.ID || d.id">{{ d.Name || d.name }} ({{ d.Code || d.code }})</option>
				</select>
			</div>

			<div v-if="currentVersion" class="flex items-center gap-3 text-xs">
				<span class="text-slate-500">Version: <strong class="text-slate-900">v{{ currentVersion.VersionNo ?? currentVersion.versionNo }}</strong></span>
				<span
					class="px-2 py-0.5 rounded-full text-[10px] font-semibold uppercase"
					:class="(currentVersion.Status || currentVersion.status) === 'published' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-amber-50 text-amber-700 border border-amber-200'"
				>
					{{ currentVersion.Status || currentVersion.status }}
				</span>
			</div>
		</div>

		<!-- Form Sections & Fields Builder -->
		<div class="space-y-4">
			<div v-for="(section, sIndex) in schema.sections" :key="section.id" class="bg-white rounded-xl border border-slate-200 p-5 shadow-xs space-y-4">
				<div class="flex items-center justify-between border-b border-slate-100 pb-3">
					<div class="flex-1 max-w-md">
						<input
							v-model="section.title"
							placeholder="Section Title..."
							class="text-sm font-bold text-slate-900 w-full focus:outline-none hover:bg-slate-50 px-2 py-1 rounded"
						/>
						<input
							v-model="section.description"
							placeholder="Optional description / instructions..."
							class="text-xs text-slate-500 w-full focus:outline-none hover:bg-slate-50 px-2 py-0.5 rounded mt-0.5"
						/>
					</div>
					<div class="flex items-center gap-2">
						<button
							@click="addFieldToSection(section)"
							class="text-[11px] font-medium text-indigo-600 hover:text-indigo-800 bg-indigo-50 hover:bg-indigo-100 px-2.5 py-1 rounded transition cursor-pointer"
						>
							+ Add Field
						</button>
						<button
							v-if="schema.sections.length > 1"
							@click="removeSection(sIndex)"
							class="text-slate-400 hover:text-rose-600 p-1 rounded"
						>
							<Trash2 class="w-4 h-4" />
						</button>
					</div>
				</div>

				<!-- Section Fields List -->
				<div class="space-y-2">
					<div
						v-for="(field, fIndex) in section.fields"
						:key="field.id"
						class="p-3 border border-slate-200 rounded-lg bg-slate-50/50 flex flex-wrap items-center justify-between gap-3 text-xs"
					>
						<div class="flex items-center gap-3 flex-1 min-w-[300px]">
							<span class="font-mono text-slate-400 text-[10px] w-6">{{ fIndex + 1 }}.</span>
							<input
								v-model="field.label"
								placeholder="Field Label"
								class="bg-white border border-slate-300 rounded px-2 py-1 font-medium text-slate-800 w-44"
							/>
							<select v-model="field.type" class="bg-white border border-slate-300 rounded px-2 py-1 text-slate-700">
								<option value="text">Text Input</option>
								<option value="number">Number</option>
								<option value="date">Date</option>
								<option value="dropdown">Dropdown Options</option>
								<option value="phone">Phone Number</option>
								<option value="email">Email Address</option>
								<option value="textarea">Textarea (Multi-line)</option>
								<option value="checkbox">Checkbox</option>
								<option value="file_upload">File Upload</option>
							</select>
							<input
								v-model="field.placeholder"
								placeholder="Placeholder..."
								class="bg-white border border-slate-200 rounded px-2 py-1 text-slate-600 flex-1 text-[11px]"
							/>
						</div>

						<div class="flex items-center gap-3">
							<label class="flex items-center gap-1.5 cursor-pointer text-slate-600 text-[11px]">
								<input type="checkbox" v-model="field.required" class="rounded border-slate-300 text-indigo-600" />
								<span>Required</span>
							</label>

							<button @click="removeField(section, fIndex)" class="text-slate-400 hover:text-rose-600 p-1">
								<X class="w-3.5 h-3.5" />
							</button>
						</div>
					</div>
				</div>
			</div>

			<button
				@click="addSection"
				class="w-full py-3 border-2 border-dashed border-slate-200 hover:border-indigo-400 hover:bg-indigo-50/20 text-slate-500 hover:text-indigo-600 rounded-xl text-xs font-semibold transition flex items-center justify-center gap-1.5 cursor-pointer"
			>
				<Plus class="w-4 h-4" />
				<span>Add Form Section</span>
			</button>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import type { FormDefinitionDTO, FormVersionDTO, FormSchemaDTO, FormSectionDTO } from '@/types/common'
import { ArrowLeft, CheckCircle2, Plus, Trash2, X, RotateCcw } from 'lucide-vue-next'

const services = getServices()
const toast = useToast()

const definitions = ref<FormDefinitionDTO[]>([])
const selectedFormDefId = ref('')
const currentVersion = ref<FormVersionDTO | null>(null)

const saving = ref(false)
const publishing = ref(false)

function getStandardInstitutionalSchema(): FormSchemaDTO {
	return {
		title: 'Official Admission Application Form',
		description: 'Official admission intake form for student registration, parents details and institutional financials.',
		sections: [
			{
				id: 'sec_office_financials',
				title: 'SECTION 1: OFFICE USE & FINANCIALS',
				description: 'Institutional enrollment details, student identification and fee schedule',
				fields: [
					{
						id: 'b_form_cnic',
						label: 'B-Form / CNIC No.',
						type: 'text',
						placeholder: 'e.g. 34101-1234567-1',
						required: false
					},
					{
						id: 'adm_fee',
						label: 'Admission Fee (PKR)',
						type: 'number',
						placeholder: 'e.g. 2000',
						required: false
					},
					{
						id: 'monthly_fee',
						label: 'Monthly Tuition Fee (PKR)',
						type: 'number',
						placeholder: 'e.g. 2500',
						required: false
					},
					{
						id: 'annual_fee',
						label: 'Annual Charges (PKR)',
						type: 'number',
						placeholder: 'e.g. 3500',
						required: false
					},
					{
						id: 'paper_fund',
						label: 'Paper / Exam Fund (PKR)',
						type: 'number',
						placeholder: 'e.g. 1000',
						required: false
					},
					{
						id: 'id_card_fee',
						label: 'ID Card Fee (PKR)',
						type: 'number',
						placeholder: 'e.g. 500',
						required: false
					}
				]
			},
			{
				id: 'sec_student_personal',
				title: 'SECTION 2: STUDENT PERSONAL INFO',
				description: 'Student personal particulars, medical history and previous academic record',
				fields: [
					{
						id: 'medical_history',
						label: 'Medical History / Special Needs',
						type: 'text',
						placeholder: 'e.g. None / Asthma / Allergies',
						required: false
					},
					{
						id: 'home_address',
						label: 'Home / Residential Address',
						type: 'text',
						placeholder: 'House #, Street #, Area / City',
						required: false
					},
					{
						id: 'previous_school',
						label: 'Previous School Attended',
						type: 'text',
						placeholder: 'Name of previous school',
						required: false
					},
					{
						id: 'last_grade_result',
						label: 'Last Grade / Result',
						type: 'text',
						placeholder: 'e.g. Grade 1 - Passed',
						required: false
					},
					{
						id: 'blood_group',
						label: 'Blood Group',
						type: 'dropdown',
						placeholder: 'Select Blood Group',
						required: false,
						options: ['A+', 'A-', 'B+', 'B-', 'O+', 'O-', 'AB+', 'AB-']
					},
					{
						id: 'religion',
						label: 'Religion',
						type: 'text',
						placeholder: 'e.g. Islam / Christianity',
						required: false
					}
				]
			},
			{
				id: 'sec_parents_details',
				title: 'SECTION 3: PARENTS DETAILS',
				description: 'Father and mother contact details, occupations and qualifications',
				fields: [
					{
						id: 'father_name',
						label: "Father's Name",
						type: 'text',
						placeholder: "Father's full name",
						required: false
					},
					{
						id: 'father_whatsapp',
						label: 'Mobile (WhatsApp)',
						type: 'phone',
						placeholder: '0321 7404703',
						required: false
					},
					{
						id: 'father_cnic',
						label: 'Father / Guardian CNIC',
						type: 'text',
						placeholder: 'e.g. 34101-1234567-1',
						required: false
					},
					{
						id: 'father_qualification',
						label: 'Father Qualification',
						type: 'text',
						placeholder: 'e.g. Matric / Graduate / Master',
						required: false
					},
					{
						id: 'father_occupation',
						label: 'Father Occupation',
						type: 'text',
						placeholder: 'e.g. Business / Government / Private',
						required: false
					},
					{
						id: 'father_office_address',
						label: 'Office Address',
						type: 'text',
						placeholder: 'Office or business address',
						required: false
					},
					{
						id: 'father_office_phone',
						label: 'Office Phone',
						type: 'phone',
						placeholder: 'e.g. 042-1234567',
						required: false
					},
					{
						id: 'mother_name',
						label: "Mother's Name",
						type: 'text',
						placeholder: "Mother's full name",
						required: false
					},
					{
						id: 'mother_phone',
						label: 'Mother Mobile Number',
						type: 'phone',
						placeholder: '0320 7721922',
						required: false
					},
					{
						id: 'mother_qualification',
						label: 'Mother Qualification',
						type: 'text',
						placeholder: 'e.g. Matric / Graduate',
						required: false
					},
					{
						id: 'mother_occupation',
						label: 'Mother Occupation',
						type: 'text',
						placeholder: 'e.g. Housewife / Teacher',
						required: false
					}
				]
			},
			{
				id: 'sec_siblings_reference',
				title: 'SECTION 4: SIBLINGS & REFERENCE',
				description: 'Siblings currently enrolled and emergency contact reference',
				fields: [
					{
						id: 'reference_name',
						label: 'Reference Name',
						type: 'text',
						placeholder: 'Name of emergency reference',
						required: false
					},
					{
						id: 'reference_phone',
						label: 'Reference Phone',
						type: 'phone',
						placeholder: '0300 1234567',
						required: false
					},
					{
						id: 'sibling_in_school',
						label: 'Sibling in School?',
						type: 'text',
						placeholder: 'Yes / No',
						required: false
					},
					{
						id: 'sibling_class',
						label: 'Sibling Class',
						type: 'text',
						placeholder: 'e.g. Class 2',
						required: false
					},
					{
						id: 'sibling_since',
						label: 'Enrolled Since',
						type: 'text',
						placeholder: 'e.g. 2024',
						required: false
					}
				]
			}
		]
	}
}

const schema = ref<FormSchemaDTO>(getStandardInstitutionalSchema())

function loadInstitutionalTemplate() {
	schema.value = getStandardInstitutionalSchema()
	toast.success('Loaded complete institutional template (4 sections, 28 fields). Click "Publish Version" to activate.')
}

async function loadDefinitions() {
	const res = await services.admissions.listFormDefinitions()
	let defs: any[] = []
	if (res && res.success && Array.isArray(res.data)) {
		defs = res.data
	}

	// Auto-provision default institutional admission form if none exists
	if (defs.length === 0) {
		await services.admissions.ensureDefaultAdmissionForm()
		const reloadRes = await services.admissions.listFormDefinitions()
		if (reloadRes && reloadRes.success && Array.isArray(reloadRes.data)) {
			defs = reloadRes.data
		}
	}

	definitions.value = defs
	if (defs.length > 0) {
		const first = defs[0]
		selectedFormDefId.value = first.ID || first.id
		await loadPublishedVersion(selectedFormDefId.value)
	}
}

async function loadPublishedVersion(defId: string) {
	if (!defId) return
	const res = await services.admissions.getPublishedFormVersion(defId)
	if (res && res.success && res.data) {
		currentVersion.value = res.data
		const sc = res.data.Schema || res.data.schema
		if (sc && sc.sections && sc.sections.length >= 4) {
			schema.value = sc
		} else {
			// Legacy or incomplete schema detected (< 4 sections) - automatically upgrade to full template
			schema.value = getStandardInstitutionalSchema()
			toast.info('Standard institutional admission template loaded with complete fieldset.')
		}
	} else {
		currentVersion.value = null
		schema.value = getStandardInstitutionalSchema()
	}
}

function onFormDefChange() {
	if (selectedFormDefId.value) {
		loadPublishedVersion(selectedFormDefId.value)
	}
}

function addSection() {
	const id = `section_${Date.now()}`
	schema.value.sections.push({
		id,
		title: 'New Section',
		description: '',
		fields: []
	})
}

function removeSection(index: number) {
	schema.value.sections.splice(index, 1)
}

function addFieldToSection(section: FormSectionDTO) {
	const id = `field_${Date.now()}`
	section.fields.push({
		id,
		label: 'New Field',
		type: 'text',
		placeholder: '',
		required: false
	})
}

function removeField(section: FormSectionDTO, index: number) {
	section.fields.splice(index, 1)
}

async function saveDraftVersion() {
	if (!selectedFormDefId.value) return
	saving.value = true
	try {
		const res = await services.admissions.createFormVersion({
			formDefinitionId: selectedFormDefId.value,
			schema: schema.value
		})
		if (res.success) {
			toast.success('Form draft saved successfully!')
			await loadPublishedVersion(selectedFormDefId.value)
		} else {
			toast.error(res.error?.message || 'Failed to save draft')
		}
	} catch (e: any) {
		toast.error(e.message || 'Error saving draft')
	} finally {
		saving.value = false
	}
}

async function publishCurrentVersion() {
	if (!selectedFormDefId.value) return
	publishing.value = true
	try {
		// First ensure draft is saved
		const draftRes = await services.admissions.createFormVersion({
			formDefinitionId: selectedFormDefId.value,
			schema: schema.value
		})
		if (draftRes.success && draftRes.data) {
			const pubRes = await services.admissions.publishFormVersion({
				versionId: draftRes.data
			})
			if (pubRes.success) {
				toast.success('Admission form version published successfully!')
				await loadPublishedVersion(selectedFormDefId.value)
			} else {
				toast.error(pubRes.error?.message || 'Failed to publish form version')
			}
		} else {
			toast.error(draftRes.error?.message || 'Failed to create form version')
		}
	} catch (e: any) {
		toast.error(e.message || 'Error publishing version')
	} finally {
		publishing.value = false
	}
}

onMounted(() => {
	loadDefinitions()
})
</script>

