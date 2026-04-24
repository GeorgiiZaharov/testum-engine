package login

import (
	"context"

	"go.uber.org/zap"

	userrepo "testum-engine/app/internal/repository/user"
)

type UseCase struct {
	factory repoFactory
	ldap    ldapAdapter
	auth    authService
	log     *zap.Logger
}

func NewUseCase(factory repoFactory, ldap ldapAdapter, auth authService, log *zap.Logger) *UseCase {
	return &UseCase{
		factory: factory,
		ldap:    ldap,
		auth:    auth,
		log:     log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req AuthLoginRequest) (AuthLoginResponse, error) {
	if req.Login == "" || req.Password == "" {
		return AuthLoginResponse{}, ErrInvalidInput
	}

	var response AuthLoginResponse

	// 1. LDAP auth check (outside tx — external system)
	if err := uc.ldap.ValidatePassword(ctx, req.Login, req.Password); err != nil {
		uc.log.Error("ldap auth failed", zap.Error(err))
		return AuthLoginResponse{}, ErrUnauthorized
	}

	var userID int

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 2. LDAP user info
		ldapUser, err := uc.ldap.GetInfo(ctx, req.Login)
		if err != nil {
			uc.log.Error("ldap get info failed", zap.Error(err))
			return ErrLDAPFailed
		}

		// 3. Upsert user
		userID, err = r.User.Upsert(ctx, userrepo.CreateUserParams{
			Login: ldapUser.Login,
			Name:  ldapUser.Name,
			Mail:  ldapUser.Mail,
			Group: ldapUser.Group,
		})
		if err != nil {
			uc.log.Error("user upsert failed", zap.Error(err))
			return err
		}

		return nil
	})

	if err != nil {
		return AuthLoginResponse{}, err
	}

	// 4. Generate tokens
	access, err := uc.auth.GenerateAccess(userID)
	if err != nil {
		uc.log.Error("generate access token failed", zap.Error(err))
		return AuthLoginResponse{}, ErrTokenGenerate
	}

	refresh, err := uc.auth.GenerateRefresh(userID)
	if err != nil {
		uc.log.Error("generate refresh token failed", zap.Error(err))
		return AuthLoginResponse{}, ErrTokenGenerate
	}

	response.AccessToken = access
	response.RefreshToken = refresh

	return response, nil
}
