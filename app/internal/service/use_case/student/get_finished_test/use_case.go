package getfinishedtest

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

func (uc *UseCase) Execute(ctx context.Context, req GetFinishedTestRequest) (GetFinishedTestResponse, error) {
	if req.UserID <= 0 {
		return GetFinishedTestResponse{}, ErrInvalidInput
	}

	var response GetFinishedTestResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		tests, err := r.StudentTest.GetFinishedTests(ctx, req.UserID)
		if err != nil {
			uc.log.Error("failed to get finished tests",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
			)
			return err
		}

		result := make([]StudentFinishTest, 0, len(tests))

		for _, t := range tests {
			result = append(result, StudentFinishTest{
				ID:               t.ID,
				Name:             t.Name,
				LecturerName:     t.LecturerName,
				CntQuestions:     t.CntQuestions,
				CntHardQuestions: t.CntHardQuestions,
				Mark:             t.Mark,
				SuccessRate:      t.SuccessRate,
				DateStart:        t.DateStart,
				DateEnd:          t.DateEnd,
			})
		}

		response.FinishedTests = result
		return nil
	})

	if err != nil {
		return GetFinishedTestResponse{}, err
	}

	return response, nil
}
