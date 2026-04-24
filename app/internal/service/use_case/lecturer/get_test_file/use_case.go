package gettestfile

import (
	"context"
	"go.uber.org/zap"
)

type UseCase struct {
	factory repoFactory
	storage StorageAdapter
	log     *zap.Logger
}

func NewUseCase(factory repoFactory, storage StorageAdapter, log *zap.Logger) *UseCase {
	return &UseCase{
		factory: factory,
		storage: storage,
		log:     log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req GetTestFileRequest) (GetTestFileResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 {
		return GetTestFileResponse{}, ErrForbidden
	}

	var response GetTestFileResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа лектора
		hasAccess, err := r.Access.HasLecturerAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check access", zap.Error(err))
			return err
		}

		if !hasAccess {
			return ErrForbidden
		}

		// 2. Получение теста
		test, err := r.LecturerTest.GetByID(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get test", zap.Error(err))
			return err
		}

		// 3. Проверка наличия файла
		if test.FileName == "" {
			return ErrFileNotFound
		}

		// 4. Получение файла из хранилища
		file, err := uc.storage.GetFile(test.FileName)
		if err != nil || file == nil {
			uc.log.Error("failed to get file from storage", zap.Error(err))
			return ErrFileNotFound
		}

		// 5. Формирование ответа
		response = GetTestFileResponse{
			File: file,
		}

		return nil
	})

	if err != nil {
		return GetTestFileResponse{}, err
	}

	return response, nil
}
