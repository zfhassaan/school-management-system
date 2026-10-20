package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/documents"
	"school-management/internal/domain/shared"
)

type DocumentsRepository struct {
	q Querier
}

func (r *DocumentsRepository) Get(ctx context.Context, id shared.ID) (*documents.Document, error) {
	row, err := r.q.GetDocument(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("document")
		}
		return nil, shared.Internal(err)
	}
	return toDomainDoc(row), nil
}

func (r *DocumentsRepository) Insert(ctx context.Context, d *documents.Document) error {
	now := d.CreatedAt.UTC().Format(time.RFC3339Nano)
	var studentID sql.NullString
	if d.StudentID != nil {
		studentID = sql.NullString{String: d.StudentID.String(), Valid: true}
	}
	var appID sql.NullString
	if d.ApplicationID != nil {
		appID = sql.NullString{String: d.ApplicationID.String(), Valid: true}
	}

	err := r.q.InsertStudentDocument(ctx, InsertStudentDocumentParams{
		ID:               d.ID.String(),
		StudentID:        studentID,
		ApplicationID:    appID,
		DocumentTypeCode: d.DocumentTypeCode,
		FileName:         d.FileName,
		MimeType:         d.MimeType,
		SizeBytes:        d.SizeBytes,
		Sha256:           d.SHA256,
		StorageKey:       d.StorageKey,
		UploadedBy:       d.UploadedBy.String(),
		CreatedAt:        now,
		UpdatedAt:        now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	d.Revision = 1
	return nil
}

func (r *DocumentsRepository) ListForStudent(ctx context.Context, studentID shared.ID) ([]documents.Document, error) {
	rows, err := r.q.ListStudentDocuments(ctx, sql.NullString{String: studentID.String(), Valid: true})
	if err != nil {
		return nil, shared.Internal(err)
	}
	docs := make([]documents.Document, len(rows))
	for i, row := range rows {
		docs[i] = *toDomainDoc(row)
	}
	return docs, nil
}

func (r *DocumentsRepository) ListForApplication(ctx context.Context, applicationID shared.ID) ([]documents.Document, error) {
	rows, err := r.q.ListApplicationDocuments(ctx, sql.NullString{String: applicationID.String(), Valid: true})
	if err != nil {
		return nil, shared.Internal(err)
	}
	docs := make([]documents.Document, len(rows))
	for i, row := range rows {
		docs[i] = *toDomainDoc(row)
	}
	return docs, nil
}

func toDomainDoc(row StudentDocument) *documents.Document {
	var studentID *shared.ID
	if row.StudentID.Valid {
		id := shared.ID(row.StudentID.String)
		studentID = &id
	}
	var appID *shared.ID
	if row.ApplicationID.Valid {
		id := shared.ID(row.ApplicationID.String)
		appID = &id
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

	return &documents.Document{
		ID:               shared.ID(row.ID),
		StudentID:        studentID,
		ApplicationID:    appID,
		DocumentTypeCode: row.DocumentTypeCode,
		FileName:         row.FileName,
		MimeType:         row.MimeType,
		SizeBytes:        row.SizeBytes,
		SHA256:           row.Sha256,
		StorageKey:       row.StorageKey,
		VersionNo:        int(row.VersionNo),
		UploadedBy:       shared.ID(row.UploadedBy),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Revision:         row.RowVersion,
	}
}
