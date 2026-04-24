package gettestresult

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

func (uc *UseCase) Execute(ctx context.Context, req GetTestResultRequest) (GetTestResultResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 {
		return GetTestResultResponse{}, ErrInvalidInput
	}

	var response GetTestResultResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {
		// 1. Проверка доступа
		hasAccess, err := r.Access.HasStudentAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check student access", zap.Error(err))
			return err
		}

		if !hasAccess {
			return ErrAccessDenied
		}

		// 2. Получение результатов теста
		testResult, err := r.Result.GetStudentResult(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to get student result", zap.Error(err))
			return ErrResultNotFound
		}

		// 3. Формирование ответа
		response.Mark = testResult.Mark
		response.SuccessRate = testResult.SuccessRate
		response.DateStart = testResult.DateStart
		response.DateEnd = testResult.DateEnd

		return nil
	})

	if err != nil {
		return GetTestResultResponse{}, err
	}

	return response, nil
}
