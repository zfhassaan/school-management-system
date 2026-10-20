package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type UserRepository struct {
	q Querier
}

func (r *UserRepository) Get(ctx context.Context, id shared.ID) (*identity.User, error) {
	row, err := r.q.GetUser(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("user")
		}
		return nil, shared.Internal(err)
	}
	u := toDomainUser(row)
	roleIDs, err := r.getUserRoleIDs(ctx, id)
	if err != nil {
		return nil, err
	}
	u.RoleIDs = roleIDs
	return u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*identity.User, error) {
	row, err := r.q.GetUserByUsername(ctx, identity.NormalizeUsername(username))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("user")
		}
		return nil, shared.Internal(err)
	}
	u := toDomainUser(row)
	roleIDs, err := r.getUserRoleIDs(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.RoleIDs = roleIDs
	return u, nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.q.CountUsers(ctx)
	if err != nil {
		return 0, shared.Internal(err)
	}
	return count, nil
}

func (r *UserRepository) List(ctx context.Context) ([]identity.User, error) {
	rows, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, shared.Internal(err)
	}
	users := make([]identity.User, len(rows))
	for i, row := range rows {
		u := toDomainUser(row)
		roleIDs, err := r.getUserRoleIDs(ctx, u.ID)
		if err != nil {
			return nil, err
		}
		u.RoleIDs = roleIDs
		users[i] = *u
	}
	return users, nil
}

func (r *UserRepository) Insert(ctx context.Context, u *identity.User) error {
	var lockedUntil sql.NullString
	if u.LockedUntil != nil {
		lockedUntil = sql.NullString{String: u.LockedUntil.UTC().Format(time.RFC3339Nano), Valid: true}
	}
	var lastLogin sql.NullString
	if u.LastLoginAt != nil {
		lastLogin = sql.NullString{String: u.LastLoginAt.UTC().Format(time.RFC3339Nano), Valid: true}
	}
	var email sql.NullString
	if u.Email != "" {
		email = sql.NullString{String: u.Email, Valid: true}
	}

	mustChange := int64(0)
	if u.MustChangePassword {
		mustChange = 1
	}

	err := r.q.InsertUser(ctx, InsertUserParams{
		ID:                 u.ID.String(),
		Username:           identity.NormalizeUsername(u.Username),
		DisplayName:        u.DisplayName,
		Email:              email,
		PasswordHash:       u.PasswordHash,
		PasswordChangedAt:  u.PasswordChangedAt.UTC().Format(time.RFC3339Nano),
		MustChangePassword: mustChange,
		Status:             string(u.Status),
		FailedLoginCount:   int64(u.FailedLoginCount),
		LockoutCount:       int64(u.LockoutCount),
		LockedUntil:        lockedUntil,
		LastLoginAt:        lastLogin,
		CreatedAt:          u.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:          u.UpdatedAt.UTC().Format(time.RFC3339Nano),
		RowVersion:         1,
	})
	if err != nil {
		return shared.Internal(err)
	}
	u.Revision = 1
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *identity.User) error {
	var lockedUntil sql.NullString
	if u.LockedUntil != nil {
		lockedUntil = sql.NullString{String: u.LockedUntil.UTC().Format(time.RFC3339Nano), Valid: true}
	}
	var lastLogin sql.NullString
	if u.LastLoginAt != nil {
		lastLogin = sql.NullString{String: u.LastLoginAt.UTC().Format(time.RFC3339Nano), Valid: true}
	}
	var email sql.NullString
	if u.Email != "" {
		email = sql.NullString{String: u.Email, Valid: true}
	}

	mustChange := int64(0)
	if u.MustChangePassword {
		mustChange = 1
	}

	err := r.q.UpdateUser(ctx, UpdateUserParams{
		Username:           identity.NormalizeUsername(u.Username),
		DisplayName:        u.DisplayName,
		Email:              email,
		PasswordHash:       u.PasswordHash,
		PasswordChangedAt:  u.PasswordChangedAt.UTC().Format(time.RFC3339Nano),
		MustChangePassword: mustChange,
		Status:             string(u.Status),
		FailedLoginCount:   int64(u.FailedLoginCount),
		LockoutCount:       int64(u.LockoutCount),
		LockedUntil:        lockedUntil,
		LastLoginAt:        lastLogin,
		UpdatedAt:          u.UpdatedAt.UTC().Format(time.RFC3339Nano),
		ID:                 u.ID.String(),
		RowVersion:         u.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	u.Revision++
	return nil
}

func (r *UserRepository) SetRoles(ctx context.Context, userID shared.ID, roleIDs []shared.ID, assignedBy *shared.ID, at time.Time) error {
	if err := r.q.ClearUserRoles(ctx, userID.String()); err != nil {
		return shared.Internal(err)
	}
	var assigner sql.NullString
	if assignedBy != nil {
		assigner = sql.NullString{String: assignedBy.String(), Valid: true}
	}
	assignedAt := at.UTC().Format(time.RFC3339Nano)
	for _, rid := range roleIDs {
		err := r.q.InsertUserRole(ctx, InsertUserRoleParams{
			UserID:     userID.String(),
			RoleID:     rid.String(),
			AssignedBy: assigner,
			AssignedAt: assignedAt,
		})
		if err != nil {
			return shared.Internal(err)
		}
	}
	return nil
}

func (r *UserRepository) EffectivePermissions(ctx context.Context, userID shared.ID) ([]identity.Permission, error) {
	codes, err := r.q.GetEffectivePermissions(ctx, userID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	perms := make([]identity.Permission, len(codes))
	for i, c := range codes {
		perms[i] = identity.Permission(c)
	}
	return perms, nil
}

func (r *UserRepository) RoleCodes(ctx context.Context, userID shared.ID) ([]string, error) {
	codes, err := r.q.GetUserRoleCodes(ctx, userID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	return codes, nil
}

func (r *UserRepository) CountActiveWithRole(ctx context.Context, roleCode string) (int64, error) {
	count, err := r.q.CountActiveWithRole(ctx, roleCode)
	if err != nil {
		return 0, shared.Internal(err)
	}
	return count, nil
}

func (r *UserRepository) getUserRoleIDs(ctx context.Context, userID shared.ID) ([]shared.ID, error) {
	roles, err := r.q.GetUserRoles(ctx, userID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	rids := make([]shared.ID, len(roles))
	for i, role := range roles {
		rids[i] = shared.ID(role.ID)
	}
	return rids, nil
}

func toDomainUser(row User) *identity.User {
	var lockedUntil *time.Time
	if row.LockedUntil.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.LockedUntil.String)
		lockedUntil = &t
	}
	var lastLogin *time.Time
	if row.LastLoginAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.LastLoginAt.String)
		lastLogin = &t
	}
	pwChanged, _ := time.Parse(time.RFC3339Nano, row.PasswordChangedAt)
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

	return &identity.User{
		ID:                 shared.ID(row.ID),
		Username:           row.Username,
		DisplayName:        row.DisplayName,
		Email:              row.Email.String,
		PasswordHash:       row.PasswordHash,
		PasswordChangedAt:  pwChanged,
		MustChangePassword: row.MustChangePassword == 1,
		Status:             identity.UserStatus(row.Status),
		FailedLoginCount:   int(row.FailedLoginCount),
		LockoutCount:       int(row.LockoutCount),
		LockedUntil:        lockedUntil,
		LastLoginAt:        lastLogin,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		Revision:           row.RowVersion,
	}
}

