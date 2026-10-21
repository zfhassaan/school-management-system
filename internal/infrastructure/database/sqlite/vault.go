package sqlite

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// VaultMagicHeader identifies authenticated encrypted database files
	VaultMagicHeader = "SMS_VAULT_AES256GCM_V1\n"
	pbkdf2Iterations = 100000
	saltLength       = 16
	nonceLength      = 12
	defaultAppSecret = "FISS-SMS-SECURE-VAULT-AUTHENTICATOR-MASTER-KEY-2026-TIGHT-ENCRYPTION"
)

// VaultManager handles transparent encryption and decryption of the SQLite database at rest.
type VaultManager struct {
	mu          sync.Mutex
	plainDBPath string
	encDBPath   string
	secretKey   string
}

// NewVaultManager creates a VaultManager for the given plain SQLite database path.
// The encrypted database file is stored at <plainDBPath>.enc.
func NewVaultManager(plainDBPath string) *VaultManager {
	encPath := plainDBPath + ".enc"
	return &VaultManager{
		plainDBPath: plainDBPath,
		encDBPath:   encPath,
		secretKey:   defaultAppSecret,
	}
}

// EncryptedPath returns the path to the encrypted database file at rest.
func (v *VaultManager) EncryptedPath() string {
	return v.encDBPath
}

// PlainPath returns the path to the working database file.
func (v *VaultManager) PlainPath() string {
	return v.plainDBPath
}

// PrepareWorkingDatabase ensures the working database is ready before opening SQLite.
// If an encrypted vault exists, it decrypts it to the working file.
// If only a plain database exists, it encrypts it.
func (v *VaultManager) PrepareWorkingDatabase() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	encExists := fileExists(v.encDBPath)
	plainExists := fileExists(v.plainDBPath)

	if encExists && !plainExists {
		// Decrypt vault to working DB
		if err := v.decryptFile(v.encDBPath, v.plainDBPath); err != nil {
			return fmt.Errorf("failed to decrypt database vault: %w", err)
		}
		return nil
	}

	if encExists && plainExists {
		// Both exist: check which one is newer (in case of previous unexpected termination)
		encInfo, err1 := os.Stat(v.encDBPath)
		plainInfo, err2 := os.Stat(v.plainDBPath)
		if err1 == nil && err2 == nil {
			if plainInfo.ModTime().After(encInfo.ModTime()) {
				// Working file is newer, sync to encrypted vault
				_ = v.encryptFile(v.plainDBPath, v.encDBPath)
			} else {
				// Encrypted vault is newer, restore from vault
				_ = v.decryptFile(v.encDBPath, v.plainDBPath)
			}
		}
		return nil
	}

	if !encExists && plainExists {
		// First-time encryption migration: encrypt existing plain DB
		if err := v.encryptFile(v.plainDBPath, v.encDBPath); err != nil {
			return fmt.Errorf("failed to encrypt initial database: %w", err)
		}
		return nil
	}

	// Neither exists: new database initialization will create it
	return nil
}

// SyncEncrypted copies current database state to the encrypted vault.
func (v *VaultManager) SyncEncrypted() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !fileExists(v.plainDBPath) {
		return nil
	}

	return v.encryptFile(v.plainDBPath, v.encDBPath)
}

// FinalizeShutdown syncs the encrypted vault and securely removes the unencrypted working database.
func (v *VaultManager) FinalizeShutdown() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !fileExists(v.plainDBPath) {
		return nil
	}

	// Encrypt latest state
	if err := v.encryptFile(v.plainDBPath, v.encDBPath); err != nil {
		return fmt.Errorf("failed to encrypt database during shutdown: %w", err)
	}

	// Securely remove plain database and auxiliary WAL files
	_ = os.Remove(v.plainDBPath)
	_ = os.Remove(v.plainDBPath + "-wal")
	_ = os.Remove(v.plainDBPath + "-shm")

	return nil
}

// ExportDecrypted writes a clean, unencrypted SQLite database to destinationPath.
func (v *VaultManager) ExportDecrypted(destinationPath string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return err
	}

	// If working plain DB is present, copy directly
	if fileExists(v.plainDBPath) {
		return copyFile(v.plainDBPath, destinationPath)
	}

	// Otherwise decrypt directly from vault to destinationPath
	if fileExists(v.encDBPath) {
		return v.decryptFile(v.encDBPath, destinationPath)
	}

	return errors.New("no database file found to export")
}

// RestoreFromBackup replaces the database with the given backup file (either plain or encrypted).
func (v *VaultManager) RestoreFromBackup(backupPath string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	if !fileExists(backupPath) {
		return errors.New("backup file does not exist")
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	isEncrypted := strings.HasPrefix(string(data), VaultMagicHeader)

	if isEncrypted {
		// Write to encDBPath, and decrypt to plainDBPath
		if err := os.WriteFile(v.encDBPath, data, 0o600); err != nil {
			return err
		}
		return v.decryptFile(v.encDBPath, v.plainDBPath)
	}

	// Validate plain SQLite header
	if len(data) < 16 || !strings.HasPrefix(string(data[:16]), "SQLite format 3\x00") {
		return errors.New("invalid SQLite database format: unrecognized file header")
	}

	// Write to plainDBPath and re-encrypt
	if err := os.WriteFile(v.plainDBPath, data, 0o600); err != nil {
		return err
	}
	return v.encryptFile(v.plainDBPath, v.encDBPath)
}

// encryptFile encrypts src file to dst file using AES-256-GCM + PBKDF2
func (v *VaultManager) encryptFile(srcPath, dstPath string) error {
	plaintext, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read source file failed: %w", err)
	}

	// Generate 16-byte random salt
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("generate salt failed: %w", err)
	}

	// Derive 32-byte AES-256 key using PBKDF2-SHA256 (100,000 iterations)
	key := pbkdf2.Key([]byte(v.secretKey), salt, pbkdf2Iterations, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM failed: %w", err)
	}

	// Generate 12-byte random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("generate nonce failed: %w", err)
	}

	// Encrypt and seal
	ciphertext := gcm.Seal(nil, nonce, plaintext, []byte(VaultMagicHeader))

	// Assemble vault payload: Header (23 bytes) + Salt (16 bytes) + Nonce (12 bytes) + Ciphertext+Tag
	var out []byte
	out = append(out, []byte(VaultMagicHeader)...)
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	// Write atomically via temporary file
	tmpPath := dstPath + ".tmp"
	if err := os.WriteFile(tmpPath, out, 0o600); err != nil {
		return fmt.Errorf("write encrypted temp file failed: %w", err)
	}

	if err := os.Rename(tmpPath, dstPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("atomic rename encrypted file failed: %w", err)
	}

	return nil
}

// decryptFile decrypts an AES-256-GCM vault to dstPath
func (v *VaultManager) decryptFile(srcPath, dstPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read encrypted vault failed: %w", err)
	}

	magicLen := len(VaultMagicHeader)
	minLen := magicLen + saltLength + nonceLength
	if len(data) < minLen {
		return errors.New("encrypted database file is corrupted or too small")
	}

	if string(data[:magicLen]) != VaultMagicHeader {
		return errors.New("file is not a valid SMS encrypted database vault")
	}

	salt := data[magicLen : magicLen+saltLength]
	nonce := data[magicLen+saltLength : magicLen+saltLength+nonceLength]
	ciphertext := data[magicLen+saltLength+nonceLength:]

	// Derive key
	key := pbkdf2.Key([]byte(v.secretKey), salt, pbkdf2Iterations, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher failed: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM failed: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, []byte(VaultMagicHeader))
	if err != nil {
		return fmt.Errorf("authentication/decryption failed (tampered data or invalid key): %w", err)
	}

	// Write decrypted database atomically
	tmpPath := dstPath + ".tmp"
	if err := os.WriteFile(tmpPath, plaintext, 0o600); err != nil {
		return fmt.Errorf("write decrypted temp file failed: %w", err)
	}

	if err := os.Rename(tmpPath, dstPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("atomic rename decrypted file failed: %w", err)
	}

	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
