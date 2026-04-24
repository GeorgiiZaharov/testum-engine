package app

import (
	studenttestattempthandler "testum-engine/app/internal/handler/student/test_attempt"

	getbasetasksuc "testum-engine/app/internal/service/use_case/student/get_base_tasks"
	gethardtasksuc "testum-engine/app/internal/service/use_case/student/get_hard_tasks"
	postbasetasksuc "testum-engine/app/internal/service/use_case/student/post_base_answers"
	posthardtasksuc "testum-engine/app/internal/service/use_case/student/post_hard_answers"
)

func buildStudentTestAttempt(container *Container) *studenttestattempthandler.Handler {

	getHardTasksUC := gethardtasksuc.NewUseCase(
		gethardtasksuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	getBaseTasksUC := getbasetasksuc.NewUseCase(
		getbasetasksuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	postHardTasksUC := posthardtasksuc.NewUseCase(
		posthardtasksuc.NewFactory(container.DB, container.Logger),
		*container.AnswerService,
		*container.ResultService,
		container.Logger,
	)
	postBaseTasksUC := postbasetasksuc.NewUseCase(
		postbasetasksuc.NewFactory(container.DB, container.Logger),
		*container.AnswerService,
		*container.ResultService,
		container.Logger,
	)
	return studenttestattempthandler.New(getHardTasksUC, getBaseTasksUC, postHardTasksUC, postBaseTasksUC)
}
