<template>
	<div
		v-if="isOpen"
		class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs transition-opacity duration-200"
		@click.self="emit('close')"
	>
		<div
			class="relative w-full max-w-lg overflow-hidden bg-white rounded-2xl shadow-2xl border border-slate-200 transform transition-all duration-200"
		>
			<!-- Gradient Header Banner -->
			<div class="h-28 bg-gradient-to-r from-slate-900 via-indigo-950 to-slate-900 relative p-6 flex items-start justify-between">
				<div class="flex items-center gap-2 text-white/80">
					<Sparkles class="w-4 h-4 text-cyan-400" />
					<span class="text-xs font-semibold tracking-wider uppercase text-cyan-200">Architect & Developer</span>
				</div>
				<button
					@click="emit('close')"
					class="p-1 rounded-lg text-white/60 hover:text-white hover:bg-white/10 transition"
					title="Close"
				>
					<X class="w-5 h-5" />
				</button>
			</div>

			<!-- Profile Avatar & Details Overlap -->
			<div class="px-6 pb-6 pt-0 relative">
				<div class="flex items-end justify-between -mt-10 mb-4">
					<div class="relative p-1.5 bg-white rounded-2xl shadow-md border border-slate-100">
						<DeveloperLogo sizeClass="w-16 h-16" />
					</div>
					<div class="flex items-center gap-2">
						<button
							@click="openUrl('https://github.com/zfhassaan')"
							class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-slate-900 hover:bg-slate-800 text-white flex items-center gap-1.5 transition shadow-xs"
						>
							<Github class="w-3.5 h-3.5" />
							<span>GitHub</span>
							<ExternalLink class="w-3 h-3 text-slate-400" />
						</button>
						<button
							@click="openUrl('https://zfhassaan.github.io')"
							class="px-3 py-1.5 rounded-lg text-xs font-semibold bg-indigo-600 hover:bg-indigo-700 text-white flex items-center gap-1.5 transition shadow-xs"
						>
							<Globe class="w-3.5 h-3.5" />
							<span>Portfolio</span>
							<ExternalLink class="w-3 h-3 text-indigo-200" />
						</button>
					</div>
				</div>

				<!-- Identity Text -->
				<div>
					<h3 class="text-lg font-bold text-slate-900 flex items-center gap-2">
						Syed Faizan Hassaan
						<span class="text-xs px-2 py-0.5 rounded-full font-mono font-medium bg-slate-100 text-slate-600 border border-slate-200">@zfhassaan</span>
					</h3>
					<p class="text-xs text-indigo-600 font-medium mt-0.5">Full Stack Software Architect & Systems Engineer</p>
					<p class="text-xs text-slate-500 mt-2 leading-relaxed">
						Creator and lead architect of the School Management System — engineered from the ground up as a resilient offline-first platform featuring exact integer currency arithmetic, strict sequential auditability, and clean hexagonal boundaries.
					</p>
				</div>

				<!-- Direct Links Card -->
				<div class="mt-4 p-3 bg-slate-50 rounded-xl border border-slate-200/80 space-y-2">
					<div class="flex items-center justify-between text-xs">
						<span class="text-slate-500 font-medium">GitHub Repository & Profile:</span>
						<button
							@click="openUrl('https://github.com/zfhassaan')"
							class="text-indigo-600 hover:text-indigo-800 font-semibold font-mono flex items-center gap-1 hover:underline"
						>
							github.com/zfhassaan
							<ExternalLink class="w-3 h-3" />
						</button>
					</div>
					<div class="flex items-center justify-between text-xs border-t border-slate-200/60 pt-2">
						<span class="text-slate-500 font-medium">Personal Portfolio & Work:</span>
						<button
							@click="openUrl('https://zfhassaan.github.io')"
							class="text-indigo-600 hover:text-indigo-800 font-semibold font-mono flex items-center gap-1 hover:underline"
						>
							zfhassaan.github.io
							<ExternalLink class="w-3 h-3" />
						</button>
					</div>
				</div>

				<!-- Architecture Highlights -->
				<div class="mt-4 pt-4 border-t border-slate-100 grid grid-cols-2 gap-2 text-[11px]">
					<div class="p-2 rounded-lg bg-indigo-50/50 border border-indigo-100 text-slate-700">
						<span class="block font-bold text-indigo-900">Backend Core</span>
						Go 1.26 • Hexagonal • SQLite WAL
					</div>
					<div class="p-2 rounded-lg bg-emerald-50/50 border border-emerald-100 text-slate-700">
						<span class="block font-bold text-emerald-900">Frontend Presentation</span>
						Vue 3 • Wails v2 • Tailwind v4
					</div>
				</div>

				<!-- Footer note -->
				<div class="mt-4 text-center">
					<p class="text-[10px] text-slate-400">
						School Management System Desktop • Built with precision &amp; care
					</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { X, Sparkles, Github, Globe, ExternalLink } from 'lucide-vue-next'
import DeveloperLogo from '@/components/DeveloperLogo.vue'
import { openExternalUrl } from '@/utils/browser'

defineProps<{
	isOpen: boolean
}>()

const emit = defineEmits<{
	(e: 'close'): void
}>()

function openUrl(url: string) {
	openExternalUrl(url)
}
</script>
