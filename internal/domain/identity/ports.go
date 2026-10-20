package identity

import (
	"context"
	"time"

	"school-management/internal/domain/shared"
)

// UserRepository persists users. Methods map to real use cases.
type UserRepository interface {
	Get(ctx context.Context, id shared.ID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Count(ctx context.Context) (int64, error)
	List(ctx context.Context) ([]User, error)
	Insert(ctx context.Context, u *User) error
	// Update persists all mutable fields using optimistic locking on Revision.
	// On success u.Revision is incremented.
	Update(ctx context.Context, u *User) error
	SetRoles(ctx context.Context, userID shared.ID, roleIDs []shared.ID, assignedBy *shared.ID, at time.Time) error
	// EffectivePermissions returns the union of permissions over all roles.
	EffectivePermissions(ctx context.Context, userID shared.ID) ([]Permission, error)
	// RoleCodes returns the codes of roles assigned to the user.
	RoleCodes(ctx context.Context, userID shared.ID) ([]string, error)
	// CountActiveWithRole counts active users holding the role (protects the last super admin).
	CountActiveWithRole(ctx context.Context, roleCode string) (int64, error)
}

// RoleRepository persists roles and the permission catalog.
type RoleRepository interface {
	Get(ctx context.Context, id shared.ID) (*Role, error)
	GetByCode(ctx context.Context, code string) (*Role, error)
	List(ctx context.Context) ([]Role, error)
	Insert(ctx context.Context, r *Role, at time.Time) error
	Update(ctx context.Context, r *Role, at time.Time) error
	Delete(ctx context.Context, id shared.ID) error
	SetPermissions(ctx context.Context, roleID shared.ID, perms []Permission) error
	UpsertPermissionCatalog(ctx context.Context, catalog []PermissionInfo) error
}

// RecoveryCode is a one-time account recovery secret (stored hashed).
type RecoveryCode struct {
	ID       shared.ID
	UserID   shared.ID
	CodeHash string
	UsedAt   *time.Time
}

type RecoveryCodeRepository interface {
	ReplaceForUser(ctx context.Context, userID shared.ID, codes []RecoveryCode, at time.Time) error
	ListUnused(ctx context.Context, userID shared.ID) ([]RecoveryCode, error)
	MarkUsed(ctx context.Context, id shared.ID, at time.Time) error
}

// PasswordHasher hashes and verifies secrets (Argon2id in production).
type PasswordHasher interface {
	Hash(secret string) (string, error)
	// Verify reports whether secret matches encoded and whether the hash should
	// be upgraded because parameters changed.
	Verify(secret, encoded string) (ok bool, needsRehash bool, err error)
}

// Session is an authenticated session. ID is the SHA-256 of the bearer token;
// the raw token is never stored.
type Session struct {
	ID                string
	UserID            shared.ID
	DeviceID          shared.ID
	CreatedAt         time.Time
	LastSeenAt        time.Time
	IdleExpiresAt     time.Time
	AbsoluteExpiresAt time.Time
}

// Expired reports whether the session is no longer valid at now.
func (s Session) Expired(now time.Time) bool {
	return !now.Before(s.IdleExpiresAt) || !now.Before(s.AbsoluteExpiresAt)
}

// SessionPolicy defines session lifetimes.
type SessionPolicy struct {
	IdleTimeout     time.Duration
	AbsoluteTimeout time.Duration
}

var DefaultSessionPolicy = SessionPolicy{IdleTimeout: 30 * time.Minute, AbsoluteTimeout: 12 * time.Hour}

// SessionStore stores sessions (in memory for desktop, database for HTTP).
type SessionStore interface {
	Create(ctx context.Context, s Session) error
	Get(ctx context.Context, id string) (*Session, error)
	Touch(ctx context.Context, id string, lastSeen, idleExpires time.Time) error
	Revoke(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID shared.ID) error
}

