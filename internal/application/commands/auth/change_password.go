package auth

import (
	"context"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordHandler struct {
	uow      uow.UnitOfWork
	hasher   identity.PasswordHasher
	clock    shared.Clock
	recorder auditlog.Recorder
}

func NewChangePasswordHandler(
	uow uow.UnitOfWork,
	hasher identity.PasswordHasher,
	clock shared.Clock,
	ids shared.IDGenerator,
) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		uow:      uow,
		hasher:   hasher,
		clock:    clock,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *ChangePasswordHandler) Execute(ctx context.Context, in ChangePasswordInput) error {
	actor, err := authz.RequireAuthenticated(ctx)
	if err != nil {
		return err
	}

	var v shared.Validator
	identity.DefaultPasswordPolicy.Validate(&v, "newPassword", in.NewPassword, actor.Username)
	if err := v.Err(); err != nil {
		return err
	}

	newHash, err := h.hasher.Hash(in.NewPassword)
	if err != nil {
		return shared.Internal(err)
	}

	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.Get(ctx, actor.UserID)
		if err != nil {
			return err
		}

		// Verify old password
		match, _, err := h.hasher.Verify(in.OldPassword, u.PasswordHash)
		if err != nil || !match {
			return shared.Forbidden("incorrect current password")
		}

		u.PasswordHash = newHash
		u.PasswordChangedAt = now
		u.MustChangePassword = false
		u.UpdatedAt = now

		if err := r.Users.Update(ctx, u); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "user.password_changed",
			EntityType: "user",
			EntityID:   u.ID.String(),
			Reason:     "password changed by user",
		})
	})
}

