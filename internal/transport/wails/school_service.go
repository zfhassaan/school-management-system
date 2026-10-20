package wails

import (
	"context"

	schoolCmd "school-management/internal/application/commands/school"
	schoolQuery "school-management/internal/application/queries/school"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/school"
	"school-management/internal/transport/dto"
)

type SchoolService struct {
	sessions    *SessionManager
	cmdHandler  *schoolCmd.SchoolCommandsHandler
	queryHandle *schoolQuery.SchoolQueriesHandler
}

func NewSchoolService(
	sessions *SessionManager,
	cmdHandler *schoolCmd.SchoolCommandsHandler,
	queryHandle *schoolQuery.SchoolQueriesHandler,
) *SchoolService {
	return &SchoolService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *SchoolService) GetSchool(ctx context.Context) dto.Response[*school.School] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetSchool(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*school.School](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *SchoolService) UpdateSchoolProfile(ctx context.Context, in schoolCmd.UpdateSchoolProfileInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateProfile(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *SchoolService) GetSettings(ctx context.Context) dto.Response[*school.Settings] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetSettings(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*school.Settings](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *SchoolService) UpdateSettings(ctx context.Context, in schoolCmd.UpdateSchoolSettingsInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateSettings(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

