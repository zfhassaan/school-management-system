package fees

import (
	"context"
	"strings"
	"time"

	"school-management/internal/domain/shared"
)

// --- Fee Heads ---

type FeeHead struct {
	ID                 shared.ID
	SchoolID           shared.ID
	Code               string
	Name               string
	Description        string
	IsRecurring        bool
	DefaultAmountMinor int64
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Revision           int64
}

func (h *FeeHead) Validate() error {
	var v shared.Validator
	v.Check(strings.TrimSpace(h.Code) != "", "code", "fee head code is required")
	v.Check(strings.TrimSpace(h.Name) != "", "name", "fee head name is required")
	v.Check(h.DefaultAmountMinor >= 0, "defaultAmountMinor", "amount cannot be negative")
	return v.Err()
}

// --- Fee Structures & Versions ---

type StructureVersionStatus string

const (
	VersionDraft     StructureVersionStatus = "draft"
	VersionPublished StructureVersionStatus = "published"
	VersionRetired   StructureVersionStatus = "retired"
)

type FeeFrequency string

const (
	FrequencyMonthly FeeFrequency = "monthly"
	FrequencyTermly  FeeFrequency = "termly"
	FrequencyAnnual  FeeFrequency = "annual"
	FrequencyOneTime FeeFrequency = "one_time"
)

type FeeStructure struct {
	ID             shared.ID
	SchoolID       shared.ID
	AcademicYearID shared.ID
	ClassID        shared.ID
	CampusID       *shared.ID
	Code           string
	Name           string
	Description    string
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Revision       int64
}

type FeeStructureVersion struct {
	ID               shared.ID
	FeeStructureID   shared.ID
	VersionNo        int
	Status           StructureVersionStatus
	Changelog        string
	CreatedByID      shared.ID
	PublishedByID    *shared.ID
	PublishedAt      *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Revision         int64
}

type FeeStructureItem struct {
	ID                     shared.ID
	FeeStructureVersionID  shared.ID
	FeeHeadID              shared.ID
	AmountMinor            int64
	Frequency              FeeFrequency
	DueDay                 int
	IsOptional             bool
	CreatedAt              time.Time
	UpdatedAt              time.Time
	// Expanded view fields
	FeeHeadCode            string
	FeeHeadName            string
}

// --- Discount Policies & Student Discounts ---

type DiscountType string

const (
	DiscountPercentage DiscountType = "percentage"
	DiscountFixed      DiscountType = "fixed"
)

type DiscountPolicy struct {
	ID               shared.ID
	SchoolID         shared.ID
	Code             string
	Name             string
	Description      string
	DiscountType     DiscountType
	RateBasisPoints  int // 10000 = 100.00%
	AmountMinor      int64
	FeeHeadID        *shared.ID
	RequiresApproval bool
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Revision         int64
}

type StudentDiscountStatus string

const (
	DiscountPending  StudentDiscountStatus = "pending"
	DiscountApproved StudentDiscountStatus = "approved"
	DiscountRejected StudentDiscountStatus = "rejected"
	DiscountRevoked  StudentDiscountStatus = "revoked"
)

type StudentDiscount struct {
	ID                 shared.ID
	StudentID          shared.ID
	DiscountPolicyID   shared.ID
	AcademicYearID     shared.ID
	Status             StudentDiscountStatus
	ApprovedByID       *shared.ID
	ApprovedAt         *time.Time
	RejectionReason    string
	ValidFrom          shared.Date
	ValidTo            *shared.Date
	CreatedByID        shared.ID
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Revision           int64
	// Expanded view fields
	PolicyCode         string
	PolicyName         string
	DiscountType       DiscountType
	RateBasisPoints    int
	PolicyAmountMinor  int64
	PolicyFeeHeadID    *shared.ID
	StudentFirstName   string
	StudentLastName    string
	AdmissionNo        string
}

// CalculateDiscount returns the integer discount amount in minor currency units for a line item.
func CalculateDiscount(itemAmountMinor int64, discountType DiscountType, basisPoints int, fixedMinor int64) int64 {
	if discountType == DiscountPercentage {
		if basisPoints <= 0 {
			return 0
		}
		if basisPoints >= 10000 {
			return itemAmountMinor
		}
		// Exact integer arithmetic
		discount := (itemAmountMinor * int64(basisPoints)) / 10000
		if discount > itemAmountMinor {
			return itemAmountMinor
		}
		return discount
	}

	if fixedMinor >= itemAmountMinor {
		return itemAmountMinor
	}
	return fixedMinor
}

// --- Invoicing ---

type InvoiceStatus string

const (
	InvoiceDraft         InvoiceStatus = "draft"
	InvoiceIssued        InvoiceStatus = "issued"
	InvoicePartiallyPaid InvoiceStatus = "partially_paid"
	InvoicePaid          InvoiceStatus = "paid"
	InvoiceOverdue       InvoiceStatus = "overdue"
	InvoiceVoid          InvoiceStatus = "void"
)

type Invoice struct {
	ID                  shared.ID
	SchoolID            shared.ID
	StudentID           shared.ID
	AcademicYearID      shared.ID
	InvoiceNo           string
	IssueDate           shared.Date
	DueDate             shared.Date
	BillingPeriod       string // e.g. "2026-09"
	SubtotalAmountMinor int64
	DiscountAmountMinor int64
	TotalAmountMinor    int64
	PaidAmountMinor     int64
	BalanceMinor        int64
	Status              InvoiceStatus
	VoidReason          string
	VoidedByID          *shared.ID
	VoidedAt            *time.Time
	CreatedByID         shared.ID
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Revision            int64
	// Expanded view fields
	StudentFirstName    string
	StudentLastName     string
	AdmissionNo         string
}

type InvoiceItem struct {
	ID                  shared.ID
	InvoiceID           shared.ID
	FeeHeadID           shared.ID
	FeeHeadName         string
	Description         string
	OriginalAmountMinor int64
	DiscountAmountMinor int64
	NetAmountMinor      int64
	CreatedAt           time.Time
}

type InvoiceFilter struct {
	SchoolID       shared.ID
	Status         string
	AcademicYearID string
	StudentID      string
	Limit          int
}

// --- Repository Interface ---

type Repository interface {
	// Fee Heads
	ListFeeHeads(ctx context.Context, schoolID shared.ID) ([]FeeHead, error)
	GetFeeHead(ctx context.Context, id shared.ID) (*FeeHead, error)
	InsertFeeHead(ctx context.Context, h *FeeHead) error
	UpdateFeeHead(ctx context.Context, h *FeeHead) error

	// Fee Structures
	ListFeeStructures(ctx context.Context, schoolID shared.ID, academicYearID shared.ID) ([]FeeStructure, error)
	GetFeeStructure(ctx context.Context, id shared.ID) (*FeeStructure, error)
	GetFeeStructureForClass(ctx context.Context, academicYearID shared.ID, classID shared.ID) (*FeeStructure, error)
	InsertFeeStructure(ctx context.Context, s *FeeStructure) error

	// Fee Structure Versions & Items
	ListFeeStructureVersions(ctx context.Context, feeStructureID shared.ID) ([]FeeStructureVersion, error)
	GetFeeStructureVersion(ctx context.Context, id shared.ID) (*FeeStructureVersion, error)
	GetPublishedFeeStructureVersion(ctx context.Context, feeStructureID shared.ID) (*FeeStructureVersion, error)
	InsertFeeStructureVersion(ctx context.Context, v *FeeStructureVersion) error
	PublishFeeStructureVersion(ctx context.Context, id shared.ID, publishedByID shared.ID, publishedAt time.Time) error
	DeleteFeeStructureVersion(ctx context.Context, id shared.ID) error
	ListFeeStructureItems(ctx context.Context, versionID shared.ID) ([]FeeStructureItem, error)
	InsertFeeStructureItem(ctx context.Context, item *FeeStructureItem) error

	// Discount Policies
	ListDiscountPolicies(ctx context.Context, schoolID shared.ID) ([]DiscountPolicy, error)
	GetDiscountPolicy(ctx context.Context, id shared.ID) (*DiscountPolicy, error)
	InsertDiscountPolicy(ctx context.Context, p *DiscountPolicy) error
	UpdateDiscountPolicy(ctx context.Context, p *DiscountPolicy) error

	// Student Discounts
	ListStudentDiscounts(ctx context.Context, studentID shared.ID, academicYearID shared.ID) ([]StudentDiscount, error)
	ListPendingStudentDiscounts(ctx context.Context, schoolID shared.ID) ([]StudentDiscount, error)
	GetStudentDiscount(ctx context.Context, id shared.ID) (*StudentDiscount, error)
	InsertStudentDiscount(ctx context.Context, d *StudentDiscount) error
	UpdateStudentDiscountStatus(ctx context.Context, id shared.ID, status StudentDiscountStatus, approvedByID *shared.ID, approvedAt *time.Time, rejectionReason string) error

	// Invoices
	ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, error)
	GetInvoice(ctx context.Context, id shared.ID) (*Invoice, error)
	GetInvoiceByNo(ctx context.Context, schoolID shared.ID, invoiceNo string) (*Invoice, error)
	InsertInvoice(ctx context.Context, inv *Invoice) error
	UpdateInvoiceStatus(ctx context.Context, id shared.ID, status InvoiceStatus, paidAmountMinor int64, balanceMinor int64) error
	VoidInvoice(ctx context.Context, id shared.ID, reason string, voidedByID shared.ID, voidedAt time.Time) error
	ListInvoiceItems(ctx context.Context, invoiceID shared.ID) ([]InvoiceItem, error)
	InsertInvoiceItem(ctx context.Context, item *InvoiceItem) error
	CountInvoicesInPeriod(ctx context.Context, schoolID shared.ID, billingPeriod string) (int64, error)
}

