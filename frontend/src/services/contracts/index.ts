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
	DocumentDTO,
	FeeHeadDTO,
	FeeStructureDTO,
	FeeStructureVersionDTO,
	FeeStructureVersionDetailDTO,
	DiscountPolicyDTO,
	StudentDiscountDTO,
	InvoiceDTO,
	InvoiceDetailDTO,
	InvoiceFilterDTO,
	PaymentDTO,
	PaymentDetailDTO,
	StudentCreditBalanceInfoDTO,
	FinancialAdjustmentDTO,
	RecordPaymentResultDTO
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
	updateTaskbarIcon(base64Data: string): Promise<ResponseEnvelope<boolean>>
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

export interface IFeeService {
	listFeeHeads(): Promise<ResponseEnvelope<FeeHeadDTO[]>>
	createFeeHead(input: {
		code: string
		name: string
		description: string
		isRecurring: boolean
		defaultAmountMinor: number
	}): Promise<ResponseEnvelope<string>>
	updateFeeHead(input: {
		id: string
		name: string
		description: string
		isRecurring: boolean
		defaultAmountMinor: number
		isActive: boolean
	}): Promise<ResponseEnvelope<boolean>>

	listFeeStructures(academicYearId: string): Promise<ResponseEnvelope<FeeStructureDTO[]>>
	createFeeStructure(input: {
		academicYearId: string
		classId: string
		campusId?: string
		code: string
		name: string
		description: string
	}): Promise<ResponseEnvelope<string>>
	listFeeStructureVersions(feeStructureId: string): Promise<ResponseEnvelope<FeeStructureVersionDTO[]>>
	getFeeStructureVersionDetail(versionId: string): Promise<ResponseEnvelope<FeeStructureVersionDetailDTO>>
	createFeeStructureVersion(input: {
		feeStructureId: string
		changelog: string
		items: {
			feeHeadId: string
			amountMinor: number
			frequency: string
			dueDay: number
			isOptional: boolean
		}[]
	}): Promise<ResponseEnvelope<string>>
	publishFeeStructureVersion(versionId: string): Promise<ResponseEnvelope<boolean>>

	listDiscountPolicies(): Promise<ResponseEnvelope<DiscountPolicyDTO[]>>
	createDiscountPolicy(input: {
		code: string
		name: string
		description: string
		discountType: string
		rateBasisPoints?: number
		amountMinor?: number
		feeHeadId?: string
		requiresApproval: boolean
	}): Promise<ResponseEnvelope<string>>
	listStudentDiscounts(studentId: string, academicYearId: string): Promise<ResponseEnvelope<StudentDiscountDTO[]>>
	listPendingDiscounts(): Promise<ResponseEnvelope<StudentDiscountDTO[]>>
	applyStudentDiscount(input: {
		studentId: string
		discountPolicyId: string
		academicYearId: string
		validFrom: string
		validTo?: string
	}): Promise<ResponseEnvelope<string>>
	decideStudentDiscount(input: {
		discountId: string
		approved: boolean
		rejectionReason?: string
	}): Promise<ResponseEnvelope<boolean>>
}

export interface IInvoiceService {
	listInvoices(filter: InvoiceFilterDTO): Promise<ResponseEnvelope<InvoiceDTO[]>>
	getInvoice(id: string): Promise<ResponseEnvelope<InvoiceDetailDTO>>
	generateInvoicesBatch(input: {
		academicYearId: string
		classId: string
		billingPeriod: string
		issueDate: string
		dueDate: string
	}): Promise<ResponseEnvelope<string[]>>
	voidInvoice(input: {
		invoiceId: string
		reason: string
	}): Promise<ResponseEnvelope<boolean>>
}

export interface IPaymentService {
	recordPayment(input: {
		studentId: string
		paymentDate: string
		paymentMethod: string
		amountMinor: number
		referenceNo?: string
		payerName?: string
		payerPhone?: string
		notes?: string
	}): Promise<ResponseEnvelope<RecordPaymentResultDTO>>
	reversePayment(input: {
		paymentId: string
		reason: string
	}): Promise<ResponseEnvelope<boolean>>
	getPayment(id: string): Promise<ResponseEnvelope<PaymentDetailDTO>>
	listPayments(filter: {
		studentId?: string
		status?: string
		paymentMethod?: string
		limit?: number
	}): Promise<ResponseEnvelope<PaymentDTO[]>>
	getStudentCreditBalance(studentId: string): Promise<ResponseEnvelope<StudentCreditBalanceInfoDTO>>
	requestFinancialAdjustment(input: {
		studentId: string
		invoiceId?: string
		invoiceItemId?: string
		adjustmentType: string
		amountMinor: number
		reason: string
	}): Promise<ResponseEnvelope<string>>
	decideFinancialAdjustment(input: {
		adjustmentId: string
		approved: boolean
		decisionNotes?: string
	}): Promise<ResponseEnvelope<boolean>>
	listFinancialAdjustments(filter: {
		studentId?: string
		status?: string
		limit?: number
	}): Promise<ResponseEnvelope<FinancialAdjustmentDTO[]>>
	getFinancialAdjustment(id: string): Promise<ResponseEnvelope<FinancialAdjustmentDTO>>
}

export interface DatabaseStatusDTO {
	isEncrypted: boolean
	encryptionAlgorithm: string
	keyDerivation: string
	vaultPath: string
	workingPath: string
	vaultSizeBytes: number
	lastEncryptedAt: string
}

export interface IBackupService {
	getDatabaseStatus(): Promise<ResponseEnvelope<DatabaseStatusDTO>>
	exportDecryptedDatabase(input: { adminPassword?: string; customSavePath?: string }): Promise<ResponseEnvelope<string>>
	exportEncryptedDatabase(input: { adminPassword?: string; customSavePath?: string }): Promise<ResponseEnvelope<string>>
	restoreDatabase(input: { adminPassword?: string; backupFilePath: string }): Promise<ResponseEnvelope<boolean>>
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
	fees: IFeeService
	invoices: IInvoiceService
	payments: IPaymentService
	backup: IBackupService
}

