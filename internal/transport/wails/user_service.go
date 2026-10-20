package wails

import (
	"context"

	userCmd "school-management/internal/application/commands/users"
	userQuery "school-management/internal/application/queries/users"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/transport/dto"
)

type UserService struct {
	sessions    *SessionManager
	cmdHandler  *userCmd.UserCommandsHandler
	queryHandle *userQuery.UserQueriesHandler
}

func NewUserService(
	sessions *SessionManager,
	cmdHandler *userCmd.UserCommandsHandler,
	queryHandle *userQuery.UserQueriesHandler,
) *UserService {
	return &UserService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *UserService) ListUsers(ctx context.Context) dto.Response[[]userQuery.UserListItemDTO] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListUsers(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]userQuery.UserListItemDTO](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *UserService) GetUser(ctx context.Context, id string) dto.Response[*userQuery.UserListItemDTO] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetUser(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*userQuery.UserListItemDTO](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *UserService) CreateUser(ctx context.Context, in userCmd.CreateUserInput) dto.Response[string] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateUser(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[string](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *UserService) UpdateUser(ctx context.Context, in userCmd.UpdateUserInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateUser(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *UserService) SetUserStatus(ctx context.Context, id string, active bool) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.SetUserStatus(ctx, id, active)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *UserService) AdminResetPassword(ctx context.Context, in userCmd.AdminResetPasswordInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.AdminResetPassword(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

