import type {
	ResponseEnvelope,
	LoginResult,
	SetupStatusDTO,
	FirstRunSetupInput,
	FirstRunSetupResult,
	UserListItemDTO,
	RoleDTO,
	PermissionInfoDTO,
	SchoolProfileDTO,
	SchoolSettingsDTO,
	AuditLogDTO
} from '@/types/common'

export interface IAuthService {
	checkSetupStatus(): Promise<ResponseEnvelope<SetupStatusDTO>>
	firstRunSetup(input: FirstRunSetupInput): Promise<ResponseEnvelope<FirstRunSetupResult>>
	login(credentials: { username: string; password: string }): Promise<ResponseEnvelope<LoginResult>>
	logout(): Promise<ResponseEnvelope<boolean>>
	changePassword(input: { oldPassword: string; newPassword: string }): Promise<ResponseEnvelope<boolean>>
	getCurrentUser(): Promise<ResponseEnvelope<LoginResult>>
}

export interface IUserService {
	listUsers(): Promise<ResponseEnvelope<UserListItemDTO[]>>
	getUser(id: string): Promise<ResponseEnvelope<UserListItemDTO>>
	createUser(input: {
		username: string
		displayName: string
		email: string
		password: string
		roleIds: string[]
	}): Promise<ResponseEnvelope<string>>
	updateUser(input: {
		id: string
		displayName: string
		email: string
		roleIds: string[]
	}): Promise<ResponseEnvelope<boolean>>
	setUserStatus(id: string, active: boolean): Promise<ResponseEnvelope<boolean>>
	adminResetPassword(input: { userId: string; newPassword: string }): Promise<ResponseEnvelope<boolean>>
}

export interface IRoleService {
	listRoles(): Promise<ResponseEnvelope<RoleDTO[]>>
	getRole(id: string): Promise<ResponseEnvelope<RoleDTO>>
	listPermissions(): Promise<ResponseEnvelope<PermissionInfoDTO[]>>
	createRole(input: {
		code: string
		name: string
		description: string
		permissions: string[]
	}): Promise<ResponseEnvelope<string>>
	updateRole(input: {
		id: string
		name: string
		description: string
		permissions: string[]
	}): Promise<ResponseEnvelope<boolean>>
	deleteRole(id: string): Promise<ResponseEnvelope<boolean>>
}

export interface ISchoolService {
	getSchool(): Promise<ResponseEnvelope<SchoolProfileDTO>>
	updateSchoolProfile(input: any): Promise<ResponseEnvelope<boolean>>
	getSettings(): Promise<ResponseEnvelope<SchoolSettingsDTO>>
	updateSettings(input: any): Promise<ResponseEnvelope<boolean>>
}

export interface IAuditService {
	listLogs(filter: {
		entityType?: string
		entityId?: string
		limit?: number
	}): Promise<ResponseEnvelope<AuditLogDTO[]>>
}

export interface AppServices {
	auth: IAuthService
	users: IUserService
	roles: IRoleService
	school: ISchoolService
	audit: IAuditService
}

