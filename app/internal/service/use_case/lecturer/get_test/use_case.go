package gettest

import (
	"context"
	"errors"

	"go.uber.org/zap"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
	taskrepo "testum-engine/app/internal/repository/task"
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

func (uc *UseCase) Execute(ctx context.Context, req GetTestRequest) (GetTestResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 {
		return GetTestResponse{}, ErrInvalidInput
	}

	var response GetTestResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа лектора
		hasAccess, err := r.Access.HasLecturerAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check access", zap.Error(err))
			return err
		}

		if !hasAccess {
			return ErrAccessDenied
		}

		// 2. Получение теста
		test, err := r.LecturerTest.GetByID(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get test", zap.Error(err))
			return ErrNotFound
		}

		// 3. Группы (EXPLICIT COPY)
		rawGroups, err := r.LecturerTest.GetGroups(ctx, req.TestID, 0)
		if errors.Is(err, lecturertestrepo.ErrTestNotFound) {
			rawGroups = []lecturertestrepo.GroupInfo{}
		} else if err != nil {
			uc.log.Error("failed to get groups", zap.Error(err))
			return err
		}

		groups := make([]GroupInfo, 0, len(rawGroups))
		for _, g := range rawGroups {
			groups = append(groups, GroupInfo{
				GroupName:    g.GroupName,
				MembersCount: g.MembersCount,
			})
		}

		// 4. Hard tasks (EXPLICIT COPY)
		rawHardTasks, err := r.Task.GetHardTasks(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get hard tasks", zap.Error(err))
			return err
		}

		hardTasks := make([]Task, 0, len(rawHardTasks))
		for _, t := range rawHardTasks {
			hardTasks = append(hardTasks, Task{
				Text:             t.Text,
				ImageURL:         t.ImageURL,
				IsHard:           t.IsHard,
				IsMultipleChoice: isMultipleChoice(t.Answers),
				Answers:          mapAnswers(t.Answers),
			})
		}

		// 5. Base tasks (EXPLICIT COPY)
		rawBaseTasks, err := r.Task.GetBaseTasks(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get base tasks", zap.Error(err))
			return err
		}

		baseTasks := make([]Task, 0, len(rawBaseTasks))
		for _, t := range rawBaseTasks {
			baseTasks = append(baseTasks, Task{
				Text:             t.Text,
				ImageURL:         t.ImageURL,
				IsHard:           t.IsHard,
				IsMultipleChoice: isMultipleChoice(t.Answers),
				Answers:          mapAnswers(t.Answers),
			})
		}

		// 6. Response mapping
		response = GetTestResponse{
			ID:               test.ID,
			Name:             test.Name,
			CntQuestions:     test.CntQuestions,
			CntHardQuestions: test.CntHardQuestions,
			Groups:           groups,
			HardTasks:        hardTasks,
			BaseTasks:        baseTasks,
		}

		return nil
	})

	if err != nil {
		return GetTestResponse{}, err
	}

	return response, nil
}
func mapAnswers(raw []taskrepo.Answer) []Answer {
	res := make([]Answer, 0, len(raw))

	for _, a := range raw {
		res = append(res, Answer{
			Text:      a.Text,
			ImageURL:  a.ImageURL,
			IsCorrect: a.IsCorrect,
		})
	}

	return res
}

func isMultipleChoice(raw []taskrepo.Answer) bool {
	cntTrueAnswers := 0

	for _, a := range raw {
		if a.IsCorrect {
			cntTrueAnswers += 1
		}
	}

	return cntTrueAnswers > 1
}
