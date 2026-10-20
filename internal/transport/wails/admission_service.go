package wails

import (
	"context"

	admCmd "school-management/internal/application/commands/admissions"
	admQuery "school-management/internal/application/queries/admissions"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/admissions"
	"school-management/internal/domain/shared"
	"school-management/internal/transport/dto"
)

type AdmissionService struct {
	sessions    *SessionManager
	cmdHandler  *admCmd.AdmissionCommandsHandler
	queryHandle *admQuery.AdmissionQueriesHandler
}

func NewAdmissionService(
	sessions *SessionManager,
	cmdHandler *admCmd.AdmissionCommandsHandler,
	queryHandle *admQuery.AdmissionQueriesHandler,
) *AdmissionService {
	return &AdmissionService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *AdmissionService) ListFormDefinitions(ctx context.Context) dto.Response[[]admissions.FormDefinition] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListFormDefinitions(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]admissions.FormDefinition](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AdmissionService) CreateFormDefinition(ctx context.Context, in admCmd.CreateFormDefInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateFormDefinition(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *AdmissionService) ListFormVersions(ctx context.Context, formDefID shared.ID) dto.Response[[]admissions.FormVersion] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListFormVersions(ctx, formDefID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]admissions.FormVersion](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AdmissionService) GetPublishedFormVersion(ctx context.Context, formDefID shared.ID) dto.Response[*admissions.FormVersion] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetPublishedFormVersion(ctx, formDefID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*admissions.FormVersion](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AdmissionService) CreateFormVersion(ctx context.Context, in admCmd.CreateFormVersionInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateFormVersion(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *AdmissionService) PublishFormVersion(ctx context.Context, in admCmd.PublishFormVersionInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.PublishFormVersion(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *AdmissionService) ListApplications(ctx context.Context, status string, yearID string, limit int) dto.Response[[]admissions.Application] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListApplications(ctx, status, yearID, limit)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]admissions.Application](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AdmissionService) GetApplication(ctx context.Context, id shared.ID) dto.Response[*admissions.Application] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetApplication(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*admissions.Application](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AdmissionService) SubmitApplication(ctx context.Context, in admCmd.SubmitApplicationInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.SubmitApplication(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *AdmissionService) UpdateStatus(ctx context.Context, in admCmd.UpdateApplicationStatusInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateStatus(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *AdmissionService) EnrollApplication(ctx context.Context, in admCmd.EnrollApplicationInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.EnrollApplication(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}
