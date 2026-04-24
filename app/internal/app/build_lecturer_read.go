package app

import (
	lecturerreadhandler "testum-engine/app/internal/handler/lecturer/read"

	getgroupsuc "testum-engine/app/internal/service/use_case/lecturer/get_groups"
	gettestuc "testum-engine/app/internal/service/use_case/lecturer/get_test"
	getresultuc "testum-engine/app/internal/service/use_case/lecturer/get_test_result"
	gettestsuc "testum-engine/app/internal/service/use_case/lecturer/get_tests"
)

func buildLecturerRead(container *Container) *lecturerreadhandler.Handler {
	getTestsUC := gettestsuc.NewUseCase(
		gettestsuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	getTestUC := gettestuc.NewUseCase(
		gettestuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	getGroupsUC := getgroupsuc.NewUseCase(
		getgroupsuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	getRestultUC := getresultuc.NewUseCase(
		getresultuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	return lecturerreadhandler.New(getTestsUC, getTestUC, getGroupsUC, getRestultUC)
}
