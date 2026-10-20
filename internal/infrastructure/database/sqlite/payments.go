package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/fees"
	"school-management/internal/domain/payments"
	"school-management/internal/domain/shared"
)

type PaymentsRepository struct {
	q Querier
}

func (r *PaymentsRepository) NextReceiptNumber(ctx context.Context, schoolID shared.ID, year int) (string, error) {
	val, err := r.q.GetAndIncrementReceiptSeq(ctx, GetAndIncrementReceiptSeqParams{
		SchoolID: schoolID.String(),
		Year:     int64(year),
	})
	if err != nil {
		return "", shared.Internal(err)
	}
	return payments.FormatReceiptNo(year, val), nil
}

func (r *PaymentsRepository) InsertPayment(ctx context.Context, p *payments.Payment) error {
	err := r.q.InsertPayment(ctx, InsertPaymentParams{
		ID:                     p.ID.String(),
		SchoolID:               p.SchoolID.String(),
		StudentID:              p.StudentID.String(),
		ReceiptNo:              p.ReceiptNo,
		PaymentDate:            p.PaymentDate.String(),
		PaymentMethod:          string(p.PaymentMethod),
		AmountMinor:            p.AmountMinor,
		AllocatedAmountMinor:   p.AllocatedAmountMinor,
		UnallocatedAmountMinor: p.UnallocatedAmountMinor,
		ReferenceNo:            p.ReferenceNo,
		PayerName:              p.PayerName,
		PayerPhone:             p.PayerPhone,
		Status:                 string(p.Status),
		RecordedByID:           p.RecordedByID.String(),
		RecordedAt:             p.RecordedAt.Format(time.RFC3339Nano),
		Notes:                  p.Notes,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) UpdatePaymentAllocation(ctx context.Context, paymentID shared.ID, allocatedMinor, unallocatedMinor int64) error {
	err := r.q.UpdatePaymentAllocation(ctx, UpdatePaymentAllocationParams{
		AllocatedAmountMinor:   allocatedMinor,
		UnallocatedAmountMinor: unallocatedMinor,
		ID:                     paymentID.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) UpdatePaymentStatus(ctx context.Context, paymentID shared.ID, status payments.PaymentStatus) error {
	err := r.q.UpdatePaymentStatus(ctx, UpdatePaymentStatusParams{
		Status: string(status),
		ID:     paymentID.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) GetPaymentByID(ctx context.Context, id shared.ID) (*payments.PaymentWithDetails, error) {
	row, err := r.q.GetPaymentByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("payment")
		}
		return nil, shared.Internal(err)
	}
	pDate, _ := shared.ParseDate(row.PaymentDate)
	rAt, _ := time.Parse(time.RFC3339Nano, row.RecordedAt)
	return &payments.PaymentWithDetails{
		Payment: payments.Payment{
			ID:                     shared.ID(row.ID),
			SchoolID:               shared.ID(row.SchoolID),
			StudentID:              shared.ID(row.StudentID),
			ReceiptNo:              row.ReceiptNo,
			PaymentDate:            pDate,
			PaymentMethod:          payments.PaymentMethod(row.PaymentMethod),
			AmountMinor:            row.AmountMinor,
			AllocatedAmountMinor:   row.AllocatedAmountMinor,
			UnallocatedAmountMinor: row.UnallocatedAmountMinor,
			ReferenceNo:            row.ReferenceNo,
			PayerName:              row.PayerName,
			PayerPhone:             row.PayerPhone,
			Status:                 payments.PaymentStatus(row.Status),
			RecordedByID:           shared.ID(row.RecordedByID),
			RecordedAt:             rAt,
			Notes:                  row.Notes,
		},
		StudentFirstName: row.StudentFirstName,
		StudentLastName:  row.StudentLastName,
		AdmissionNo:      row.AdmissionNo,
		RecordedByName:   row.RecordedByName,
	}, nil
}

func (r *PaymentsRepository) ListPayments(ctx context.Context, filter payments.PaymentFilter) ([]*payments.PaymentWithDetails, error) {
	studentIDStr := ""
	if filter.StudentID != nil {
		studentIDStr = filter.StudentID.String()
	}
	statusStr := ""
	if filter.Status != nil {
		statusStr = string(*filter.Status)
	}
	methodStr := ""
	if filter.PaymentMethod != nil {
		methodStr = string(*filter.PaymentMethod)
	}
	limit := int64(100)
	if filter.Limit > 0 {
		limit = int64(filter.Limit)
	}

	rows, err := r.q.ListPayments(ctx, ListPaymentsParams{
		SchoolID:      filter.SchoolID.String(),
		StudentID:     studentIDStr,
		Status:        statusStr,
		PaymentMethod: methodStr,
		Limit:         limit,
	})
	if err != nil {
		return nil, shared.Internal(err)
	}

	res := make([]*payments.PaymentWithDetails, len(rows))
	for i, row := range rows {
		pDate, _ := shared.ParseDate(row.PaymentDate)
		rAt, _ := time.Parse(time.RFC3339Nano, row.RecordedAt)
		res[i] = &payments.PaymentWithDetails{
			Payment: payments.Payment{
				ID:                     shared.ID(row.ID),
				SchoolID:               shared.ID(row.SchoolID),
				StudentID:              shared.ID(row.StudentID),
				ReceiptNo:              row.ReceiptNo,
				PaymentDate:            pDate,
				PaymentMethod:          payments.PaymentMethod(row.PaymentMethod),
				AmountMinor:            row.AmountMinor,
				AllocatedAmountMinor:   row.AllocatedAmountMinor,
				UnallocatedAmountMinor: row.UnallocatedAmountMinor,
				ReferenceNo:            row.ReferenceNo,
				PayerName:              row.PayerName,
				PayerPhone:             row.PayerPhone,
				Status:                 payments.PaymentStatus(row.Status),
				RecordedByID:           shared.ID(row.RecordedByID),
				RecordedAt:             rAt,
				Notes:                  row.Notes,
			},
			StudentFirstName: row.StudentFirstName,
			StudentLastName:  row.StudentLastName,
			AdmissionNo:      row.AdmissionNo,
			RecordedByName:   row.RecordedByName,
		}
	}
	return res, nil
}

func (r *PaymentsRepository) InsertAllocation(ctx context.Context, a *payments.PaymentAllocation) error {
	var itemIDStr sql.NullString
	if a.InvoiceItemID != nil {
		itemIDStr = sql.NullString{String: a.InvoiceItemID.String(), Valid: true}
	}
	err := r.q.InsertPaymentAllocation(ctx, InsertPaymentAllocationParams{
		ID:                   a.ID.String(),
		PaymentID:            a.PaymentID.String(),
		InvoiceID:            a.InvoiceID.String(),
		InvoiceItemID:        itemIDStr,
		AllocatedAmountMinor: a.AllocatedAmountMinor,
		AllocatedAt:          a.AllocatedAt.Format(time.RFC3339Nano),
		AllocatedByID:        a.AllocatedByID.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) ListAllocationsByPayment(ctx context.Context, paymentID shared.ID) ([]*payments.PaymentAllocation, error) {
	rows, err := r.q.ListAllocationsByPayment(ctx, paymentID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	res := make([]*payments.PaymentAllocation, len(rows))
	for i, row := range rows {
		aAt, _ := time.Parse(time.RFC3339Nano, row.AllocatedAt)
		var itemID *shared.ID
		if row.InvoiceItemID.Valid {
			id := shared.ID(row.InvoiceItemID.String)
			itemID = &id
		}
		res[i] = &payments.PaymentAllocation{
			ID:                   shared.ID(row.ID),
			PaymentID:            shared.ID(row.PaymentID),
			InvoiceID:            shared.ID(row.InvoiceID),
			InvoiceItemID:        itemID,
			AllocatedAmountMinor: row.AllocatedAmountMinor,
			AllocatedAt:          aAt,
			AllocatedByID:        shared.ID(row.AllocatedByID),
			InvoiceNo:            row.InvoiceNo,
			BillingPeriod:        row.BillingPeriod,
		}
	}
	return res, nil
}

func (r *PaymentsRepository) ListAllocationsByInvoice(ctx context.Context, invoiceID shared.ID) ([]*payments.PaymentAllocation, error) {
	rows, err := r.q.ListAllocationsByInvoice(ctx, invoiceID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	res := make([]*payments.PaymentAllocation, len(rows))
	for i, row := range rows {
		aAt, _ := time.Parse(time.RFC3339Nano, row.AllocatedAt)
		pDate, _ := shared.ParseDate(row.PaymentDate)
		var itemID *shared.ID
		if row.InvoiceItemID.Valid {
			id := shared.ID(row.InvoiceItemID.String)
			itemID = &id
		}
		res[i] = &payments.PaymentAllocation{
			ID:                   shared.ID(row.ID),
			PaymentID:            shared.ID(row.PaymentID),
			InvoiceID:            shared.ID(row.InvoiceID),
			InvoiceItemID:        itemID,
			AllocatedAmountMinor: row.AllocatedAmountMinor,
			AllocatedAt:          aAt,
			AllocatedByID:        shared.ID(row.AllocatedByID),
			ReceiptNo:            row.ReceiptNo,
			PaymentDate:          pDate,
			PaymentMethod:        payments.PaymentMethod(row.PaymentMethod),
		}
	}
	return res, nil
}

func (r *PaymentsRepository) InsertReversal(ctx context.Context, rv *payments.PaymentReversal) error {
	err := r.q.InsertPaymentReversal(ctx, InsertPaymentReversalParams{
		ID:           rv.ID.String(),
		PaymentID:    rv.PaymentID.String(),
		ReversalDate: rv.ReversalDate.String(),
		Reason:       rv.Reason,
		ReversedByID: rv.ReversedByID.String(),
		ReversedAt:   rv.ReversedAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) GetReversal(ctx context.Context, paymentID shared.ID) (*payments.PaymentReversal, error) {
	row, err := r.q.GetPaymentReversal(ctx, paymentID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("payment reversal")
		}
		return nil, shared.Internal(err)
	}
	rDate, _ := shared.ParseDate(row.ReversalDate)
	rAt, _ := time.Parse(time.RFC3339Nano, row.ReversedAt)
	return &payments.PaymentReversal{
		ID:             shared.ID(row.ID),
		PaymentID:      shared.ID(row.PaymentID),
		ReversalDate:   rDate,
		Reason:         row.Reason,
		ReversedByID:   shared.ID(row.ReversedByID),
		ReversedAt:     rAt,
		ReversedByName: row.ReversedByName,
	}, nil
}

func (r *PaymentsRepository) GetCreditBalance(ctx context.Context, schoolID, studentID shared.ID) (*payments.CreditBalance, error) {
	row, err := r.q.GetCreditBalanceByStudent(ctx, GetCreditBalanceByStudentParams{
		SchoolID:  schoolID.String(),
		StudentID: studentID.String(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("credit balance")
		}
		return nil, shared.Internal(err)
	}
	uAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	cAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	return &payments.CreditBalance{
		ID:          shared.ID(row.ID),
		SchoolID:    shared.ID(row.SchoolID),
		StudentID:   shared.ID(row.StudentID),
		AmountMinor: row.AmountMinor,
		UpdatedAt:   uAt,
		CreatedAt:   cAt,
	}, nil
}

func (r *PaymentsRepository) UpsertCreditBalance(ctx context.Context, cb *payments.CreditBalance) error {
	err := r.q.UpsertCreditBalance(ctx, UpsertCreditBalanceParams{
		ID:          cb.ID.String(),
		SchoolID:    cb.SchoolID.String(),
		StudentID:   cb.StudentID.String(),
		AmountMinor: cb.AmountMinor,
		UpdatedAt:   cb.UpdatedAt.Format(time.RFC3339Nano),
		CreatedAt:   cb.CreatedAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) InsertCreditLedger(ctx context.Context, entry *payments.CreditLedgerEntry) error {
	var refID sql.NullString
	if entry.ReferenceID != nil {
		refID = sql.NullString{String: *entry.ReferenceID, Valid: true}
	}
	err := r.q.InsertCreditLedgerEntry(ctx, InsertCreditLedgerEntryParams{
		ID:                entry.ID.String(),
		CreditBalanceID:   entry.CreditBalanceID.String(),
		StudentID:         entry.StudentID.String(),
		AmountMinor:       entry.AmountMinor,
		BalanceAfterMinor: entry.BalanceAfterMinor,
		ReferenceType:     string(entry.ReferenceType),
		ReferenceID:       refID,
		Notes:             entry.Notes,
		CreatedAt:         entry.CreatedAt.Format(time.RFC3339Nano),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) ListCreditLedger(ctx context.Context, studentID shared.ID) ([]*payments.CreditLedgerEntry, error) {
	rows, err := r.q.ListCreditLedger(ctx, studentID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	res := make([]*payments.CreditLedgerEntry, len(rows))
	for i, row := range rows {
		cAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		var refID *string
		if row.ReferenceID.Valid {
			s := row.ReferenceID.String
			refID = &s
		}
		res[i] = &payments.CreditLedgerEntry{
			ID:                shared.ID(row.ID),
			CreditBalanceID:   shared.ID(row.CreditBalanceID),
			StudentID:         shared.ID(row.StudentID),
			AmountMinor:       row.AmountMinor,
			BalanceAfterMinor: row.BalanceAfterMinor,
			ReferenceType:     payments.CreditReferenceType(row.ReferenceType),
			ReferenceID:       refID,
			Notes:             row.Notes,
			CreatedAt:         cAt,
		}
	}
	return res, nil
}

func (r *PaymentsRepository) InsertAdjustment(ctx context.Context, fa *payments.FinancialAdjustment) error {
	var invID sql.NullString
	if fa.InvoiceID != nil {
		invID = sql.NullString{String: fa.InvoiceID.String(), Valid: true}
	}
	var itemID sql.NullString
	if fa.InvoiceItemID != nil {
		itemID = sql.NullString{String: fa.InvoiceItemID.String(), Valid: true}
	}
	err := r.q.InsertFinancialAdjustment(ctx, InsertFinancialAdjustmentParams{
		ID:             fa.ID.String(),
		SchoolID:       fa.SchoolID.String(),
		StudentID:      fa.StudentID.String(),
		InvoiceID:      invID,
		InvoiceItemID:  itemID,
		AdjustmentType: string(fa.AdjustmentType),
		AmountMinor:    fa.AmountMinor,
		Reason:         fa.Reason,
		Status:         string(fa.Status),
		RequestedByID:  fa.RequestedByID.String(),
		RequestedAt:    fa.RequestedAt.Format(time.RFC3339Nano),
		DecisionNotes:  fa.DecisionNotes,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) GetAdjustmentByID(ctx context.Context, id shared.ID) (*payments.FinancialAdjustment, error) {
	row, err := r.q.GetFinancialAdjustmentByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("financial adjustment")
		}
		return nil, shared.Internal(err)
	}

	rAt, _ := time.Parse(time.RFC3339Nano, row.RequestedAt)
	var dAt *time.Time
	if row.DecidedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.DecidedAt.String)
		dAt = &t
	}
	var decID *shared.ID
	if row.DecidedByID.Valid {
		id := shared.ID(row.DecidedByID.String)
		decID = &id
	}
	var invID *shared.ID
	if row.InvoiceID.Valid {
		id := shared.ID(row.InvoiceID.String)
		invID = &id
	}
	var itemID *shared.ID
	if row.InvoiceItemID.Valid {
		id := shared.ID(row.InvoiceItemID.String)
		itemID = &id
	}
	var decidedByName *string
	if row.DecidedByName.Valid {
		decidedByName = &row.DecidedByName.String
	}
	var invoiceNo *string
	if row.InvoiceNo.Valid {
		invoiceNo = &row.InvoiceNo.String
	}

	return &payments.FinancialAdjustment{
		ID:               shared.ID(row.ID),
		SchoolID:         shared.ID(row.SchoolID),
		StudentID:        shared.ID(row.StudentID),
		InvoiceID:        invID,
		InvoiceItemID:    itemID,
		AdjustmentType:   payments.AdjustmentType(row.AdjustmentType),
		AmountMinor:      row.AmountMinor,
		Reason:           row.Reason,
		Status:           payments.AdjustmentStatus(row.Status),
		RequestedByID:    shared.ID(row.RequestedByID),
		RequestedAt:      rAt,
		DecidedByID:      decID,
		DecidedAt:        dAt,
		DecisionNotes:    row.DecisionNotes,
		StudentFirstName: row.StudentFirstName,
		StudentLastName:  row.StudentLastName,
		AdmissionNo:      row.AdmissionNo,
		RequestedByName:  row.RequestedByName,
		DecidedByName:    decidedByName,
		InvoiceNo:        invoiceNo,
	}, nil
}

func (r *PaymentsRepository) UpdateAdjustmentStatus(ctx context.Context, id shared.ID, status payments.AdjustmentStatus, decidedByID shared.ID, decidedAt time.Time, notes string) error {
	err := r.q.UpdateFinancialAdjustmentStatus(ctx, UpdateFinancialAdjustmentStatusParams{
		Status:        string(status),
		DecidedByID:   sql.NullString{String: decidedByID.String(), Valid: true},
		DecidedAt:     sql.NullString{String: decidedAt.Format(time.RFC3339Nano), Valid: true},
		DecisionNotes: notes,
		ID:            id.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *PaymentsRepository) ListAdjustments(ctx context.Context, filter payments.AdjustmentFilter) ([]*payments.FinancialAdjustment, error) {
	studentIDStr := ""
	if filter.StudentID != nil {
		studentIDStr = filter.StudentID.String()
	}
	statusStr := ""
	if filter.Status != nil {
		statusStr = string(*filter.Status)
	}
	limit := int64(100)
	if filter.Limit > 0 {
		limit = int64(filter.Limit)
	}

	rows, err := r.q.ListFinancialAdjustments(ctx, ListFinancialAdjustmentsParams{
		SchoolID:  filter.SchoolID.String(),
		StudentID: studentIDStr,
		Status:    statusStr,
		Limit:     limit,
	})
	if err != nil {
		return nil, shared.Internal(err)
	}

	res := make([]*payments.FinancialAdjustment, len(rows))
	for i, row := range rows {
		rAt, _ := time.Parse(time.RFC3339Nano, row.RequestedAt)
		var dAt *time.Time
		if row.DecidedAt.Valid {
			t, _ := time.Parse(time.RFC3339Nano, row.DecidedAt.String)
			dAt = &t
		}
		var decID *shared.ID
		if row.DecidedByID.Valid {
			id := shared.ID(row.DecidedByID.String)
			decID = &id
		}
		var invID *shared.ID
		if row.InvoiceID.Valid {
			id := shared.ID(row.InvoiceID.String)
			invID = &id
		}
		var itemID *shared.ID
		if row.InvoiceItemID.Valid {
			id := shared.ID(row.InvoiceItemID.String)
			itemID = &id
		}
		var decidedByName *string
		if row.DecidedByName.Valid {
			decidedByName = &row.DecidedByName.String
		}
		var invoiceNo *string
		if row.InvoiceNo.Valid {
			invoiceNo = &row.InvoiceNo.String
		}

		res[i] = &payments.FinancialAdjustment{
			ID:               shared.ID(row.ID),
			SchoolID:         shared.ID(row.SchoolID),
			StudentID:        shared.ID(row.StudentID),
			InvoiceID:        invID,
			InvoiceItemID:    itemID,
			AdjustmentType:   payments.AdjustmentType(row.AdjustmentType),
			AmountMinor:      row.AmountMinor,
			Reason:           row.Reason,
			Status:           payments.AdjustmentStatus(row.Status),
			RequestedByID:    shared.ID(row.RequestedByID),
			RequestedAt:      rAt,
			DecidedByID:      decID,
			DecidedAt:        dAt,
			DecisionNotes:    row.DecisionNotes,
			StudentFirstName: row.StudentFirstName,
			StudentLastName:  row.StudentLastName,
			AdmissionNo:      row.AdmissionNo,
			RequestedByName:  row.RequestedByName,
			DecidedByName:    decidedByName,
			InvoiceNo:        invoiceNo,
		}
	}
	return res, nil
}

func (r *PaymentsRepository) GetUnpaidInvoicesForStudent(ctx context.Context, schoolID, studentID shared.ID) ([]*fees.Invoice, error) {
	rows, err := r.q.GetUnpaidInvoicesForStudent(ctx, GetUnpaidInvoicesForStudentParams{
		SchoolID:  schoolID.String(),
		StudentID: studentID.String(),
	})
	if err != nil {
		return nil, shared.Internal(err)
	}

	res := make([]*fees.Invoice, len(rows))
	for i, row := range rows {
		issueDate, _ := shared.ParseDate(row.IssueDate)
		dueDate, _ := shared.ParseDate(row.DueDate)
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		var voidedAt *time.Time
		if row.VoidedAt.Valid {
			t, _ := time.Parse(time.RFC3339Nano, row.VoidedAt.String)
			voidedAt = &t
		}

		res[i] = &fees.Invoice{
			ID:                  shared.ID(row.ID),
			SchoolID:            shared.ID(row.SchoolID),
			StudentID:           shared.ID(row.StudentID),
			AcademicYearID:      shared.ID(row.AcademicYearID),
			InvoiceNo:           row.InvoiceNo,
			IssueDate:           issueDate,
			DueDate:             dueDate,
			BillingPeriod:       row.BillingPeriod,
			SubtotalAmountMinor: row.SubtotalAmountMinor,
			DiscountAmountMinor: row.DiscountAmountMinor,
			TotalAmountMinor:    row.TotalAmountMinor,
			PaidAmountMinor:     row.PaidAmountMinor,
			BalanceMinor:        row.BalanceMinor,
			Status:              fees.InvoiceStatus(row.Status),
			VoidReason:          row.VoidReason,
			CreatedByID:         shared.ID(row.CreatedBy),
			CreatedAt:           createdAt,
			UpdatedAt:           updatedAt,
			VoidedAt:            voidedAt,
			Revision:            row.RowVersion,
		}
	}
	return res, nil
}
