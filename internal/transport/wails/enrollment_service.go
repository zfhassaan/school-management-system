package wails

import (
	"context"

	enrCmd "school-management/internal/application/commands/enrollments"
	enrQuery "school-management/internal/application/queries/enrollments"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
	"school-management/internal/transport/dto"
)

type EnrollmentService struct {
	sessions    *SessionManager
	cmdHandler  *enrCmd.EnrollmentCommandsHandler
	queryHandle *enrQuery.EnrollmentQueriesHandler
}

func NewEnrollmentService(
	sessions *SessionManager,
	cmdHandler *enrCmd.EnrollmentCommandsHandler,
	queryHandle *enrQuery.EnrollmentQueriesHandler,
) *EnrollmentService {
	return &EnrollmentService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *EnrollmentService) EnrollStudent(ctx context.Context, in enrCmd.EnrollStudentInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.EnrollStudent(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *EnrollmentService) GetActive(ctx context.Context, studentID shared.ID) dto.Response[*students.Enrollment] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetActiveEnrollment(ctx, studentID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*students.Enrollment](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *EnrollmentService) ListHistory(ctx context.Context, studentID shared.ID) dto.Response[[]students.Enrollment] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListEnrollmentHistory(ctx, studentID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]students.Enrollment](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *EnrollmentService) ListSectionRoster(ctx context.Context, sectionID shared.ID) dto.Response[[]students.Enrollment] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListSectionRoster(ctx, sectionID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]students.Enrollment](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *EnrollmentService) ExecutePromotionBatch(ctx context.Context, in enrCmd.ExecutePromotionBatchInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.ExecutePromotionBatch(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *EnrollmentService) GetPromotionBatch(ctx context.Context, id shared.ID) dto.Response[*students.PromotionBatch] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetPromotionBatch(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*students.PromotionBatch](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

