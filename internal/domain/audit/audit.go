// Package audit defines the append-only audit trail.
package audit

import (
	"context"
	"time"

	"school-management/internal/domain/shared"
)

// Source identifies the channel that produced an audited action.
type Source string

const (
	SourceDesktop Source = "desktop"
	SourceHTTP    Source = "http"
	SourceSystem  Source = "system"
	SourceSync    Source = "sync"
)

// Entry is an immutable audit record.
type Entry struct {
	ID            shared.ID
	OccurredAt    time.Time
	ActorUserID   *shared.ID
	DeviceID      shared.ID
	Action        string // e.g. "user.created"
	EntityType    string
	EntityID      string
	CorrelationID string
	Reason        string
	Before        any // serialized as JSON; nil when not applicable
	After         any
	Source        Source
}

// Writer appends audit entries inside the caller's transaction.
type Writer interface {
	Append(ctx context.Context, e Entry) error
}

// Filter selects audit entries for the audit log screen.
type Filter struct {
	EntityType string
	EntityID   string
	ActorID    *shared.ID
	From, To   *time.Time
	// Cursor is the ID of the last entry of the previous page (UUIDv7 ordering).
	Cursor string
	Limit  int
}

