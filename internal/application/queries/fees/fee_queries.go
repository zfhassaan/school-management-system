package fees

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/fees"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type FeeQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewFeeQueriesHandler(uow uow.UnitOfWork) *FeeQueriesHandler {
	return &FeeQueriesHandler{uow: uow}
}

// --- Fee Heads ---

func (h *FeeQueriesHandler) ListFeeHeads(ctx context.Context) ([]fees.FeeHead, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res []fees.FeeHead
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		list, err := r.Fees.ListFeeHeads(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *FeeQueriesHandler) GetFeeHead(ctx context.Context, id shared.ID) (*fees.FeeHead, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res *fees.FeeHead
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		fh, err := r.Fees.GetFeeHead(ctx, id)
		if err != nil {
			return err
		}
		res = fh
		return nil
	})
	return res, err
}

// --- Fee Structures & Versions ---

func (h *FeeQueriesHandler) ListFeeStructures(ctx context.Context, academicYearID shared.ID) ([]fees.FeeStructure, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res []fees.FeeStructure
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		list, err := r.Fees.ListFeeStructures(ctx, sch.ID, academicYearID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *FeeQueriesHandler) GetFeeStructure(ctx context.Context, id shared.ID) (*fees.FeeStructure, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res *fees.FeeStructure
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		fs, err := r.Fees.GetFeeStructure(ctx, id)
		if err != nil {
			return err
		}
		res = fs
		return nil
	})
	return res, err
}

type FeeStructureVersionDetail struct {
	Version fees.FeeStructureVersion `json:"version"`
	Items   []fees.FeeStructureItem  `json:"items"`
}

func (h *FeeQueriesHandler) GetFeeStructureVersionDetail(ctx context.Context, versionID shared.ID) (*FeeStructureVersionDetail, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res *FeeStructureVersionDetail
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		ver, err := r.Fees.GetFeeStructureVersion(ctx, versionID)
		if err != nil {
			return err
		}
		items, err := r.Fees.ListFeeStructureItems(ctx, versionID)
		if err != nil {
			return err
		}
		res = &FeeStructureVersionDetail{
			Version: *ver,
			Items:   items,
		}
		return nil
	})
	return res, err
}

func (h *FeeQueriesHandler) ListFeeStructureVersions(ctx context.Context, feeStructureID shared.ID) ([]fees.FeeStructureVersion, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res []fees.FeeStructureVersion
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		list, err := r.Fees.ListFeeStructureVersions(ctx, feeStructureID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

// --- Discount Policies ---

func (h *FeeQueriesHandler) ListDiscountPolicies(ctx context.Context) ([]fees.DiscountPolicy, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res []fees.DiscountPolicy
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		list, err := r.Fees.ListDiscountPolicies(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *FeeQueriesHandler) ListStudentDiscounts(ctx context.Context, studentID shared.ID, academicYearID shared.ID) ([]fees.StudentDiscount, error) {
	if _, err := authz.Require(ctx, identity.PermFeesView); err != nil {
		return nil, err
	}
	var res []fees.StudentDiscount
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		list, err := r.Fees.ListStudentDiscounts(ctx, studentID, academicYearID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *FeeQueriesHandler) ListPendingDiscounts(ctx context.Context) ([]fees.StudentDiscount, error) {
	if _, err := authz.Require(ctx, identity.PermFeesApprove); err != nil {
		return nil, err
	}
	var res []fees.StudentDiscount
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		list, err := r.Fees.ListPendingStudentDiscounts(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

// --- Invoices ---

type InvoiceDetail struct {
	Invoice fees.Invoice       `json:"invoice"`
	Items   []fees.InvoiceItem `json:"items"`
}

func (h *FeeQueriesHandler) ListInvoices(ctx context.Context, filter fees.InvoiceFilter) ([]fees.Invoice, error) {
	if _, err := authz.Require(ctx, identity.PermInvoicesView); err != nil {
		return nil, err
	}
	var res []fees.Invoice
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		filter.SchoolID = sch.ID
		list, err := r.Fees.ListInvoices(ctx, filter)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *FeeQueriesHandler) GetInvoice(ctx context.Context, id shared.ID) (*InvoiceDetail, error) {
	if _, err := authz.Require(ctx, identity.PermInvoicesView); err != nil {
		return nil, err
	}
	var res *InvoiceDetail
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		inv, err := r.Fees.GetInvoice(ctx, id)
		if err != nil {
			return err
		}
		items, err := r.Fees.ListInvoiceItems(ctx, id)
		if err != nil {
			return err
		}
		res = &InvoiceDetail{
			Invoice: *inv,
			Items:   items,
		}
		return nil
	})
	return res, err
}
