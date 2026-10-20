<template>
	<div class="min-h-screen w-screen flex items-center justify-center bg-slate-100 p-4">
		<div class="max-w-md w-full bg-white rounded-2xl shadow-xl border border-slate-200/80 p-8">
			<!-- Header -->
			<div class="text-center mb-8">
				<div class="w-12 h-12 bg-indigo-600 text-white rounded-xl mx-auto flex items-center justify-center font-bold text-lg mb-3 shadow-md shadow-indigo-200">
					{{ schoolInitials }}
				</div>
				<h1 class="text-xl font-bold text-slate-900">{{ authStore.school?.name || 'Sign in to your School' }}</h1>
				<p class="text-xs text-slate-500 mt-1">{{ authStore.school?.code ? `${authStore.school.code} • Management Portal` : 'Single-School Management System (Offline-Ready)' }}</p>
			</div>

			<!-- Error Alert -->
			<div v-if="errorMessage" class="mb-5 p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs flex items-start gap-2">
				<AlertCircle class="w-4 h-4 shrink-0 mt-0.5" />
				<div>{{ errorMessage }}</div>
			</div>

			<!-- Form -->
			<form @submit.prevent="handleLogin" class="space-y-4">
				<div>
					<label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">Username</label>
					<div class="relative">
						<input
							v-model="username"
							type="text"
							required
							autofocus
							class="w-full pl-9 pr-3 py-2.5 bg-slate-50 border border-slate-300 rounded-lg text-sm text-slate-900 focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-hidden transition"
							placeholder="admin"
						/>
						<User class="w-4 h-4 text-slate-400 absolute left-3 top-3" />
					</div>
				</div>

				<div>
					<label class="block text-xs font-semibold text-slate-700 uppercase tracking-wider mb-1">Password</label>
					<div class="relative">
						<input
							v-model="password"
							type="password"
							required
							class="w-full pl-9 pr-3 py-2.5 bg-slate-50 border border-slate-300 rounded-lg text-sm text-slate-900 focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-hidden transition"
							placeholder="••••••••••••"
						/>
						<Lock class="w-4 h-4 text-slate-400 absolute left-3 top-3" />
					</div>
				</div>

				<div class="pt-2">
					<button
						type="submit"
						:disabled="isLoading"
						class="w-full py-2.5 px-4 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-medium rounded-lg text-sm transition shadow-sm flex items-center justify-center gap-2 cursor-pointer"
					>
						<Loader2 v-if="isLoading" class="w-4 h-4 animate-spin" />
						<span>{{ isLoading ? 'Authenticating...' : 'Sign In' }}</span>
					</button>
				</div>
			</form>

			<div class="mt-6 pt-4 border-t border-slate-100 flex items-center justify-between text-[11px] text-slate-400">
				<span>Offline-First Monolith</span>
				<span>SQLite WAL Engine</span>
			</div>
		</div>

		<!-- Floating Bottom-Right Developer Badge -->
		<DeveloperBadge @open-modal="showDeveloperModal = true" />

		<!-- Developer About Modal -->
		<DeveloperAboutModal :isOpen="showDeveloperModal" @close="showDeveloperModal = false" />
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import DeveloperBadge from '@/components/DeveloperBadge.vue'
import DeveloperAboutModal from '@/components/DeveloperAboutModal.vue'
import { User, Lock, AlertCircle, Loader2 } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const showDeveloperModal = ref(false)

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const isLoading = ref(false)

const schoolInitials = computed(() => {
	const code = authStore.school?.code
	if (code && code !== 'OFFLINE') {
		return code.length <= 3 ? code.toUpperCase() : code.substring(0, 2).toUpperCase()
	}
	const name = authStore.school?.name || 'SMS'
	const words = name.trim().split(/\s+/)
	if (words.length >= 2) {
		return (words[0][0] + words[1][0]).toUpperCase()
	}
	return name.substring(0, 2).toUpperCase()
})

onMounted(async () => {
	await authStore.checkSetup()
})

async function handleLogin() {
	errorMessage.value = ''
	isLoading.value = true
	try {
		const res = await authStore.login(username.value, password.value)
		if (res.success) {
			router.push('/')
		} else {
			errorMessage.value = res.error || 'Invalid credentials'
		}
	} catch (e: any) {
		errorMessage.value = e.message || 'Connection error'
	} finally {
		isLoading.value = false
	}
}
</script>

