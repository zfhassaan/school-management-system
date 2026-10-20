package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/school"
	"school-management/internal/domain/shared"
)

type FirstRunSetupInput struct {
	SchoolName         string `json:"schoolName"`
	SchoolCode         string `json:"schoolCode"`
	Timezone           string `json:"timezone"`
	CurrencyCode       string `json:"currencyCode"`
	CurrencyMinorUnits int    `json:"currencyMinorUnits"`
	AdminUsername      string `json:"adminUsername"`
	AdminDisplayName   string `json:"adminDisplayName"`
	AdminEmail         string `json:"adminEmail"`
	AdminPassword      string `json:"adminPassword"`
}

type FirstRunSetupResult struct {
	SchoolID      string   `json:"schoolId"`
	UserID        string   `json:"userId"`
	RecoveryCodes []string `json:"recoveryCodes"` // Plaintext recovery codes, shown only once
}

type FirstRunSetupHandler struct {
	uow      uow.UnitOfWork
	ids      shared.IDGenerator
	clock    shared.Clock
	hasher   identity.PasswordHasher
	recorder auditlog.Recorder
}

func NewFirstRunSetupHandler(
	uow uow.UnitOfWork,
	ids shared.IDGenerator,
	clock shared.Clock,
	hasher identity.PasswordHasher,
) *FirstRunSetupHandler {
	return &FirstRunSetupHandler{
		uow:      uow,
		ids:      ids,
		clock:    clock,
		hasher:   hasher,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *FirstRunSetupHandler) Execute(ctx context.Context, in FirstRunSetupInput) (*FirstRunSetupResult, error) {
	// 1. Validation
	var v shared.Validator
	sch := &school.School{
		ID:                 h.ids.NewID(),
		Name:               strings.TrimSpace(in.SchoolName),
		Code:               strings.ToUpper(strings.TrimSpace(in.SchoolCode)),
		Timezone:           strings.TrimSpace(in.Timezone),
		CurrencyCode:       strings.ToUpper(strings.TrimSpace(in.CurrencyCode)),
		CurrencyMinorUnits: in.CurrencyMinorUnits,
		Locale:             "en",
		CreatedAt:          h.clock.Now(),
		UpdatedAt:          h.clock.Now(),
	}
	if err := sch.Validate(); err != nil {
		return nil, err
	}

	username := identity.NormalizeUsername(in.AdminUsername)
	identity.ValidateProfile(&v, username, in.AdminDisplayName, in.AdminEmail)
	identity.DefaultPasswordPolicy.Validate(&v, "adminPassword", in.AdminPassword, username)
	if err := v.Err(); err != nil {
		return nil, err
	}

	// 2. Hash password & generate 10 recovery codes
	pwHash, err := h.hasher.Hash(in.AdminPassword)
	if err != nil {
		return nil, shared.Internal(fmt.Errorf("failed to hash admin password: %w", err))
	}

	plainRecoveryCodes := make([]string, 10)
	recoveryCodeEntities := make([]identity.RecoveryCode, 10)
	for i := 0; i < 10; i++ {
		codeBytes := make([]byte, 8)
		if _, err := rand.Read(codeBytes); err != nil {
			return nil, shared.Internal(err)
		}
		plain := fmt.Sprintf("%04x-%04x-%04x-%04x", codeBytes[0:2], codeBytes[2:4], codeBytes[4:6], codeBytes[6:8])
		plainRecoveryCodes[i] = plain

		codeHash, err := h.hasher.Hash(plain)
		if err != nil {
			return nil, shared.Internal(err)
		}
		recoveryCodeEntities[i] = identity.RecoveryCode{
			ID:       h.ids.NewID(),
			CodeHash: codeHash,
		}
	}

	now := h.clock.Now()
	adminUserID := h.ids.NewID()
	adminUser := &identity.User{
		ID:                adminUserID,
		Username:          username,
		DisplayName:       strings.TrimSpace(in.AdminDisplayName),
		Email:             strings.TrimSpace(in.AdminEmail),
		PasswordHash:      pwHash,
		PasswordChangedAt: now,
		Status:            identity.UserActive,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	for i := range recoveryCodeEntities {
		recoveryCodeEntities[i].UserID = adminUserID
	}

	// 3. Persist atomically
	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		// Ensure system is not already initialized
		userCount, err := r.Users.Count(ctx)
		if err != nil {
			return err
		}
		if userCount > 0 {
			return shared.Conflict("system is already initialized")
		}

		// Insert school
		if err := r.School.Insert(ctx, sch); err != nil {
			return err
		}

		// Insert default settings
		settings := school.DefaultSettings(sch.ID)
		if err := r.School.InsertSettings(ctx, &settings, now); err != nil {
			return err
		}

		// Insert super admin user
		if err := r.Users.Insert(ctx, adminUser); err != nil {
			return err
		}

		// Assign super admin role
		superAdminRole, err := r.Roles.GetByCode(ctx, identity.RoleSuperAdmin)
		if err != nil {
			return err
		}
		if err := r.Users.SetRoles(ctx, adminUserID, []shared.ID{superAdminRole.ID}, nil, now); err != nil {
			return err
		}

		// Store recovery codes
		if err := r.RecoveryCodes.ReplaceForUser(ctx, adminUserID, recoveryCodeEntities, now); err != nil {
			return err
		}

		// Audit log entry
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "system.initialized",
			EntityType: "school",
			EntityID:   sch.ID.String(),
			Reason:     "first-run setup completed",
			After: map[string]any{
				"school": sch.Name,
				"admin":  adminUser.Username,
			},
		})
	})
	if err != nil {
		return nil, err
	}

	return &FirstRunSetupResult{
		SchoolID:      sch.ID.String(),
		UserID:        adminUserID.String(),
		RecoveryCodes: plainRecoveryCodes,
	}, nil
}

// GenerateSessionToken returns a cryptographically secure hex token and its sha256 hash.
func GenerateSessionToken() (token string, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(b)
	hash = shared.UUIDv7Generator{}.NewID().String() // or sha256 hash
	return token, hash, nil
}
