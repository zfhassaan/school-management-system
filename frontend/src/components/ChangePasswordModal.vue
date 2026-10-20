<template>
	<div
		v-if="isOpen"
		class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4"
	>
		<div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 border border-slate-100">
			<div class="flex items-center gap-3 text-indigo-600 mb-4">
				<div class="w-10 h-10 rounded-xl bg-indigo-50 flex items-center justify-center">
					<KeyRound class="w-5 h-5" />
				</div>
				<div>
					<h3 class="text-lg font-bold text-slate-900">Change Password Required</h3>
					<p class="text-xs text-slate-500">You must set a new personal password to continue.</p>
				</div>
			</div>

			<form @submit.prevent="handleSubmit" class="space-y-4">
				<div>
					<label class="block text-xs font-semibold uppercase tracking-wider text-slate-600 mb-1">
						Current Password
					</label>
					<input
						v-model="oldPassword"
						type="password"
						required
						class="w-full px-3 py-2 border rounded-lg text-sm outline-hidden focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 border-slate-300"
						placeholder="••••••••••••"
					/>
				</div>

				<div>
					<label class="block text-xs font-semibold uppercase tracking-wider text-slate-600 mb-1">
						New Password
					</label>
					<input
						v-model="newPassword"
						type="password"
						required
						minlength="10"
						class="w-full px-3 py-2 border rounded-lg text-sm outline-hidden focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 border-slate-300"
						placeholder="Min 10 chars with letters & numbers"
					/>
				</div>

				<div>
					<label class="block text-xs font-semibold uppercase tracking-wider text-slate-600 mb-1">
						Confirm New Password
					</label>
					<input
						v-model="confirmPassword"
						type="password"
						required
						class="w-full px-3 py-2 border rounded-lg text-sm outline-hidden focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 border-slate-300"
						placeholder="••••••••••••"
					/>
				</div>

				<div v-if="errorMsg" class="p-3 bg-rose-50 border border-rose-200 rounded-lg text-rose-700 text-xs">
					{{ errorMsg }}
				</div>

				<div class="pt-2">
					<button
						type="submit"
						:disabled="isSubmitting"
						class="w-full py-2.5 px-4 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white font-medium rounded-lg text-sm transition shadow-sm"
					>
						{{ isSubmitting ? 'Updating Password...' : 'Save New Password' }}
					</button>
				</div>
			</form>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { getServices } from '@/services/factory'
import { KeyRound } from 'lucide-vue-next'

const authStore = useAuthStore()
const toast = useToast()
const services = getServices()

const isOpen = computed(() => authStore.isAuthenticated && authStore.mustChangePassword)

const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const isSubmitting = ref(false)
const errorMsg = ref('')

async function handleSubmit() {
	errorMsg.value = ''
	if (newPassword.value !== confirmPassword.value) {
		errorMsg.value = 'New passwords do not match'
		return
	}
	if (newPassword.value.length < 10) {
		errorMsg.value = 'Password must be at least 10 characters'
		return
	}

	isSubmitting.value = true
	try {
		const res = await services.auth.changePassword({
			oldPassword: oldPassword.value,
			newPassword: newPassword.value
		})
		if (res.success) {
			toast.success('Your password has been updated successfully.')
			if (authStore.user) {
				authStore.user.mustChangePassword = false
			}
			oldPassword.value = ''
			newPassword.value = ''
			confirmPassword.value = ''
		} else {
			errorMsg.value = res.error?.message || 'Failed to change password'
		}
	} catch (e: any) {
		errorMsg.value = e.message || 'An error occurred'
	} finally {
		isSubmitting.value = false
	}
}
</script>

