package payments

import (
	"context"
	"errors"
	"fmt"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/fees"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/payments"
	"school-management/internal/domain/shared"
)

type PaymentCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewPaymentCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *PaymentCommandsHandler {
	return &PaymentCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

// --- Inputs ---

type RecordPaymentInput struct {
	StudentID     shared.ID              `json:"studentId"`
	PaymentDate   string                 `json:"paymentDate"`
	PaymentMethod payments.PaymentMethod `json:"paymentMethod"`
	AmountMinor   int64                  `json:"amountMinor"`
	ReferenceNo   string                 `json:"referenceNo"`
	PayerName     string                 `json:"payerName"`
	PayerPhone    string                 `json:"payerPhone"`
	Notes         string                 `json:"notes"`
}

type RecordPaymentResult struct {
	PaymentID   shared.ID `json:"paymentId"`
	ReceiptNo   string    `json:"receiptNo"`
	Allocated   int64     `json:"allocatedMinor"`
	Unallocated int64     `json:"unallocatedMinor"`
}

type ReversePaymentInput struct {
	PaymentID shared.ID `json:"paymentId"`
	Reason    string    `json:"reason"`
}

type RequestFinancialAdjustmentInput struct {
	StudentID      shared.ID               `json:"studentId"`
	InvoiceID      *shared.ID              `json:"invoiceId,omitempty"`
	InvoiceItemID  *shared.ID              `json:"invoiceItemId,omitempty"`
	AdjustmentType payments.AdjustmentType `json:"adjustmentType"`
	AmountMinor    int64                   `json:"amountMinor"`
	Reason         string                  `json:"reason"`
}

type DecideFinancialAdjustmentInput struct {
	AdjustmentID  shared.ID `json:"adjustmentId"`
	Approved      bool      `json:"approved"`
	DecisionNotes string    `json:"decisionNotes"`
}

// --- Handlers ---

func (h *PaymentCommandsHandler) RecordPayment(ctx context.Context, in RecordPaymentInput) (*RecordPaymentResult, error) {
	actor, err := authz.Require(ctx, identity.PermPaymentsRecord)
	if err != nil {
		return nil, err
	}

	if in.StudentID.IsZero() {
		return nil, shared.Validation("studentId is required", nil)
	}
	if in.AmountMinor <= 0 {
		return nil, shared.Validation("payment amount must be greater than zero", nil)
	}
	if in.PaymentMethod == "" {
		in.PaymentMethod = payments.MethodCash
	}

	now := h.clock.Now()
	pDate, err := shared.ParseDate(in.PaymentDate)
	if err != nil {
		pDate = shared.DateOf(now)
	}

	var result RecordPaymentResult

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		// 1. Generate sequential receipt number
		receiptNo, err := r.Payments.NextReceiptNumber(ctx, sch.ID, pDate.Year)
		if err != nil {
			return fmt.Errorf("failed to generate receipt number: %w", err)
		}

		paymentID := h.ids.NewID()

		payment := &payments.Payment{
			ID:                     paymentID,
			SchoolID:               sch.ID,
			StudentID:              in.StudentID,
			ReceiptNo:              receiptNo,
			PaymentDate:            pDate,
			PaymentMethod:          in.PaymentMethod,
			AmountMinor:            in.AmountMinor,
			AllocatedAmountMinor:   0,
			UnallocatedAmountMinor: in.AmountMinor,
			ReferenceNo:            in.ReferenceNo,
			PayerName:              in.PayerName,
			PayerPhone:             in.PayerPhone,
			Status:                 payments.PaymentPosted,
			RecordedByID:           actor.UserID,
			RecordedAt:             now,
			Notes:                  in.Notes,
		}

		if err := r.Payments.InsertPayment(ctx, payment); err != nil {
			return fmt.Errorf("failed to insert payment: %w", err)
		}

		// 2. Fetch outstanding unpaid invoices for the student (FIFO order)
		unpaidInvoices, err := r.Payments.GetUnpaidInvoicesForStudent(ctx, sch.ID, in.StudentID)
		if err != nil {
			return fmt.Errorf("failed to query unpaid invoices: %w", err)
		}

		// 3. Plan FIFO Allocation
		plan, unallocated := payments.PlanFIFOAllocation(unpaidInvoices, in.AmountMinor)
		allocatedMinor := in.AmountMinor - unallocated

		for _, alloc := range plan {
			allocationID := h.ids.NewID()
			allocRecord := &payments.PaymentAllocation{
				ID:                   allocationID,
				PaymentID:            paymentID,
				InvoiceID:            alloc.InvoiceID,
				InvoiceItemID:        alloc.InvoiceItemID,
				AllocatedAmountMinor: alloc.AllocatedAmount,
				AllocatedAt:          now,
				AllocatedByID:        actor.UserID,
			}
			if err := r.Payments.InsertAllocation(ctx, allocRecord); err != nil {
				return fmt.Errorf("failed to insert allocation: %w", err)
			}

			// Update invoice state
			if err := r.Fees.UpdateInvoiceStatus(ctx, alloc.InvoiceID, alloc.NewStatus, alloc.NewPaidAmount, alloc.NewBalance); err != nil {
				return fmt.Errorf("failed to update invoice balance: %w", err)
			}
		}

		// 4. Update payment allocation totals
		if err := r.Payments.UpdatePaymentAllocation(ctx, paymentID, allocatedMinor, unallocated); err != nil {
			return fmt.Errorf("failed to update payment allocation: %w", err)
		}

		// 5. Overpayment / unallocated amount is credited to student credit balance
		if unallocated > 0 {
			var currentCredit int64
			cb, err := r.Payments.GetCreditBalance(ctx, sch.ID, in.StudentID)
			if err != nil {
				if !errors.Is(err, shared.ErrNotFound) {
					return fmt.Errorf("failed to fetch student credit balance: %w", err)
				}
				cb = &payments.CreditBalance{
					ID:          h.ids.NewID(),
					SchoolID:    sch.ID,
					StudentID:   in.StudentID,
					AmountMinor: 0,
					CreatedAt:   now,
				}
			} else {
				currentCredit = cb.AmountMinor
			}

			newCredit := currentCredit + unallocated
			cb.AmountMinor = newCredit
			cb.UpdatedAt = now

			if err := r.Payments.UpsertCreditBalance(ctx, cb); err != nil {
				return fmt.Errorf("failed to update credit balance: %w", err)
			}

			refStr := paymentID.String()
			ledger := &payments.CreditLedgerEntry{
				ID:                h.ids.NewID(),
				CreditBalanceID:   cb.ID,
				StudentID:         in.StudentID,
				AmountMinor:       unallocated,
				BalanceAfterMinor: newCredit,
				ReferenceType:     payments.CreditOverpayment,
				ReferenceID:       &refStr,
				Notes:             fmt.Sprintf("Overpayment on receipt %s", receiptNo),
				CreatedAt:         now,
			}
			if err := r.Payments.InsertCreditLedger(ctx, ledger); err != nil {
				return fmt.Errorf("failed to record credit ledger entry: %w", err)
			}
		}

		result = RecordPaymentResult{
			PaymentID:   paymentID,
			ReceiptNo:   receiptNo,
			Allocated:   allocatedMinor,
			Unallocated: unallocated,
		}

		// Audit Log
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "payment.recorded",
			EntityType: "payment",
			EntityID:   paymentID.String(),
			Reason:     "student payment recorded",
			After: map[string]any{
				"receiptNo":   receiptNo,
				"studentId":   in.StudentID,
				"amountMinor": in.AmountMinor,
				"allocated":   allocatedMinor,
				"unallocated": unallocated,
			},
		})
	})

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (h *PaymentCommandsHandler) ReversePayment(ctx context.Context, in ReversePaymentInput) error {
	actor, err := authz.Require(ctx, identity.PermPaymentsReverse)
	if err != nil {
		return err
	}

	if in.PaymentID.IsZero() {
		return shared.Validation("paymentId is required", nil)
	}
	if in.Reason == "" {
		return shared.Validation("reason is required for payment reversal", nil)
	}

	now := h.clock.Now()
	today := shared.DateOf(now)

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		p, err := r.Payments.GetPaymentByID(ctx, in.PaymentID)
		if err != nil {
			return err
		}
		if p.Status == payments.PaymentReversed {
			return shared.Conflict("payment has already been reversed")
		}

		// 1. Fetch all allocations made by this payment
		allocations, err := r.Payments.ListAllocationsByPayment(ctx, in.PaymentID)
		if err != nil {
			return fmt.Errorf("failed to list payment allocations: %w", err)
		}

		// 2. Revert invoice balances for each allocation
		for _, alloc := range allocations {
			inv, err := r.Fees.GetInvoice(ctx, alloc.InvoiceID)
			if err != nil {
				return fmt.Errorf("failed to fetch invoice %s: %w", alloc.InvoiceID, err)
			}

			// Subtract payment allocation from invoice paid amount
			newPaid := inv.PaidAmountMinor - alloc.AllocatedAmountMinor
			newBalance := inv.BalanceMinor + alloc.AllocatedAmountMinor
			if newPaid < 0 {
				newPaid = 0
			}

			newStatus := fees.InvoiceIssued
			if newPaid > 0 {
				newStatus = fees.InvoicePartiallyPaid
			}

			if err := r.Fees.UpdateInvoiceStatus(ctx, alloc.InvoiceID, newStatus, newPaid, newBalance); err != nil {
				return fmt.Errorf("failed to revert invoice %s: %w", alloc.InvoiceID, err)
			}
		}

		// 3. If payment had unallocated amount placed in credit balance, deduct it
		if p.UnallocatedAmountMinor > 0 {
			cb, err := r.Payments.GetCreditBalance(ctx, sch.ID, p.StudentID)
			if err == nil && cb != nil {
				newCredit := cb.AmountMinor - p.UnallocatedAmountMinor
				if newCredit < 0 {
					newCredit = 0
				}
				cb.AmountMinor = newCredit
				cb.UpdatedAt = now
				_ = r.Payments.UpsertCreditBalance(ctx, cb)

				refStr := in.PaymentID.String()
				_ = r.Payments.InsertCreditLedger(ctx, &payments.CreditLedgerEntry{
					ID:                h.ids.NewID(),
					CreditBalanceID:   cb.ID,
					StudentID:         p.StudentID,
					AmountMinor:       -p.UnallocatedAmountMinor,
					BalanceAfterMinor: newCredit,
					ReferenceType:     payments.CreditRefund,
					ReferenceID:       &refStr,
					Notes:             fmt.Sprintf("Reversal of payment receipt %s", p.ReceiptNo),
					CreatedAt:         now,
				})
			}
		}

		// 4. Update payment status to reversed
		if err := r.Payments.UpdatePaymentStatus(ctx, in.PaymentID, payments.PaymentReversed); err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		// 5. Insert reversal record
		reversal := &payments.PaymentReversal{
			ID:           h.ids.NewID(),
			PaymentID:    in.PaymentID,
			ReversalDate: today,
			Reason:       in.Reason,
			ReversedByID: actor.UserID,
			ReversedAt:   now,
		}
		if err := r.Payments.InsertReversal(ctx, reversal); err != nil {
			return fmt.Errorf("failed to insert reversal record: %w", err)
		}

		// Audit Log
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "payment.reversed",
			EntityType: "payment",
			EntityID:   in.PaymentID.String(),
			Reason:     in.Reason,
			After: map[string]any{
				"receiptNo": p.ReceiptNo,
				"reason":    in.Reason,
			},
		})
	})
}

func (h *PaymentCommandsHandler) RequestFinancialAdjustment(ctx context.Context, in RequestFinancialAdjustmentInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermAdjustmentsRequest)
	if err != nil {
		return "", err
	}

	if in.StudentID.IsZero() {
		return "", shared.Validation("studentId is required", nil)
	}
	if in.AmountMinor <= 0 {
		return "", shared.Validation("adjustment amount must be greater than zero", nil)
	}
	if in.Reason == "" {
		return "", shared.Validation("reason is required", nil)
	}

	adjID := h.ids.NewID()
	now := h.clock.Now()

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		fa := &payments.FinancialAdjustment{
			ID:             adjID,
			SchoolID:       sch.ID,
			StudentID:      in.StudentID,
			InvoiceID:      in.InvoiceID,
			InvoiceItemID:  in.InvoiceItemID,
			AdjustmentType: in.AdjustmentType,
			AmountMinor:    in.AmountMinor,
			Reason:         in.Reason,
			Status:         payments.AdjPending,
			RequestedByID:  actor.UserID,
			RequestedAt:    now,
		}
		if err := r.Payments.InsertAdjustment(ctx, fa); err != nil {
			return fmt.Errorf("failed to insert adjustment: %w", err)
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "adjustment.requested",
			EntityType: "adjustment",
			EntityID:   adjID.String(),
			Reason:     in.Reason,
			After: map[string]any{
				"studentId": in.StudentID,
				"amount":    in.AmountMinor,
				"type":      in.AdjustmentType,
			},
		})
	})

	if err != nil {
		return "", err
	}
	return adjID, nil
}

func (h *PaymentCommandsHandler) DecideFinancialAdjustment(ctx context.Context, in DecideFinancialAdjustmentInput) error {
	actor, err := authz.Require(ctx, identity.PermAdjustmentsApprove)
	if err != nil {
		return err
	}

	if in.AdjustmentID.IsZero() {
		return shared.Validation("adjustmentId is required", nil)
	}

	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		adj, err := r.Payments.GetAdjustmentByID(ctx, in.AdjustmentID)
		if err != nil {
			return err
		}
		if adj.Status != payments.AdjPending {
			return shared.Conflict("financial adjustment has already been decided")
		}

		targetStatus := payments.AdjRejected
		if in.Approved {
			targetStatus = payments.AdjApproved

			// Apply financial impact if approved
			if adj.InvoiceID != nil {
				inv, err := r.Fees.GetInvoice(ctx, *adj.InvoiceID)
				if err != nil {
					return fmt.Errorf("failed to fetch invoice %s: %w", *adj.InvoiceID, err)
				}

				if adj.AdjustmentType == payments.AdjWaiver || adj.AdjustmentType == payments.AdjWriteOff {
					newBalance := inv.BalanceMinor - adj.AmountMinor
					if newBalance < 0 {
						newBalance = 0
					}
					var newStatus fees.InvoiceStatus = inv.Status
					if newBalance == 0 {
						newStatus = fees.InvoicePaid
					}
					if err := r.Fees.UpdateInvoiceStatus(ctx, inv.ID, newStatus, inv.PaidAmountMinor, newBalance); err != nil {
						return fmt.Errorf("failed to adjust invoice balance: %w", err)
					}
				}
			}
		}

		if err := r.Payments.UpdateAdjustmentStatus(ctx, in.AdjustmentID, targetStatus, actor.UserID, now, in.DecisionNotes); err != nil {
			return fmt.Errorf("failed to update adjustment status: %w", err)
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "adjustment.decided",
			EntityType: "adjustment",
			EntityID:   in.AdjustmentID.String(),
			Reason:     in.DecisionNotes,
			After: map[string]any{
				"approved": in.Approved,
				"notes":    in.DecisionNotes,
			},
		})
	})
}
