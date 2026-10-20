package uow

import (
	"context"

	"school-management/internal/application/services/outbox"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/admissions"
	"school-management/internal/domain/audit"
	"school-management/internal/domain/documents"
	"school-management/internal/domain/fees"
	"school-management/internal/domain/guardians"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/school"
	"school-management/internal/domain/students"
)

// Repos exposes transaction-bound repositories.
type Repos struct {
	Users         identity.UserRepository
	Roles         identity.RoleRepository
	RecoveryCodes identity.RecoveryCodeRepository
	School        school.Repository
	Academics     academics.Repository
	Students      students.Repository
	Enrollments   students.EnrollmentRepository
	Guardians     guardians.Repository
	Admissions    admissions.Repository
	Documents     documents.Repository
	Fees          fees.Repository
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
	Write(ctx context.Context, fn func(ctx context.Context, r Repos) error) error
	Read(ctx context.Context, fn func(ctx context.Context, r Repos) error) error
}
