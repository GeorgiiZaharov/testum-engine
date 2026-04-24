package giveaccess

import (
	"context"

	"go.uber.org/zap"
)

type UseCase struct {
	factory repoFactory
	log     *zap.Logger
}

func NewUseCase(factory repoFactory, log *zap.Logger) *UseCase {
	return &UseCase{
		factory: factory,
		log:     log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req GiveAccessRequest) (GiveAccessResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 || req.Group == "" {
		return GiveAccessResponse{}, ErrInvalidInput
	}

	var response GiveAccessResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа лектора к тесту
		ok, err := r.Access.HasLecturerAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check lecturer access",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
				zap.Int("test_id", req.TestID),
			)
			return err
		}

		if !ok {
			return ErrAccessDenied
		}

		// 2. Выдача доступа группе
		ok, err = r.Access.GiveAccess(ctx, req.TestID, req.Group)
		if !ok {
			uc.log.Error("failed to give access",
				zap.Error(err),
				zap.Int("test_id", req.TestID),
				zap.String("group", req.Group),
			)
			return err
		}

		response.Success = true
		return nil
	})

	if err != nil {
		return GiveAccessResponse{}, err
	}

	return response, nil
}
