package app

import (
	studentgettesthandler "testum-engine/app/internal/handler/student/get_tests"

	getactivetestuc "testum-engine/app/internal/service/use_case/student/get_active_test"
	getfinishedtestuc "testum-engine/app/internal/service/use_case/student/get_finished_test"
)

func buildStudentGetTests(container *Container) *studentgettesthandler.Handler {
	getActiveTestUC := getactivetestuc.NewUseCase(
		getactivetestuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	getFinishedTestUC := getfinishedtestuc.NewUseCase(
		getfinishedtestuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	return studentgettesthandler.New(getActiveTestUC, getFinishedTestUC)
}
