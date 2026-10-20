package academics

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type AcademicQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewAcademicQueriesHandler(uow uow.UnitOfWork) *AcademicQueriesHandler {
	return &AcademicQueriesHandler{uow: uow}
}

func (h *AcademicQueriesHandler) ListAcademicYears(ctx context.Context) ([]academics.AcademicYear, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsView); err != nil {
		return nil, err
	}
	var res []academics.AcademicYear
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		years, err := r.Academics.ListAcademicYears(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = years
		return nil
	})
	return res, err
}

func (h *AcademicQueriesHandler) GetActiveAcademicYear(ctx context.Context) (*academics.AcademicYear, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsView); err != nil {
		return nil, err
	}
	var res *academics.AcademicYear
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		ay, err := r.Academics.GetActiveAcademicYear(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = ay
		return nil
	})
	return res, err
}

func (h *AcademicQueriesHandler) ListCampuses(ctx context.Context) ([]academics.Campus, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsView); err != nil {
		return nil, err
	}
	var res []academics.Campus
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		campuses, err := r.Academics.ListCampuses(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = campuses
		return nil
	})
	return res, err
}

func (h *AcademicQueriesHandler) ListClasses(ctx context.Context) ([]academics.Class, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsView); err != nil {
		return nil, err
	}
	var res []academics.Class
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		classes, err := r.Academics.ListClasses(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = classes
		return nil
	})
	return res, err
}

func (h *AcademicQueriesHandler) ListSectionsByYear(ctx context.Context, yearID shared.ID) ([]academics.Section, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsView); err != nil {
		return nil, err
	}
	var res []academics.Section
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sections, err := r.Academics.ListSectionsByYear(ctx, yearID)
		if err != nil {
			return err
		}
		res = sections
		return nil
	})
	return res, err
}

func (h *AcademicQueriesHandler) ListTerms(ctx context.Context, yearID shared.ID) ([]academics.AcademicTerm, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsView); err != nil {
		return nil, err
	}
	var res []academics.AcademicTerm
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		terms, err := r.Academics.ListAcademicTerms(ctx, yearID)
		if err != nil {
			return err
		}
		res = terms
		return nil
	})
	return res, err
}

