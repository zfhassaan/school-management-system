export interface ResponseEnvelope<T> {
	success: boolean
	data?: T
	error?: {
		code: string
		message: string
		fields?: Record<string, string>
		correlationId?: string
	}
}

export interface UserDTO {
	id: string
	username: string
	displayName: string
	email: string
	roles: string[]
	permissions: string[]
	mustChangePassword: boolean
}

export interface SchoolDTO {
	id: string
	name: string
	code: string
	timezone: string
	currencyCode: string
	currencyMinorUnits: number
	dateFormat: string
}

export interface LoginResult {
	sessionToken: string
	user: UserDTO
	school: SchoolDTO
}

export interface SetupStatusDTO {
	isInitialized: boolean
	schoolName?: string
}

export interface FirstRunSetupInput {
	schoolName: string
	schoolCode: string
	timezone: string
	currencyCode: string
	currencyMinorUnits: number
	adminUsername: string
	adminDisplayName: string
	adminEmail: string
	adminPassword: string
}

export interface FirstRunSetupResult {
	schoolId: string
	userId: string
	recoveryCodes: string[]
}

export interface UserListItemDTO {
	id: string
	username: string
	displayName: string
	email: string
	status: string
	roles: string[]
	roleIds: string[]
	mustChangePassword: boolean
	lastLoginAt?: string
	createdAt: string
}

export interface RoleDTO {
	id: string
	code: string
	name: string
	description: string
	isSystem: boolean
	permissions: string[]
	userCount: number
}

export interface PermissionInfoDTO {
	code: string
	module: string
	description: string
	sensitive: boolean
}

export interface SchoolSettingsDTO {
	schoolId: string
	admissionNoFormat: string
	invoicePrefix: string
	receiptPrefix: string
	roundingUnitMinor: number
	roundingMode: string
	defaultDueDay: number
	lateFeeGraceDays: number
	discountCombination: string
	requireSecondApprover: boolean
	dateFormat: string
}

export interface SchoolProfileDTO {
	id: string
	name: string
	legalName?: string
	code: string
	address: {
		line1?: string
		line2?: string
		city?: string
		region?: string
		postalCode?: string
		country?: string
	}
	phone?: string
	email?: string
	website?: string
	timezone: string
	currencyCode: string
	currencyMinorUnits: number
	locale: string
}

export interface AuditLogDTO {
	id: string
	occurredAt: string
	actorUserId?: string
	deviceId: string
	action: string
	entityType: string
	entityId: string
	correlationId: string
	reason: string
	before?: any
	after?: any
	source: string
}

