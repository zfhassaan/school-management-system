<template>
	<div class="fixed bottom-16 right-4 z-50 flex flex-col gap-2 max-w-md w-full pointer-events-none">
		<transition-group
			enter-active-class="transform ease-out duration-300 transition"
			enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
			enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
			leave-active-class="transition ease-in duration-100"
			leave-from-class="opacity-100"
			leave-to-class="opacity-0"
		>
			<div
				v-for="toast in toasts"
				:key="toast.id"
				class="pointer-events-auto p-4 rounded-xl border shadow-lg flex items-start gap-3 bg-white"
				:class="{
					'border-emerald-200 text-emerald-950 bg-emerald-50/90': toast.type === 'success',
					'border-rose-200 text-rose-950 bg-rose-50/90': toast.type === 'error',
					'border-sky-200 text-sky-950 bg-sky-50/90': toast.type === 'info',
					'border-amber-200 text-amber-950 bg-amber-50/90': toast.type === 'warning'
				}"
			>
				<div class="mt-0.5 shrink-0">
					<CheckCircle2 v-if="toast.type === 'success'" class="w-5 h-5 text-emerald-600" />
					<AlertCircle v-else-if="toast.type === 'error'" class="w-5 h-5 text-rose-600" />
					<Info v-else class="w-5 h-5 text-sky-600" />
				</div>
				<div class="flex-1">
					<h4 v-if="toast.title" class="text-sm font-semibold">{{ toast.title }}</h4>
					<p class="text-sm mt-0.5 opacity-90">{{ toast.message }}</p>
				</div>
				<button
					@click="remove(toast.id)"
					class="shrink-0 text-slate-400 hover:text-slate-600 p-1"
				>
					<X class="w-4 h-4" />
				</button>
			</div>
		</transition-group>
	</div>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'
import { CheckCircle2, AlertCircle, Info, X } from 'lucide-vue-next'

const { toasts, remove } = useToast()
</script>

