package payments

import (
	"context"
	"errors"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/payments"
	"school-management/internal/domain/shared"
)

type PaymentQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewPaymentQueriesHandler(u uow.UnitOfWork) *PaymentQueriesHandler {
	return &PaymentQueriesHandler{uow: u}
}

type PaymentDetail struct {
	Payment     payments.PaymentWithDetails   `json:"payment"`
	Allocations []*payments.PaymentAllocation `json:"allocations"`
	Reversal    *payments.PaymentReversal     `json:"reversal,omitempty"`
}

type StudentCreditBalanceInfo struct {
	BalanceMinor int64                         `json:"balanceMinor"`
	Ledger       []*payments.CreditLedgerEntry `json:"ledger"`
}

func (h *PaymentQueriesHandler) GetPayment(ctx context.Context, id shared.ID) (*PaymentDetail, error) {
	if _, err := authz.Require(ctx, identity.PermPaymentsView); err != nil {
		return nil, err
	}

	var detail PaymentDetail
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		p, err := r.Payments.GetPaymentByID(ctx, id)
		if err != nil {
			return err
		}
		detail.Payment = *p

		allocations, err := r.Payments.ListAllocationsByPayment(ctx, id)
		if err != nil {
			return err
		}
		detail.Allocations = allocations

		if p.Status == payments.PaymentReversed {
			rev, err := r.Payments.GetReversal(ctx, id)
			if err == nil {
				detail.Reversal = rev
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (h *PaymentQueriesHandler) ListPayments(ctx context.Context, filter payments.PaymentFilter) ([]*payments.PaymentWithDetails, error) {
	if _, err := authz.Require(ctx, identity.PermPaymentsView); err != nil {
		return nil, err
	}

	var list []*payments.PaymentWithDetails
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		filter.SchoolID = sch.ID

		res, err := r.Payments.ListPayments(ctx, filter)
		if err != nil {
			return err
		}
		list = res
		return nil
	})

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (h *PaymentQueriesHandler) GetStudentCreditBalance(ctx context.Context, studentID shared.ID) (*StudentCreditBalanceInfo, error) {
	if _, err := authz.Require(ctx, identity.PermPaymentsView); err != nil {
		return nil, err
	}

	var info StudentCreditBalanceInfo
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		cb, err := r.Payments.GetCreditBalance(ctx, sch.ID, studentID)
		if err != nil {
			if !errors.Is(err, shared.ErrNotFound) {
				return err
			}
			info.BalanceMinor = 0
		} else {
			info.BalanceMinor = cb.AmountMinor
		}

		ledger, err := r.Payments.ListCreditLedger(ctx, studentID)
		if err != nil {
			return err
		}
		info.Ledger = ledger
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (h *PaymentQueriesHandler) ListFinancialAdjustments(ctx context.Context, filter payments.AdjustmentFilter) ([]*payments.FinancialAdjustment, error) {
	if _, err := authz.Require(ctx, identity.PermAdjustmentsRequest); err != nil {
		if _, err2 := authz.Require(ctx, identity.PermAdjustmentsApprove); err2 != nil {
			return nil, err
		}
	}

	var list []*payments.FinancialAdjustment
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		filter.SchoolID = sch.ID

		res, err := r.Payments.ListAdjustments(ctx, filter)
		if err != nil {
			return err
		}
		list = res
		return nil
	})

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (h *PaymentQueriesHandler) GetFinancialAdjustment(ctx context.Context, id shared.ID) (*payments.FinancialAdjustment, error) {
	if _, err := authz.Require(ctx, identity.PermAdjustmentsRequest); err != nil {
		if _, err2 := authz.Require(ctx, identity.PermAdjustmentsApprove); err2 != nil {
			return nil, err
		}
	}

	var adj *payments.FinancialAdjustment
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		res, err := r.Payments.GetAdjustmentByID(ctx, id)
		if err != nil {
			return err
		}
		adj = res
		return nil
	})

	if err != nil {
		return nil, err
	}
	return adj, nil
}

