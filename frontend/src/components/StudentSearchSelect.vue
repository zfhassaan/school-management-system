<template>
	<div class="space-y-1 relative" ref="containerRef">
		<label v-if="label" class="block text-xs font-medium text-slate-700">
			{{ label }} <span v-if="required" class="text-rose-500">*</span>
		</label>

		<!-- Selected Student Card View -->
		<div
			v-if="selectedStudent && !isSearching"
			class="flex items-center justify-between p-2.5 bg-indigo-50/60 border border-indigo-200 rounded-lg text-xs"
		>
			<div class="flex items-center gap-2.5 min-w-0">
				<div class="w-7 h-7 rounded-full bg-indigo-600 text-white font-bold text-[10px] flex items-center justify-center shrink-0">
					{{ (selectedStudent.FirstName?.[0] || '') + (selectedStudent.LastName?.[0] || '') }}
				</div>
				<div class="truncate">
					<p class="font-bold text-slate-900 truncate">
						{{ selectedStudent.FirstName }} {{ selectedStudent.LastName }}
					</p>
					<p class="text-[10px] font-mono font-medium text-indigo-700">
						ID: {{ selectedStudent.AdmissionNo }}
					</p>
				</div>
			</div>
			<div class="flex items-center gap-1.5 shrink-0 ml-2">
				<button
					type="button"
					@click="startSearch"
					class="px-2 py-1 text-[11px] font-medium text-indigo-600 hover:text-indigo-800 hover:bg-white rounded border border-transparent hover:border-indigo-200 transition cursor-pointer"
				>
					Change
				</button>
				<button
					type="button"
					@click="clearSelection"
					class="p-1 text-slate-400 hover:text-rose-600 rounded transition cursor-pointer"
					title="Clear selection"
				>
					<X class="w-3.5 h-3.5" />
				</button>
			</div>
		</div>

		<!-- Search Input Bar & Dropdown -->
		<div v-else class="relative">
			<div class="relative">
				<Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-2.5" />
				<input
					ref="inputRef"
					v-model="searchQuery"
					type="text"
					:placeholder="placeholder"
					@focus="isOpen = true"
					@keydown.esc="isOpen = false"
					@keydown.down.prevent="navigateDown"
					@keydown.up.prevent="navigateUp"
					@keydown.enter.prevent="selectHighlighted"
					class="w-full pl-8.5 pr-8 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 bg-white"
				/>
				<button
					v-if="searchQuery"
					type="button"
					@click="searchQuery = ''"
					class="absolute right-2.5 top-2.5 text-slate-400 hover:text-slate-600 cursor-pointer"
				>
					<X class="w-3 h-3" />
				</button>
			</div>

			<!-- Floating Search Results Dropdown -->
			<div
				v-if="isOpen"
				class="absolute left-0 right-0 top-full mt-1 bg-white rounded-xl border border-slate-200 shadow-xl z-50 max-h-56 overflow-y-auto divide-y divide-slate-100 animate-in fade-in zoom-in-95 duration-100"
			>
				<div v-if="filteredStudents.length === 0" class="p-3 text-center text-slate-400 text-xs">
					No students found matching "<span class="font-medium text-slate-600">{{ searchQuery }}</span>"
				</div>

				<div
					v-for="(s, idx) in filteredStudents"
					:key="s.ID"
					@click="selectStudent(s)"
					@mouseenter="highlightedIndex = idx"
					class="p-2.5 flex items-center justify-between cursor-pointer transition"
					:class="highlightedIndex === idx ? 'bg-indigo-50/80 text-indigo-900' : 'hover:bg-slate-50 text-slate-700'"
				>
					<div class="flex items-center gap-2.5 min-w-0">
						<div class="w-6 h-6 rounded-full bg-slate-200 text-slate-700 font-bold text-[9px] flex items-center justify-center shrink-0">
							{{ (s.FirstName?.[0] || '') + (s.LastName?.[0] || '') }}
						</div>
						<div class="truncate">
							<p class="font-semibold text-xs text-slate-900 truncate">
								{{ s.FirstName }} {{ s.LastName }}
							</p>
							<p class="text-[10px] text-slate-500 font-mono">
								Admission: {{ s.AdmissionNo }}
							</p>
						</div>
					</div>
					<span class="text-[10px] font-mono font-medium px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 shrink-0 ml-2">
						{{ s.AdmissionNo }}
					</span>
				</div>
			</div>
		</div>

		<!-- Hidden native input to fulfill HTML5 required validation -->
		<input
			type="text"
			:value="modelValue"
			:required="required"
			tabindex="-1"
			class="opacity-0 pointer-events-none absolute bottom-0 left-0 h-0 w-0"
		/>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import type { StudentDTO } from '@/types/common'
import { Search, X } from 'lucide-vue-next'

const props = withDefaults(
	defineProps<{
		modelValue?: string
		students?: StudentDTO[]
		label?: string
		required?: boolean
		placeholder?: string
	}>(),
	{
		modelValue: '',
		students: () => [],
		label: '',
		required: false,
		placeholder: 'Search student by name or admission ID...'
	}
)

const emit = defineEmits<{
	(e: 'update:modelValue', value: string): void
	(e: 'select', student: StudentDTO | null): void
}>()

const containerRef = ref<HTMLElement | null>(null)
const inputRef = ref<HTMLInputElement | null>(null)

const searchQuery = ref('')
const isOpen = ref(false)
const isSearching = ref(false)
const highlightedIndex = ref(0)

const selectedStudent = computed(() => {
	if (!props.modelValue) return null
	return props.students.find(s => s.ID === props.modelValue) || null
})

const filteredStudents = computed(() => {
	const q = searchQuery.value.trim().toLowerCase()
	if (!q) return props.students.slice(0, 50)
	return props.students.filter(s => {
		const name = `${s.FirstName} ${s.LastName}`.toLowerCase()
		const adm = (s.AdmissionNo || '').toLowerCase()
		const id = (s.ID || '').toLowerCase()
		return name.includes(q) || adm.includes(q) || id.includes(q)
	}).slice(0, 50)
})

function selectStudent(student: StudentDTO) {
	emit('update:modelValue', student.ID)
	emit('select', student)
	isOpen.value = false
	isSearching.value = false
	searchQuery.value = ''
}

function clearSelection() {
	emit('update:modelValue', '')
	emit('select', null)
	isSearching.value = true
	searchQuery.value = ''
	nextTick(() => {
		inputRef.value?.focus()
		isOpen.value = true
	})
}

function startSearch() {
	isSearching.value = true
	nextTick(() => {
		inputRef.value?.focus()
		isOpen.value = true
	})
}

function navigateDown() {
	if (!isOpen.value) {
		isOpen.value = true
		return
	}
	if (highlightedIndex.value < filteredStudents.value.length - 1) {
		highlightedIndex.value++
	}
}

function navigateUp() {
	if (highlightedIndex.value > 0) {
		highlightedIndex.value--
	}
}

function selectHighlighted() {
	if (filteredStudents.value.length > 0 && highlightedIndex.value < filteredStudents.value.length) {
		selectStudent(filteredStudents.value[highlightedIndex.value])
	}
}

function handleClickOutside(e: MouseEvent) {
	if (containerRef.value && !containerRef.value.contains(e.target as Node)) {
		isOpen.value = false
		if (selectedStudent.value) {
			isSearching.value = false
		}
	}
}

onMounted(() => {
	document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
	document.removeEventListener('click', handleClickOutside)
})

watch(() => props.modelValue, (newVal) => {
	if (newVal) {
		isSearching.value = false
	}
})
</script>
