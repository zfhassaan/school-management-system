// Package authz carries the authenticated actor through context and enforces
// permissions inside application services (never only in the UI).
package authz

import (
	"context"

	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

// Actor is the authenticated principal executing a use case.
type Actor struct {
	UserID      shared.ID
	Username    string
	DisplayName string
	DeviceID    shared.ID
	SessionID   string
	Roles       []string
	Permissions identity.PermissionSet
	// MustChangePassword restricts the actor to password change only.
	MustChangePassword bool
	// System marks internal jobs (startup sync, scheduled backups).
	System bool
}

type actorKey struct{}

// WithActor returns a context carrying the actor.
func WithActor(ctx context.Context, a Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, a)
}

// ActorFrom returns the actor or an unauthenticated error.
func ActorFrom(ctx context.Context) (Actor, error) {
	a, ok := ctx.Value(actorKey{}).(Actor)
	if !ok {
		return Actor{}, shared.Unauthenticated("sign in required")
	}
	return a, nil
}

// SystemActor returns an actor for internal jobs with all permissions.
func SystemActor(deviceID shared.ID) Actor {
	perms := make([]identity.Permission, len(identity.Catalog))
	for i, p := range identity.Catalog {
		perms[i] = p.Code
	}
	return Actor{Username: "system", DisplayName: "System", DeviceID: deviceID, System: true, Permissions: identity.NewPermissionSet(perms)}
}

// Require returns the actor if it holds every listed permission.
func Require(ctx context.Context, perms ...identity.Permission) (Actor, error) {
	a, err := ActorFrom(ctx)
	if err != nil {
		return Actor{}, err
	}
	if a.MustChangePassword {
		return Actor{}, shared.Forbidden("password change required")
	}
	for _, p := range perms {
		if !a.Permissions.Has(p) {
			return Actor{}, shared.Forbidden("missing permission " + string(p))
		}
	}
	return a, nil
}

// RequireAuthenticated returns the actor without a permission check; used for
// self-service operations (change own password, sign out).
func RequireAuthenticated(ctx context.Context) (Actor, error) {
	return ActorFrom(ctx)
}

// UserIDPtr returns a pointer to the actor's user id, or nil for system actors.
func (a Actor) UserIDPtr() *shared.ID {
	if a.UserID.IsZero() {
		return nil
	}
	id := a.UserID
	return &id
}

