import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getServices } from '@/services/factory'
import type { UserDTO, SchoolDTO, FirstRunSetupInput } from '@/types/common'

export const useAuthStore = defineStore('auth', () => {
	const services = getServices()

	const isInitialized = ref<boolean | null>(null)
	const isAuthenticated = ref<boolean>(false)
	const user = ref<UserDTO | null>(null)
	const school = ref<SchoolDTO | null>(null)
	const token = ref<string | null>(null)
	const isLoading = ref<boolean>(false)

	const mustChangePassword = computed(() => !!user.value?.mustChangePassword)

	const permissionsSet = computed(() => {
		return new Set(user.value?.permissions || [])
	})

	function hasPermission(permission: string): boolean {
		if (!user.value) return false
		if (user.value.roles.includes('super_admin')) return true
		return permissionsSet.value.has(permission)
	}

	async function checkSetup(): Promise<boolean> {
		try {
			const res = await services.auth.checkSetupStatus()
			if (res.success && res.data) {
				isInitialized.value = res.data.isInitialized
				return res.data.isInitialized
			}
			isInitialized.value = false
			return false
		} catch (e) {
			console.error('Failed to check setup status:', e)
			isInitialized.value = false
			return false
		}
	}

	async function login(username: string, password: string) {
		isLoading.value = true
		try {
			const res = await services.auth.login({ username, password })
			if (res.success && res.data) {
				token.value = res.data.sessionToken
				user.value = res.data.user
				school.value = res.data.school
				isAuthenticated.value = true
				return { success: true }
			}
			return {
				success: false,
				error: res.error?.message || 'Login failed'
			}
		} finally {
			isLoading.value = false
		}
	}

	async function firstRunSetup(input: FirstRunSetupInput) {
		isLoading.value = true
		try {
			const res = await services.auth.firstRunSetup(input)
			if (res.success && res.data) {
				isInitialized.value = true
				return {
					success: true,
					data: res.data
				}
			}
			return {
				success: false,
				error: res.error?.message || 'Setup failed'
			}
		} finally {
			isLoading.value = false
		}
	}

	async function logout() {
		try {
			await services.auth.logout()
		} finally {
			token.value = null
			user.value = null
			isAuthenticated.value = false
		}
	}

	return {
		isInitialized,
		isAuthenticated,
		user,
		school,
		token,
		isLoading,
		mustChangePassword,
		hasPermission,
		checkSetup,
		login,
		firstRunSetup,
		logout
	}
})

