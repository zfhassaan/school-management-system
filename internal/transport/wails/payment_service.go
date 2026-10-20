package wails

import (
	"context"

	paymentCmd "school-management/internal/application/commands/payments"
	paymentQuery "school-management/internal/application/queries/payments"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/payments"
	"school-management/internal/domain/shared"
	"school-management/internal/transport/dto"
)

type PaymentService struct {
	sessions    *SessionManager
	cmdHandler  *paymentCmd.PaymentCommandsHandler
	queryHandle *paymentQuery.PaymentQueriesHandler
}

func NewPaymentService(
	sessions *SessionManager,
	cmdHandler *paymentCmd.PaymentCommandsHandler,
	queryHandle *paymentQuery.PaymentQueriesHandler,
) *PaymentService {
	return &PaymentService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *PaymentService) RecordPayment(ctx context.Context, in paymentCmd.RecordPaymentInput) dto.Response[*paymentCmd.RecordPaymentResult] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.cmdHandler.RecordPayment(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*paymentCmd.RecordPaymentResult](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *PaymentService) ReversePayment(ctx context.Context, in paymentCmd.ReversePaymentInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.ReversePayment(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *PaymentService) GetPayment(ctx context.Context, id shared.ID) dto.Response[*paymentQuery.PaymentDetail] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetPayment(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*paymentQuery.PaymentDetail](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *PaymentService) ListPayments(ctx context.Context, filter payments.PaymentFilter) dto.Response[[]*payments.PaymentWithDetails] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListPayments(ctx, filter)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]*payments.PaymentWithDetails](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *PaymentService) GetStudentCreditBalance(ctx context.Context, studentID shared.ID) dto.Response[*paymentQuery.StudentCreditBalanceInfo] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetStudentCreditBalance(ctx, studentID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*paymentQuery.StudentCreditBalanceInfo](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *PaymentService) RequestFinancialAdjustment(ctx context.Context, in paymentCmd.RequestFinancialAdjustmentInput) dto.Response[shared.ID] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.RequestFinancialAdjustment(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *PaymentService) DecideFinancialAdjustment(ctx context.Context, in paymentCmd.DecideFinancialAdjustmentInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.DecideFinancialAdjustment(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *PaymentService) ListFinancialAdjustments(ctx context.Context, filter payments.AdjustmentFilter) dto.Response[[]*payments.FinancialAdjustment] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListFinancialAdjustments(ctx, filter)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]*payments.FinancialAdjustment](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *PaymentService) GetFinancialAdjustment(ctx context.Context, id shared.ID) dto.Response[*payments.FinancialAdjustment] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetFinancialAdjustment(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*payments.FinancialAdjustment](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}
