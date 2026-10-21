package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	_ "time/tzdata"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"school-management/frontend"
	acadCmd "school-management/internal/application/commands/academics"
	admCmd "school-management/internal/application/commands/admissions"
	authCmd "school-management/internal/application/commands/auth"
	docCmd "school-management/internal/application/commands/documents"
	enrCmd "school-management/internal/application/commands/enrollments"
	feeCmd "school-management/internal/application/commands/fees"
	guardCmd "school-management/internal/application/commands/guardians"
	paymentCmd "school-management/internal/application/commands/payments"
	rolesCmd "school-management/internal/application/commands/roles"
	schoolCmd "school-management/internal/application/commands/school"
	studCmd "school-management/internal/application/commands/students"
	usersCmd "school-management/internal/application/commands/users"
	acadQuery "school-management/internal/application/queries/academics"
	admQuery "school-management/internal/application/queries/admissions"
	auditQuery "school-management/internal/application/queries/audit"
	authQuery "school-management/internal/application/queries/auth"
	docQuery "school-management/internal/application/queries/documents"
	enrQuery "school-management/internal/application/queries/enrollments"
	feeQuery "school-management/internal/application/queries/fees"
	guardQuery "school-management/internal/application/queries/guardians"
	paymentQuery "school-management/internal/application/queries/payments"
	rolesQuery "school-management/internal/application/queries/roles"
	schoolQuery "school-management/internal/application/queries/school"
	studQuery "school-management/internal/application/queries/students"
	usersQuery "school-management/internal/application/queries/users"
	"school-management/internal/domain/shared"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/infrastructure/logging"
	"school-management/internal/infrastructure/storage"
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

	blobDir := filepath.Join(filepath.Dir(cfg.Database.Path), "blobs")
	blobStore, err := storage.NewFilesystemBlobStore(blobDir)
	if err != nil {
		logger.Error("failed to initialize blob store", "error", err)
		os.Exit(1)
	}

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
	acadCmdHandler := acadCmd.NewAcademicCommandsHandler(uow, clock, ids)
	studCmdHandler := studCmd.NewStudentCommandsHandler(uow, clock, ids)
	guardCmdHandler := guardCmd.NewGuardianCommandsHandler(uow, clock, ids)
	enrCmdHandler := enrCmd.NewEnrollmentCommandsHandler(uow, clock, ids)
	admCmdHandler := admCmd.NewAdmissionCommandsHandler(uow, clock, ids)
	docCmdHandler := docCmd.NewDocumentCommandsHandler(uow, blobStore, clock, ids)
	feeCmdHandler := feeCmd.NewFeeCommandsHandler(uow, clock, ids)
	paymentCmdHandler := paymentCmd.NewPaymentCommandsHandler(uow, clock, ids)

	authQueries := authQuery.NewAuthQueriesHandler(uow)
	usersQueries := usersQuery.NewUserQueriesHandler(uow)
	rolesQueries := rolesQuery.NewRoleQueriesHandler(uow)
	schoolQueries := schoolQuery.NewSchoolQueriesHandler(uow)
	auditQueries := auditQuery.NewAuditQueriesHandler(uow)
	acadQueries := acadQuery.NewAcademicQueriesHandler(uow)
	studQueries := studQuery.NewStudentQueriesHandler(uow)
	guardQueries := guardQuery.NewGuardianQueriesHandler(uow)
	enrQueries := enrQuery.NewEnrollmentQueriesHandler(uow)
	admQueries := admQuery.NewAdmissionQueriesHandler(uow)
	docQueries := docQuery.NewDocumentQueriesHandler(uow, blobStore)
	feeQueries := feeQuery.NewFeeQueriesHandler(uow)
	paymentQueries := paymentQuery.NewPaymentQueriesHandler(uow)

	// Desktop Session Manager and Transport Services
	sessionMgr := wailsTransport.NewSessionManager(ids)
	authSvc := wailsTransport.NewAuthService(sessionMgr, setupHandler, loginHandler, logoutHandler, changePwHandler, authQueries, db.DeviceID())
	userSvc := wailsTransport.NewUserService(sessionMgr, usersCmdHandler, usersQueries)
	roleSvc := wailsTransport.NewRoleService(sessionMgr, rolesCmdHandler, rolesQueries)
	schoolSvc := wailsTransport.NewSchoolService(sessionMgr, schoolCmdHandler, schoolQueries)
	auditSvc := wailsTransport.NewAuditService(sessionMgr, auditQueries)
	acadSvc := wailsTransport.NewAcademicService(sessionMgr, acadCmdHandler, acadQueries)
	studSvc := wailsTransport.NewStudentService(sessionMgr, studCmdHandler, studQueries)
	guardSvc := wailsTransport.NewGuardianService(sessionMgr, guardCmdHandler, guardQueries)
	enrSvc := wailsTransport.NewEnrollmentService(sessionMgr, enrCmdHandler, enrQueries)
	admSvc := wailsTransport.NewAdmissionService(sessionMgr, admCmdHandler, admQueries)
	docSvc := wailsTransport.NewDocumentService(sessionMgr, docCmdHandler, docQueries)
	feeSvc := wailsTransport.NewFeeService(sessionMgr, feeCmdHandler, feeQueries)
	invoiceSvc := wailsTransport.NewInvoiceService(sessionMgr, feeCmdHandler, feeQueries)
	paymentSvc := wailsTransport.NewPaymentService(sessionMgr, paymentCmdHandler, paymentQueries)

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
			acadSvc,
			studSvc,
			guardSvc,
			enrSvc,
			admSvc,
			docSvc,
			feeSvc,
			invoiceSvc,
			paymentSvc,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
	})

	if err != nil {
		logger.Error("application crashed", "error", err)
	}
}
