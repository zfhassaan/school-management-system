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
	legalName?: string
	code: string
	timezone: string
	currencyCode: string
	currencyMinorUnits: number
	dateFormat?: string
	phone?: string
	email?: string
	website?: string
	address?: {
		line1?: string
		line2?: string
		city?: string
		region?: string
		postalCode?: string
		country?: string
	}
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

// Phase 3: Academics, Students, Admissions, Documents

export interface AcademicYearDTO {
	ID: string
	SchoolID: string
	Code: string
	Name: string
	StartsOn: { Year: number; Month: number; Day: number } | string
	EndsOn: { Year: number; Month: number; Day: number } | string
	Status: 'planned' | 'active' | 'closed'
}

export interface AcademicTermDTO {
	ID: string
	AcademicYearID: string
	Sequence: number
	Name: string
	StartsOn: { Year: number; Month: number; Day: number } | string
	EndsOn: { Year: number; Month: number; Day: number } | string
}

export interface CampusDTO {
	ID: string
	SchoolID: string
	Code: string
	Name: string
	Address?: string
	IsActive: boolean
}

export interface ClassDTO {
	ID: string
	SchoolID: string
	Code: string
	Name: string
	LevelOrder: number
	NextClassID?: string
	IsActive: boolean
}

export interface SectionDTO {
	ID: string
	ClassID: string
	AcademicYearID: string
	CampusID?: string
	Name: string
	Capacity?: number
	ClassTeacherUserID?: string
	IsActive: boolean
}

export interface StudentDTO {
	ID: string
	SchoolID: string
	AdmissionNo: string
	FirstName: string
	MiddleName?: string
	LastName: string
	DateOfBirth: { Year: number; Month: number; Day: number } | string
	Gender: 'male' | 'female' | 'other' | 'unspecified'
	NationalID?: string
	Address?: any
	PhotoDocID?: string
	AdmissionDate: { Year: number; Month: number; Day: number } | string
	Status: 'active' | 'inactive' | 'graduated' | 'withdrawn' | 'alumni'
	Notes?: string
	CreatedAt: string
	UpdatedAt: string
}

export interface StudentFilterDTO {
	schoolId?: string
	status?: string
	classId?: string
	sectionId?: string
	searchQuery?: string
	cursorName?: string
	cursorId?: string
	limit?: number
}

export interface StatusHistoryDTO {
	ID: string
	StudentID: string
	FromStatus: string
	ToStatus: string
	EffectiveOn: any
	Reason: string
	ChangedBy: string
	ChangedAt: string
}

export interface GuardianDTO {
	ID: string
	SchoolID: string
	FirstName: string
	LastName: string
	PhonePrimary: string
	PhoneSecondary?: string
	Email?: string
	Occupation?: string
	NationalID?: string
	Address?: any
}

export interface StudentGuardianLinkDTO {
	StudentID: string
	GuardianID: string
	Relationship: string
	IsPrimary: boolean
	IsEmergencyContact: boolean
	ReceivesBilling: boolean
	CanPickup: boolean
	Priority: number
	FirstName?: string
	LastName?: string
	PhonePrimary?: string
	Email?: string
}

export interface EnrollmentDTO {
	ID: string
	StudentID: string
	AcademicYearID: string
	ClassID: string
	SectionID?: string
	RollNo?: string
	Status: string
	EnrolledOn: any
	EndedOn?: any
	ClassName?: string
	SectionName?: string
	YearCode?: string
}

export interface PromotionBatchDTO {
	ID: string
	FromAcademicYearID: string
	ToAcademicYearID: string
	Status: string
	SummaryJSON: string
	ExecutedAt?: string
}

export interface FormFieldDTO {
	id: string
	label: string
	type: 'text' | 'number' | 'date' | 'dropdown' | 'radio' | 'checkbox' | 'phone' | 'email' | 'address' | 'file_upload' | 'textarea'
	placeholder?: string
	required: boolean
	options?: string[]
	dependsOn?: {
		fieldId: string
		equals: string
	}
}

export interface FormSectionDTO {
	id: string
	title: string
	description?: string
	fields: FormFieldDTO[]
}

export interface FormSchemaDTO {
	title: string
	description?: string
	sections: FormSectionDTO[]
	requiredDocuments?: string[]
}

export interface FormDefinitionDTO {
	ID: string
	SchoolID: string
	Code: string
	Name: string
	Description: string
	IsActive: boolean
	id?: string
	schoolId?: string
	code?: string
	name?: string
	description?: string
	isActive?: boolean
}

export interface FormVersionDTO {
	ID: string
	FormDefinitionID: string
	VersionNo: number
	Status: 'draft' | 'published' | 'retired'
	Schema: FormSchemaDTO
	SchemaHash: string
	PublishedAt?: string
	id?: string
	formDefinitionId?: string
	versionNo?: number
	status?: 'draft' | 'published' | 'retired' | string
	schema?: FormSchemaDTO
}

export interface ApplicationDTO {
	ID: string
	SchoolID: string
	ApplicationNo: string
	FormVersionID: string
	AcademicYearID: string
	AppliedClassID: string
	Status: 'draft' | 'submitted' | 'under_review' | 'approved' | 'rejected' | 'enrolled' | 'withdrawn'
	ApplicantFirstName: string
	ApplicantLastName: string
	ApplicantDOB?: any
	ApplicantGender: string
	GuardianPhoneNormalized: string
	StudentID?: string
	Response: Record<string, any>
	SubmittedAt?: string
	DecidedAt?: string
	DecisionReason?: string
	ClassName?: string
	YearCode?: string
	GuardianName?: string
	AdmissionNo?: string
	RollNo?: string
}

export interface DocumentDTO {
	ID: string
	StudentID?: string
	ApplicationID?: string
	DocumentTypeCode: string
	FileName: string
	MimeType: string
	SizeBytes: number
	SHA256: string
	CreatedAt: string
}

// Phase 4: Fees & Invoicing

export interface FeeHeadDTO {
	ID: string
	SchoolID: string
	Code: string
	Name: string
	Description: string
	IsRecurring: boolean
	DefaultAmountMinor: number
	IsActive: boolean
}

export interface FeeStructureDTO {
	ID: string
	SchoolID: string
	AcademicYearID: string
	ClassID: string
	CampusID?: string
	Code: string
	Name: string
	Description: string
	IsActive: boolean
}

export interface FeeStructureVersionDTO {
	ID: string
	FeeStructureID: string
	VersionNo: number
	Status: 'draft' | 'published' | 'retired'
	Changelog: string
	CreatedByID: string
	PublishedByID?: string
	PublishedAt?: string
}

export interface FeeStructureItemDTO {
	ID: string
	FeeStructureVersionID: string
	FeeHeadID: string
	AmountMinor: number
	Frequency: 'monthly' | 'termly' | 'annual' | 'one_time'
	DueDay: number
	IsOptional: boolean
	FeeHeadCode?: string
	FeeHeadName?: string
}

export interface FeeStructureVersionDetailDTO {
	version: FeeStructureVersionDTO
	items: FeeStructureItemDTO[]
}

export interface DiscountPolicyDTO {
	ID: string
	SchoolID: string
	Code: string
	Name: string
	Description: string
	DiscountType: 'percentage' | 'fixed'
	RateBasisPoints: number
	AmountMinor: number
	FeeHeadID?: string
	RequiresApproval: boolean
	IsActive: boolean
}

export interface StudentDiscountDTO {
	ID: string
	StudentID: string
	DiscountPolicyID: string
	AcademicYearID: string
	Status: 'pending' | 'approved' | 'rejected' | 'revoked'
	ApprovedByID?: string
	ApprovedAt?: string
	RejectionReason: string
	ValidFrom: any
	ValidTo?: any
	PolicyCode?: string
	PolicyName?: string
	DiscountType?: string
	RateBasisPoints?: number
	PolicyAmountMinor?: number
	StudentFirstName?: string
	StudentLastName?: string
	AdmissionNo?: string
}

export interface InvoiceDTO {
	ID: string
	SchoolID: string
	StudentID: string
	AcademicYearID: string
	InvoiceNo: string
	IssueDate: any
	DueDate: any
	BillingPeriod: string
	SubtotalAmountMinor: number
	DiscountAmountMinor: number
	TotalAmountMinor: number
	PaidAmountMinor: number
	BalanceMinor: number
	Status: 'draft' | 'issued' | 'partially_paid' | 'paid' | 'overdue' | 'void'
	VoidReason?: string
	StudentFirstName?: string
	StudentLastName?: string
	AdmissionNo?: string
	CreatedAt: string
}

export interface InvoiceItemDTO {
	ID: string
	InvoiceID: string
	FeeHeadID: string
	FeeHeadName: string
	Description: string
	OriginalAmountMinor: number
	DiscountAmountMinor: number
	NetAmountMinor: number
}

export interface InvoiceDetailDTO {
	invoice: InvoiceDTO
	items: InvoiceItemDTO[]
}

export interface InvoiceFilterDTO {
	schoolId?: string
	status?: string
	academicYearId?: string
	studentId?: string
	limit?: number
}

// Phase 5: Payments, Allocations, Credits & Financial Adjustments

export interface PaymentDTO {
	ID: string
	SchoolID: string
	StudentID: string
	ReceiptNo: string
	PaymentDate: any
	PaymentMethod: 'cash' | 'bank_transfer' | 'cheque' | 'pos_card' | 'mobile_money' | 'credit_balance'
	AmountMinor: number
	AllocatedAmountMinor: number
	UnallocatedAmountMinor: number
	ReferenceNo: string
	PayerName: string
	PayerPhone: string
	Status: 'posted' | 'reversed'
	RecordedByID: string
	RecordedByName?: string
	RecordedAt: string
	Notes: string
	StudentFirstName?: string
	StudentLastName?: string
	AdmissionNo?: string
}

export interface PaymentAllocationDTO {
	ID: string
	PaymentID: string
	InvoiceID: string
	InvoiceItemID?: string
	AllocatedAmountMinor: number
	AllocatedAt: string
	AllocatedByID: string
	InvoiceNo?: string
	BillingPeriod?: string
	ReceiptNo?: string
	PaymentDate?: any
	PaymentMethod?: string
}

export interface PaymentReversalDTO {
	ID: string
	PaymentID: string
	ReversalDate: any
	Reason: string
	ReversedByID: string
	ReversedByName?: string
	ReversedAt: string
}

export interface PaymentDetailDTO {
	payment: PaymentDTO
	allocations: PaymentAllocationDTO[]
	reversal?: PaymentReversalDTO
}

export interface CreditLedgerEntryDTO {
	ID: string
	CreditBalanceID: string
	StudentID: string
	AmountMinor: number
	BalanceAfterMinor: number
	ReferenceType: 'overpayment' | 'advance' | 'invoice_application' | 'refund' | 'adjustment'
	ReferenceID?: string
	Notes: string
	CreatedAt: string
}

export interface StudentCreditBalanceInfoDTO {
	balanceMinor: number
	ledger: CreditLedgerEntryDTO[]
}

export interface FinancialAdjustmentDTO {
	ID: string
	SchoolID: string
	StudentID: string
	InvoiceID?: string
	InvoiceItemID?: string
	AdjustmentType: 'waiver' | 'write_off' | 'debit_correction' | 'credit_correction'
	AmountMinor: number
	Reason: string
	Status: 'pending' | 'approved' | 'rejected'
	RequestedByID: string
	RequestedByName?: string
	RequestedAt: string
	DecidedByID?: string
	DecidedByName?: string
	DecidedAt?: string
	DecisionNotes: string
	StudentFirstName?: string
	StudentLastName?: string
	AdmissionNo?: string
	InvoiceNo?: string
}

export interface RecordPaymentResultDTO {
	paymentId: string
	receiptNo: string
	allocatedMinor: number
	unallocatedMinor: number
}




