package wails

import (
	feeCmd "school-management/internal/application/commands/fees"
	feeQuery "school-management/internal/application/queries/fees"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/fees"
	"school-management/internal/domain/shared"
	"school-management/internal/transport/dto"
)

type FeeService struct {
	sessions    *SessionManager
	cmdHandler  *feeCmd.FeeCommandsHandler
	queryHandle *feeQuery.FeeQueriesHandler
}

func NewFeeService(
	sessions *SessionManager,
	cmdHandler *feeCmd.FeeCommandsHandler,
	queryHandle *feeQuery.FeeQueriesHandler,
) *FeeService {
	return &FeeService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

// --- Fee Heads ---

func (s *FeeService) ListFeeHeads() dto.Response[[]fees.FeeHead] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListFeeHeads(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]fees.FeeHead](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *FeeService) CreateFeeHead(in feeCmd.CreateFeeHeadInput) dto.Response[shared.ID] {
	ctx := s.sessions.Context()
	id, err := s.cmdHandler.CreateFeeHead(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *FeeService) UpdateFeeHead(in feeCmd.UpdateFeeHeadInput) dto.Response[bool] {
	ctx := s.sessions.Context()
	err := s.cmdHandler.UpdateFeeHead(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

// --- Fee Structures & Versions ---

func (s *FeeService) ListFeeStructures(academicYearID shared.ID) dto.Response[[]fees.FeeStructure] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListFeeStructures(ctx, academicYearID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]fees.FeeStructure](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *FeeService) CreateFeeStructure(in feeCmd.CreateFeeStructureInput) dto.Response[shared.ID] {
	ctx := s.sessions.Context()
	id, err := s.cmdHandler.CreateFeeStructure(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *FeeService) ListFeeStructureVersions(feeStructureID shared.ID) dto.Response[[]fees.FeeStructureVersion] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListFeeStructureVersions(ctx, feeStructureID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]fees.FeeStructureVersion](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *FeeService) GetFeeStructureVersionDetail(versionID shared.ID) dto.Response[*feeQuery.FeeStructureVersionDetail] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.GetFeeStructureVersionDetail(ctx, versionID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*feeQuery.FeeStructureVersionDetail](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *FeeService) CreateFeeStructureVersion(in feeCmd.CreateFeeStructureVersionInput) dto.Response[shared.ID] {
	ctx := s.sessions.Context()
	id, err := s.cmdHandler.CreateFeeStructureVersion(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *FeeService) PublishFeeStructureVersion(versionID shared.ID) dto.Response[bool] {
	ctx := s.sessions.Context()
	err := s.cmdHandler.PublishFeeStructureVersion(ctx, versionID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

// --- Discount Policies & Student Discounts ---

func (s *FeeService) ListDiscountPolicies() dto.Response[[]fees.DiscountPolicy] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListDiscountPolicies(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]fees.DiscountPolicy](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *FeeService) CreateDiscountPolicy(in feeCmd.CreateDiscountPolicyInput) dto.Response[shared.ID] {
	ctx := s.sessions.Context()
	id, err := s.cmdHandler.CreateDiscountPolicy(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *FeeService) ListStudentDiscounts(studentID shared.ID, academicYearID shared.ID) dto.Response[[]fees.StudentDiscount] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListStudentDiscounts(ctx, studentID, academicYearID)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]fees.StudentDiscount](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *FeeService) ListPendingDiscounts() dto.Response[[]fees.StudentDiscount] {
	ctx := s.sessions.Context()
	res, err := s.queryHandle.ListPendingDiscounts(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]fees.StudentDiscount](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *FeeService) ApplyStudentDiscount(in feeCmd.ApplyStudentDiscountInput) dto.Response[shared.ID] {
	ctx := s.sessions.Context()
	id, err := s.cmdHandler.ApplyStudentDiscount(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[shared.ID](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *FeeService) DecideStudentDiscount(in feeCmd.DecideStudentDiscountInput) dto.Response[bool] {
	ctx := s.sessions.Context()
	err := s.cmdHandler.DecideStudentDiscount(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

