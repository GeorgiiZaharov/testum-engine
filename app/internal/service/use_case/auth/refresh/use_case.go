package refresh

import (
	"context"

	"go.uber.org/zap"
)

type UseCase struct {
	factory repoFactory
	auth    authService
	log     *zap.Logger
}

func NewUseCase(factory repoFactory, auth authService, log *zap.Logger) *UseCase {
	return &UseCase{
		factory: factory,
		auth:    auth,
		log:     log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req AuthRefreshRequest) (AuthRefreshResponse, error) {
	if req.UserID <= 0 {
		return AuthRefreshResponse{}, ErrInvalidUserID
	}

	var response AuthRefreshResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка существования пользователя
		_, err := r.User.GetByID(ctx, req.UserID)
		if err != nil {
			uc.log.Error("failed to get user", zap.Error(err))
			return err
		}

		// 2. Генерация access token
		access, err := uc.auth.GenerateAccess(req.UserID)
		if err != nil {
			uc.log.Error("failed to generate access token", zap.Error(err))
			return ErrAuthFailed
		}

		// 3. Генерация refresh token
		refresh, err := uc.auth.GenerateRefresh(req.UserID)
		if err != nil {
			uc.log.Error("failed to generate refresh token", zap.Error(err))
			return ErrAuthFailed
		}

		response.AccessToken = access
		response.RefreshToken = refresh

		return nil
	})

	if err != nil {
		return AuthRefreshResponse{}, err
	}

	return response, nil
}
