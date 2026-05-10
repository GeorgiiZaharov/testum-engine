package app

import (
	lecturerwritehandler "testum-engine/app/internal/handler/lecturer/write"

	deleteattemptuc "testum-engine/app/internal/service/use_case/lecturer/delete_attempt"
	deletetestuc "testum-engine/app/internal/service/use_case/lecturer/delete_test"
	giveaccessuc "testum-engine/app/internal/service/use_case/lecturer/give_access"
	takeaccessuc "testum-engine/app/internal/service/use_case/lecturer/take_access"
)

func buildLecturerWrite(container *Container) *lecturerwritehandler.Handler {
	deleteUC := deletetestuc.NewUseCase(
		deletetestuc.NewFactory(container.DB, container.Logger),
		container.Storage,
		container.Logger,
	)
	giveAccessUC := giveaccessuc.NewUseCase(
		giveaccessuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	takeAccessUC := takeaccessuc.NewUseCase(
		takeaccessuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	deleteAttemptUC := deleteattemptuc.NewUseCase(
		deleteattemptuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)

	return lecturerwritehandler.New(deleteUC, giveAccessUC, takeAccessUC, deleteAttemptUC)
}
