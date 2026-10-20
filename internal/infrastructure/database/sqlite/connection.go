package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	migrations "school-management/database/migrations/sqlite"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/platform/config"
)

type DB struct {
	writerDB *sql.DB
	readerDB *sql.DB
	logger   *slog.Logger
	deviceID shared.ID
}

func Open(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*DB, error) {
	// Build Writer DSN with IMMEDIATE lock and WAL pragmas
	// Note: ModernC sqlite driver accepts parameters formatted as url query params
	writerParams := url.Values{}
	writerParams.Add("_pragma", "foreign_keys(1)")
	writerParams.Add("_pragma", "journal_mode(wal)")
	writerParams.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", cfg.Database.BusyTimeout))
	writerParams.Add("_pragma", "synchronous(normal)")
	writerParams.Add("_txlock", "immediate")

	writerDSN := fmt.Sprintf("file:%s?%s", cfg.Database.Path, writerParams.Encode())
	writerDB, err := sql.Open("sqlite", writerDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open writer sqlite connection: %w", err)
	}

	// Strictly 1 writer connection to prevent SQLITE_BUSY concurrency conflicts
	writerDB.SetMaxOpenConns(1)
	writerDB.SetMaxIdleConns(1)
	writerDB.SetConnMaxLifetime(0)

	// Build Reader DSN
	readerParams := url.Values{}
	readerParams.Add("_pragma", "foreign_keys(1)")
	readerParams.Add("_pragma", "journal_mode(wal)")
	readerParams.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", cfg.Database.BusyTimeout))
	readerParams.Add("_pragma", "synchronous(normal)")
	readerParams.Add("_pragma", "query_only(1)")

	readerDSN := fmt.Sprintf("file:%s?%s", cfg.Database.Path, readerParams.Encode())
	readerDB, err := sql.Open("sqlite", readerDSN)
	if err != nil {
		_ = writerDB.Close()
		return nil, fmt.Errorf("failed to open reader sqlite connection: %w", err)
	}

	maxReadConns := cfg.Database.MaxReadConns
	if maxReadConns <= 0 {
		maxReadConns = 4
	}
	readerDB.SetMaxOpenConns(maxReadConns)
	readerDB.SetMaxIdleConns(maxReadConns)
	readerDB.SetConnMaxLifetime(0)

	// Verify connectivity
	if err := writerDB.PingContext(ctx); err != nil {
		_ = writerDB.Close()
		_ = readerDB.Close()
		return nil, fmt.Errorf("failed to ping sqlite writer database: %w", err)
	}
	if err := readerDB.PingContext(ctx); err != nil {
		_ = writerDB.Close()
		_ = readerDB.Close()
		return nil, fmt.Errorf("failed to ping sqlite reader database: %w", err)
	}

	db := &DB{
		writerDB: writerDB,
		readerDB: readerDB,
		logger:   logger,
	}

	// Run migrations
	if err := db.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	// Initialize local device and permission catalog
	if err := db.initDeviceAndCatalog(ctx, cfg.App); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to init device and permission catalog: %w", err)
	}

	return db, nil
}

func (db *DB) Migrate(ctx context.Context) error {
	db.logger.Info("running embedded sqlite database migrations")
	goose.SetBaseFS(migrations.SQLiteFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.UpContext(ctx, db.writerDB, "."); err != nil {
		return fmt.Errorf("goose up failed: %w", err)
	}
	db.logger.Info("database migrations up to date")
	return nil
}

func (db *DB) initDeviceAndCatalog(ctx context.Context, appCfg config.AppConfig) error {
	queries := New(db.writerDB)

	// Synchronize permission catalog
	for _, perm := range identity.Catalog {
		sensitive := int64(0)
		if perm.Sensitive {
			sensitive = 1
		}
		err := queries.UpsertPermission(ctx, UpsertPermissionParams{
			Code:        string(perm.Code),
			Module:      perm.Module,
			Description: perm.Description,
			IsSensitive: sensitive,
		})
		if err != nil {
			return fmt.Errorf("failed to sync permission %s: %w", perm.Code, err)
		}
	}

	// Synchronize system roles
	for _, sysRole := range identity.SystemRoles() {
		r, err := queries.GetRoleByCode(ctx, sysRole.Code)
		var roleID string
		now := time.Now().UTC().Format(time.RFC3339Nano)
		if err != nil {
			// Role doesn't exist, create it
			roleID = shared.UUIDv7Generator{}.NewID().String()
			err = queries.InsertRole(ctx, InsertRoleParams{
				ID:          roleID,
				Code:        sysRole.Code,
				Name:        sysRole.Name,
				Description: sql.NullString{String: sysRole.Description, Valid: true},
				IsSystem:    1,
				CreatedAt:   now,
				UpdatedAt:   now,
				RowVersion:  1,
			})
			if err != nil {
				return fmt.Errorf("failed to create system role %s: %w", sysRole.Code, err)
			}
			for _, p := range sysRole.Permissions {
				if err := queries.InsertRolePermission(ctx, InsertRolePermissionParams{
					RoleID:         roleID,
					PermissionCode: string(p),
				}); err != nil {
					return fmt.Errorf("failed to assign permission %s to role %s: %w", p, sysRole.Code, err)
				}
			}
		} else {
			roleID = r.ID
			// If it's super_admin, guarantee it has all catalog permissions
			if sysRole.Code == identity.RoleSuperAdmin {
				_ = queries.ClearRolePermissions(ctx, roleID)
				for _, p := range sysRole.Permissions {
					_ = queries.InsertRolePermission(ctx, InsertRolePermissionParams{
						RoleID:         roleID,
						PermissionCode: string(p),
					})
				}
			}
		}
	}

	// Register or fetch local device
	dev, err := queries.GetLocalDevice(ctx)
	if err != nil {
		deviceID := shared.UUIDv7Generator{}.NewID().String()
		now := time.Now().UTC().Format(time.RFC3339Nano)
		code := appCfg.DeviceCode
		if code == "" {
			code = "DSK1"
		}
		name := appCfg.DeviceName
		if name == "" {
			name = "Primary Desktop"
		}
		err = queries.InsertLocalDevice(ctx, InsertLocalDeviceParams{
			ID:           deviceID,
			DeviceCode:   code,
			Name:         name,
			Kind:         "desktop",
			PublicKey:    sql.NullString{},
			RegisteredAt: now,
			LastSeenAt:   sql.NullString{String: now, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to register local device: %w", err)
		}
		db.deviceID = shared.ID(deviceID)
	} else {
		db.deviceID = shared.ID(dev.ID)
	}

	return nil
}

func (db *DB) DeviceID() shared.ID {
	return db.deviceID
}

func (db *DB) Close() error {
	var err1, err2 error
	if db.readerDB != nil {
		err1 = db.readerDB.Close()
	}
	if db.writerDB != nil {
		err2 = db.writerDB.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}

func (db *DB) WriterDB() *sql.DB {
	return db.writerDB
}

func (db *DB) ReaderDB() *sql.DB {
	return db.readerDB
}
