package audit

import (
	"context"
	"time"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/audit"
	"school-management/internal/domain/identity"
)

type AuditLogDTO struct {
	ID            string  `json:"id"`
	OccurredAt    string  `json:"occurredAt"`
	ActorUserID   *string `json:"actorUserId,omitempty"`
	DeviceID      string  `json:"deviceId"`
	Action        string  `json:"action"`
	EntityType    string  `json:"entityType"`
	EntityID      string  `json:"entityId"`
	CorrelationID string  `json:"correlationId"`
	Reason        string  `json:"reason"`
	Before        any     `json:"before,omitempty"`
	After         any     `json:"after,omitempty"`
	Source        string  `json:"source"`
}

type AuditQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewAuditQueriesHandler(uow uow.UnitOfWork) *AuditQueriesHandler {
	return &AuditQueriesHandler{uow: uow}
}

func (h *AuditQueriesHandler) ListLogs(ctx context.Context, f audit.Filter) ([]AuditLogDTO, error) {
	if _, err := authz.Require(ctx, identity.PermAuditView); err != nil {
		return nil, err
	}

	var results []AuditLogDTO
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		entries, err := r.AuditReader.List(ctx, f)
		if err != nil {
			return err
		}
		results = make([]AuditLogDTO, len(entries))
		for i, e := range entries {
			var actorID *string
			if e.ActorUserID != nil {
				s := e.ActorUserID.String()
				actorID = &s
			}
			results[i] = AuditLogDTO{
				ID:            e.ID.String(),
				OccurredAt:    e.OccurredAt.Format(time.RFC3339),
				ActorUserID:   actorID,
				DeviceID:      e.DeviceID.String(),
				Action:        e.Action,
				EntityType:    e.EntityType,
				EntityID:      e.EntityID,
				CorrelationID: e.CorrelationID,
				Reason:        e.Reason,
				Before:        e.Before,
				After:         e.After,
				Source:        string(e.Source),
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

