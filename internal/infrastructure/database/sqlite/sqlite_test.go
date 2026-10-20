package sqlite_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/school"
	"school-management/internal/domain/shared"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/platform/config"
	"school-management/internal/platform/security"
)

func TestSQLiteIntegration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sms-sqlite-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.DefaultConfig(tempDir)
	cfg.Database.Path = filepath.Join(tempDir, "test.db")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	ctx := context.Background()
	db, err := sqlite.Open(ctx, &cfg, logger)
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	defer db.Close()

	if db.DeviceID().IsZero() {
		t.Errorf("expected local device ID to be initialized")
	}

	uowInstance := sqlite.NewUnitOfWork(db)

	// 1. Verify system roles and permission catalog were seeded
	err = uowInstance.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		roles, err := r.Roles.List(ctx)
		if err != nil {
			return err
		}
		if len(roles) == 0 {
			t.Errorf("expected seeded system roles, got none")
		}
		superAdmin, err := r.Roles.GetByCode(ctx, identity.RoleSuperAdmin)
		if err != nil {
			t.Errorf("failed to get super admin role: %v", err)
		} else if len(superAdmin.Permissions) < len(identity.Catalog) {
			t.Errorf("super admin expected all catalog permissions, got %d", len(superAdmin.Permissions))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed read transaction: %v", err)
	}

	// 2. Insert School and Super Admin user in a single write transaction
	hasher := security.NewArgon2Hasher(security.Argon2Params{
		Memory:      1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})
	hash, err := hasher.Hash("SuperSecret123!")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	schoolID := shared.UUIDv7Generator{}.NewID()
	userID := shared.UUIDv7Generator{}.NewID()
	now := time.Now().UTC()

	err = uowInstance.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		superAdminRole, err := r.Roles.GetByCode(ctx, identity.RoleSuperAdmin)
		if err != nil {
			return err
		}

		sch := &school.School{
			ID:                 schoolID,
			Name:               "Horizon International Academy",
			LegalName:          "Horizon Academy Trust",
			Code:               "HORIZON-01",
			Timezone:           "Asia/Karachi",
			CurrencyCode:       "PKR",
			CurrencyMinorUnits: 2,
			Locale:             "en",
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := r.School.Insert(ctx, sch); err != nil {
			return err
		}

		settings := school.DefaultSettings(schoolID)
		if err := r.School.InsertSettings(ctx, &settings, now); err != nil {
			return err
		}

		u := &identity.User{
			ID:                userID,
			Username:          "superadmin",
			DisplayName:       "Headmaster",
			Email:             "admin@horizon.edu",
			PasswordHash:      hash,
			PasswordChangedAt: now,
			Status:            identity.UserActive,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		if err := r.Users.Insert(ctx, u); err != nil {
			return err
		}

		return r.Users.SetRoles(ctx, userID, []shared.ID{superAdminRole.ID}, nil, now)
	})
	if err != nil {
		t.Fatalf("failed write transaction: %v", err)
	}

	// 3. Verify user retrieval and effective permissions
	err = uowInstance.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.GetByUsername(ctx, "SuperAdmin") // case-insensitive test
		if err != nil {
			return err
		}
		if u.DisplayName != "Headmaster" {
			t.Errorf("got %q, want 'Headmaster'", u.DisplayName)
		}
		perms, err := r.Users.EffectivePermissions(ctx, u.ID)
		if err != nil {
			return err
		}
		if len(perms) < len(identity.Catalog) {
			t.Errorf("expected all permissions, got %d", len(perms))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("verification read failed: %v", err)
	}
}


