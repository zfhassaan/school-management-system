package documents

import (
	"context"
	"io"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/documents"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type DocumentQueriesHandler struct {
	uow       uow.UnitOfWork
	blobStore documents.BlobStore
}

func NewDocumentQueriesHandler(uow uow.UnitOfWork, blobStore documents.BlobStore) *DocumentQueriesHandler {
	return &DocumentQueriesHandler{uow: uow, blobStore: blobStore}
}

func (h *DocumentQueriesHandler) ListDocumentsForStudent(ctx context.Context, studentID shared.ID) ([]documents.Document, error) {
	if _, err := authz.Require(ctx, identity.PermDocumentsView); err != nil {
		return nil, err
	}
	var res []documents.Document
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		docs, err := r.Documents.ListForStudent(ctx, studentID)
		if err != nil {
			return err
		}
		res = docs
		return nil
	})
	return res, err
}

func (h *DocumentQueriesHandler) ListDocumentsForApplication(ctx context.Context, appID shared.ID) ([]documents.Document, error) {
	if _, err := authz.Require(ctx, identity.PermDocumentsView); err != nil {
		return nil, err
	}
	var res []documents.Document
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		docs, err := r.Documents.ListForApplication(ctx, appID)
		if err != nil {
			return err
		}
		res = docs
		return nil
	})
	return res, err
}

func (h *DocumentQueriesHandler) GetDocumentContent(ctx context.Context, id shared.ID) ([]byte, *documents.Document, error) {
	if _, err := authz.Require(ctx, identity.PermDocumentsView); err != nil {
		return nil, nil, err
	}
	var doc *documents.Document
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		d, err := r.Documents.Get(ctx, id)
		if err != nil {
			return err
		}
		doc = d
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	rc, err := h.blobStore.Open(ctx, doc.StorageKey)
	if err != nil {
		return nil, nil, shared.Internal(err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, nil, shared.Internal(err)
	}
	return data, doc, nil
}
