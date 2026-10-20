import type {
	IAuthService,
	IUserService,
	IRoleService,
	ISchoolService,
	IAuditService,
	IAcademicService,
	IStudentService,
	IGuardianService,
	IEnrollmentService,
	IAdmissionService,
	IDocumentService,
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
				AcademicService?: any
				StudentService?: any
				GuardianService?: any
				EnrollmentService?: any
				AdmissionService?: any
				DocumentService?: any
			}
		}
	}
}

function getBinding(serviceName: string) {
	const svc = (window.go?.wails as any)?.[serviceName]
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

export class WailsAcademicService implements IAcademicService {
	async listAcademicYears() {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListAcademicYears()
	}
	async getActiveAcademicYear() {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetActiveAcademicYear()
	}
	async createAcademicYear(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateAcademicYear(input)
	}
	async updateAcademicYear(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateAcademicYear(input)
	}
	async setTerms(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.SetTerms(input)
	}
	async listTerms(yearId: string) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListTerms(yearId)
	}
	async listCampuses() {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListCampuses()
	}
	async createCampus(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateCampus(input)
	}
	async listClasses() {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListClasses()
	}
	async createClass(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateClass(input)
	}
	async updateClass(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateClass(input)
	}
	async listSectionsByYear(yearId: string) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListSectionsByYear(yearId)
	}
	async createSection(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateSection(input)
	}
	async updateSection(input: any) {
		const svc = getBinding('AcademicService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateSection(input)
	}
}

export class WailsStudentService implements IStudentService {
	async getStudent(id: string) {
		const svc = getBinding('StudentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetStudent(id)
	}
	async listStudents(filter: any) {
		const svc = getBinding('StudentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListStudents(filter)
	}
	async registerStudent(input: any) {
		const svc = getBinding('StudentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.RegisterStudent(input)
	}
	async updateStudent(input: any) {
		const svc = getBinding('StudentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateStudent(input)
	}
	async changeStatus(input: any) {
		const svc = getBinding('StudentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ChangeStatus(input)
	}
	async listStatusHistory(studentId: string) {
		const svc = getBinding('StudentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListStatusHistory(studentId)
	}
}

export class WailsGuardianService implements IGuardianService {
	async getGuardian(id: string) {
		const svc = getBinding('GuardianService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetGuardian(id)
	}
	async createGuardian(input: any) {
		const svc = getBinding('GuardianService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateGuardian(input)
	}
	async updateGuardian(input: any) {
		const svc = getBinding('GuardianService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateGuardian(input)
	}
	async listForStudent(studentId: string) {
		const svc = getBinding('GuardianService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListForStudent(studentId)
	}
	async setStudentGuardians(input: any) {
		const svc = getBinding('GuardianService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.SetStudentGuardians(input)
	}
	async findSiblings(studentId: string) {
		const svc = getBinding('GuardianService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.FindSiblings(studentId)
	}
}

export class WailsEnrollmentService implements IEnrollmentService {
	async enrollStudent(input: any) {
		const svc = getBinding('EnrollmentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.EnrollStudent(input)
	}
	async getActive(studentId: string) {
		const svc = getBinding('EnrollmentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetActive(studentId)
	}
	async listHistory(studentId: string) {
		const svc = getBinding('EnrollmentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListHistory(studentId)
	}
	async listSectionRoster(sectionId: string) {
		const svc = getBinding('EnrollmentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListSectionRoster(sectionId)
	}
	async executePromotionBatch(input: any) {
		const svc = getBinding('EnrollmentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ExecutePromotionBatch(input)
	}
	async getPromotionBatch(id: string) {
		const svc = getBinding('EnrollmentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetPromotionBatch(id)
	}
}

export class WailsAdmissionService implements IAdmissionService {
	async listFormDefinitions() {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListFormDefinitions()
	}
	async createFormDefinition(input: any) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateFormDefinition(input)
	}
	async listFormVersions(formDefId: string) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListFormVersions(formDefId)
	}
	async getPublishedFormVersion(formDefId: string) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetPublishedFormVersion(formDefId)
	}
	async createFormVersion(input: any) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.CreateFormVersion(input)
	}
	async publishFormVersion(input: any) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.PublishFormVersion(input)
	}
	async listApplications(status: string, yearId: string, limit: number) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListApplications(status, yearId, limit)
	}
	async getApplication(id: string) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetApplication(id)
	}
	async submitApplication(input: any) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.SubmitApplication(input)
	}
	async updateStatus(input: any) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UpdateStatus(input)
	}
	async enrollApplication(input: any) {
		const svc = getBinding('AdmissionService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.EnrollApplication(input)
	}
}

export class WailsDocumentService implements IDocumentService {
	async uploadDocument(input: any) {
		const svc = getBinding('DocumentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.UploadDocument(input)
	}
	async listForStudent(studentId: string) {
		const svc = getBinding('DocumentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListForStudent(studentId)
	}
	async listForApplication(applicationId: string) {
		const svc = getBinding('DocumentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.ListForApplication(applicationId)
	}
	async getContent(id: string) {
		const svc = getBinding('DocumentService')
		if (!svc) return { success: false, error: { code: 'unavailable', message: 'Wails runtime unavailable' } }
		return await svc.GetContent(id)
	}
}

export function createWailsServices(): AppServices {
	return {
		auth: new WailsAuthService(),
		users: new WailsUserService(),
		roles: new WailsRoleService(),
		school: new WailsSchoolService(),
		audit: new WailsAuditService(),
		academics: new WailsAcademicService(),
		students: new WailsStudentService(),
		guardians: new WailsGuardianService(),
		enrollments: new WailsEnrollmentService(),
		admissions: new WailsAdmissionService(),
		documents: new WailsDocumentService()
	}
}

