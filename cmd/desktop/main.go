package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"school-management/frontend"
	authCmd "school-management/internal/application/commands/auth"
	rolesCmd "school-management/internal/application/commands/roles"
	schoolCmd "school-management/internal/application/commands/school"
	usersCmd "school-management/internal/application/commands/users"
	auditQuery "school-management/internal/application/queries/audit"
	authQuery "school-management/internal/application/queries/auth"
	rolesQuery "school-management/internal/application/queries/roles"
	schoolQuery "school-management/internal/application/queries/school"
	usersQuery "school-management/internal/application/queries/users"
	"school-management/internal/domain/shared"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/infrastructure/logging"
	"school-management/internal/platform/config"
	"school-management/internal/platform/security"
	wailsTransport "school-management/internal/transport/wails"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger := logging.SetupLogger(cfg.Logging)
	logger.Info("starting school management desktop application", "env", cfg.App.Environment)

	ctx := context.Background()
	db, err := sqlite.Open(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to initialize sqlite database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	uow := sqlite.NewUnitOfWork(db)
	sessionStore := sqlite.NewSessionStore(db)
	clock := shared.SystemClock{}
	ids := shared.UUIDv7Generator{}
	hasher := security.NewDefaultArgon2Hasher()

	// Application Handlers
	setupHandler := authCmd.NewFirstRunSetupHandler(uow, ids, clock, hasher)
	loginHandler := authCmd.NewLoginHandler(uow, sessionStore, ids, clock, hasher, db.DeviceID())
	logoutHandler := authCmd.NewLogoutHandler(sessionStore)
	changePwHandler := authCmd.NewChangePasswordHandler(uow, hasher, clock, ids)

	usersCmdHandler := usersCmd.NewUserCommandsHandler(uow, hasher, clock, ids)
	rolesCmdHandler := rolesCmd.NewRoleCommandsHandler(uow, clock, ids)
	schoolCmdHandler := schoolCmd.NewSchoolCommandsHandler(uow, clock, ids)

	authQueries := authQuery.NewAuthQueriesHandler(uow)
	usersQueries := usersQuery.NewUserQueriesHandler(uow)
	rolesQueries := rolesQuery.NewRoleQueriesHandler(uow)
	schoolQueries := schoolQuery.NewSchoolQueriesHandler(uow)
	auditQueries := auditQuery.NewAuditQueriesHandler(uow)

	// Desktop Session Manager and Transport Services
	sessionMgr := wailsTransport.NewSessionManager(ids)
	authSvc := wailsTransport.NewAuthService(sessionMgr, setupHandler, loginHandler, logoutHandler, changePwHandler, authQueries, db.DeviceID())
	userSvc := wailsTransport.NewUserService(sessionMgr, usersCmdHandler, usersQueries)
	roleSvc := wailsTransport.NewRoleService(sessionMgr, rolesCmdHandler, rolesQueries)
	schoolSvc := wailsTransport.NewSchoolService(sessionMgr, schoolCmdHandler, schoolQueries)
	auditSvc := wailsTransport.NewAuditService(sessionMgr, auditQueries)

	err = wails.Run(&options.App{
		Title:             "School Management System",
		Width:             1366,
		Height:            768,
		MinWidth:          1024,
		MinHeight:         700,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 248, G: 250, B: 252, A: 255},
		AssetServer: &assetserver.Options{
			Assets: frontend.Assets,
		},
		Bind: []any{
			authSvc,
			userSvc,
			roleSvc,
			schoolSvc,
			auditSvc,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			BackdropType:         windows.Mica,
			Theme:                windows.SystemDefault,
		},
	})

	if err != nil {
		logger.Error("application execution failed", "error", err)
	}
}

