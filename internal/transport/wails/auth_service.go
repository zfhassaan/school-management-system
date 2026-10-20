package wails

import (
	"context"

	authCmd "school-management/internal/application/commands/auth"
	authQuery "school-management/internal/application/queries/auth"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/reqctx"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/transport/dto"
)

type AuthService struct {
	sessions        *SessionManager
	setupHandler    *authCmd.FirstRunSetupHandler
	loginHandler    *authCmd.LoginHandler
	logoutHandler   *authCmd.LogoutHandler
	changePwHandler *authCmd.ChangePasswordHandler
	authQueries     *authQuery.AuthQueriesHandler
	deviceID        shared.ID
}

func NewAuthService(
	sessions *SessionManager,
	setupHandler *authCmd.FirstRunSetupHandler,
	loginHandler *authCmd.LoginHandler,
	logoutHandler *authCmd.LogoutHandler,
	changePwHandler *authCmd.ChangePasswordHandler,
	authQueries *authQuery.AuthQueriesHandler,
	deviceID shared.ID,
) *AuthService {
	return &AuthService{
		sessions:        sessions,
		setupHandler:    setupHandler,
		loginHandler:    loginHandler,
		logoutHandler:   logoutHandler,
		changePwHandler: changePwHandler,
		authQueries:     authQueries,
		deviceID:        deviceID,
	}
}

func (s *AuthService) CheckSetupStatus(ctx context.Context) dto.Response[*authQuery.SetupStatusDTO] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.authQueries.CheckSetupStatus(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*authQuery.SetupStatusDTO](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AuthService) FirstRunSetup(ctx context.Context, in authCmd.FirstRunSetupInput) dto.Response[*authCmd.FirstRunSetupResult] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.setupHandler.Execute(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*authCmd.FirstRunSetupResult](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

func (s *AuthService) Login(ctx context.Context, in authCmd.LoginInput) dto.Response[*authCmd.LoginResult] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.loginHandler.Execute(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*authCmd.LoginResult](code, msg, fields, reqctx.CorrelationID(ctx))
	}

	// Update desktop session manager with active actor
	uid, _ := shared.ParseID(res.User.ID)
	perms := make([]identity.Permission, len(res.User.Permissions))
	for i, p := range res.User.Permissions {
		perms[i] = identity.Permission(p)
	}

	actor := &authz.Actor{
		UserID:             uid,
		Username:           res.User.Username,
		DisplayName:        res.User.DisplayName,
		DeviceID:           s.deviceID,
		SessionID:          res.SessionToken,
		Roles:              res.User.Roles,
		Permissions:        identity.NewPermissionSet(perms),
		MustChangePassword: res.User.MustChangePassword,
	}
	s.sessions.SetActor(actor)

	return dto.Ok(res)
}

func (s *AuthService) Logout(ctx context.Context) dto.Response[bool] {
	actor := s.sessions.CurrentActor()
	if actor != nil {
		_ = s.logoutHandler.Execute(ctx, actor.SessionID)
		s.sessions.ClearActor()
	}
	return dto.Ok(true)
}

func (s *AuthService) ChangePassword(ctx context.Context, in authCmd.ChangePasswordInput) dto.Response[bool] {
	ctx = s.sessions.WithContext(ctx)
	err := s.changePwHandler.Execute(ctx, in)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[bool](code, msg, fields, reqctx.CorrelationID(ctx))
	}

	// Update actor mustChangePassword to false
	actor := s.sessions.CurrentActor()
	if actor != nil {
		actor.MustChangePassword = false
		s.sessions.SetActor(actor)
	}

	return dto.Ok(true)
}

func (s *AuthService) GetCurrentUser(ctx context.Context) dto.Response[*authCmd.LoginResult] {
	ctx = s.sessions.WithContext(ctx)
	res, err := s.authQueries.GetCurrentUser(ctx)
	if err != nil {
		code, msg, fields := MapError(err)
		return dto.Fail[*authCmd.LoginResult](code, msg, fields, reqctx.CorrelationID(ctx))
	}
	return dto.Ok(res)
}

