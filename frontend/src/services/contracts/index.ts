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
	AuditLogDTO,
	AcademicYearDTO,
	AcademicTermDTO,
	CampusDTO,
	ClassDTO,
	SectionDTO,
	StudentDTO,
	StudentFilterDTO,
	StatusHistoryDTO,
	GuardianDTO,
	StudentGuardianLinkDTO,
	EnrollmentDTO,
	PromotionBatchDTO,
	FormDefinitionDTO,
	FormVersionDTO,
	ApplicationDTO,
	DocumentDTO
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

export interface IAcademicService {
	listAcademicYears(): Promise<ResponseEnvelope<AcademicYearDTO[]>>
	getActiveAcademicYear(): Promise<ResponseEnvelope<AcademicYearDTO>>
	createAcademicYear(input: any): Promise<ResponseEnvelope<string>>
	updateAcademicYear(input: any): Promise<ResponseEnvelope<boolean>>
	setTerms(input: any): Promise<ResponseEnvelope<boolean>>
	listTerms(yearId: string): Promise<ResponseEnvelope<AcademicTermDTO[]>>
	listCampuses(): Promise<ResponseEnvelope<CampusDTO[]>>
	createCampus(input: any): Promise<ResponseEnvelope<string>>
	listClasses(): Promise<ResponseEnvelope<ClassDTO[]>>
	createClass(input: any): Promise<ResponseEnvelope<string>>
	updateClass(input: any): Promise<ResponseEnvelope<boolean>>
	listSectionsByYear(yearId: string): Promise<ResponseEnvelope<SectionDTO[]>>
	createSection(input: any): Promise<ResponseEnvelope<string>>
	updateSection(input: any): Promise<ResponseEnvelope<boolean>>
}

export interface IStudentService {
	getStudent(id: string): Promise<ResponseEnvelope<StudentDTO>>
	listStudents(filter: StudentFilterDTO): Promise<ResponseEnvelope<StudentDTO[]>>
	registerStudent(input: any): Promise<ResponseEnvelope<string>>
	updateStudent(input: any): Promise<ResponseEnvelope<boolean>>
	changeStatus(input: any): Promise<ResponseEnvelope<boolean>>
	listStatusHistory(studentId: string): Promise<ResponseEnvelope<StatusHistoryDTO[]>>
}

export interface IGuardianService {
	getGuardian(id: string): Promise<ResponseEnvelope<GuardianDTO>>
	createGuardian(input: any): Promise<ResponseEnvelope<string>>
	updateGuardian(input: any): Promise<ResponseEnvelope<boolean>>
	listForStudent(studentId: string): Promise<ResponseEnvelope<StudentGuardianLinkDTO[]>>
	setStudentGuardians(input: any): Promise<ResponseEnvelope<boolean>>
	findSiblings(studentId: string): Promise<ResponseEnvelope<StudentDTO[]>>
}

export interface IEnrollmentService {
	enrollStudent(input: any): Promise<ResponseEnvelope<string>>
	getActive(studentId: string): Promise<ResponseEnvelope<EnrollmentDTO>>
	listHistory(studentId: string): Promise<ResponseEnvelope<EnrollmentDTO[]>>
	listSectionRoster(sectionId: string): Promise<ResponseEnvelope<EnrollmentDTO[]>>
	executePromotionBatch(input: any): Promise<ResponseEnvelope<string>>
	getPromotionBatch(id: string): Promise<ResponseEnvelope<PromotionBatchDTO>>
}

export interface IAdmissionService {
	listFormDefinitions(): Promise<ResponseEnvelope<FormDefinitionDTO[]>>
	createFormDefinition(input: any): Promise<ResponseEnvelope<string>>
	listFormVersions(formDefId: string): Promise<ResponseEnvelope<FormVersionDTO[]>>
	getPublishedFormVersion(formDefId: string): Promise<ResponseEnvelope<FormVersionDTO>>
	createFormVersion(input: any): Promise<ResponseEnvelope<string>>
	publishFormVersion(input: any): Promise<ResponseEnvelope<boolean>>
	listApplications(status: string, yearId: string, limit: number): Promise<ResponseEnvelope<ApplicationDTO[]>>
	getApplication(id: string): Promise<ResponseEnvelope<ApplicationDTO>>
	submitApplication(input: any): Promise<ResponseEnvelope<string>>
	updateStatus(input: any): Promise<ResponseEnvelope<boolean>>
	enrollApplication(input: any): Promise<ResponseEnvelope<string>>
}

export interface IDocumentService {
	uploadDocument(input: any): Promise<ResponseEnvelope<string>>
	listForStudent(studentId: string): Promise<ResponseEnvelope<DocumentDTO[]>>
	listForApplication(applicationId: string): Promise<ResponseEnvelope<DocumentDTO[]>>
	getContent(id: string): Promise<ResponseEnvelope<{ document: DocumentDTO; dataBase64: string }>>
}

export interface AppServices {
	auth: IAuthService
	users: IUserService
	roles: IRoleService
	school: ISchoolService
	audit: IAuditService
	academics: IAcademicService
	students: IStudentService
	guardians: IGuardianService
	enrollments: IEnrollmentService
	admissions: IAdmissionService
	documents: IDocumentService
}

