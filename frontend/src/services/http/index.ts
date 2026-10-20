import type {
	IAuthService,
	IUserService,
	IRoleService,
	ISchoolService,
	IAuditService,
	AppServices
} from '@/services/contracts'

export class HttpAuthService implements IAuthService {
	async checkSetupStatus() {
		return fetch('/api/v1/auth/setup-status').then(r => r.json())
	}
	async firstRunSetup(input: any) {
		return fetch('/api/v1/auth/first-run-setup', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json())
	}
	async login(credentials: any) {
		return fetch('/api/v1/auth/login', { method: 'POST', body: JSON.stringify(credentials) }).then(r => r.json())
	}
	async logout() {
		return fetch('/api/v1/auth/logout', { method: 'POST' }).then(r => r.json())
	}
	async changePassword(input: any) {
		return fetch('/api/v1/auth/change-password', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json())
	}
	async getCurrentUser() {
		return fetch('/api/v1/auth/me').then(r => r.json())
	}
}

export class HttpUserService implements IUserService {
	async listUsers() {
		return fetch('/api/v1/users').then(r => r.json())
	}
	async getUser(id: string) {
		return fetch(`/api/v1/users/${id}`).then(r => r.json())
	}
	async createUser(input: any) {
		return fetch('/api/v1/users', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json())
	}
	async updateUser(input: any) {
		return fetch(`/api/v1/users/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json())
	}
	async setUserStatus(id: string, active: boolean) {
		return fetch(`/api/v1/users/${id}/status`, { method: 'PATCH', body: JSON.stringify({ active }) }).then(r => r.json())
	}
	async adminResetPassword(input: any) {
		return fetch(`/api/v1/users/${input.userId}/reset-password`, { method: 'POST', body: JSON.stringify(input) }).then(r => r.json())
	}
}

export class HttpRoleService implements IRoleService {
	async listRoles() {
		return fetch('/api/v1/roles').then(r => r.json())
	}
	async getRole(id: string) {
		return fetch(`/api/v1/roles/${id}`).then(r => r.json())
	}
	async listPermissions() {
		return fetch('/api/v1/roles/permissions').then(r => r.json())
	}
	async createRole(input: any) {
		return fetch('/api/v1/roles', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json())
	}
	async updateRole(input: any) {
		return fetch(`/api/v1/roles/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json())
	}
	async deleteRole(id: string) {
		return fetch(`/api/v1/roles/${id}`, { method: 'DELETE' }).then(r => r.json())
	}
}

export class HttpSchoolService implements ISchoolService {
	async getSchool() {
		return fetch('/api/v1/school').then(r => r.json())
	}
	async updateSchoolProfile(input: any) {
		return fetch('/api/v1/school', { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json())
	}
	async getSettings() {
		return fetch('/api/v1/school/settings').then(r => r.json())
	}
	async updateSettings(input: any) {
		return fetch('/api/v1/school/settings', { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json())
	}
}

export class HttpAuditService implements IAuditService {
	async listLogs(filter: any) {
		const params = new URLSearchParams(filter)
		return fetch(`/api/v1/audit?${params.toString()}`).then(r => r.json())
	}
}

export function createHttpServices(): AppServices {
	return {
		auth: new HttpAuthService(),
		users: new HttpUserService(),
		roles: new HttpRoleService(),
		school: new HttpSchoolService(),
		audit: new HttpAuditService()
	}
}

