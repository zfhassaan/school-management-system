<template>
	<aside
		aria-label="Developer Attribution"
		class="fixed bottom-3.5 right-4 z-30 flex items-center gap-2.5 py-1.5 px-3 bg-white/90 hover:bg-white backdrop-blur-md rounded-xl border border-slate-200/90 shadow-xs hover:shadow-md transition-all duration-200 select-none group"
	>
		<!-- Logo & Name (Clickable to open Developer Architecture Modal) -->
		<button
			type="button"
			@click="$emit('open-modal')"
			class="flex items-center gap-2 text-left cursor-pointer group/author focus:outline-none"
			title="Architect & Developer: Hassaan Ali (@zfhassaan)"
		>
			<DeveloperLogo sizeClass="w-5 h-5 group-hover/author:scale-105 transition-transform" />
			<div class="pr-2 border-r border-slate-200/80">
				<p class="text-[11px] font-bold text-slate-800 leading-tight group-hover/author:text-indigo-600 transition-colors">
					Hassaan Ali
				</p>
				<p class="text-[9px] font-mono font-medium text-indigo-600">
					@zfhassaan
				</p>
			</div>
		</button>

		<!-- External Profile Quick Links -->
		<div class="flex items-center gap-1">
			<button
				type="button"
				@click="openUrl('https://github.com/zfhassaan')"
				class="p-1 rounded-md text-slate-400 hover:text-slate-800 hover:bg-slate-100 transition cursor-pointer"
				title="GitHub: https://github.com/zfhassaan"
			>
				<Github class="w-3.5 h-3.5" />
			</button>
			<button
				type="button"
				@click="openUrl('https://zfhassaan.github.io')"
				class="p-1 rounded-md text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 transition cursor-pointer"
				title="Portfolio: https://zfhassaan.github.io"
			>
				<Globe class="w-3.5 h-3.5" />
			</button>
		</div>
	</aside>
</template>

<script setup lang="ts">
import DeveloperLogo from '@/components/DeveloperLogo.vue'
import { openExternalUrl } from '@/utils/browser'
import { Github, Globe } from 'lucide-vue-next'

defineEmits<{
	(e: 'open-modal'): void
}>()

function openUrl(url: string) {
	openExternalUrl(url)
}
</script>
