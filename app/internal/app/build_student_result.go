package app

import (
	studentresulthandler "testum-engine/app/internal/handler/student/result"

	gettestresultuc "testum-engine/app/internal/service/use_case/student/get_test_result"
)

func buildStudentResult(container *Container) *studentresulthandler.Handler {
	getTestResultUC := gettestresultuc.NewUseCase(
		gettestresultuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	return studentresulthandler.New(getTestResultUC)
}
