package integration_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	authCmd "school-management/internal/application/commands/auth"
	rolesCmd "school-management/internal/application/commands/roles"
	schoolCmd "school-management/internal/application/commands/school"
	usersCmd "school-management/internal/application/commands/users"
	authQuery "school-management/internal/application/queries/auth"
	usersQuery "school-management/internal/application/queries/users"
	"school-management/internal/application/services/authz"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/platform/config"
	"school-management/internal/platform/security"
)

func TestAuthAndUserLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sms-auth-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.DefaultConfig(tempDir)
	cfg.Database.Path = filepath.Join(tempDir, "auth_test.db")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	ctx := context.Background()
	db, err := sqlite.Open(ctx, &cfg, logger)
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	defer db.Close()

	uow := sqlite.NewUnitOfWork(db)
	sessionStore := sqlite.NewSessionStore(db)
	clock := shared.SystemClock{}
	ids := shared.UUIDv7Generator{}
	hasher := security.NewArgon2Hasher(security.Argon2Params{
		Memory:      1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})

	setupHandler := authCmd.NewFirstRunSetupHandler(uow, ids, clock, hasher)
	loginHandler := authCmd.NewLoginHandler(uow, sessionStore, ids, clock, hasher, db.DeviceID())
	logoutHandler := authCmd.NewLogoutHandler(sessionStore)
	changePwHandler := authCmd.NewChangePasswordHandler(uow, hasher, clock, ids)

	authQueries := authQuery.NewAuthQueriesHandler(uow)
	usersHandler := usersCmd.NewUserCommandsHandler(uow, hasher, clock, ids)
	userQueries := usersQuery.NewUserQueriesHandler(uow)
	_ = rolesCmd.NewRoleCommandsHandler(uow, clock, ids)
	_ = schoolCmd.NewSchoolCommandsHandler(uow, clock, ids)

	// Step 1: Check setup before initialization
	status, err := authQueries.CheckSetupStatus(ctx)
	if err != nil {
		t.Fatalf("CheckSetupStatus failed: %v", err)
	}
	if status.IsInitialized {
		t.Errorf("expected isInitialized=false before setup")
	}

	// Step 2: First-run setup
	setupRes, err := setupHandler.Execute(ctx, authCmd.FirstRunSetupInput{
		SchoolName:         "Al-Huda Grammar School",
		SchoolCode:         "ALHUDA",
		Timezone:           "Asia/Karachi",
		CurrencyCode:       "PKR",
		CurrencyMinorUnits: 2,
		AdminUsername:      "admin",
		AdminDisplayName:   "Administrator",
		AdminEmail:         "admin@alhuda.edu",
		AdminPassword:      "SecureAdminPass123!",
	})
	if err != nil {
		t.Fatalf("FirstRunSetup failed: %v", err)
	}
	if len(setupRes.RecoveryCodes) != 10 {
		t.Errorf("expected 10 recovery codes, got %d", len(setupRes.RecoveryCodes))
	}

	// Step 3: Check setup after initialization
	status, err = authQueries.CheckSetupStatus(ctx)
	if err != nil || !status.IsInitialized {
		t.Fatalf("expected isInitialized=true, got %v", status)
	}
	if status.SchoolName != "Al-Huda Grammar School" {
		t.Errorf("expected school name 'Al-Huda Grammar School', got %q", status.SchoolName)
	}

	// Step 4: Login with wrong password
	_, err = loginHandler.Execute(ctx, authCmd.LoginInput{
		Username: "admin",
		Password: "WrongPassword!",
	})
	if err == nil {
		t.Fatalf("expected error on wrong password")
	}

	// Step 5: Login with correct credentials
	loginRes, err := loginHandler.Execute(ctx, authCmd.LoginInput{
		Username: "admin",
		Password: "SecureAdminPass123!",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loginRes.SessionToken == "" {
		t.Errorf("expected session token")
	}
	if loginRes.User.Username != "admin" {
		t.Errorf("got username %q, want 'admin'", loginRes.User.Username)
	}
	if len(loginRes.User.Permissions) == 0 {
		t.Errorf("expected permissions for super admin")
	}

	// Step 6: Create authenticated context as the logged-in super admin
	adminUserID, _ := shared.ParseID(loginRes.User.ID)
	adminPerms := make([]identity.Permission, len(loginRes.User.Permissions))
	for i, p := range loginRes.User.Permissions {
		adminPerms[i] = identity.Permission(p)
	}
	authedCtx := authz.WithActor(ctx, authz.Actor{
		UserID:      adminUserID,
		Username:    loginRes.User.Username,
		DisplayName: loginRes.User.DisplayName,
		DeviceID:    db.DeviceID(),
		SessionID:   loginRes.SessionToken,
		Roles:       loginRes.User.Roles,
		Permissions: identity.NewPermissionSet(adminPerms),
	})

	// Step 7: GetCurrentUser query
	currUser, err := authQueries.GetCurrentUser(authedCtx)
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	}
	if currUser.User.DisplayName != "Administrator" {
		t.Errorf("got %q, want 'Administrator'", currUser.User.DisplayName)
	}

	// Step 8: Create a second user (e.g. an Accountant)
	newUserID, err := usersHandler.CreateUser(authedCtx, usersCmd.CreateUserInput{
		Username:    "accountant1",
		DisplayName: "Ali Ahmed",
		Email:       "ali@alhuda.edu",
		Password:    "AccountantPass123!",
		RoleIDs:     []string{},
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if newUserID == "" {
		t.Fatalf("expected new user id")
	}

	// Step 9: List users
	userList, err := userQueries.ListUsers(authedCtx)
	if err != nil {
		t.Fatalf("ListUsers failed: %v", err)
	}
	if len(userList) != 2 {
		t.Errorf("expected 2 users, got %d", len(userList))
	}

	// Step 10: Change admin password
	err = changePwHandler.Execute(authedCtx, authCmd.ChangePasswordInput{
		OldPassword: "SecureAdminPass123!",
		NewPassword: "BrandNewPassword123!",
	})
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	// Step 11: Sign out
	err = logoutHandler.Execute(ctx, loginRes.SessionToken)
	if err != nil {
		t.Fatalf("logout failed: %v", err)
	}
}

