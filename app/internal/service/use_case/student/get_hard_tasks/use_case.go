package gethardtasks

import (
	"context"
	"errors"

	"go.uber.org/zap"

	resultrepo "testum-engine/app/internal/repository/result"
	studenttestrepo "testum-engine/app/internal/repository/student_test"
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

func (uc *UseCase) Execute(ctx context.Context, req GetHardTasksRequest) (GetHardTasksResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 {
		return GetHardTasksResponse{}, ErrInvalidInput
	}

	var response GetHardTasksResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {
		// 1. Проверка доступа
		ok, err := r.Access.HasStudentAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check student access", zap.Error(err))
			return err
		}

		if !ok {
			return ErrAccessDenied
		}

		// 3. Проверка завершён ли тест
		_, err = r.Result.GetStudentResult(ctx, req.UserID, req.TestID)
		if errors.Is(err, resultrepo.ErrResultNotFound) {
			_, err = r.StudentTest.StartTest(ctx, studenttestrepo.StartTestParams{
				UserID: req.UserID,
				TestID: req.TestID,
			})
			if err != nil {
				uc.log.Error("failed to start test",
					zap.Error(err),
					zap.Int("user_id", req.UserID),
					zap.Int("test_id", req.TestID),
				)
				return err
			}
		} else if err != nil {
			uc.log.Error("failed to get student result",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
				zap.Int("test_id", req.TestID),
			)
			return err
		}

		// 2. Получение сложных заданий
		tasks, err := r.Task.GetHardTasks(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get hard tasks", zap.Error(err))
			return ErrGetHardTasks
		}

		hardTasks := make([]Task, 0, len(tasks))
		for _, t := range tasks {
			answers := make([]Answer, 0, len(t.Answers))
			trueAnswersCnt := 0

			for _, a := range t.Answers {
				answers = append(answers, Answer{
					ID:       a.ID,
					Text:     a.Text,
					ImageURL: a.ImageURL,
				})
				if a.IsCorrect {
					trueAnswersCnt += 1
				}
			}

			isMultipleChoice := false
			if trueAnswersCnt > 1 {
				isMultipleChoice = true
			}
			hardTasks = append(hardTasks, Task{
				ID:               t.ID,
				Text:             t.Text,
				ImageURL:         t.ImageURL,
				IsHard:           t.IsHard,
				IsMultipleChoice: isMultipleChoice,
				Answers:          answers,
			})
		}

		response.HardTasks = hardTasks
		return nil
	})

	if err != nil {
		return GetHardTasksResponse{}, err
	}

	return response, nil
}
