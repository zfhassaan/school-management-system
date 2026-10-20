package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type SessionStore struct {
	db *DB
}

func NewSessionStore(db *DB) *SessionStore {
	return &SessionStore{db: db}
}

func (s *SessionStore) Create(ctx context.Context, session identity.Session) error {
	q := New(s.db.writerDB)
	var deviceID sql.NullString
	if !session.DeviceID.IsZero() {
		deviceID = sql.NullString{String: session.DeviceID.String(), Valid: true}
	}
	return q.CreateSession(ctx, CreateSessionParams{
		ID:                session.ID,
		UserID:            session.UserID.String(),
		TokenHash:         session.ID, // ID is sha256 token hash
		DeviceID:          deviceID,
		CreatedAt:         session.CreatedAt.UTC().Format(time.RFC3339Nano),
		LastSeenAt:        session.LastSeenAt.UTC().Format(time.RFC3339Nano),
		IdleExpiresAt:     session.IdleExpiresAt.UTC().Format(time.RFC3339Nano),
		AbsoluteExpiresAt: session.AbsoluteExpiresAt.UTC().Format(time.RFC3339Nano),
	})
}

func (s *SessionStore) Get(ctx context.Context, id string) (*identity.Session, error) {
	q := New(s.db.readerDB)
	row, err := q.GetSession(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("session")
		}
		return nil, shared.Internal(err)
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	lastSeen, _ := time.Parse(time.RFC3339Nano, row.LastSeenAt)
	idleExpires, _ := time.Parse(time.RFC3339Nano, row.IdleExpiresAt)
	absExpires, _ := time.Parse(time.RFC3339Nano, row.AbsoluteExpiresAt)

	return &identity.Session{
		ID:                row.ID,
		UserID:            shared.ID(row.UserID),
		DeviceID:          shared.ID(row.DeviceID.String),
		CreatedAt:         createdAt,
		LastSeenAt:        lastSeen,
		IdleExpiresAt:     idleExpires,
		AbsoluteExpiresAt: absExpires,
	}, nil
}

func (s *SessionStore) Touch(ctx context.Context, id string, lastSeen, idleExpires time.Time) error {
	q := New(s.db.writerDB)
	return q.TouchSession(ctx, TouchSessionParams{
		ID:            id,
		LastSeenAt:    lastSeen.UTC().Format(time.RFC3339Nano),
		IdleExpiresAt: idleExpires.UTC().Format(time.RFC3339Nano),
	})
}

func (s *SessionStore) Revoke(ctx context.Context, id string) error {
	q := New(s.db.writerDB)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return q.RevokeSession(ctx, RevokeSessionParams{
		ID:        id,
		RevokedAt: sql.NullString{String: now, Valid: true},
	})
}

func (s *SessionStore) RevokeAllForUser(ctx context.Context, userID shared.ID) error {
	q := New(s.db.writerDB)
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return q.RevokeAllUserSessions(ctx, RevokeAllUserSessionsParams{
		UserID:    userID.String(),
		RevokedAt: sql.NullString{String: now, Valid: true},
	})
}

