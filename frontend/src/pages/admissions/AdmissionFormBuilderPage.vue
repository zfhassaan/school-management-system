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
					<option v-for="d in definitions" :key="d.ID" :value="d.ID">{{ d.Name }} ({{ d.Code }})</option>
				</select>
			</div>

			<div v-if="currentVersion" class="flex items-center gap-3 text-xs">
				<span class="text-slate-500">Version: <strong class="text-slate-900">v{{ currentVersion.VersionNo }}</strong></span>
				<span
					class="px-2 py-0.5 rounded-full text-[10px] font-semibold uppercase"
					:class="currentVersion.Status === 'published' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-amber-50 text-amber-700 border border-amber-200'"
				>
					{{ currentVersion.Status }}
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
import type { FormDefinitionDTO, FormVersionDTO, FormSchemaDTO, FormSectionDTO } from '@/types/common'
import { ArrowLeft, CheckCircle2, Plus, Trash2, X } from 'lucide-vue-next'

const services = getServices()

const definitions = ref<FormDefinitionDTO[]>([])
const selectedFormDefId = ref('')
const currentVersion = ref<FormVersionDTO | null>(null)

const saving = ref(false)
const publishing = ref(false)

const schema = ref<FormSchemaDTO>({
	title: 'Standard Admission Application',
	description: 'Please submit accurate student and guardian details',
	sections: [
		{
			id: 'general_info',
			title: 'General Information',
			description: 'Previous schooling and language background',
			fields: [
				{
					id: 'prev_school',
					label: 'Previous School Attended',
					type: 'text',
					placeholder: 'Name of last school',
					required: false
				},
				{
					id: 'primary_language',
					label: 'Primary Language at Home',
					type: 'text',
					placeholder: 'e.g. English',
					required: true
				}
			]
		}
	]
})

async function loadDefinitions() {
	const res = await services.admissions.listFormDefinitions()
	if (res.success && res.data) {
		definitions.value = res.data
		if (definitions.value.length > 0) {
			selectedFormDefId.value = definitions.value[0].ID
			await loadPublishedVersion(selectedFormDefId.value)
		} else {
			// Create a default form definition if none exists
			const created = await services.admissions.createFormDefinition({
				code: 'STD-ADMISSION',
				name: 'General Admission Form',
				description: 'Standard single-school admission intake'
			})
			if (created.success && created.data) {
				await loadDefinitions()
			}
		}
	}
}

async function loadPublishedVersion(defId: string) {
	const res = await services.admissions.getPublishedFormVersion(defId)
	if (res.success && res.data) {
		currentVersion.value = res.data
		if (res.data.Schema && res.data.Schema.sections) {
			schema.value = res.data.Schema
		}
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
			await loadPublishedVersion(selectedFormDefId.value)
		}
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
				await loadPublishedVersion(selectedFormDefId.value)
			}
		}
	} finally {
		publishing.value = false
	}
}

onMounted(() => {
	loadDefinitions()
})
</script>
