package users

import (
	"context"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type CreateUserInput struct {
	Username    string   `json:"username"`
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	RoleIDs     []string `json:"roleIds"`
}

type UpdateUserInput struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email"`
	RoleIDs     []string `json:"roleIds"`
}

type AdminResetPasswordInput struct {
	UserID      string `json:"userId"`
	NewPassword string `json:"newPassword"`
}

type UserCommandsHandler struct {
	uow      uow.UnitOfWork
	hasher   identity.PasswordHasher
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewUserCommandsHandler(
	uow uow.UnitOfWork,
	hasher identity.PasswordHasher,
	clock shared.Clock,
	ids shared.IDGenerator,
) *UserCommandsHandler {
	return &UserCommandsHandler{
		uow:      uow,
		hasher:   hasher,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *UserCommandsHandler) CreateUser(ctx context.Context, in CreateUserInput) (string, error) {
	if _, err := authz.Require(ctx, identity.PermUsersManage); err != nil {
		return "", err
	}

	var v shared.Validator
	username := identity.NormalizeUsername(in.Username)
	identity.ValidateProfile(&v, username, in.DisplayName, in.Email)
	identity.DefaultPasswordPolicy.Validate(&v, "password", in.Password, username)
	if err := v.Err(); err != nil {
		return "", err
	}

	pwHash, err := h.hasher.Hash(in.Password)
	if err != nil {
		return "", shared.Internal(err)
	}

	now := h.clock.Now()
	newUserID := h.ids.NewID()
	user := &identity.User{
		ID:                 newUserID,
		Username:           username,
		DisplayName:        strings.TrimSpace(in.DisplayName),
		Email:              strings.TrimSpace(in.Email),
		PasswordHash:       pwHash,
		PasswordChangedAt:  now,
		MustChangePassword: true, // New users must change password upon first login
		Status:             identity.UserActive,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	roleIDs := make([]shared.ID, len(in.RoleIDs))
	for i, rid := range in.RoleIDs {
		id, err := shared.ParseID(rid)
		if err != nil {
			return "", err
		}
		roleIDs[i] = id
	}

	actor, _ := authz.ActorFrom(ctx)

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		// Ensure username is unique
		if _, err := r.Users.GetByUsername(ctx, username); err == nil {
			return shared.Conflict("username already exists")
		}

		if err := r.Users.Insert(ctx, user); err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			if err := r.Users.SetRoles(ctx, newUserID, roleIDs, actor.UserIDPtr(), now); err != nil {
				return err
			}
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "user.created",
			EntityType: "user",
			EntityID:   newUserID.String(),
			Reason:     "new user account created",
			After: map[string]any{
				"username":    username,
				"displayName": user.DisplayName,
			},
		})
	})

	if err != nil {
		return "", err
	}
	return newUserID.String(), nil
}

func (h *UserCommandsHandler) UpdateUser(ctx context.Context, in UpdateUserInput) error {
	if _, err := authz.Require(ctx, identity.PermUsersManage); err != nil {
		return err
	}

	uid, err := shared.ParseID(in.ID)
	if err != nil {
		return err
	}

	var v shared.Validator
	dn := strings.TrimSpace(in.DisplayName)
	v.Check(dn != "", "displayName", "displayName is required")
	if err := v.Err(); err != nil {
		return err
	}

	now := h.clock.Now()
	actor, _ := authz.ActorFrom(ctx)

	roleIDs := make([]shared.ID, len(in.RoleIDs))
	for i, rid := range in.RoleIDs {
		id, err := shared.ParseID(rid)
		if err != nil {
			return err
		}
		roleIDs[i] = id
	}

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.Get(ctx, uid)
		if err != nil {
			return err
		}

		before := map[string]any{
			"displayName": u.DisplayName,
			"email":       u.Email,
		}

		u.DisplayName = dn
		u.Email = strings.TrimSpace(in.Email)
		u.UpdatedAt = now

		if err := r.Users.Update(ctx, u); err != nil {
			return err
		}

		if err := r.Users.SetRoles(ctx, uid, roleIDs, actor.UserIDPtr(), now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "user.updated",
			EntityType: "user",
			EntityID:   uid.String(),
			Reason:     "user profile updated",
			Before:     before,
			After: map[string]any{
				"displayName": u.DisplayName,
				"email":       u.Email,
			},
		})
	})
}

func (h *UserCommandsHandler) SetUserStatus(ctx context.Context, id string, active bool) error {
	if _, err := authz.Require(ctx, identity.PermUsersManage); err != nil {
		return err
	}
	uid, err := shared.ParseID(id)
	if err != nil {
		return err
	}

	actor, _ := authz.ActorFrom(ctx)
	if actor.UserID == uid && !active {
		return shared.Forbidden("you cannot disable your own account")
	}

	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.Get(ctx, uid)
		if err != nil {
			return err
		}

		// Protect last super admin from being disabled
		if !active {
			count, err := r.Users.CountActiveWithRole(ctx, identity.RoleSuperAdmin)
			if err != nil {
				return err
			}
			roles, _ := r.Users.RoleCodes(ctx, uid)
			isSuperAdmin := false
			for _, rc := range roles {
				if rc == identity.RoleSuperAdmin {
					isSuperAdmin = true
					break
				}
			}
			if isSuperAdmin && count <= 1 {
				return shared.Forbidden("cannot disable the only active Super Administrator")
			}
		}

		if active {
			u.Status = identity.UserActive
			u.Unlock()
		} else {
			u.Status = identity.UserDisabled
		}
		u.UpdatedAt = now

		if err := r.Users.Update(ctx, u); err != nil {
			return err
		}

		action := "user.disabled"
		if active {
			action = "user.enabled"
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     action,
			EntityType: "user",
			EntityID:   uid.String(),
			Reason:     "administrative status change",
		})
	})
}

func (h *UserCommandsHandler) AdminResetPassword(ctx context.Context, in AdminResetPasswordInput) error {
	if _, err := authz.Require(ctx, identity.PermUsersManage); err != nil {
		return err
	}
	uid, err := shared.ParseID(in.UserID)
	if err != nil {
		return err
	}

	var v shared.Validator
	identity.DefaultPasswordPolicy.Validate(&v, "newPassword", in.NewPassword, "")
	if err := v.Err(); err != nil {
		return err
	}

	pwHash, err := h.hasher.Hash(in.NewPassword)
	if err != nil {
		return shared.Internal(err)
	}

	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.Get(ctx, uid)
		if err != nil {
			return err
		}

		u.PasswordHash = pwHash
		u.PasswordChangedAt = now
		u.MustChangePassword = true
		u.Unlock()
		u.UpdatedAt = now

		if err := r.Users.Update(ctx, u); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "user.password_reset_by_admin",
			EntityType: "user",
			EntityID:   uid.String(),
			Reason:     "administrator reset password",
		})
	})
}

