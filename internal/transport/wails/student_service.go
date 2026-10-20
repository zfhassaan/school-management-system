package wails

import (
	"context"

	studCmd "school-management/internal/application/commands/students"
	studQuery "school-management/internal/application/queries/students"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
	"school-management/internal/transport/dto"
)

type StudentService struct {
	sessions    *SessionManager
	cmdHandler  *studCmd.StudentCommandsHandler
	queryHandle *studQuery.StudentQueriesHandler
}

func NewStudentService(
	sessions *SessionManager,
	cmdHandler *studCmd.StudentCommandsHandler,
	queryHandle *studQuery.StudentQueriesHandler,
) *StudentService {
	return &StudentService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *StudentService) GetStudent(ctx context.Context, id shared.ID) dto.Response[*students.Student] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetStudent(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*students.Student](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *StudentService) ListStudents(ctx context.Context, filter students.StudentFilter) dto.Response[[]students.Student] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListStudents(ctx, filter)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]students.Student](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *StudentService) RegisterStudent(ctx context.Context, in studCmd.RegisterStudentInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.RegisterStudent(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *StudentService) UpdateStudent(ctx context.Context, in studCmd.UpdateStudentInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateStudent(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *StudentService) ChangeStatus(ctx context.Context, in studCmd.ChangeStudentStatusInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.ChangeStudentStatus(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *StudentService) ListStatusHistory(ctx context.Context, studentID shared.ID) dto.Response[[]students.StatusHistoryEntry] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListStatusHistory(ctx, studentID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]students.StatusHistoryEntry](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}
