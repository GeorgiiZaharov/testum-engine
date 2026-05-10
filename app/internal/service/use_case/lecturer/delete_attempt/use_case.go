package deleteattempt

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

func (uc *UseCase) Execute(ctx context.Context, req DeleteAttemptRequest) (DeleteAttemptResponse, error) {
	if req.LecturerID <= 0 || req.UserID <= 0 || req.TestID <= 0 {
		return DeleteAttemptResponse{}, ErrInvalidInput
	}

	var resp DeleteAttemptResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа лектора
		hasAccess, err := r.Access.HasLecturerAccess(ctx, req.LecturerID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check lecturer access", zap.Error(err))
			return err
		}

		if !hasAccess {
			return ErrAccessDenied
		}

		// 2. Удаление попытки студента (answers + student_tests)
		err = r.Result.DeleteAttempt(ctx, req.TestID, req.UserID)
		if err != nil {
			uc.log.Error("failed to delete attempt",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
				zap.Int("test_id", req.TestID),
			)
			return ErrDeleteFailed
		}

		resp.Success = true
		return nil
	})

	if err != nil {
		return DeleteAttemptResponse{}, err
	}

	return resp, nil
}
