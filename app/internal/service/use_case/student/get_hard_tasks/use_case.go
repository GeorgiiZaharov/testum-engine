package gethardtasks

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

		// 2. Получение сложных заданий
		tasks, err := r.Task.GetHardTasks(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get hard tasks", zap.Error(err))
			return ErrGetHardTasks
		}

		// 3. Маппинг задач
		hardTasks := make([]Task, 0, len(tasks))
		for _, t := range tasks {
			answers := make([]Answer, 0, len(t.Answers))
			for _, a := range t.Answers {
				answers = append(answers, Answer{
					Text:     a.Text,
					ImageURL: a.ImageURL,
				})
			}

			hardTasks = append(hardTasks, Task{
				Text:     t.Text,
				ImageURL: t.ImageURL,
				IsHard:   t.IsHard,
				Answers:  answers,
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
