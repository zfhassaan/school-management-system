package wails

import (
	feeCmd "school-management/internal/application/commands/fees"
	feeQuery "school-management/internal/application/queries/fees"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/fees"
	"school-management/internal/domain/shared"
	"school-management/internal/transport/dto"
)

type InvoiceService struct {
	sessions    *SessionManager
	cmdHandler  *feeCmd.FeeCommandsHandler
	queryHandle *feeQuery.FeeQueriesHandler
}

func NewInvoiceService(
	sessions *SessionManager,
	cmdHandler *feeCmd.FeeCommandsHandler,
	queryHandle *feeQuery.FeeQueriesHandler,
) *InvoiceService {
	return &InvoiceService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *InvoiceService) ListInvoices(filter fees.InvoiceFilter) dto.Response[[]fees.Invoice] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListInvoices(ctx, filter)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]fees.Invoice](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *InvoiceService) GetInvoice(id shared.ID) dto.Response[*feeQuery.InvoiceDetail] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.GetInvoice(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*feeQuery.InvoiceDetail](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *InvoiceService) GenerateInvoicesBatch(in feeCmd.GenerateInvoicesBatchInput) dto.Response[[]shared.ID] {
	ctx := s.sessions.Context()
	ids, err := s.cmdHandler.GenerateInvoicesBatch(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(ids)
}

func (s *InvoiceService) VoidInvoice(in feeCmd.VoidInvoiceInput) dto.Response[bool] {
	ctx := s.sessions.Context()
	err := s.cmdHandler.VoidInvoice(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

