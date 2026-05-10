package app

import (
	lecturerfilehandler "testum-engine/app/internal/handler/lecturer/file"

	gettestfileuc "testum-engine/app/internal/service/use_case/lecturer/get_test_file"
	uploadpictureuc "testum-engine/app/internal/service/use_case/lecturer/upload_picture"
	uploadtestuc "testum-engine/app/internal/service/use_case/lecturer/upload_test"
)

func buildLecturerFile(container *Container) *lecturerfilehandler.Handler {
	uploadTestUC := uploadtestuc.NewUseCase(
		uploadtestuc.NewFactory(container.DB, container.Logger),
		container.Storage,
		container.ValidationService,
		container.LatexValidationService,
		container.Logger,
	)
	getTestFileUC := gettestfileuc.NewUseCase(
		gettestfileuc.NewFactory(container.DB, container.Logger),
		container.Storage,
		container.Logger,
	)
	uploadPictureUC := uploadpictureuc.NewUseCase(
		uploadpictureuc.NewFactory(container.DB, container.Logger),
		container.Storage,
		container.Logger,
		container.BasePictureURL,
	)
	return lecturerfilehandler.New(uploadTestUC, getTestFileUC, uploadPictureUC)
}
