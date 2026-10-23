package fees

import (
	"context"
	"fmt"
	"strings"
	"time"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/fees"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/payments"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type CreateFeeHeadInput struct {
	Code               string `json:"code"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	IsRecurring        bool   `json:"isRecurring"`
	DefaultAmountMinor int64  `json:"defaultAmountMinor"`
}

type UpdateFeeHeadInput struct {
	ID                 shared.ID `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	IsRecurring        bool      `json:"isRecurring"`
	DefaultAmountMinor int64     `json:"defaultAmountMinor"`
	IsActive           bool      `json:"isActive"`
}

type FeeStructureItemInput struct {
	FeeHeadID   shared.ID         `json:"feeHeadId"`
	AmountMinor int64             `json:"amountMinor"`
	Frequency   fees.FeeFrequency `json:"frequency"`
	DueDay      int               `json:"dueDay"`
	IsOptional  bool              `json:"isOptional"`
}

type CreateFeeStructureInput struct {
	AcademicYearID shared.ID `json:"academicYearId"`
	ClassID        shared.ID `json:"classId"`
	CampusID       *shared.ID `json:"campusId,omitempty"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
}

type CreateFeeStructureVersionInput struct {
	FeeStructureID shared.ID               `json:"feeStructureId"`
	Changelog      string                  `json:"changelog"`
	Items          []FeeStructureItemInput `json:"items"`
}

type CreateDiscountPolicyInput struct {
	Code             string            `json:"code"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	DiscountType     fees.DiscountType `json:"discountType"`
	RateBasisPoints  int               `json:"rateBasisPoints"`
	AmountMinor      int64             `json:"amountMinor"`
	FeeHeadID        *shared.ID        `json:"feeHeadId,omitempty"`
	RequiresApproval bool              `json:"requiresApproval"`
}

type ApplyStudentDiscountInput struct {
	StudentID        shared.ID   `json:"studentId"`
	DiscountPolicyID shared.ID   `json:"discountPolicyId"`
	AcademicYearID   shared.ID   `json:"academicYearId"`
	ValidFrom        shared.Date `json:"validFrom"`
	ValidTo          *shared.Date `json:"validTo,omitempty"`
}

type DecideStudentDiscountInput struct {
	DiscountID      shared.ID `json:"discountId"`
	Approved        bool      `json:"approved"`
	RejectionReason string    `json:"rejectionReason"`
}

type GenerateInvoicesBatchInput struct {
	AcademicYearID shared.ID   `json:"academicYearId"`
	ClassID        *shared.ID  `json:"classId,omitempty"` // nil = all classes in academic year
	BillingPeriod  string      `json:"billingPeriod"`  // e.g. "2026-09"
	IssueDate      shared.Date `json:"issueDate"`
	DueDate        shared.Date `json:"dueDate"`
}

type VoidInvoiceInput struct {
	InvoiceID shared.ID `json:"invoiceId"`
	Reason    string    `json:"reason"`
}

type FeeCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewFeeCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *FeeCommandsHandler {
	return &FeeCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

// --- Fee Heads ---

func (h *FeeCommandsHandler) CreateFeeHead(ctx context.Context, in CreateFeeHeadInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermFeesEdit); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		fh := &fees.FeeHead{
			ID:                 id,
			SchoolID:           sch.ID,
			Code:               strings.TrimSpace(in.Code),
			Name:               strings.TrimSpace(in.Name),
			Description:        strings.TrimSpace(in.Description),
			IsRecurring:        in.IsRecurring,
			DefaultAmountMinor: in.DefaultAmountMinor,
			IsActive:           true,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := fh.Validate(); err != nil {
			return err
		}
		if err := r.Fees.InsertFeeHead(ctx, fh); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "fee_head.created",
			EntityType: "fee_head",
			EntityID:   id.String(),
			Reason:     "fee head created",
			After:      map[string]any{"name": fh.Name, "code": fh.Code},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *FeeCommandsHandler) UpdateFeeHead(ctx context.Context, in UpdateFeeHeadInput) error {
	if _, err := authz.Require(ctx, identity.PermFeesEdit); err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		fh, err := r.Fees.GetFeeHead(ctx, in.ID)
		if err != nil {
			return err
		}
		fh.Name = strings.TrimSpace(in.Name)
		fh.Description = strings.TrimSpace(in.Description)
		fh.IsRecurring = in.IsRecurring
		fh.DefaultAmountMinor = in.DefaultAmountMinor
		fh.IsActive = in.IsActive
		fh.UpdatedAt = now

		if err := fh.Validate(); err != nil {
			return err
		}
		if err := r.Fees.UpdateFeeHead(ctx, fh); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "fee_head.updated",
			EntityType: "fee_head",
			EntityID:   in.ID.String(),
			Reason:     "fee head updated",
			After:      map[string]any{"name": fh.Name, "isActive": fh.IsActive},
		})
	})
}

// --- Fee Structures & Versions ---

func (h *FeeCommandsHandler) CreateFeeStructure(ctx context.Context, in CreateFeeStructureInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermFeesEdit); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		fs := &fees.FeeStructure{
			ID:             id,
			SchoolID:       sch.ID,
			AcademicYearID: in.AcademicYearID,
			ClassID:        in.ClassID,
			CampusID:       in.CampusID,
			Code:           strings.TrimSpace(in.Code),
			Name:           strings.TrimSpace(in.Name),
			Description:    strings.TrimSpace(in.Description),
			IsActive:       true,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := r.Fees.InsertFeeStructure(ctx, fs); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "fee_structure.created",
			EntityType: "fee_structure",
			EntityID:   id.String(),
			Reason:     "fee structure created",
			After:      map[string]any{"name": fs.Name, "code": fs.Code},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *FeeCommandsHandler) CreateFeeStructureVersion(ctx context.Context, in CreateFeeStructureVersionInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermFeesEdit)
	if err != nil {
		return "", err
	}
	versionID := h.ids.NewID()
	now := h.clock.Now()

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		versions, err := r.Fees.ListFeeStructureVersions(ctx, in.FeeStructureID)
		if err != nil {
			return err
		}
		verNo := len(versions) + 1

		v := &fees.FeeStructureVersion{
			ID:             versionID,
			FeeStructureID: in.FeeStructureID,
			VersionNo:      verNo,
			Status:         fees.VersionDraft,
			Changelog:      strings.TrimSpace(in.Changelog),
			CreatedByID:    actor.UserID,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := r.Fees.InsertFeeStructureVersion(ctx, v); err != nil {
			return err
		}

		for _, item := range in.Items {
			fsi := &fees.FeeStructureItem{
				ID:                    h.ids.NewID(),
				FeeStructureVersionID: versionID,
				FeeHeadID:             item.FeeHeadID,
				AmountMinor:           item.AmountMinor,
				Frequency:             item.Frequency,
				DueDay:                item.DueDay,
				IsOptional:            item.IsOptional,
				CreatedAt:             now,
				UpdatedAt:             now,
			}
			if err := r.Fees.InsertFeeStructureItem(ctx, fsi); err != nil {
				return err
			}
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "fee_structure_version.created",
			EntityType: "fee_structure_version",
			EntityID:   versionID.String(),
			Reason:     "fee structure draft version created",
			After:      map[string]any{"versionNo": verNo, "itemCount": len(in.Items)},
		})
	})
	if err != nil {
		return "", err
	}
	return versionID, nil
}

func (h *FeeCommandsHandler) PublishFeeStructureVersion(ctx context.Context, versionID shared.ID) error {
	actor, err := authz.Require(ctx, identity.PermFeesPublish)
	if err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		v, err := r.Fees.GetFeeStructureVersion(ctx, versionID)
		if err != nil {
			return err
		}
		if v.Status == fees.VersionPublished {
			return nil
		}

		if err := r.Fees.PublishFeeStructureVersion(ctx, versionID, actor.UserID, now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "fee_structure_version.published",
			EntityType: "fee_structure_version",
			EntityID:   versionID.String(),
			Reason:     "fee structure version published and locked",
			After:      map[string]any{"versionNo": v.VersionNo},
		})
	})
}

func (h *FeeCommandsHandler) DeleteFeeStructureVersion(ctx context.Context, versionID shared.ID) error {
	_, err := authz.Require(ctx, identity.PermFeesEdit)
	if err != nil {
		return err
	}

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		v, err := r.Fees.GetFeeStructureVersion(ctx, versionID)
		if err != nil {
			return err
		}
		if v.Status != fees.VersionDraft {
			return shared.Validation("only draft fee structure versions can be deleted", map[string]string{"status": "version is not in draft status"})
		}

		if err := r.Fees.DeleteFeeStructureVersion(ctx, versionID); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "fee_structure_version.deleted",
			EntityType: "fee_structure_version",
			EntityID:   versionID.String(),
			Reason:     "fee structure draft version deleted",
			Before:     map[string]any{"versionNo": v.VersionNo, "feeStructureId": v.FeeStructureID.String()},
		})
	})
}

// --- Discount Policies & Student Discounts ---

func (h *FeeCommandsHandler) CreateDiscountPolicy(ctx context.Context, in CreateDiscountPolicyInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermFeesEdit); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		dp := &fees.DiscountPolicy{
			ID:               id,
			SchoolID:         sch.ID,
			Code:             strings.TrimSpace(in.Code),
			Name:             strings.TrimSpace(in.Name),
			Description:      strings.TrimSpace(in.Description),
			DiscountType:     in.DiscountType,
			RateBasisPoints:  in.RateBasisPoints,
			AmountMinor:      in.AmountMinor,
			FeeHeadID:        in.FeeHeadID,
			RequiresApproval: in.RequiresApproval,
			IsActive:         true,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := r.Fees.InsertDiscountPolicy(ctx, dp); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "discount_policy.created",
			EntityType: "discount_policy",
			EntityID:   id.String(),
			Reason:     "discount policy created",
			After:      map[string]any{"name": dp.Name, "code": dp.Code},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *FeeCommandsHandler) ApplyStudentDiscount(ctx context.Context, in ApplyStudentDiscountInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermFeesAssign)
	if err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		policy, err := r.Fees.GetDiscountPolicy(ctx, in.DiscountPolicyID)
		if err != nil {
			return err
		}

		status := fees.DiscountPending
		var approvedBy *shared.ID
		var approvedAt *time.Time
		if !policy.RequiresApproval {
			status = fees.DiscountApproved
			approvedBy = &actor.UserID
			approvedAt = &now
		}

		sd := &fees.StudentDiscount{
			ID:               id,
			StudentID:        in.StudentID,
			DiscountPolicyID: in.DiscountPolicyID,
			AcademicYearID:   in.AcademicYearID,
			Status:           status,
			ApprovedByID:     approvedBy,
			ApprovedAt:       approvedAt,
			ValidFrom:        in.ValidFrom,
			ValidTo:          in.ValidTo,
			CreatedByID:      actor.UserID,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := r.Fees.InsertStudentDiscount(ctx, sd); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "student_discount.applied",
			EntityType: "student_discount",
			EntityID:   id.String(),
			Reason:     "discount assigned to student",
			After:      map[string]any{"status": string(status), "studentId": in.StudentID.String()},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *FeeCommandsHandler) DecideStudentDiscount(ctx context.Context, in DecideStudentDiscountInput) error {
	actor, err := authz.Require(ctx, identity.PermFeesApprove)
	if err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		status := fees.DiscountApproved
		var appByID *shared.ID = &actor.UserID
		var appAt *time.Time = &now
		if !in.Approved {
			status = fees.DiscountRejected
			appByID = nil
			appAt = nil
		}

		if err := r.Fees.UpdateStudentDiscountStatus(ctx, in.DiscountID, status, appByID, appAt, in.RejectionReason); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "student_discount.decided",
			EntityType: "student_discount",
			EntityID:   in.DiscountID.String(),
			Reason:     "discount approval decision made",
			After:      map[string]any{"status": string(status), "reason": in.RejectionReason},
		})
	})
}

// --- Invoicing ---

func (h *FeeCommandsHandler) GenerateInvoicesBatch(ctx context.Context, in GenerateInvoicesBatchInput) ([]shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermInvoicesGenerate)
	if err != nil {
		return nil, err
	}
	now := h.clock.Now()
	var generatedIDs []shared.ID

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		// Find target students
		var studentsToInvoice []students.Student
		if in.ClassID != nil {
			// Find active students in class
			secList, err := r.Academics.ListSectionsByYear(ctx, in.AcademicYearID)
			if err != nil {
				return err
			}
			for _, sec := range secList {
				if sec.ClassID == *in.ClassID {
					roster, err := r.Enrollments.ListSectionRoster(ctx, sec.ID)
					if err != nil {
						return err
					}
					for _, enr := range roster {
						s, err := r.Students.Get(ctx, enr.StudentID)
						if err == nil && s.Status == students.StatusActive {
							studentsToInvoice = append(studentsToInvoice, *s)
						}
					}
				}
			}
		} else {
			// All active students in school
			list, err := r.Students.ListCursor(ctx, students.StudentFilter{
				SchoolID: sch.ID,
				Status:   string(students.StatusActive),
				Limit:    1000,
			})
			if err != nil {
				return err
			}
			studentsToInvoice = list
		}

		if len(studentsToInvoice) == 0 {
			return shared.Validation("no active students found for invoice generation", nil)
		}

		countInPeriod, _ := r.Fees.CountInvoicesInPeriod(ctx, sch.ID, in.BillingPeriod)

		for i, student := range studentsToInvoice {
			// Get student active enrollment
			enr, err := r.Enrollments.GetActive(ctx, student.ID)
			if err != nil || enr == nil {
				continue
			}

			// Get published fee structure for class
			feeStructure, err := r.Fees.GetFeeStructureForClass(ctx, in.AcademicYearID, enr.ClassID)
			if err != nil {
				continue
			}
			publishedVer, err := r.Fees.GetPublishedFeeStructureVersion(ctx, feeStructure.ID)
			if err != nil {
				continue
			}
			items, err := r.Fees.ListFeeStructureItems(ctx, publishedVer.ID)
			if err != nil || len(items) == 0 {
				continue
			}

			// Get approved student discounts
			studentDiscounts, _ := r.Fees.ListStudentDiscounts(ctx, student.ID, in.AcademicYearID)
			var approvedDiscounts []fees.StudentDiscount
			for _, sd := range studentDiscounts {
				if sd.Status == fees.DiscountApproved {
					approvedDiscounts = append(approvedDiscounts, sd)
				}
			}

			invoiceID := h.ids.NewID()
			invoiceSeqNo := countInPeriod + int64(i) + 1
			invoiceNo := fmt.Sprintf("INV-%d-%05d", now.Year(), invoiceSeqNo)

			var subtotalMinor int64 = 0
			var totalDiscountMinor int64 = 0
			var invoiceLineItems []fees.InvoiceItem

			for _, item := range items {
				var itemDiscountMinor int64 = 0
				for _, sd := range approvedDiscounts {
					// Check if discount policy applies to this specific fee head or general
					if sd.PolicyFeeHeadID == nil || *sd.PolicyFeeHeadID == item.FeeHeadID {
						d := fees.CalculateDiscount(item.AmountMinor, sd.DiscountType, sd.RateBasisPoints, sd.PolicyAmountMinor)
						itemDiscountMinor += d
					}
				}
				if itemDiscountMinor > item.AmountMinor {
					itemDiscountMinor = item.AmountMinor
				}

				netAmount := item.AmountMinor - itemDiscountMinor
				subtotalMinor += item.AmountMinor
				totalDiscountMinor += itemDiscountMinor

				invoiceLineItems = append(invoiceLineItems, fees.InvoiceItem{
					ID:                  h.ids.NewID(),
					InvoiceID:           invoiceID,
					FeeHeadID:           item.FeeHeadID,
					FeeHeadName:         item.FeeHeadName,
					Description:         fmt.Sprintf("%s (%s)", item.FeeHeadName, in.BillingPeriod),
					OriginalAmountMinor: item.AmountMinor,
					DiscountAmountMinor: itemDiscountMinor,
					NetAmountMinor:      netAmount,
					CreatedAt:           now,
				})
			}

			totalMinor := subtotalMinor - totalDiscountMinor
			inv := &fees.Invoice{
				ID:                  invoiceID,
				SchoolID:            sch.ID,
				StudentID:           student.ID,
				AcademicYearID:      in.AcademicYearID,
				InvoiceNo:           invoiceNo,
				IssueDate:           in.IssueDate,
				DueDate:             in.DueDate,
				BillingPeriod:       in.BillingPeriod,
				SubtotalAmountMinor: subtotalMinor,
				DiscountAmountMinor: totalDiscountMinor,
				TotalAmountMinor:    totalMinor,
				PaidAmountMinor:     0,
				BalanceMinor:        totalMinor,
				Status:              fees.InvoiceDraft,
				CreatedByID:         actor.UserID,
				CreatedAt:           now,
				UpdatedAt:           now,
			}

			if err := r.Fees.InsertInvoice(ctx, inv); err != nil {
				return fmt.Errorf("failed to insert invoice: %w", err)
			}

			for _, li := range invoiceLineItems {
				if err := r.Fees.InsertInvoiceItem(ctx, &li); err != nil {
					return fmt.Errorf("failed to insert invoice line item: %w", err)
				}
			}

			// Finalize invoice status from draft to issued
			paidAmount := int64(0)
			balance := totalMinor
			status := fees.InvoiceIssued

			// Check if student has advance credit balance to auto-allocate
			cb, err := r.Payments.GetCreditBalance(ctx, sch.ID, student.ID)
			if err == nil && cb != nil && cb.AmountMinor > 0 && totalMinor > 0 {
				creditToApply := cb.AmountMinor
				if creditToApply > totalMinor {
					creditToApply = totalMinor
				}

				paidAmount = creditToApply
				balance = totalMinor - creditToApply
				if balance == 0 {
					status = fees.InvoicePaid
				} else {
					status = fees.InvoicePartiallyPaid
				}

				// Deduct from student's credit balance
				cb.AmountMinor -= creditToApply
				cb.UpdatedAt = now
				_ = r.Payments.UpsertCreditBalance(ctx, cb)

				// Record credit ledger deduction
				invRef := invoiceID.String()
				_ = r.Payments.InsertCreditLedger(ctx, &payments.CreditLedgerEntry{
					ID:                h.ids.NewID(),
					CreditBalanceID:   cb.ID,
					StudentID:         student.ID,
					AmountMinor:       -creditToApply,
					BalanceAfterMinor: cb.AmountMinor,
					ReferenceType:     payments.CreditInvoiceApplication,
					ReferenceID:       &invRef,
					Notes:             fmt.Sprintf("Auto-applied advance credit to %s", invoiceNo),
					CreatedAt:         now,
				})

				// If there are payments with unallocated funds, create allocation records
				paymentsList, _ := r.Payments.ListPayments(ctx, payments.PaymentFilter{
					SchoolID:  sch.ID,
					StudentID: &student.ID,
					Limit:     100,
				})
				remainingToAllocate := creditToApply
				for _, p := range paymentsList {
					if p.UnallocatedAmountMinor > 0 && remainingToAllocate > 0 {
						thisAlloc := p.UnallocatedAmountMinor
						if thisAlloc > remainingToAllocate {
							thisAlloc = remainingToAllocate
						}
						_ = r.Payments.InsertAllocation(ctx, &payments.PaymentAllocation{
							ID:                   h.ids.NewID(),
							PaymentID:            p.ID,
							InvoiceID:            invoiceID,
							AllocatedAmountMinor: thisAlloc,
							AllocatedAt:          now,
							AllocatedByID:        actor.UserID,
						})
						_ = r.Payments.UpdatePaymentAllocation(ctx, p.ID, p.AllocatedAmountMinor+thisAlloc, p.UnallocatedAmountMinor-thisAlloc)
						remainingToAllocate -= thisAlloc
					}
				}
			}

			if err := r.Fees.UpdateInvoiceStatus(ctx, invoiceID, status, paidAmount, balance); err != nil {
				return fmt.Errorf("failed to issue invoice: %w", err)
			}

			generatedIDs = append(generatedIDs, invoiceID)
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "invoices.batch_generated",
			EntityType: "invoice_batch",
			EntityID:   in.BillingPeriod,
			Reason:     "monthly invoice run",
			After:      map[string]any{"count": len(generatedIDs), "billingPeriod": in.BillingPeriod},
		})
	})
	if err != nil {
		return nil, err
	}
	return generatedIDs, nil
}

func (h *FeeCommandsHandler) VoidInvoice(ctx context.Context, in VoidInvoiceInput) error {
	actor, err := authz.Require(ctx, identity.PermInvoicesVoid)
	if err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		inv, err := r.Fees.GetInvoice(ctx, in.InvoiceID)
		if err != nil {
			return err
		}
		if inv.Status == fees.InvoicePaid {
			return shared.Validation("cannot void an invoice that has been paid", nil)
		}

		if err := r.Fees.VoidInvoice(ctx, in.InvoiceID, in.Reason, actor.UserID, now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "invoice.voided",
			EntityType: "invoice",
			EntityID:   in.InvoiceID.String(),
			Reason:     in.Reason,
			Before:     map[string]any{"status": string(inv.Status)},
			After:      map[string]any{"status": string(fees.InvoiceVoid)},
		})
	})
}
