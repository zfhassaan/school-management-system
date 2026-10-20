package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FilesystemBlobStore struct {
	baseDir string
}

func NewFilesystemBlobStore(baseDir string) (*FilesystemBlobStore, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create blobstore base directory: %w", err)
	}
	return &FilesystemBlobStore{baseDir: baseDir}, nil
}

func (s *FilesystemBlobStore) Save(ctx context.Context, r io.Reader) (storageKey string, sha256Hex string, sizeBytes int64, err error) {
	// Create temporary file
	tmpFile, err := os.CreateTemp(s.baseDir, "upload-*")
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath) // clean up if not renamed
	}()

	hasher := sha256.New()
	mw := io.MultiWriter(tmpFile, hasher)

	n, err := io.Copy(mw, r)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to write stream: %w", err)
	}
	_ = tmpFile.Close()

	hashBytes := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)

	// Content-addressed sharding: ab/cd/<hash>
	subDir := filepath.Join(s.baseDir, hashHex[:2], hashHex[2:4])
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		return "", "", 0, fmt.Errorf("failed to create shard directory: %w", err)
	}

	finalKey := filepath.Join(hashHex[:2], hashHex[2:4], hashHex)
	finalPath := filepath.Join(s.baseDir, finalKey)

	// Atomic rename or replace
	if err := os.Rename(tmpPath, finalPath); err != nil {
		// If file already exists, it is identical content
		if !os.IsExist(err) {
			return "", "", 0, fmt.Errorf("failed to move blob to final path: %w", err)
		}
	}

	return finalKey, hashHex, n, nil
}

func (s *FilesystemBlobStore) Open(ctx context.Context, storageKey string) (io.ReadCloser, error) {
	cleanPath := filepath.Clean(filepath.Join(s.baseDir, storageKey))
	return os.Open(cleanPath)
}

func (s *FilesystemBlobStore) Delete(ctx context.Context, storageKey string) error {
	cleanPath := filepath.Clean(filepath.Join(s.baseDir, storageKey))
	return os.Remove(cleanPath)
}
