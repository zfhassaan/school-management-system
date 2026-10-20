package wails

import (
	auditQuery "school-management/internal/application/queries/audit"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/audit"
	"school-management/internal/transport/dto"
)

type AuditService struct {
	sessions    *SessionManager
	queryHandle *auditQuery.AuditQueriesHandler
}

func NewAuditService(
	sessions *SessionManager,
	queryHandle *auditQuery.AuditQueriesHandler,
) *AuditService {
	return &AuditService{
		sessions:    sessions,
		queryHandle: queryHandle,
	}
}

func (s *AuditService) ListLogs(f audit.Filter) dto.Response[[]auditQuery.AuditLogDTO] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListLogs(ctx, f)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]auditQuery.AuditLogDTO](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

