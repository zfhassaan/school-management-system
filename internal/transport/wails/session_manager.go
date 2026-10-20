package wails

import (
	"context"
	"sync"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/audit"
	"school-management/internal/domain/shared"
)

// SessionManager manages desktop in-memory session and injects Actor into Context.
type SessionManager struct {
	mu           sync.RWMutex
	currentActor *authz.Actor
	ids          shared.IDGenerator
}

func NewSessionManager(ids shared.IDGenerator) *SessionManager {
	return &SessionManager{ids: ids}
}

func (m *SessionManager) SetActor(actor *authz.Actor) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentActor = actor
}

func (m *SessionManager) ClearActor() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentActor = nil
}

func (m *SessionManager) CurrentActor() *authz.Actor {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentActor
}

// Context returns a new context decorated with current desktop Actor, correlation ID, and desktop audit source.
func (m *SessionManager) Context() context.Context {
	ctx := reqctx.New(context.Background(), m.ids, audit.SourceDesktop)
	m.mu.RLock()
	actor := m.currentActor
	m.mu.RUnlock()
	if actor != nil {
		ctx = authz.WithActor(ctx, *actor)
	}
	return ctx
}

// WithContext decorates ctx with the current desktop Actor, correlation ID, and desktop audit source.
func (m *SessionManager) WithContext(ctx context.Context) context.Context {
	ctx = reqctx.New(ctx, m.ids, audit.SourceDesktop)
	m.mu.RLock()
	actor := m.currentActor
	m.mu.RUnlock()
	if actor != nil {
		ctx = authz.WithActor(ctx, *actor)
	}
	return ctx
}

// MapError translates a domain error into the standard ErrorDTO.
func MapError(err error) (code, message string, fields map[string]string) {
	if err == nil {
		return "", "", nil
	}
	var domErr *shared.Error
	if e, ok := err.(*shared.Error); ok {
		domErr = e
	}
	if domErr != nil {
		return string(domErr.Code), domErr.Message, domErr.Fields
	}
	return string(shared.CodeInternal), err.Error(), nil
}
