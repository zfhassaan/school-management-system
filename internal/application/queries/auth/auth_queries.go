package auth

import (
	"context"

	"school-management/internal/application/commands/auth"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
)

type SetupStatusDTO struct {
	IsInitialized bool   `json:"isInitialized"`
	SchoolName    string `json:"schoolName,omitempty"`
}

type AuthQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewAuthQueriesHandler(uow uow.UnitOfWork) *AuthQueriesHandler {
	return &AuthQueriesHandler{uow: uow}
}

func (h *AuthQueriesHandler) CheckSetupStatus(ctx context.Context) (*SetupStatusDTO, error) {
	var status SetupStatusDTO
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		count, err := r.Users.Count(ctx)
		if err != nil {
			return err
		}
		status.IsInitialized = count > 0
		if status.IsInitialized {
			sch, err := r.School.Get(ctx)
			if err == nil && sch != nil {
				status.SchoolName = sch.Name
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (h *AuthQueriesHandler) GetCurrentUser(ctx context.Context) (*auth.LoginResult, error) {
	actor, err := authz.RequireAuthenticated(ctx)
	if err != nil {
		return nil, err
	}

	var res auth.LoginResult
	err = h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.Get(ctx, actor.UserID)
		if err != nil {
			return err
		}
		roles, err := r.Users.RoleCodes(ctx, u.ID)
		if err != nil {
			return err
		}
		perms, err := r.Users.EffectivePermissions(ctx, u.ID)
		if err != nil {
			return err
		}

		sch, err := r.School.Get(ctx)
		var schDTO auth.SchoolDTO
		if err == nil {
			settings, _ := r.School.GetSettings(ctx, sch.ID)
			dateFormat := "dd/MM/yyyy"
			if settings != nil {
				dateFormat = settings.DateFormat
			}
			schDTO = auth.SchoolDTO{
				ID:                 sch.ID.String(),
				Name:               sch.Name,
				Code:               sch.Code,
				Timezone:           sch.Timezone,
				CurrencyCode:       sch.CurrencyCode,
				CurrencyMinorUnits: sch.CurrencyMinorUnits,
				DateFormat:         dateFormat,
			}
		}

		permStrings := make([]string, len(perms))
		for i, p := range perms {
			permStrings[i] = string(p)
		}

		res.SessionToken = actor.SessionID
		res.User = auth.UserDTO{
			ID:                 u.ID.String(),
			Username:           u.Username,
			DisplayName:        u.DisplayName,
			Email:              u.Email,
			Roles:              roles,
			Permissions:        permStrings,
			MustChangePassword: u.MustChangePassword,
		}
		res.School = schDTO
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

