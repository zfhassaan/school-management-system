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
	const schoolLogo = ref<string | null>(typeof window !== 'undefined' ? localStorage.getItem('sms_school_logo') : null)
	const token = ref<string | null>(null)
	const isLoading = ref<boolean>(false)

	function setSchoolLogo(logoData: string | null) {
		schoolLogo.value = logoData
		if (typeof window !== 'undefined') {
			if (logoData) {
				localStorage.setItem('sms_school_logo', logoData)
			} else {
				localStorage.removeItem('sms_school_logo')
			}
		}
	}

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
				const d = res.data as any
				const initialized = !!(d.isInitialized ?? d.IsInitialized)
				isInitialized.value = initialized
				if ((d.schoolName || d.SchoolName) && !school.value?.name) {
					updateSchoolDetails({
						name: d.schoolName || d.SchoolName
					})
				}
				return initialized
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

	function updateSchoolDetails(updated: Partial<SchoolDTO>) {
		if (school.value) {
			school.value = { ...school.value, ...updated }
		} else {
			school.value = updated as SchoolDTO
		}
	}

	async function refreshSchool() {
		try {
			const res = await services.school.getSchool()
			if (res.success && res.data) {
				const d = res.data as any
				updateSchoolDetails({
					id: d.id || d.ID || (school.value?.id ?? ''),
					name: d.name || d.Name || school.value?.name || '',
					legalName: d.legalName || d.LegalName || school.value?.legalName || '',
					code: d.code || d.Code || school.value?.code || '',
					timezone: d.timezone || d.Timezone || school.value?.timezone || '',
					currencyCode: d.currencyCode || d.CurrencyCode || school.value?.currencyCode || 'PKR',
					currencyMinorUnits: d.currencyMinorUnits ?? d.CurrencyMinorUnits ?? school.value?.currencyMinorUnits ?? 2,
					phone: d.phone || d.Phone || school.value?.phone || '',
					email: d.email || d.Email || school.value?.email || '',
					website: d.website || d.Website || school.value?.website || '',
					address: d.address || d.Address || school.value?.address
				})
			}
		} catch (e) {
			console.error('Failed to refresh school details:', e)
		}
	}

	return {
		isInitialized,
		isAuthenticated,
		user,
		school,
		schoolLogo,
		setSchoolLogo,
		token,
		isLoading,
		mustChangePassword,
		hasPermission,
		checkSetup,
		login,
		firstRunSetup,
		logout,
		updateSchoolDetails,
		refreshSchool
	}
})

