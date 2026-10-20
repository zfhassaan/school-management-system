package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"school-management/internal/application/services/outbox"
	"school-management/internal/domain/shared"
)

type OutboxWriter struct {
	q        Querier
	deviceID shared.ID
}

func (w *OutboxWriter) Append(ctx context.Context, e outbox.Event) error {
	var cmdName sql.NullString
	if e.CommandName != "" {
		cmdName = sql.NullString{String: e.CommandName, Valid: true}
	}
	var baseRev, newRev sql.NullInt64
	if e.BaseRevision != nil {
		baseRev = sql.NullInt64{Int64: *e.BaseRevision, Valid: true}
	}
	if e.NewRevision != nil {
		newRev = sql.NullInt64{Int64: *e.NewRevision, Valid: true}
	}

	payloadJSON, err := json.Marshal(e.Payload)
	if err != nil {
		return shared.Internal(err)
	}

	return w.q.AppendOutboxEvent(ctx, AppendOutboxEventParams{
		OperationID:          e.OperationID.String(),
		DeviceID:             w.deviceID.String(),
		EntityType:           e.EntityType,
		EntityID:             e.EntityID,
		OperationType:        string(e.OperationType),
		CommandName:          cmdName,
		BaseRevision:         baseRev,
		NewRevision:          newRev,
		PayloadJson:          string(payloadJSON),
		PayloadSchemaVersion: int64(e.PayloadSchemaVersion),
		IdempotencyKey:       e.IdempotencyKey,
		OccurredAt:           e.OccurredAt.UTC().Format(time.RFC3339Nano),
	})
}

