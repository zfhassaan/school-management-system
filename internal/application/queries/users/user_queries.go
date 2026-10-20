package users

import (
	"context"
	"time"

	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type UserListItemDTO struct {
	ID                 string    `json:"id"`
	Username           string    `json:"username"`
	DisplayName        string    `json:"displayName"`
	Email              string    `json:"email"`
	Status             string    `json:"status"`
	Roles              []string  `json:"roles"`
	RoleIDs            []string  `json:"roleIds"`
	MustChangePassword bool      `json:"mustChangePassword"`
	LastLoginAt        *string   `json:"lastLoginAt,omitempty"`
	CreatedAt          string    `json:"createdAt"`
}

type UserQueriesHandler struct {
	uow uow.UnitOfWork
}

func NewUserQueriesHandler(uow uow.UnitOfWork) *UserQueriesHandler {
	return &UserQueriesHandler{uow: uow}
}

func (h *UserQueriesHandler) ListUsers(ctx context.Context) ([]UserListItemDTO, error) {
	if _, err := authz.Require(ctx, identity.PermUsersView); err != nil {
		return nil, err
	}

	var results []UserListItemDTO
	err := h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		users, err := r.Users.List(ctx)
		if err != nil {
			return err
		}
		results = make([]UserListItemDTO, len(users))
		for i, u := range users {
			roleCodes, _ := r.Users.RoleCodes(ctx, u.ID)
			roleIDStrings := make([]string, len(u.RoleIDs))
			for j, rid := range u.RoleIDs {
				roleIDStrings[j] = rid.String()
			}

			var lastLogin *string
			if u.LastLoginAt != nil {
				s := u.LastLoginAt.Format(time.RFC3339)
				lastLogin = &s
			}

			results[i] = UserListItemDTO{
				ID:                 u.ID.String(),
				Username:           u.Username,
				DisplayName:        u.DisplayName,
				Email:              u.Email,
				Status:             string(u.Status),
				Roles:              roleCodes,
				RoleIDs:            roleIDStrings,
				MustChangePassword: u.MustChangePassword,
				LastLoginAt:        lastLogin,
				CreatedAt:          u.CreatedAt.Format(time.RFC3339),
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (h *UserQueriesHandler) GetUser(ctx context.Context, id string) (*UserListItemDTO, error) {
	if _, err := authz.Require(ctx, identity.PermUsersView); err != nil {
		return nil, err
	}
	uid, err := shared.ParseID(id)
	if err != nil {
		return nil, err
	}

	var result *UserListItemDTO
	err = h.uow.Read(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.Get(ctx, uid)
		if err != nil {
			return err
		}
		roleCodes, _ := r.Users.RoleCodes(ctx, u.ID)
		roleIDStrings := make([]string, len(u.RoleIDs))
		for j, rid := range u.RoleIDs {
			roleIDStrings[j] = rid.String()
		}

		var lastLogin *string
		if u.LastLoginAt != nil {
			s := u.LastLoginAt.Format(time.RFC3339)
			lastLogin = &s
		}

		result = &UserListItemDTO{
			ID:                 u.ID.String(),
			Username:           u.Username,
			DisplayName:        u.DisplayName,
			Email:              u.Email,
			Status:             string(u.Status),
			Roles:              roleCodes,
			RoleIDs:            roleIDStrings,
			MustChangePassword: u.MustChangePassword,
			LastLoginAt:        lastLogin,
			CreatedAt:          u.CreatedAt.Format(time.RFC3339),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

