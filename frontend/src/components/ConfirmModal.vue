<template>
	<div v-if="show" class="fixed inset-0 bg-slate-900/50 backdrop-blur-xs flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 border border-slate-200 transform transition-all animate-in fade-in zoom-in-95 duration-150">
			<!-- Header / Icon -->
			<div class="flex items-start gap-3.5 mb-4">
				<div
					class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
					:class="{
						'bg-indigo-50 text-indigo-600': variant === 'primary',
						'bg-rose-50 text-rose-600': variant === 'danger',
						'bg-amber-50 text-amber-600': variant === 'warning',
						'bg-sky-50 text-sky-600': variant === 'info'
					}"
				>
					<AlertTriangle v-if="variant === 'warning'" class="w-5 h-5" />
					<Trash2 v-else-if="variant === 'danger'" class="w-5 h-5" />
					<Info v-else-if="variant === 'info'" class="w-5 h-5" />
					<HelpCircle v-else class="w-5 h-5" />
				</div>
				<div class="flex-1">
					<h3 class="text-base font-bold text-slate-900 leading-tight">{{ title }}</h3>
					<p v-if="message" class="text-xs text-slate-500 mt-1 leading-relaxed">{{ message }}</p>
				</div>
			</div>

			<!-- Slot for custom body if needed -->
			<div v-if="$slots.default" class="mb-5 text-xs text-slate-600">
				<slot />
			</div>

			<!-- Action buttons -->
			<div class="flex items-center justify-end gap-2.5 pt-2 border-t border-slate-100">
				<button
					type="button"
					@click="$emit('cancel')"
					class="px-4 py-2 border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition cursor-pointer"
					:disabled="loading"
				>
					{{ cancelText }}
				</button>
				<button
					type="button"
					@click="$emit('confirm')"
					class="px-4 py-2 font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
					:class="{
						'bg-indigo-600 hover:bg-indigo-700 text-white': variant === 'primary',
						'bg-rose-600 hover:bg-rose-700 text-white': variant === 'danger',
						'bg-amber-600 hover:bg-amber-700 text-white': variant === 'warning',
						'bg-sky-600 hover:bg-sky-700 text-white': variant === 'info'
					}"
					:disabled="loading"
				>
					<span v-if="loading" class="animate-spin mr-1">⏳</span>
					<span>{{ confirmText }}</span>
				</button>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { AlertTriangle, Trash2, Info, HelpCircle } from 'lucide-vue-next'

withDefaults(
	defineProps<{
		show: boolean
		title: string
		message?: string
		confirmText?: string
		cancelText?: string
		variant?: 'primary' | 'danger' | 'warning' | 'info'
		loading?: boolean
	}>(),
	{
		message: '',
		confirmText: 'Confirm',
		cancelText: 'Cancel',
		variant: 'primary',
		loading: false
	}
)

defineEmits<{
	(e: 'confirm'): void
	(e: 'cancel'): void
}>()
</script>

