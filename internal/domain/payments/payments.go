package payments

import (
	"context"
	"fmt"
	"time"

	"school-management/internal/domain/fees"
	"school-management/internal/domain/shared"
)

type PaymentMethod string

const (
	MethodCash          PaymentMethod = "cash"
	MethodBankTransfer  PaymentMethod = "bank_transfer"
	MethodCheque        PaymentMethod = "cheque"
	MethodPOSCard       PaymentMethod = "pos_card"
	MethodMobileMoney   PaymentMethod = "mobile_money"
	MethodCreditBalance PaymentMethod = "credit_balance"
)

type PaymentStatus string

const (
	PaymentPosted   PaymentStatus = "posted"
	PaymentReversed PaymentStatus = "reversed"
)

type AdjustmentType string

const (
	AdjWaiver           AdjustmentType = "waiver"
	AdjWriteOff         AdjustmentType = "write_off"
	AdjDebitCorrection  AdjustmentType = "debit_correction"
	AdjCreditCorrection AdjustmentType = "credit_correction"
)

type AdjustmentStatus string

const (
	AdjPending  AdjustmentStatus = "pending"
	AdjApproved AdjustmentStatus = "approved"
	AdjRejected AdjustmentStatus = "rejected"
)

type CreditReferenceType string

const (
	CreditOverpayment        CreditReferenceType = "overpayment"
	CreditAdvance            CreditReferenceType = "advance"
	CreditInvoiceApplication CreditReferenceType = "invoice_application"
	CreditRefund             CreditReferenceType = "refund"
	CreditAdjustment         CreditReferenceType = "adjustment"
)

// Payment represents a single recorded student payment.
type Payment struct {
	ID                     shared.ID
	SchoolID               shared.ID
	StudentID              shared.ID
	ReceiptNo              string
	PaymentDate            shared.Date
	PaymentMethod          PaymentMethod
	AmountMinor            int64
	AllocatedAmountMinor   int64
	UnallocatedAmountMinor int64
	ReferenceNo            string
	PayerName              string
	PayerPhone             string
	Status                 PaymentStatus
	RecordedByID           shared.ID
	RecordedAt             time.Time
	Notes                  string
}

type PaymentWithDetails struct {
	Payment
	StudentFirstName string
	StudentLastName  string
	AdmissionNo      string
	RecordedByName   string
}

// PaymentAllocation maps payments to invoices/items.
type PaymentAllocation struct {
	ID                   shared.ID
	PaymentID            shared.ID
	InvoiceID            shared.ID
	InvoiceItemID        *shared.ID
	AllocatedAmountMinor int64
	AllocatedAt          time.Time
	AllocatedByID        shared.ID

	InvoiceNo     string
	BillingPeriod string
	ReceiptNo     string
	PaymentDate   shared.Date
	PaymentMethod PaymentMethod
}

// PaymentReversal records a reversed payment audit event.
type PaymentReversal struct {
	ID             shared.ID
	PaymentID      shared.ID
	ReversalDate   shared.Date
	Reason         string
	ReversedByID   shared.ID
	ReversedAt     time.Time
	ReversedByName string
}

// CreditBalance tracks a student's available credit balance.
type CreditBalance struct {
	ID          shared.ID
	SchoolID    shared.ID
	StudentID   shared.ID
	AmountMinor int64
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// CreditLedgerEntry records every change to student credit balances.
type CreditLedgerEntry struct {
	ID                shared.ID
	CreditBalanceID   shared.ID
	StudentID         shared.ID
	AmountMinor       int64
	BalanceAfterMinor int64
	ReferenceType     CreditReferenceType
	ReferenceID       *string
	Notes             string
	CreatedAt         time.Time
}

// FinancialAdjustment represents fee waivers, write-offs, or balance corrections.
type FinancialAdjustment struct {
	ID             shared.ID
	SchoolID       shared.ID
	StudentID      shared.ID
	InvoiceID      *shared.ID
	InvoiceItemID  *shared.ID
	AdjustmentType AdjustmentType
	AmountMinor    int64
	Reason         string
	Status         AdjustmentStatus
	RequestedByID  shared.ID
	RequestedAt    time.Time
	DecidedByID    *shared.ID
	DecidedAt      *time.Time
	DecisionNotes  string

	StudentFirstName string
	StudentLastName  string
	AdmissionNo      string
	RequestedByName  string
	DecidedByName    *string
	InvoiceNo        *string
}

type PaymentFilter struct {
	SchoolID      shared.ID
	StudentID     *shared.ID
	Status        *PaymentStatus
	PaymentMethod *PaymentMethod
	Limit         int
}

type AdjustmentFilter struct {
	SchoolID  shared.ID
	StudentID *shared.ID
	Status    *AdjustmentStatus
	Limit     int
}

// PlannedAllocation represents an allocated portion of a payment to an invoice.
type PlannedAllocation struct {
	InvoiceID       shared.ID
	InvoiceItemID   *shared.ID
	AllocatedAmount int64
	NewPaidAmount   int64
	NewBalance      int64
	NewStatus       fees.InvoiceStatus
}

// PlanFIFOAllocation calculates how an amount will be distributed across outstanding invoices.
func PlanFIFOAllocation(unpaidInvoices []*fees.Invoice, paymentAmountMinor int64) ([]PlannedAllocation, int64) {
	remaining := paymentAmountMinor
	var plan []PlannedAllocation

	for _, inv := range unpaidInvoices {
		if remaining <= 0 {
			break
		}
		if inv.BalanceMinor <= 0 || inv.Status == fees.InvoicePaid || inv.Status == fees.InvoiceVoid {
			continue
		}

		toAllocate := remaining
		if toAllocate > inv.BalanceMinor {
			toAllocate = inv.BalanceMinor
		}

		remaining -= toAllocate
		newPaid := inv.PaidAmountMinor + toAllocate
		newBalance := inv.BalanceMinor - toAllocate

		var newStatus fees.InvoiceStatus
		if newBalance == 0 {
			newStatus = fees.InvoicePaid
		} else {
			newStatus = fees.InvoicePartiallyPaid
		}

		plan = append(plan, PlannedAllocation{
			InvoiceID:       inv.ID,
			AllocatedAmount: toAllocate,
			NewPaidAmount:   newPaid,
			NewBalance:      newBalance,
			NewStatus:       newStatus,
		})
	}

	return plan, remaining
}

// FormatReceiptNo formats a receipt number according to the school sequence standard.
func FormatReceiptNo(year int, seq int64) string {
	return fmt.Sprintf("RCP-%d-%06d", year, seq)
}

// Repository defines persistence boundaries for payments, allocations, credits, and adjustments.
type Repository interface {
	NextReceiptNumber(ctx context.Context, schoolID shared.ID, year int) (string, error)
	InsertPayment(ctx context.Context, p *Payment) error
	UpdatePaymentAllocation(ctx context.Context, paymentID shared.ID, allocatedMinor, unallocatedMinor int64) error
	UpdatePaymentStatus(ctx context.Context, paymentID shared.ID, status PaymentStatus) error
	GetPaymentByID(ctx context.Context, id shared.ID) (*PaymentWithDetails, error)
	ListPayments(ctx context.Context, filter PaymentFilter) ([]*PaymentWithDetails, error)

	InsertAllocation(ctx context.Context, a *PaymentAllocation) error
	ListAllocationsByPayment(ctx context.Context, paymentID shared.ID) ([]*PaymentAllocation, error)
	ListAllocationsByInvoice(ctx context.Context, invoiceID shared.ID) ([]*PaymentAllocation, error)

	InsertReversal(ctx context.Context, r *PaymentReversal) error
	GetReversal(ctx context.Context, paymentID shared.ID) (*PaymentReversal, error)

	GetCreditBalance(ctx context.Context, schoolID, studentID shared.ID) (*CreditBalance, error)
	UpsertCreditBalance(ctx context.Context, cb *CreditBalance) error
	InsertCreditLedger(ctx context.Context, entry *CreditLedgerEntry) error
	ListCreditLedger(ctx context.Context, studentID shared.ID) ([]*CreditLedgerEntry, error)

	InsertAdjustment(ctx context.Context, fa *FinancialAdjustment) error
	GetAdjustmentByID(ctx context.Context, id shared.ID) (*FinancialAdjustment, error)
	UpdateAdjustmentStatus(ctx context.Context, id shared.ID, status AdjustmentStatus, decidedByID shared.ID, decidedAt time.Time, notes string) error
	ListAdjustments(ctx context.Context, filter AdjustmentFilter) ([]*FinancialAdjustment, error)

	GetUnpaidInvoicesForStudent(ctx context.Context, schoolID, studentID shared.ID) ([]*fees.Invoice, error)
}
