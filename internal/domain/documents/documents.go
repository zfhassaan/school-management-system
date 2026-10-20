package documents

import (
	"context"
	"io"
	"time"

	"school-management/internal/domain/shared"
)

type Document struct {
	ID               shared.ID
	StudentID        *shared.ID
	ApplicationID    *shared.ID
	DocumentTypeCode string // e.g. "PHOTO", "BIRTH_CERT", "PREVIOUS_MARKS"
	FileName         string
	MimeType         string
	SizeBytes        int64
	SHA256           string
	StorageKey       string
	VersionNo        int
	UploadedBy       shared.ID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Revision         int64
}

// BlobStore abstracts local filesystem (now) and S3 object storage (Phase 7).
type BlobStore interface {
	Save(ctx context.Context, r io.Reader) (storageKey string, sha256Hex string, sizeBytes int64, err error)
	Open(ctx context.Context, storageKey string) (io.ReadCloser, error)
	Delete(ctx context.Context, storageKey string) error
}

type Repository interface {
	Get(ctx context.Context, id shared.ID) (*Document, error)
	Insert(ctx context.Context, d *Document) error
	ListForStudent(ctx context.Context, studentID shared.ID) ([]Document, error)
	ListForApplication(ctx context.Context, applicationID shared.ID) ([]Document, error)
}

