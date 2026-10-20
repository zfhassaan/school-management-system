package documents

import (
	"bytes"
	"context"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/documents"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type UploadDocumentInput struct {
	StudentID        *shared.ID `json:"studentId,omitempty"`
	ApplicationID    *shared.ID `json:"applicationId,omitempty"`
	DocumentTypeCode string     `json:"documentTypeCode"`
	FileName         string     `json:"fileName"`
	MimeType         string     `json:"mimeType"`
	Data             []byte     `json:"data"`
}

type DocumentCommandsHandler struct {
	uow       uow.UnitOfWork
	blobStore documents.BlobStore
	clock     shared.Clock
	ids       shared.IDGenerator
	recorder  auditlog.Recorder
}

func NewDocumentCommandsHandler(
	uow uow.UnitOfWork,
	blobStore documents.BlobStore,
	clock shared.Clock,
	ids shared.IDGenerator,
) *DocumentCommandsHandler {
	return &DocumentCommandsHandler{
		uow:       uow,
		blobStore: blobStore,
		clock:     clock,
		ids:       ids,
		recorder:  auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *DocumentCommandsHandler) UploadDocument(ctx context.Context, in UploadDocumentInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermDocumentsManage)
	if err != nil {
		return "", err
	}
	if len(in.Data) == 0 {
		return "", shared.Validation("file content cannot be empty", nil)
	}

	storageKey, sha256Hex, size, err := h.blobStore.Save(ctx, bytes.NewReader(in.Data))
	if err != nil {
		return "", shared.Internal(err)
	}

	id := h.ids.NewID()
	now := h.clock.Now()

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		doc := &documents.Document{
			ID:               id,
			StudentID:        in.StudentID,
			ApplicationID:    in.ApplicationID,
			DocumentTypeCode: strings.TrimSpace(in.DocumentTypeCode),
			FileName:         strings.TrimSpace(in.FileName),
			MimeType:         strings.TrimSpace(in.MimeType),
			SizeBytes:        size,
			SHA256:           sha256Hex,
			StorageKey:       storageKey,
			VersionNo:        1,
			UploadedBy:       actor.UserID,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		if err := r.Documents.Insert(ctx, doc); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "document.uploaded",
			EntityType: "document",
			EntityID:   id.String(),
			Reason:     "document uploaded",
			After: map[string]any{
				"fileName": doc.FileName,
				"size":     size,
			},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}
