// Package auditlog builds audit entries from request context.
package auditlog

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/audit"
	"school-management/internal/domain/shared"
)

// Recorder creates audit entries with actor, device and correlation metadata.
type Recorder struct {
	IDs   shared.IDGenerator
	Clock shared.Clock
}

// Change describes an audited change.
type Change struct {
	Action     string
	EntityType string
	EntityID   string
	Reason     string
	Before     any
	After      any
}

// Record appends an audit entry via w.
func (r Recorder) Record(ctx context.Context, w audit.Writer, c Change) error {
	e := audit.Entry{
		ID:            r.IDs.NewID(),
		OccurredAt:    r.Clock.Now(),
		Action:        c.Action,
		EntityType:    c.EntityType,
		EntityID:      c.EntityID,
		CorrelationID: reqctx.CorrelationID(ctx),
		Reason:        c.Reason,
		Before:        c.Before,
		After:         c.After,
		Source:        reqctx.Source(ctx),
	}
	if e.CorrelationID == "" {
		e.CorrelationID = e.ID.String()
	}
	if a, err := authz.ActorFrom(ctx); err == nil {
		e.ActorUserID = a.UserIDPtr()
		e.DeviceID = a.DeviceID
	}
	return w.Append(ctx, e)
}

