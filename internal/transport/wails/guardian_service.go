package wails

import (
	"context"

	guardCmd "school-management/internal/application/commands/guardians"
	guardQuery "school-management/internal/application/queries/guardians"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/guardians"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
	"school-management/internal/transport/dto"
)

type GuardianService struct {
	sessions    *SessionManager
	cmdHandler  *guardCmd.GuardianCommandsHandler
	queryHandle *guardQuery.GuardianQueriesHandler
}

func NewGuardianService(
	sessions *SessionManager,
	cmdHandler *guardCmd.GuardianCommandsHandler,
	queryHandle *guardQuery.GuardianQueriesHandler,
) *GuardianService {
	return &GuardianService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *GuardianService) GetGuardian(ctx context.Context, id shared.ID) dto.Response[*guardians.Guardian] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetGuardian(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*guardians.Guardian](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *GuardianService) CreateGuardian(ctx context.Context, in guardCmd.CreateGuardianInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateGuardian(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *GuardianService) UpdateGuardian(ctx context.Context, in guardCmd.UpdateGuardianInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateGuardian(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *GuardianService) ListForStudent(ctx context.Context, studentID shared.ID) dto.Response[[]guardians.StudentGuardianLink] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListGuardiansForStudent(ctx, studentID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]guardians.StudentGuardianLink](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *GuardianService) SetStudentGuardians(ctx context.Context, in guardCmd.SetStudentGuardiansInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.SetStudentGuardians(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *GuardianService) FindSiblings(ctx context.Context, studentID shared.ID) dto.Response[[]students.Student] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.FindSiblings(ctx, studentID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]students.Student](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}
