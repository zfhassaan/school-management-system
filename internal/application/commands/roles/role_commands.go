package roles

import (
	"context"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type CreateRoleInput struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type UpdateRoleInput struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

type RoleCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewRoleCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *RoleCommandsHandler {
	return &RoleCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *RoleCommandsHandler) CreateRole(ctx context.Context, in CreateRoleInput) (string, error) {
	if _, err := authz.Require(ctx, identity.PermRolesManage); err != nil {
		return "", err
	}

	code := strings.ToLower(strings.TrimSpace(in.Code))
	perms := make([]identity.Permission, len(in.Permissions))
	for i, p := range in.Permissions {
		perms[i] = identity.Permission(p)
	}

	var v shared.Validator
	identity.ValidateRole(&v, code, in.Name, perms)
	if err := v.Err(); err != nil {
		return "", err
	}

	newRoleID := h.ids.NewID()
	role := &identity.Role{
		ID:          newRoleID,
		Code:        code,
		Name:        strings.TrimSpace(in.Name),
		Description: strings.TrimSpace(in.Description),
		IsSystem:    false,
		Permissions: perms,
	}

	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		if _, err := r.Roles.GetByCode(ctx, code); err == nil {
			return shared.Conflict("role with this code already exists")
		}

		if err := r.Roles.Insert(ctx, role, now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "role.created",
			EntityType: "role",
			EntityID:   newRoleID.String(),
			Reason:     "new custom role created",
			After: map[string]any{
				"code": role.Code,
				"name": role.Name,
			},
		})
	})

	if err != nil {
		return "", err
	}
	return newRoleID.String(), nil
}

func (h *RoleCommandsHandler) UpdateRole(ctx context.Context, in UpdateRoleInput) error {
	if _, err := authz.Require(ctx, identity.PermRolesManage); err != nil {
		return err
	}

	rid, err := shared.ParseID(in.ID)
	if err != nil {
		return err
	}

	perms := make([]identity.Permission, len(in.Permissions))
	for i, p := range in.Permissions {
		perms[i] = identity.Permission(p)
	}

	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		existing, err := r.Roles.Get(ctx, rid)
		if err != nil {
			return err
		}

		if existing.IsSystem && existing.Code == identity.RoleSuperAdmin {
			return shared.Forbidden("super_admin role permissions cannot be modified")
		}

		existing.Name = strings.TrimSpace(in.Name)
		existing.Description = strings.TrimSpace(in.Description)
		existing.Permissions = perms

		if err := r.Roles.Update(ctx, existing, now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "role.updated",
			EntityType: "role",
			EntityID:   rid.String(),
			Reason:     "role updated",
		})
	})
}

func (h *RoleCommandsHandler) DeleteRole(ctx context.Context, id string) error {
	if _, err := authz.Require(ctx, identity.PermRolesManage); err != nil {
		return err
	}

	rid, err := shared.ParseID(id)
	if err != nil {
		return err
	}

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		role, err := r.Roles.Get(ctx, rid)
		if err != nil {
			return err
		}
		if role.IsSystem {
			return shared.Forbidden("system roles cannot be deleted")
		}
		if role.UserCount > 0 {
			return shared.Conflict("cannot delete role that is assigned to users")
		}

		if err := r.Roles.Delete(ctx, rid); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "role.deleted",
			EntityType: "role",
			EntityID:   rid.String(),
			Reason:     "custom role deleted",
		})
	})
}

