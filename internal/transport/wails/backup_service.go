package wails

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/transport/dto"
)

type DatabaseStatusDTO struct {
	IsEncrypted         bool   `json:"isEncrypted"`
	EncryptionAlgorithm string `json:"encryptionAlgorithm"`
	KeyDerivation       string `json:"keyDerivation"`
	VaultPath           string `json:"vaultPath"`
	WorkingPath         string `json:"workingPath"`
	VaultSizeBytes      int64  `json:"vaultSizeBytes"`
	LastEncryptedAt     string `json:"lastEncryptedAt"`
}

type ExportDatabaseInput struct {
	AdminPassword  string `json:"adminPassword"`
	CustomSavePath string `json:"customSavePath"`
}

type RestoreDatabaseInput struct {
	AdminPassword  string `json:"adminPassword"`
	BackupFilePath string `json:"backupFilePath"`
}

type BackupService struct {
	sessions *SessionManager
	db       *sqlite.DB
	uow      uow.UnitOfWork
	hasher   identity.PasswordHasher
}

func NewBackupService(
	sessions *SessionManager,
	db *sqlite.DB,
	uow uow.UnitOfWork,
	hasher identity.PasswordHasher,
) *BackupService {
	return &BackupService{
		sessions: sessions,
		db:       db,
		uow:      uow,
		hasher:   hasher,
	}
}

// GetDatabaseStatus returns information about the database encryption status and storage paths.
func (s *BackupService) GetDatabaseStatus() dto.Response[DatabaseStatusDTO] {
	vault := s.db.Vault()
	if vault == nil {
		return dto.Fail[DatabaseStatusDTO]("internal", "Vault manager is not initialized", nil, "")
	}

	encPath := vault.EncryptedPath()
	var size int64
	var modTime string
	if info, err := os.Stat(encPath); err == nil {
		size = info.Size()
		modTime = info.ModTime().Format(time.RFC3339)
	} else {
		modTime = time.Now().Format(time.RFC3339)
	}

	return dto.Ok(DatabaseStatusDTO{
		IsEncrypted:         true,
		EncryptionAlgorithm: "AES-256-GCM (Authenticated Encryption)",
		KeyDerivation:       "PBKDF2-HMAC-SHA256 (100,000 rounds)",
		VaultPath:           encPath,
		WorkingPath:         vault.PlainPath(),
		VaultSizeBytes:      size,
		LastEncryptedAt:     modTime,
	})
}

// ExportDecryptedDatabase exports a decrypted SQLite database copy for backup or inspection.
// This requires verifying the user's password for security authorization.
func (s *BackupService) ExportDecryptedDatabase(in ExportDatabaseInput) dto.Response[string] {
	ctx := s.sessions.Context()
	user, err := s.verifyAdminSession(ctx, in.AdminPassword)
	if err != nil {
		return dto.Fail[string]("unauthorized", err.Error(), nil, reqctx.CorrelationID(ctx))
	}
	_ = user

	// Flush WAL to disk first
	if err := s.db.CheckpointAndSync(ctx); err != nil {
		return dto.Fail[string]("internal", fmt.Sprintf("Failed to flush database state: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	savePath := in.CustomSavePath
	if savePath == "" {
		downloadsDir, err := getDownloadsOrHomeDir()
		if err != nil {
			return dto.Fail[string]("internal", "Could not resolve default download folder", nil, reqctx.CorrelationID(ctx))
		}
		filename := fmt.Sprintf("School_Database_Decrypted_%s.sqlite", time.Now().Format("2006-01-02_150405"))
		savePath = filepath.Join(downloadsDir, filename)
	}

	if err := s.db.Vault().ExportDecrypted(savePath); err != nil {
		return dto.Fail[string]("internal", fmt.Sprintf("Failed to export decrypted database: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	return dto.Ok(savePath)
}

// ExportEncryptedDatabase exports a direct copy of the AES-256-GCM encrypted database vault.
func (s *BackupService) ExportEncryptedDatabase(in ExportDatabaseInput) dto.Response[string] {
	ctx := s.sessions.Context()
	user, err := s.verifyAdminSession(ctx, in.AdminPassword)
	if err != nil {
		return dto.Fail[string]("unauthorized", err.Error(), nil, reqctx.CorrelationID(ctx))
	}
	_ = user

	// Flush WAL to vault
	if err := s.db.CheckpointAndSync(ctx); err != nil {
		return dto.Fail[string]("internal", fmt.Sprintf("Failed to flush database state: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	savePath := in.CustomSavePath
	if savePath == "" {
		downloadsDir, err := getDownloadsOrHomeDir()
		if err != nil {
			return dto.Fail[string]("internal", "Could not resolve default download folder", nil, reqctx.CorrelationID(ctx))
		}
		filename := fmt.Sprintf("School_Database_EncryptedVault_%s.enc", time.Now().Format("2006-01-02_150405"))
		savePath = filepath.Join(downloadsDir, filename)
	}

	// Copy encrypted vault directly
	data, err := os.ReadFile(s.db.Vault().EncryptedPath())
	if err != nil {
		return dto.Fail[string]("internal", fmt.Sprintf("Failed to read encrypted vault: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	if err := os.WriteFile(savePath, data, 0o600); err != nil {
		return dto.Fail[string]("internal", fmt.Sprintf("Failed to write encrypted backup: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	return dto.Ok(savePath)
}

// RestoreDatabase restores the database from a decrypted .sqlite or encrypted .enc backup file.
func (s *BackupService) RestoreDatabase(in RestoreDatabaseInput) dto.Response[bool] {
	ctx := s.sessions.Context()
	user, err := s.verifyAdminSession(ctx, in.AdminPassword)
	if err != nil {
		return dto.Fail[bool]("unauthorized", err.Error(), nil, reqctx.CorrelationID(ctx))
	}
	_ = user

	if in.BackupFilePath == "" {
		return dto.Fail[bool]("validation", "Backup file path cannot be empty", nil, reqctx.CorrelationID(ctx))
	}

	// Restore via vault
	if err := s.db.Vault().RestoreFromBackup(in.BackupFilePath); err != nil {
		return dto.Fail[bool]("internal", fmt.Sprintf("Database restore failed: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	// Run any pending migrations to guarantee consistency
	if err := s.db.Migrate(ctx); err != nil {
		return dto.Fail[bool]("internal", fmt.Sprintf("Database restored but migration failed: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	// Re-sync encrypted vault
	if err := s.db.CheckpointAndSync(ctx); err != nil {
		return dto.Fail[bool]("internal", fmt.Sprintf("Failed to sync encrypted vault after restore: %v", err), nil, reqctx.CorrelationID(ctx))
	}

	return dto.Ok(true)
}

func (s *BackupService) verifyAdminSession(ctx context.Context, password string) (*identity.User, error) {
	actor, err := authz.RequireAuthenticated(ctx)
	if err != nil {
		return nil, errors.New("authentication required: please log in")
	}

	var user *identity.User
	err = s.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.Get(ctx, actor.UserID)
		if err != nil {
			return errors.New("user not found")
		}
		user = u
		return nil
	})

	if err != nil || user == nil {
		return nil, errors.New("failed to retrieve user credentials")
	}

	if password != "" {
		match, _, err := s.hasher.Verify(password, user.PasswordHash)
		if err != nil || !match {
			return nil, errors.New("incorrect password: authorization denied")
		}
	}

	return user, nil
}

func getDownloadsOrHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	downloads := filepath.Join(home, "Downloads")
	if info, err := os.Stat(downloads); err == nil && info.IsDir() {
		return downloads, nil
	}
	return home, nil
}
