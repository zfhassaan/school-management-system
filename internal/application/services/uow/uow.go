// Package uow defines the unit-of-work port. All write use cases run inside
// UnitOfWork.Write so business rows, audit entries and outbox events commit
// atomically.
package uow

import (
	"context"

	"school-management/internal/application/services/outbox"
	"school-management/internal/domain/audit"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/school"
)

// Repos exposes transaction-bound repositories.
type Repos struct {
	Users         identity.UserRepository
	Roles         identity.RoleRepository
	RecoveryCodes identity.RecoveryCodeRepository
	School        school.Repository
	Audit         audit.Writer
	AuditReader   AuditReader
	Outbox        outbox.Writer
}

// AuditReader reads the audit log (queries only).
type AuditReader interface {
	List(ctx context.Context, f audit.Filter) ([]audit.Entry, error)
}

// UnitOfWork runs functions in database transactions.
type UnitOfWork interface {
	// Write runs fn in a serialized write transaction. Returning an error rolls back.
	Write(ctx context.Context, fn func(ctx context.Context, r Repos) error) error
	// Read runs fn against a consistent read snapshot.
	Read(ctx context.Context, fn func(ctx context.Context, r Repos) error) error
}

