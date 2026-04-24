package getactivetest

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

func (uc *UseCase) Execute(ctx context.Context, req GetActiveTestRequest) (GetActiveTestResponse, error) {
	if req.UserID <= 0 {
		return GetActiveTestResponse{}, ErrInvalidInput
	}

	var response GetActiveTestResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		tests, err := r.StudentTest.GetActiveTests(ctx, req.UserID)
		if err != nil {
			uc.log.Error("failed to get active tests",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
			)
			return err
		}

		result := make([]StudentActiveTest, 0, len(tests))

		for _, t := range tests {
			result = append(result, StudentActiveTest{
				ID:               t.ID,
				Name:             t.Name,
				LecturerName:     t.LecturerName,
				CntQuestions:     t.CntQuestions,
				CntHardQuestions: t.CntHardQuestions,
				DateStart:        t.DateStart,
			})
		}

		response.ActiveTests = result
		return nil
	})

	if err != nil {
		return GetActiveTestResponse{}, err
	}

	return response, nil
}
