package students

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type StudentQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewStudentQueriesHandler(uow uow.UnitOfWork) *StudentQueriesHandler {
	return &StudentQueriesHandler{uow: uow}
}

func (h *StudentQueriesHandler) GetStudent(ctx context.Context, id shared.ID) (*students.Student, error) {
	if _, err := authz.Require(ctx, identity.PermStudentsView); err != nil {
		return nil, err
	}
	var res *students.Student
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		s, err := r.Students.Get(ctx, id)
		if err != nil {
			return err
		}
		res = s
		return nil
	})
	return res, err
}

func (h *StudentQueriesHandler) ListStudents(ctx context.Context, filter students.StudentFilter) ([]students.Student, error) {
	if _, err := authz.Require(ctx, identity.PermStudentsView); err != nil {
		return nil, err
	}
	var res []students.Student
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		filter.SchoolID = sch.ID

		if filter.SearchQuery != "" {
			list, err := r.Students.SearchFTS(ctx, sch.ID, filter.SearchQuery, filter.Limit)
			if err != nil {
				return err
			}
			res = list
			return nil
		}

		list, err := r.Students.ListCursor(ctx, filter)
		if err != nil {
			return err
		}
		res = list
		return nil
	})
	return res, err
}

func (h *StudentQueriesHandler) ListStatusHistory(ctx context.Context, studentID shared.ID) ([]students.StatusHistoryEntry, error) {
	if _, err := authz.Require(ctx, identity.PermStudentsView); err != nil {
		return nil, err
	}
	var res []students.StatusHistoryEntry
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		hist, err := r.Students.ListStatusHistory(ctx, studentID)
		if err != nil {
			return err
		}
		res = hist
		return nil
	})
	return res, err
}

