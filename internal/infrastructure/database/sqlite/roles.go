package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type RoleRepository struct {
	q Querier
}

func (r *RoleRepository) Get(ctx context.Context, id shared.ID) (*identity.Role, error) {
	row, err := r.q.GetRole(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("role")
		}
		return nil, shared.Internal(err)
	}
	role := toDomainRole(row)
	perms, err := r.q.GetRolePermissions(ctx, id.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	role.Permissions = make([]identity.Permission, len(perms))
	for i, p := range perms {
		role.Permissions[i] = identity.Permission(p)
	}
	return role, nil
}

func (r *RoleRepository) GetByCode(ctx context.Context, code string) (*identity.Role, error) {
	row, err := r.q.GetRoleByCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("role")
		}
		return nil, shared.Internal(err)
	}
	role := toDomainRole(row)
	perms, err := r.q.GetRolePermissions(ctx, role.ID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	role.Permissions = make([]identity.Permission, len(perms))
	for i, p := range perms {
		role.Permissions[i] = identity.Permission(p)
	}
	return role, nil
}

func (r *RoleRepository) List(ctx context.Context) ([]identity.Role, error) {
	rows, err := r.q.ListRolesWithUserCount(ctx)
	if err != nil {
		return nil, shared.Internal(err)
	}
	roles := make([]identity.Role, len(rows))
	for i, row := range rows {
		role := identity.Role{
			ID:          shared.ID(row.ID),
			Code:        row.Code,
			Name:        row.Name,
			Description: row.Description.String,
			IsSystem:    row.IsSystem == 1,
			UserCount:   int(row.UserCount),
			Revision:    row.RowVersion,
		}
		perms, err := r.q.GetRolePermissions(ctx, row.ID)
		if err != nil {
			return nil, shared.Internal(err)
		}
		role.Permissions = make([]identity.Permission, len(perms))
		for j, p := range perms {
			role.Permissions[j] = identity.Permission(p)
		}
		roles[i] = role
	}
	return roles, nil
}

func (r *RoleRepository) Insert(ctx context.Context, role *identity.Role, at time.Time) error {
	now := at.UTC().Format(time.RFC3339Nano)
	err := r.q.InsertRole(ctx, InsertRoleParams{
		ID:          role.ID.String(),
		Code:        role.Code,
		Name:        role.Name,
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		IsSystem:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		RowVersion:  1,
	})
	if err != nil {
		return shared.Internal(err)
	}
	role.Revision = 1
	return r.SetPermissions(ctx, role.ID, role.Permissions)
}

func (r *RoleRepository) Update(ctx context.Context, role *identity.Role, at time.Time) error {
	now := at.UTC().Format(time.RFC3339Nano)
	err := r.q.UpdateRole(ctx, UpdateRoleParams{
		Name:        role.Name,
		Description: sql.NullString{String: role.Description, Valid: role.Description != ""},
		UpdatedAt:   now,
		ID:          role.ID.String(),
		RowVersion:  role.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	role.Revision++
	return r.SetPermissions(ctx, role.ID, role.Permissions)
}

func (r *RoleRepository) Delete(ctx context.Context, id shared.ID) error {
	return r.q.DeleteRole(ctx, id.String())
}

func (r *RoleRepository) SetPermissions(ctx context.Context, roleID shared.ID, perms []identity.Permission) error {
	if err := r.q.ClearRolePermissions(ctx, roleID.String()); err != nil {
		return shared.Internal(err)
	}
	for _, p := range perms {
		err := r.q.InsertRolePermission(ctx, InsertRolePermissionParams{
			RoleID:         roleID.String(),
			PermissionCode: string(p),
		})
		if err != nil {
			return shared.Internal(err)
		}
	}
	return nil
}

func (r *RoleRepository) UpsertPermissionCatalog(ctx context.Context, catalog []identity.PermissionInfo) error {
	for _, p := range catalog {
		sensitive := int64(0)
		if p.Sensitive {
			sensitive = 1
		}
		err := r.q.UpsertPermission(ctx, UpsertPermissionParams{
			Code:        string(p.Code),
			Module:      p.Module,
			Description: p.Description,
			IsSensitive: sensitive,
		})
		if err != nil {
			return shared.Internal(err)
		}
	}
	return nil
}

func toDomainRole(row Role) *identity.Role {
	return &identity.Role{
		ID:          shared.ID(row.ID),
		Code:        row.Code,
		Name:        row.Name,
		Description: row.Description.String,
		IsSystem:    row.IsSystem == 1,
		Revision:    row.RowVersion,
	}
}

