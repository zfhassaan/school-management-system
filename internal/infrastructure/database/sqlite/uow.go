package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"school-management/internal/application/services/uow"
)

type UnitOfWork struct {
	db *DB
}

func NewUnitOfWork(db *DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Write(ctx context.Context, fn func(ctx context.Context, r uow.Repos) error) error {
	// Execute on the dedicated single-connection writer pool with immediate lock
	tx, err := u.db.writerDB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return fmt.Errorf("failed to begin write transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	q := New(tx)
	repos := uow.Repos{
		Users:         &UserRepository{q: q},
		Roles:         &RoleRepository{q: q},
		RecoveryCodes: &RecoveryCodeRepository{q: q},
		School:        &SchoolRepository{q: q},
		Academics:     &AcademicsRepository{q: q},
		Students:      &StudentsRepository{q: q},
		Enrollments:   &EnrollmentsRepository{q: q},
		Guardians:     &GuardiansRepository{q: q},
		Admissions:    &AdmissionsRepository{q: q},
		Documents:     &DocumentsRepository{q: q},
		Fees:          &FeesRepository{q: q},
		Audit:         &AuditWriter{q: q},
		AuditReader:   &AuditReader{q: q},
		Outbox:        &OutboxWriter{q: q, deviceID: u.db.deviceID},
	}

	if err := fn(ctx, repos); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit write transaction: %w", err)
	}
	return nil
}

func (u *UnitOfWork) Read(ctx context.Context, fn func(ctx context.Context, r uow.Repos) error) error {
	// Execute on reader pool
	q := New(u.db.readerDB)
	repos := uow.Repos{
		Users:         &UserRepository{q: q},
		Roles:         &RoleRepository{q: q},
		RecoveryCodes: &RecoveryCodeRepository{q: q},
		School:        &SchoolRepository{q: q},
		Academics:     &AcademicsRepository{q: q},
		Students:      &StudentsRepository{q: q},
		Enrollments:   &EnrollmentsRepository{q: q},
		Guardians:     &GuardiansRepository{q: q},
		Admissions:    &AdmissionsRepository{q: q},
		Documents:     &DocumentsRepository{q: q},
		Fees:          &FeesRepository{q: q},
		AuditReader:   &AuditReader{q: q},
	}
	return fn(ctx, repos)
}

