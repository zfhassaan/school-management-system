package guardians

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/guardians"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type GuardianQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewGuardianQueriesHandler(uow uow.UnitOfWork) *GuardianQueriesHandler {
	return &GuardianQueriesHandler{uow: uow}
}

func (h *GuardianQueriesHandler) GetGuardian(ctx context.Context, id shared.ID) (*guardians.Guardian, error) {
	if _, err := authz.Require(ctx, identity.PermGuardiansView); err != nil {
		return nil, err
	}
	var res *guardians.Guardian
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		g, err := r.Guardians.Get(ctx, id)
		if err != nil {
			return err
		}
		res = g
		return nil
	})
	return res, err
}

func (h *GuardianQueriesHandler) ListGuardiansForStudent(ctx context.Context, studentID shared.ID) ([]guardians.StudentGuardianLink, error) {
	if _, err := authz.Require(ctx, identity.PermGuardiansView); err != nil {
		return nil, err
	}
	var res []guardians.StudentGuardianLink
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		list, err := r.Guardians.ListForStudent(ctx, studentID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *GuardianQueriesHandler) FindSiblings(ctx context.Context, studentID shared.ID) ([]students.Student, error) {
	if _, err := authz.Require(ctx, identity.PermGuardiansView); err != nil {
		return nil, err
	}
	var res []students.Student
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		list, err := r.Students.FindSiblings(ctx, studentID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

