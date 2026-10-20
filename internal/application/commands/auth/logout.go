package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"school-management/internal/domain/identity"
)

type LogoutHandler struct {
	sessionStore identity.SessionStore
}

func NewLogoutHandler(sessionStore identity.SessionStore) *LogoutHandler {
	return &LogoutHandler{sessionStore: sessionStore}
}

func (h *LogoutHandler) Execute(ctx context.Context, rawToken string) error {
	if rawToken == "" {
		return nil
	}
	hasher := sha256.New()
	hasher.Write([]byte(rawToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))
	return h.sessionStore.Revoke(ctx, tokenHash)
}

