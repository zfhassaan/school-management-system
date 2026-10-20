package school

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/school"
)

type SchoolQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewSchoolQueriesHandler(uow uow.UnitOfWork) *SchoolQueriesHandler {
	return &SchoolQueriesHandler{uow: uow}
}

func (h *SchoolQueriesHandler) GetSchool(ctx context.Context) (*school.School, error) {
	if _, err := authz.Require(ctx, identity.PermSchoolView); err != nil {
		return nil, err
	}
	var res *school.School
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		s, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		res = s
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *SchoolQueriesHandler) GetSettings(ctx context.Context) (*school.Settings, error) {
	if _, err := authz.Require(ctx, identity.PermSchoolView); err != nil {
		return nil, err
	}
	var res *school.Settings
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		s, err := r.School.GetSettings(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = s
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

