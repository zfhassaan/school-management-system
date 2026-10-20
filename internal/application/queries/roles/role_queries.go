package roles

import (
	"context"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type RoleDTO struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	IsSystem    bool     `json:"isSystem"`
	Permissions []string `json:"permissions"`
	UserCount   int      `json:"userCount"`
}

type PermissionInfoDTO struct {
	Code        string `json:"code"`
	Module      string `json:"module"`
	Description string `json:"description"`
	Sensitive   bool   `json:"sensitive"`
}

type RoleQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewRoleQueriesHandler(uow uow.UnitOfWork) *RoleQueriesHandler {
	return &RoleQueriesHandler{uow: uow}
}

func (h *RoleQueriesHandler) ListRoles(ctx context.Context) ([]RoleDTO, error) {
	if _, err := authz.Require(ctx, identity.PermUsersView); err != nil {
		return nil, err
	}

	var results []RoleDTO
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		roles, err := r.Roles.List(ctx)
		if err != nil {
			return err
		}
		results = make([]RoleDTO, len(roles))
		for i, role := range roles {
			perms := make([]string, len(role.Permissions))
			for j, p := range role.Permissions {
				perms[j] = string(p)
			}
			results[i] = RoleDTO{
				ID:          role.ID.String(),
				Code:        role.Code,
				Name:        role.Name,
				Description: role.Description,
				IsSystem:    role.IsSystem,
				Permissions: perms,
				UserCount:   role.UserCount,
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (h *RoleQueriesHandler) GetRole(ctx context.Context, id string) (*RoleDTO, error) {
	if _, err := authz.Require(ctx, identity.PermUsersView); err != nil {
		return nil, err
	}
	rid, err := shared.ParseID(id)
	if err != nil {
		return nil, err
	}

	var result *RoleDTO
	err = h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		role, err := r.Roles.Get(ctx, rid)
		if err != nil {
			return err
		}
		perms := make([]string, len(role.Permissions))
		for j, p := range role.Permissions {
			perms[j] = string(p)
		}
		result = &RoleDTO{
			ID:          role.ID.String(),
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			IsSystem:    role.IsSystem,
			Permissions: perms,
			UserCount:   role.UserCount,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (h *RoleQueriesHandler) ListPermissions(ctx context.Context) ([]PermissionInfoDTO, error) {
	if _, err := authz.Require(ctx, identity.PermUsersView); err != nil {
		return nil, err
	}
	results := make([]PermissionInfoDTO, len(identity.Catalog))
	for i, p := range identity.Catalog {
		results[i] = PermissionInfoDTO{
			Code:        string(p.Code),
			Module:      p.Module,
			Description: p.Description,
			Sensitive:   p.Sensitive,
		}
	}
	return results, nil
}

