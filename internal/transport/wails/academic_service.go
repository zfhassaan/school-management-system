package wails

import (
	"context"

	acadCmd "school-management/internal/application/commands/academics"
	acadQuery "school-management/internal/application/queries/academics"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/shared"
	"school-management/internal/transport/dto"
)

type AcademicService struct {
	sessions    *SessionManager
	cmdHandler  *acadCmd.AcademicCommandsHandler
	queryHandle *acadQuery.AcademicQueriesHandler
}

func NewAcademicService(
	sessions *SessionManager,
	cmdHandler *acadCmd.AcademicCommandsHandler,
	queryHandle *acadQuery.AcademicQueriesHandler,
) *AcademicService {
	return &AcademicService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *AcademicService) ListAcademicYears(ctx context.Context) dto.Response[[]academics.AcademicYear] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListAcademicYears(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]academics.AcademicYear](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AcademicService) GetActiveAcademicYear(ctx context.Context) dto.Response[*academics.AcademicYear] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetActiveAcademicYear(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*academics.AcademicYear](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AcademicService) CreateAcademicYear(ctx context.Context, in acadCmd.CreateAcademicYearInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateAcademicYear(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *AcademicService) UpdateAcademicYear(ctx context.Context, in acadCmd.UpdateAcademicYearInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateAcademicYear(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *AcademicService) SetTerms(ctx context.Context, in acadCmd.SetTermsInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.SetTerms(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *AcademicService) ListTerms(ctx context.Context, yearID shared.ID) dto.Response[[]academics.AcademicTerm] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListTerms(ctx, yearID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]academics.AcademicTerm](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AcademicService) ListCampuses(ctx context.Context) dto.Response[[]academics.Campus] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListCampuses(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]academics.Campus](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AcademicService) CreateCampus(ctx context.Context, in acadCmd.CreateCampusInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateCampus(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *AcademicService) ListClasses(ctx context.Context) dto.Response[[]academics.Class] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListClasses(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]academics.Class](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AcademicService) CreateClass(ctx context.Context, in acadCmd.CreateClassInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateClass(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *AcademicService) UpdateClass(ctx context.Context, in acadCmd.UpdateClassInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateClass(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *AcademicService) ListSectionsByYear(ctx context.Context, yearID shared.ID) dto.Response[[]academics.Section] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListSectionsByYear(ctx, yearID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]academics.Section](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AcademicService) CreateSection(ctx context.Context, in acadCmd.CreateSectionInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateSection(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *AcademicService) UpdateSection(ctx context.Context, in acadCmd.UpdateSectionInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateSection(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}
