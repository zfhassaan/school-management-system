package enrollments

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type EnrollmentQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewEnrollmentQueriesHandler(uow uow.UnitOfWork) *EnrollmentQueriesHandler {
	return &EnrollmentQueriesHandler{uow: uow}
}

func (h *EnrollmentQueriesHandler) GetActiveEnrollment(ctx context.Context, studentID shared.ID) (*students.Enrollment, error) {
	if _, err := authz.Require(ctx, identity.PermStudentsView); err != nil {
		return nil, err
	}
	var res *students.Enrollment
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		enr, err := r.Enrollments.GetActive(ctx, studentID)
		if err != nil {
			return err
		}
		res = enr
		return nil
	})
	return res, err
}

func (h *EnrollmentQueriesHandler) ListEnrollmentHistory(ctx context.Context, studentID shared.ID) ([]students.Enrollment, error) {
	if _, err := authz.Require(ctx, identity.PermStudentsView); err != nil {
		return nil, err
	}
	var res []students.Enrollment
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		list, err := r.Enrollments.ListHistory(ctx, studentID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *EnrollmentQueriesHandler) ListSectionRoster(ctx context.Context, sectionID shared.ID) ([]students.Enrollment, error) {
	if _, err := authz.Require(ctx, identity.PermStudentsView); err != nil {
		return nil, err
	}
	var res []students.Enrollment
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		list, err := r.Enrollments.ListSectionRoster(ctx, sectionID)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *EnrollmentQueriesHandler) GetPromotionBatch(ctx context.Context, id shared.ID) (*students.PromotionBatch, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsView); err != nil {
		return nil, err
	}
	var res *students.PromotionBatch
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		batch, err := r.Enrollments.GetPromotionBatch(ctx, id)
		if err != nil {
			return err
		}
		res = batch
		return nil
	})
	return res, err
}

