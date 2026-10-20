package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/platform/security"
)

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserDTO struct {
	ID                 string   `json:"id"`
	Username           string   `json:"username"`
	DisplayName        string   `json:"displayName"`
	Email              string   `json:"email"`
	Roles              []string `json:"roles"`
	Permissions        []string `json:"permissions"`
	MustChangePassword bool     `json:"mustChangePassword"`
}

type SchoolDTO struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Code               string `json:"code"`
	Timezone           string `json:"timezone"`
	CurrencyCode       string `json:"currencyCode"`
	CurrencyMinorUnits int    `json:"currencyMinorUnits"`
	DateFormat         string `json:"dateFormat"`
}

type LoginResult struct {
	SessionToken string    `json:"sessionToken"`
	User         UserDTO   `json:"user"`
	School       SchoolDTO `json:"school"`
}

type LoginHandler struct {
	uow          uow.UnitOfWork
	sessionStore identity.SessionStore
	ids          shared.IDGenerator
	clock        shared.Clock
	hasher       identity.PasswordHasher
	recorder     auditlog.Recorder
	deviceID     shared.ID
}

func NewLoginHandler(
	uow uow.UnitOfWork,
	sessionStore identity.SessionStore,
	ids shared.IDGenerator,
	clock shared.Clock,
	hasher identity.PasswordHasher,
	deviceID shared.ID,
) *LoginHandler {
	return &LoginHandler{
		uow:          uow,
		sessionStore: sessionStore,
		ids:          ids,
		clock:        clock,
		hasher:       hasher,
		recorder:     auditlog.Recorder{IDs: ids, Clock: clock},
		deviceID:     deviceID,
	}
}

func (h *LoginHandler) Execute(ctx context.Context, in LoginInput) (*LoginResult, error) {
	username := identity.NormalizeUsername(in.Username)
	if username == "" || in.Password == "" {
		return nil, shared.Unauthenticated("username and password are required")
	}

	var targetUser *identity.User
	var roles []string
	var perms []identity.Permission
	var schoolDTO SchoolDTO

	now := h.clock.Now()

	// Step 1: Check user credentials and status inside transaction
	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		u, err := r.Users.GetByUsername(ctx, username)
		if err != nil {
			// Prevent user enumeration by taking constant time or generic error
			return shared.Unauthenticated("invalid username or password")
		}

		if u.Status == identity.UserDisabled {
			return shared.Forbidden("account is disabled")
		}

		// Check Lockout
		if u.IsLocked(now) {
			return &shared.Error{
				Code:    shared.CodeLocked,
				Message: fmt.Sprintf("account is locked until %s", u.LockedUntil.Format(time.RFC3339)),
			}
		}

		// Verify password
		match, needsRehash, verifyErr := h.hasher.Verify(in.Password, u.PasswordHash)
		if verifyErr != nil || !match {
			// Register failure and check if newly locked
			locked := u.RegisterFailedLogin(now, identity.DefaultLockoutPolicy)
			_ = r.Users.Update(ctx, u)

			_ = h.recorder.Record(ctx, r.Audit, auditlog.Change{
				Action:     "auth.login_failed",
				EntityType: "user",
				EntityID:   u.ID.String(),
				Reason:     "invalid credentials",
			})

			if locked {
				return &shared.Error{
					Code:    shared.CodeLocked,
					Message: "account has been locked due to too many failed attempts",
				}
			}
			return shared.Unauthenticated("invalid username or password")
		}

		// Success! Clear throttling state
		u.RegisterSuccessfulLogin(now)
		if needsRehash {
			if newHash, err := h.hasher.Hash(in.Password); err == nil {
				u.PasswordHash = newHash
			}
		}
		if err := r.Users.Update(ctx, u); err != nil {
			return err
		}

		userRoles, err := r.Users.RoleCodes(ctx, u.ID)
		if err != nil {
			return err
		}
		userPerms, err := r.Users.EffectivePermissions(ctx, u.ID)
		if err != nil {
			return err
		}

		// Fetch school details
		sch, err := r.School.Get(ctx)
		if err == nil {
			settings, _ := r.School.GetSettings(ctx, sch.ID)
			dateFormat := "dd/MM/yyyy"
			if settings != nil {
				dateFormat = settings.DateFormat
			}
			schoolDTO = SchoolDTO{
				ID:                 sch.ID.String(),
				Name:               sch.Name,
				Code:               sch.Code,
				Timezone:           sch.Timezone,
				CurrencyCode:       sch.CurrencyCode,
				CurrencyMinorUnits: sch.CurrencyMinorUnits,
				DateFormat:         dateFormat,
			}
		}

		_ = h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "auth.login_success",
			EntityType: "user",
			EntityID:   u.ID.String(),
			Reason:     "successful sign-in",
		})

		targetUser = u
		roles = userRoles
		perms = userPerms
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Step 2: Create Session
	rawToken, err := security.GenerateRandomToken(32)
	if err != nil {
		return nil, shared.Internal(err)
	}
	hasher := sha256.New()
	hasher.Write([]byte(rawToken))
	tokenHash := hex.EncodeToString(hasher.Sum(nil))

	session := identity.Session{
		ID:                tokenHash,
		UserID:            targetUser.ID,
		DeviceID:          h.deviceID,
		CreatedAt:         now,
		LastSeenAt:        now,
		IdleExpiresAt:     now.Add(identity.DefaultSessionPolicy.IdleTimeout),
		AbsoluteExpiresAt: now.Add(identity.DefaultSessionPolicy.AbsoluteTimeout),
	}

	if err := h.sessionStore.Create(ctx, session); err != nil {
		return nil, shared.Internal(fmt.Errorf("failed to create session: %w", err))
	}

	permStrings := make([]string, len(perms))
	for i, p := range perms {
		permStrings[i] = string(p)
	}

	return &LoginResult{
		SessionToken: rawToken,
		User: UserDTO{
			ID:                 targetUser.ID.String(),
			Username:           targetUser.Username,
			DisplayName:        targetUser.DisplayName,
			Email:              targetUser.Email,
			Roles:              roles,
			Permissions:        permStrings,
			MustChangePassword: targetUser.MustChangePassword,
		},
		School: schoolDTO,
	}, nil
}
