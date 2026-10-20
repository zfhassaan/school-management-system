// Package outbox defines the transactional outbox port used to record
// synchronization events in the same transaction as business changes.
package outbox

import (
	"context"
	"time"

	"school-management/internal/domain/shared"
)

// OperationType classifies a synchronized change.
type OperationType string

const (
	OpCreate  OperationType = "create"
	OpUpdate  OperationType = "update"
	OpDelete  OperationType = "delete"
	OpCommand OperationType = "command"
)

// Event is a synchronization operation. The adapter fills the local device id.
type Event struct {
	OperationID          shared.ID
	EntityType           string
	EntityID             string
	OperationType        OperationType
	CommandName          string
	BaseRevision         *int64
	NewRevision          *int64
	Payload              any
	PayloadSchemaVersion int
	IdempotencyKey       string
	OccurredAt           time.Time
}

// Writer appends events within the active transaction.
type Writer interface {
	Append(ctx context.Context, e Event) error
}

