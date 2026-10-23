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
	IFeeService,
	IInvoiceService,
	IPaymentService,
	IBackupService,
	IUpdateService,
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
	async updateTaskbarIcon(base64Data: string) {
		return fetch('/api/v1/school/taskbar-icon', { method: 'POST', body: JSON.stringify({ icon: base64Data }) }).then(r => r.json())
	}
}

export class HttpAuditService implements IAuditService {
	async listLogs(filter: any) {
		const params = new URLSearchParams(filter)
		return fetch(`/api/v1/audit?${params.toString()}`).then(r => r.json())
	}
}

export class HttpAcademicService implements IAcademicService {
	async listAcademicYears() { return fetch('/api/v1/academics/years').then(r => r.json()) }
	async getActiveAcademicYear() { return fetch('/api/v1/academics/years/active').then(r => r.json()) }
	async createAcademicYear(input: any) { return fetch('/api/v1/academics/years', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async updateAcademicYear(input: any) { return fetch(`/api/v1/academics/years/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json()) }
	async setTerms(input: any) { return fetch('/api/v1/academics/terms', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listTerms(yearId: string) { return fetch(`/api/v1/academics/years/${yearId}/terms`).then(r => r.json()) }
	async listCampuses() { return fetch('/api/v1/academics/campuses').then(r => r.json()) }
	async createCampus(input: any) { return fetch('/api/v1/academics/campuses', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listClasses() { return fetch('/api/v1/academics/classes').then(r => r.json()) }
	async createClass(input: any) { return fetch('/api/v1/academics/classes', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async updateClass(input: any) { return fetch(`/api/v1/academics/classes/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json()) }
	async listSectionsByYear(yearId: string) { return fetch(`/api/v1/academics/years/${yearId}/sections`).then(r => r.json()) }
	async createSection(input: any) { return fetch('/api/v1/academics/sections', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async updateSection(input: any) { return fetch(`/api/v1/academics/sections/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json()) }
}

export class HttpStudentService implements IStudentService {
	async getStudent(id: string) { return fetch(`/api/v1/students/${id}`).then(r => r.json()) }
	async listStudents(filter: any) {
		const params = new URLSearchParams(filter)
		return fetch(`/api/v1/students?${params.toString()}`).then(r => r.json())
	}
	async registerStudent(input: any) { return fetch('/api/v1/students', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async updateStudent(input: any) { return fetch(`/api/v1/students/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json()) }
	async changeStatus(input: any) { return fetch(`/api/v1/students/${input.studentId}/status`, { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listStatusHistory(studentId: string) { return fetch(`/api/v1/students/${studentId}/history`).then(r => r.json()) }
}

export class HttpGuardianService implements IGuardianService {
	async getGuardian(id: string) { return fetch(`/api/v1/guardians/${id}`).then(r => r.json()) }
	async createGuardian(input: any) { return fetch('/api/v1/guardians', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async updateGuardian(input: any) { return fetch(`/api/v1/guardians/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json()) }
	async listForStudent(studentId: string) { return fetch(`/api/v1/students/${studentId}/guardians`).then(r => r.json()) }
	async setStudentGuardians(input: any) { return fetch(`/api/v1/students/${input.studentId}/guardians`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json()) }
	async findSiblings(studentId: string) { return fetch(`/api/v1/students/${studentId}/siblings`).then(r => r.json()) }
}

export class HttpEnrollmentService implements IEnrollmentService {
	async enrollStudent(input: any) { return fetch('/api/v1/enrollments', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async getActive(studentId: string) { return fetch(`/api/v1/students/${studentId}/enrollment`).then(r => r.json()) }
	async listHistory(studentId: string) { return fetch(`/api/v1/students/${studentId}/enrollments`).then(r => r.json()) }
	async listSectionRoster(sectionId: string) { return fetch(`/api/v1/sections/${sectionId}/roster`).then(r => r.json()) }
	async executePromotionBatch(input: any) { return fetch('/api/v1/promotions/execute', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async getPromotionBatch(id: string) { return fetch(`/api/v1/promotions/${id}`).then(r => r.json()) }
}

export class HttpAdmissionService implements IAdmissionService {
	async listFormDefinitions() { return fetch('/api/v1/admissions/forms').then(r => r.json()) }
	async ensureDefaultAdmissionForm() { return fetch('/api/v1/admissions/forms/ensure-default', { method: 'POST' }).then(r => r.json()) }
	async createFormDefinition(input: any) { return fetch('/api/v1/admissions/forms', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listFormVersions(formDefId: string) { return fetch(`/api/v1/admissions/forms/${formDefId}/versions`).then(r => r.json()) }
	async getPublishedFormVersion(formDefId: string) { return fetch(`/api/v1/admissions/forms/${formDefId}/published`).then(r => r.json()) }
	async createFormVersion(input: any) { return fetch('/api/v1/admissions/versions', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async publishFormVersion(input: any) { return fetch(`/api/v1/admissions/versions/${input.versionId}/publish`, { method: 'POST' }).then(r => r.json()) }
	async listApplications(status: string, yearId: string, limit: number) {
		const params = new URLSearchParams({ status, yearId, limit: String(limit) })
		return fetch(`/api/v1/admissions/applications?${params.toString()}`).then(r => r.json())
	}
	async getApplication(id: string) { return fetch(`/api/v1/admissions/applications/${id}`).then(r => r.json()) }
	async submitApplication(input: any) { return fetch('/api/v1/admissions/applications', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async updateStatus(input: any) { return fetch(`/api/v1/admissions/applications/${input.applicationId}/status`, { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async enrollApplication(input: any) { return fetch(`/api/v1/admissions/applications/${input.applicationId}/enroll`, { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
}

export class HttpDocumentService implements IDocumentService {
	async uploadDocument(input: any) { return fetch('/api/v1/documents', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listForStudent(studentId: string) { return fetch(`/api/v1/students/${studentId}/documents`).then(r => r.json()) }
	async listForApplication(applicationId: string) { return fetch(`/api/v1/admissions/applications/${applicationId}/documents`).then(r => r.json()) }
	async getContent(id: string) { return fetch(`/api/v1/documents/${id}/content`).then(r => r.json()) }
}

export class HttpFeeService implements IFeeService {
	async listFeeHeads() { return fetch('/api/v1/fees/heads').then(r => r.json()) }
	async createFeeHead(input: any) { return fetch('/api/v1/fees/heads', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async updateFeeHead(input: any) { return fetch(`/api/v1/fees/heads/${input.id}`, { method: 'PUT', body: JSON.stringify(input) }).then(r => r.json()) }

	async listFeeStructures(academicYearId: string) { return fetch(`/api/v1/fees/structures?academicYearId=${academicYearId}`).then(r => r.json()) }
	async createFeeStructure(input: any) { return fetch('/api/v1/fees/structures', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listFeeStructureVersions(feeStructureId: string) { return fetch(`/api/v1/fees/structures/${feeStructureId}/versions`).then(r => r.json()) }
	async getFeeStructureVersionDetail(versionId: string) { return fetch(`/api/v1/fees/versions/${versionId}`).then(r => r.json()) }
	async createFeeStructureVersion(input: any) { return fetch('/api/v1/fees/versions', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async publishFeeStructureVersion(versionId: string) { return fetch(`/api/v1/fees/versions/${versionId}/publish`, { method: 'POST' }).then(r => r.json()) }
	async deleteFeeStructureVersion(versionId: string) { return fetch(`/api/v1/fees/versions/${versionId}`, { method: 'DELETE' }).then(r => r.json()) }

	async listDiscountPolicies() { return fetch('/api/v1/fees/discounts/policies').then(r => r.json()) }
	async createDiscountPolicy(input: any) { return fetch('/api/v1/fees/discounts/policies', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listStudentDiscounts(studentId: string, academicYearId: string) { return fetch(`/api/v1/fees/discounts/students?studentId=${studentId}&academicYearId=${academicYearId}`).then(r => r.json()) }
	async listPendingDiscounts() { return fetch('/api/v1/fees/discounts/pending').then(r => r.json()) }
	async applyStudentDiscount(input: any) { return fetch('/api/v1/fees/discounts/apply', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async decideStudentDiscount(input: any) { return fetch('/api/v1/fees/discounts/decide', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
}

export class HttpInvoiceService implements IInvoiceService {
	async listInvoices(filter: any) {
		const params = new URLSearchParams()
		if (filter.status) params.set('status', filter.status)
		if (filter.academicYearId) params.set('academicYearId', filter.academicYearId)
		if (filter.studentId) params.set('studentId', filter.studentId)
		if (filter.limit) params.set('limit', String(filter.limit))
		return fetch(`/api/v1/invoices?${params.toString()}`).then(r => r.json())
	}
	async getInvoice(id: string) { return fetch(`/api/v1/invoices/${id}`).then(r => r.json()) }
	async generateInvoicesBatch(input: any) { return fetch('/api/v1/invoices/batch', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async voidInvoice(input: any) { return fetch(`/api/v1/invoices/${input.invoiceId}/void`, { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
}

export class HttpPaymentService implements IPaymentService {
	async recordPayment(input: any) { return fetch('/api/v1/payments', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async reversePayment(input: any) { return fetch(`/api/v1/payments/${input.paymentId}/reverse`, { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async getPayment(id: string) { return fetch(`/api/v1/payments/${id}`).then(r => r.json()) }
	async listPayments(filter: any) {
		const params = new URLSearchParams()
		if (filter.studentId) params.set('studentId', filter.studentId)
		if (filter.status) params.set('status', filter.status)
		if (filter.paymentMethod) params.set('paymentMethod', filter.paymentMethod)
		if (filter.limit) params.set('limit', String(filter.limit))
		return fetch(`/api/v1/payments?${params.toString()}`).then(r => r.json())
	}
	async getStudentCreditBalance(studentId: string) { return fetch(`/api/v1/students/${studentId}/credit-balance`).then(r => r.json()) }
	async requestFinancialAdjustment(input: any) { return fetch('/api/v1/adjustments', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async decideFinancialAdjustment(input: any) { return fetch(`/api/v1/adjustments/${input.adjustmentId}/decide`, { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async listFinancialAdjustments(filter: any) {
		const params = new URLSearchParams()
		if (filter.studentId) params.set('studentId', filter.studentId)
		if (filter.status) params.set('status', filter.status)
		if (filter.limit) params.set('limit', String(filter.limit))
		return fetch(`/api/v1/adjustments?${params.toString()}`).then(r => r.json())
	}
	async getFinancialAdjustment(id: string) { return fetch(`/api/v1/adjustments/${id}`).then(r => r.json()) }
}

class HttpBackupService implements IBackupService {
	async getDatabaseStatus() { return fetch('/api/v1/backup/status').then(r => r.json()) }
	async exportDecryptedDatabase(input: any) { return fetch('/api/v1/backup/export-decrypted', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async exportEncryptedDatabase(input: any) { return fetch('/api/v1/backup/export-encrypted', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
	async restoreDatabase(input: any) { return fetch('/api/v1/backup/restore', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
}

class HttpUpdateService implements IUpdateService {
	async getVersion() { return { success: true, data: '1.0.0' } }
	async checkForUpdates(customRepo?: string) { return fetch(`/api/v1/updates/check?repo=${customRepo || ''}`).then(r => r.json()) }
	async downloadAndApplyUpdate(input: any) { return fetch('/api/v1/updates/apply', { method: 'POST', body: JSON.stringify(input) }).then(r => r.json()) }
}

export function createHttpServices(): AppServices {
	return {
		auth: new HttpAuthService(),
		users: new HttpUserService(),
		roles: new HttpRoleService(),
		school: new HttpSchoolService(),
		audit: new HttpAuditService(),
		academics: new HttpAcademicService(),
		students: new HttpStudentService(),
		guardians: new HttpGuardianService(),
		enrollments: new HttpEnrollmentService(),
		admissions: new HttpAdmissionService(),
		documents: new HttpDocumentService(),
		fees: new HttpFeeService(),
		invoices: new HttpInvoiceService(),
		payments: new HttpPaymentService(),
		backup: new HttpBackupService(),
		updates: new HttpUpdateService()
	}
}
