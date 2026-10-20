import type {
	IAuthService,
	IUserService,
	IRoleService,
	ISchoolService,
	IAuditService,
	AppServices
} from '@/services/contracts'

// Declare Wails runtime bindings on window
declare global {
	interface Window {
		go?: {
			wails?: {
				AuthService?: any
				UserService?: any
				RoleService?: any
				SchoolService?: any
				AuditService?: any
			}
		}
	}
}

function getBinding(serviceName: 'AuthService' | 'UserService' | 'RoleService' | 'SchoolService' | 'AuditService') {
	const svc = window.go?.wails?.[serviceName]
	if (!svc) {
		console.warn(`[Wails] Service ${serviceName} is not available on window.go.wails`)
		return null
	}
	return svc
}

export class WailsAuthService implements IAuthService {
	async checkSetupStatus() {
		const svc = getBinding('AuthService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CheckSetupStatus()
	}
	async firstRunSetup(input: any) {
		const svc = getBinding('AuthService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.FirstRunSetup(input)
	}
	async login(credentials: any) {
		const svc = getBinding('AuthService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.Login(credentials)
	}
	async logout() {
		const svc = getBinding('AuthService')
		if (!svc) return { success: true }
		return await svc.Logout()
	}
	async changePassword(input: any) {
		const svc = getBinding('AuthService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ChangePassword(input)
	}
	async getCurrentUser() {
		const svc = getBinding('AuthService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetCurrentUser()
	}
}

export class WailsUserService implements IUserService {
	async listUsers() {
		const svc = getBinding('UserService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListUsers()
	}
	async getUser(id: string) {
		const svc = getBinding('UserService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetUser(id)
	}
	async createUser(input: any) {
		const svc = getBinding('UserService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateUser(input)
	}
	async updateUser(input: any) {
		const svc = getBinding('UserService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateUser(input)
	}
	async setUserStatus(id: string, active: boolean) {
		const svc = getBinding('UserService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.SetUserStatus(id, active)
	}
	async adminResetPassword(input: any) {
		const svc = getBinding('UserService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.AdminResetPassword(input)
	}
}

export class WailsRoleService implements IRoleService {
	async listRoles() {
		const svc = getBinding('RoleService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListRoles()
	}
	async getRole(id: string) {
		const svc = getBinding('RoleService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetRole(id)
	}
	async listPermissions() {
		const svc = getBinding('RoleService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListPermissions()
	}
	async createRole(input: any) {
		const svc = getBinding('RoleService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateRole(input)
	}
	async updateRole(input: any) {
		const svc = getBinding('RoleService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateRole(input)
	}
	async deleteRole(id: string) {
		const svc = getBinding('RoleService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.DeleteRole(id)
	}
}

export class WailsSchoolService implements ISchoolService {
	async getSchool() {
		const svc = getBinding('SchoolService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetSchool()
	}
	async updateSchoolProfile(input: any) {
		const svc = getBinding('SchoolService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateSchoolProfile(input)
	}
	async getSettings() {
		const svc = getBinding('SchoolService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetSettings()
	}
	async updateSettings(input: any) {
		const svc = getBinding('SchoolService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateSettings(input)
	}
}

export class WailsAuditService implements IAuditService {
	async listLogs(filter: any) {
		const svc = getBinding('AuditService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListLogs(filter)
	}
}

export function createWailsServices(): AppServices {
	return {
		auth: new WailsAuthService(),
		users: new WailsUserService(),
		roles: new WailsRoleService(),
		school: new WailsSchoolService(),
		audit: new WailsAuditService()
	}
}

