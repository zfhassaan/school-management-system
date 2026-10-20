package wails

import (
	"encoding/base64"

	docCmd "school-management/internal/application/commands/documents"
	docQuery "school-management/internal/application/queries/documents"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/documents"
	"school-management/internal/domain/shared"
	"school-management/internal/transport/dto"
)

type DocumentService struct {
	sessions    *SessionManager
	cmdHandler  *docCmd.DocumentCommandsHandler
	queryHandle *docQuery.DocumentQueriesHandler
}

func NewDocumentService(
	sessions *SessionManager,
	cmdHandler *docCmd.DocumentCommandsHandler,
	queryHandle *docQuery.DocumentQueriesHandler,
) *DocumentService {
	return &DocumentService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

type UploadDocumentRequest struct {
	StudentID        *shared.ID `json:"studentId,omitempty"`
	ApplicationID    *shared.ID `json:"applicationId,omitempty"`
	DocumentTypeCode string     `json:"documentTypeCode"`
	FileName         string     `json:"fileName"`
	MimeType         string     `json:"mimeType"`
	DataBase64       string     `json:"dataBase64"`
}

type DocumentContentResponse struct {
	Document   *documents.Document `json:"document"`
	DataBase64 string              `json:"dataBase64"`
}

func (s *DocumentService) UploadDocument(in UploadDocumentRequest) dto.Response[shared.ID] {
	ctx := s.sessions.Context()
	data, err := base64.StdEncoding.DecodeString(in.DataBase64)
	if err != nil {
		code, msg, fields := MapError(shared.Validation("invalid base64 document content", nil))
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}

	id, err := s.cmdHandler.UploadDocument(ctx, docCmd.UploadDocumentInput{
		StudentID:        in.StudentID,
		ApplicationID:    in.ApplicationID,
		DocumentTypeCode: in.DocumentTypeCode,
		FileName:         in.FileName,
		MimeType:         in.MimeType,
		Data:             data,
	})
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *DocumentService) ListForStudent(studentID shared.ID) dto.Response[[]documents.Document] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListDocumentsForStudent(ctx, studentID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]documents.Document](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *DocumentService) ListForApplication(applicationID shared.ID) dto.Response[[]documents.Document] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListDocumentsForApplication(ctx, applicationID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]documents.Document](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *DocumentService) GetContent(id shared.ID) dto.Response[*DocumentContentResponse] {
	ctx := s.sessions.Context()
	data, doc, err := s.queryHandle.GetDocumentContent(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*DocumentContentResponse](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(&DocumentContentResponse{
		Document:   doc,
		DataBase64: base64.StdEncoding.EncodeToString(data),
	})
}

