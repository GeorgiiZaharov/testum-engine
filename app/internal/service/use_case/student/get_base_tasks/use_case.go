package getbasetasks

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

func (uc *UseCase) Execute(ctx context.Context, req GetBaseTasksRequest) (GetBaseTasksResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 {
		return GetBaseTasksResponse{}, ErrInvalidInput
	}

	var response GetBaseTasksResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа
		ok, err := r.Access.HasStudentAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check student access",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
				zap.Int("test_id", req.TestID),
			)
			return err
		}

		if !ok {
			return ErrAccessDenied
		}

		// 2. Проверка завершён ли тест
		result, err := r.Result.GetStudentResult(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to get student result",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
				zap.Int("test_id", req.TestID),
			)
			return err
		}

		if result.DateEnd != nil {
			return ErrTestCompleted
		}

		// 4. Получение базовых задач
		tasks, err := r.Task.GetBaseTasks(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get base tasks",
				zap.Error(err),
				zap.Int("test_id", req.TestID),
			)
			return err
		}

		// 5. Маппинг
		resultTasks := make([]Task, 0, len(tasks))

		for _, t := range tasks {
			answers := make([]Answer, 0, len(t.Answers))

			for _, a := range t.Answers {
				answers = append(answers, Answer{
					ID:       a.ID,
					Text:     a.Text,
					ImageURL: a.ImageURL,
				})
			}

			resultTasks = append(resultTasks, Task{
				ID:       t.ID,
				Text:     t.Text,
				ImageURL: t.ImageURL,
				IsHard:   t.IsHard,
				Answers:  answers,
			})
		}

		response.BaseTasks = resultTasks
		return nil
	})

	if err != nil {
		return GetBaseTasksResponse{}, err
	}

	return response, nil
}
