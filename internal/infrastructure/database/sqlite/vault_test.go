package sqlite

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestVaultManagerEncryptionRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	plainPath := filepath.Join(tmpDir, "test.db")
	vault := NewVaultManager(plainPath)

	// Create dummy sqlite content
	sampleContent := []byte("SQLite format 3\x00DummyDatabaseContentForTesting1234567890")
	if err := os.WriteFile(plainPath, sampleContent, 0o600); err != nil {
		t.Fatalf("failed to write plain file: %v", err)
	}

	// Encrypt
	if err := vault.SyncEncrypted(); err != nil {
		t.Fatalf("failed to encrypt: %v", err)
	}

	// Check encrypted file exists
	encPath := vault.EncryptedPath()
	encData, err := os.ReadFile(encPath)
	if err != nil {
		t.Fatalf("failed to read encrypted file: %v", err)
	}

	// Must NOT contain "SQLite format 3"
	if bytes.Contains(encData, []byte("SQLite format 3")) {
		t.Fatalf("encrypted file leaked sqlite header in plaintext!")
	}

	// Must contain magic vault header
	if !bytes.HasPrefix(encData, []byte(VaultMagicHeader)) {
		t.Fatalf("encrypted file missing vault header!")
	}

	// Remove plain file and decrypt via PrepareWorkingDatabase
	_ = os.Remove(plainPath)
	if err := vault.PrepareWorkingDatabase(); err != nil {
		t.Fatalf("failed to decrypt working db: %v", err)
	}

	decryptedData, err := os.ReadFile(plainPath)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(decryptedData, sampleContent) {
		t.Fatalf("decrypted data does not match original! Got %s", string(decryptedData))
	}

	// Test ExportDecrypted
	exportPath := filepath.Join(tmpDir, "export.sqlite")
	if err := vault.ExportDecrypted(exportPath); err != nil {
		t.Fatalf("failed to export decrypted: %v", err)
	}

	exportedData, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("failed to read exported file: %v", err)
	}
	if !bytes.Equal(exportedData, sampleContent) {
		t.Fatalf("exported data does not match original!")
	}
}

