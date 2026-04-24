package app

import (
	adminhandler "testum-engine/app/internal/handler/admin"

	createlectureruc "testum-engine/app/internal/service/use_case/admin/create_lecturer"
	deletelectureruc "testum-engine/app/internal/service/use_case/admin/delete_lecturer"
	getlecturersuc "testum-engine/app/internal/service/use_case/admin/get_lecturers"
)

func buildAdmin(container *Container) *adminhandler.Handler {
	createLecturerUC := createlectureruc.NewUseCase(
		createlectureruc.NewFactory(container.DB, container.Logger),
		*container.LDAP,
		container.Logger,
	)
	deleteLecturerUC := deletelectureruc.NewUseCase(
		deletelectureruc.NewFactory(container.DB, container.Logger),
		container.Storage,
		container.Logger,
	)
	getLectureresUC := getlecturersuc.NewUseCase(
		getlecturersuc.NewFactory(container.DB, container.Logger),
		container.Logger,
	)
	return adminhandler.New(
		createLecturerUC,
		deleteLecturerUC,
		getLectureresUC,
	)
}
