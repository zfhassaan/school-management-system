package wails

import (
	"context"

	roleCmd "school-management/internal/application/commands/roles"
	roleQuery "school-management/internal/application/queries/roles"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/transport/dto"
)

type RoleService struct {
	sessions    *SessionManager
	cmdHandler  *roleCmd.RoleCommandsHandler
	queryHandle *roleQuery.RoleQueriesHandler
}

func NewRoleService(
	sessions *SessionManager,
	cmdHandler *roleCmd.RoleCommandsHandler,
	queryHandle *roleQuery.RoleQueriesHandler,
) *RoleService {
	return &RoleService{
		sessions:    sessions,
		cmdHandler:  cmdHandler,
		queryHandle: queryHandle,
	}
}

func (s *RoleService) ListRoles(ctx context.Context) dto.Response[[]roleQuery.RoleDTO] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListRoles(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]roleQuery.RoleDTO](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *RoleService) GetRole(ctx context.Context, id string) dto.Response[*roleQuery.RoleDTO] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.GetRole(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*roleQuery.RoleDTO](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *RoleService) ListPermissions(ctx context.Context) dto.Response[[]roleQuery.PermissionInfoDTO] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.queryHandle.ListPermissions(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[[]roleQuery.PermissionInfoDTO](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *RoleService) CreateRole(ctx context.Context, in roleCmd.CreateRoleInput) dto.Response[string] {
	ctx = s.sessions.WithContext(ctx)
	id, err := s.cmdHandler.CreateRole(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[string](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(id)
}

func (s *RoleService) UpdateRole(ctx context.Context, in roleCmd.UpdateRoleInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.UpdateRole(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.cmdHandler.DeleteRole(ctx, id)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(true)
}

