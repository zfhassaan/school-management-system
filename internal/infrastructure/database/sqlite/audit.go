package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"school-management/internal/domain/audit"
	"school-management/internal/domain/shared"
)

type AuditWriter struct {
	q Querier
}

func (w *AuditWriter) Append(ctx context.Context, e audit.Entry) error {
	var actorID sql.NullString
	if e.ActorUserID != nil && !e.ActorUserID.IsZero() {
		actorID = sql.NullString{String: e.ActorUserID.String(), Valid: true}
	}
	var deviceID sql.NullString
	if !e.DeviceID.IsZero() {
		deviceID = sql.NullString{String: e.DeviceID.String(), Valid: true}
	}
	var entityID sql.NullString
	if e.EntityID != "" {
		entityID = sql.NullString{String: e.EntityID, Valid: true}
	}
	var reason sql.NullString
	if e.Reason != "" {
		reason = sql.NullString{String: e.Reason, Valid: true}
	}

	var beforeJSON, afterJSON sql.NullString
	if e.Before != nil {
		b, err := json.Marshal(e.Before)
		if err == nil {
			beforeJSON = sql.NullString{String: string(b), Valid: true}
		}
	}
	if e.After != nil {
		b, err := json.Marshal(e.After)
		if err == nil {
			afterJSON = sql.NullString{String: string(b), Valid: true}
		}
	}

	params := AppendAuditLogParams{
		ID:            e.ID.String(),
		OccurredAt:    e.OccurredAt.UTC().Format(time.RFC3339Nano),
		ActorUserID:   actorID,
		DeviceID:      deviceID,
		Action:        e.Action,
		EntityType:    e.EntityType,
		EntityID:      entityID,
		CorrelationID: e.CorrelationID,
		Reason:        reason,
		BeforeJson:    beforeJSON,
		AfterJson:     afterJSON,
		Source:        string(e.Source),
	}

	err := w.q.AppendAuditLog(ctx, params)
	if err != nil && actorID.Valid && strings.Contains(err.Error(), "FOREIGN KEY") {
		// If actor user id was deleted/reseeded while in-memory session stayed open, retry with NULL actor
		params.ActorUserID = sql.NullString{Valid: false}
		err = w.q.AppendAuditLog(ctx, params)
	}
	return err
}

type AuditReader struct {
	q Querier
}

func (r *AuditReader) List(ctx context.Context, f audit.Filter) ([]audit.Entry, error) {
	limit := int64(f.Limit)
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	var actorID sql.NullString
	if f.ActorID != nil && !f.ActorID.IsZero() {
		actorID = sql.NullString{String: f.ActorID.String(), Valid: true}
	}

	rows, err := r.q.ListAuditLogs(ctx, ListAuditLogsParams{
		Column1:     f.EntityType,
		EntityType:  f.EntityType,
		Column3:     f.EntityID,
		EntityID:    sql.NullString{String: f.EntityID, Valid: f.EntityID != ""},
		Column5:     actorID,
		ActorUserID: actorID,
		Limit:       limit,
	})
	if err != nil {
		return nil, shared.Internal(err)
	}

	entries := make([]audit.Entry, len(rows))
	for i, row := range rows {
		occurredAt, _ := time.Parse(time.RFC3339Nano, row.OccurredAt)
		var actorUserID *shared.ID
		if row.ActorUserID.Valid {
			id := shared.ID(row.ActorUserID.String)
			actorUserID = &id
		}
		var before, after any
		if row.BeforeJson.Valid {
			_ = json.Unmarshal([]byte(row.BeforeJson.String), &before)
		}
		if row.AfterJson.Valid {
			_ = json.Unmarshal([]byte(row.AfterJson.String), &after)
		}

		entries[i] = audit.Entry{
			ID:            shared.ID(row.ID),
			OccurredAt:    occurredAt,
			ActorUserID:   actorUserID,
			DeviceID:      shared.ID(row.DeviceID.String),
			Action:        row.Action,
			EntityType:    row.EntityType,
			EntityID:      row.EntityID.String,
			CorrelationID: row.CorrelationID,
			Reason:        row.Reason.String,
			Before:        before,
			After:         after,
			Source:        audit.Source(row.Source),
		}
	}

	return entries, nil
}

