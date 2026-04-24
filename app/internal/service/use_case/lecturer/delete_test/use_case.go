package deletetest

import (
	"context"

	"go.uber.org/zap"
)

type UseCase struct {
	factory repoFactory
	storage storageAdapter
	log     *zap.Logger
}

func NewUseCase(factory repoFactory, storage storageAdapter, log *zap.Logger) *UseCase {
	return &UseCase{
		factory: factory,
		storage: storage,
		log:     log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req DeleteTestRequest) (DeleteTestResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 {
		return DeleteTestResponse{}, ErrInvalidInput
	}

	var response DeleteTestResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа лектора к тесту
		hasAccess, err := r.Access.HasLecturerAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check lecturer access", zap.Error(err))
			return err
		}

		if !hasAccess {
			return ErrAccessDenied
		}

		// 2. Получение информации о тесте
		test, err := r.LecturerTest.GetByID(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to get test", zap.Error(err))
			return ErrTestNotFound
		}

		// 3. Удаление файла из storage
		if test.FileName != "" {
			if err := uc.storage.DeleteFile(test.FileName); err != nil {
				uc.log.Error("failed to delete file from storage", zap.Error(err))
				return ErrStorageFailed
			}
		}

		// 4. Удаление теста (каскад: задания, ответы, доступы — в репозитории/БД)
		ok, err := r.LecturerTest.Delete(ctx, req.TestID)
		if err != nil {
			uc.log.Error("failed to delete test", zap.Error(err))
			return err
		}

		if !ok {
			return ErrTestNotFound
		}

		response.Success = true
		return nil
	})

	if err != nil {
		return DeleteTestResponse{}, err
	}

	return response, nil
}
