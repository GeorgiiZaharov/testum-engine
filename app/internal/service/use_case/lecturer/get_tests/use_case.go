package gettests

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

func (uc *UseCase) Execute(ctx context.Context, req GetTestsRequest) (GetTestsResponse, error) {
	if req.UserID <= 0 {
		return GetTestsResponse{}, ErrInvalidInput
	}

	var response GetTestsResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Получение тестов лектора
		tests, err := r.LecturerTest.GetByLecturer(ctx, req.UserID)
		if err != nil {
			uc.log.Error("failed to get tests by lecturer", zap.Error(err))
			return err
		}

		// 2. Явный маппинг
		mapped := make([]TestInfo, 0, len(tests))

		for _, t := range tests {
			mapped = append(mapped, TestInfo{
				ID:               t.ID,
				Name:             t.Name,
				CntQuestions:     t.CntQuestions,
				CntHardQuestions: t.CntHardQuestions,
				Groups:           t.Groups,
				DateCreated:      t.DateCreated,
			})
		}

		response = GetTestsResponse{
			Tests: mapped,
		}

		return nil
	})

	if err != nil {
		return GetTestsResponse{}, err
	}

	return response, nil
}
