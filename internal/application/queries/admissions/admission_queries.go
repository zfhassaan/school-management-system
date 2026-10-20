package admissions

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/admissions"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type AdmissionQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewAdmissionQueriesHandler(uow uow.UnitOfWork) *AdmissionQueriesHandler {
	return &AdmissionQueriesHandler{uow: uow}
}

func (h *AdmissionQueriesHandler) ListFormDefinitions(ctx context.Context) ([]admissions.FormDefinition, error) {
	if _, err := authz.Require(ctx, identity.PermAdmissionsView); err != nil {
		return nil, err
	}
	var res []admissions.FormDefinition
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		defs, err := r.Admissions.ListFormDefinitions(ctx, sch.ID)
		if err != nil {
			return err
		}
		res = defs
		return nil
	})
	return res, err
}

func (h *AdmissionQueriesHandler) GetPublishedFormVersion(ctx context.Context, formDefID shared.ID) (*admissions.FormVersion, error) {
	if _, err := authz.Require(ctx, identity.PermAdmissionsView); err != nil {
		return nil, err
	}
	var res *admissions.FormVersion
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		v, err := r.Admissions.GetPublishedVersion(ctx, formDefID)
		if err != nil {
			return err
		}
		res = v
		return nil
	})
	return res, err
}

func (h *AdmissionQueriesHandler) ListFormVersions(ctx context.Context, formDefID shared.ID) ([]admissions.FormVersion, error) {
	if _, err := authz.Require(ctx, identity.PermAdmissionsView); err != nil {
		return nil, err
	}
	var res []admissions.FormVersion
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		versions, err := r.Admissions.ListFormVersions(ctx, formDefID)
		if err != nil {
			return err
		}
		res = versions
		return nil
	})
	return res, err
}

func (h *AdmissionQueriesHandler) ListApplications(ctx context.Context, status string, yearID string, limit int) ([]admissions.Application, error) {
	if _, err := authz.Require(ctx, identity.PermAdmissionsView); err != nil {
		return nil, err
	}
	var res []admissions.Application
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		apps, err := r.Admissions.ListApplications(ctx, sch.ID, status, yearID, limit)
		if err != nil {
			return err
		}
		res = apps
		return nil
	})
	return res, err
}

func (h *AdmissionQueriesHandler) GetApplication(ctx context.Context, id shared.ID) (*admissions.Application, error) {
	if _, err := authz.Require(ctx, identity.PermAdmissionsView); err != nil {
		return nil, err
	}
	var res *admissions.Application
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		app, err := r.Admissions.GetApplication(ctx, id)
		if err != nil {
			return err
		}
		res = app
		return nil
	})
	return res, err
}
