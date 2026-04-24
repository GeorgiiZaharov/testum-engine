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
	if req.UserID <= 0 || req.TestID <= 0 || req.Group == "" {
		return GetTestResultResponse{}, ErrInvalidInput
	}

	var response GetTestResultResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа
		hasAccess, err := r.Access.HasLecturerAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check access", zap.Error(err))
			return err
		}

		if !hasAccess {
			return ErrForbidden
		}

		// 2. Получение результатов группы
		results, err := r.Result.GetGroupResult(ctx, req.TestID, req.Group, req.Year)
		if err != nil {
			uc.log.Error("failed to get group results", zap.Error(err))
			return err
		}

		// 3. Маппинг результатов
		mapped := make([]StudentResult, 0, len(results))

		for _, r := range results {
			mapped = append(mapped, StudentResult{
				UserID:      r.UserID,
				Name:        r.Name,
				Login:       r.Login,
				Mail:        r.Mail,
				Mark:        r.Mark,
				SuccessRate: r.SuccessRate,
				DateStart:   r.DateStart,
				DateEnd:     r.DateEnd,
			})
		}

		// 4. Response
		response = GetTestResultResponse{
			Results: mapped,
		}

		return nil
	})

	if err != nil {
		return GetTestResultResponse{}, err
	}

	return response, nil
}
