package deletelecturer

import (
	"context"
	"errors"

	"go.uber.org/zap"

	userrepo "testum-engine/app/internal/repository/user"
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

func (uc *UseCase) Execute(ctx context.Context, req DeleteLecturerRequest) (DeleteLecturerResponse, error) {
	if req.AdminID <= 0 || req.LecturerID <= 0 {
		return DeleteLecturerResponse{}, ErrInvalidInput
	}

	var response DeleteLecturerResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка администратора
		admin, err := r.User.GetByID(ctx, req.AdminID)
		if err != nil {
			uc.log.Error("failed to get admin", zap.Error(err))
			return err
		}

		if !isAdmin(admin) {
			return ErrForbidden
		}

		// 2. Получение лектора
		lecturer, err := r.User.GetByID(ctx, req.LecturerID)
		if errors.Is(err, userrepo.ErrUserNotFound) {
			return ErrNotFound
		} else if err != nil {
			uc.log.Error("failed to get lecturer", zap.Error(err))
			return err
		}

		if !lecturer.IsLecturer {
			return ErrNotLecturer
		}

		// 3. Получение файлов тестов
		files, err := r.File.GetAllTestFiles(ctx, lecturer.ID)
		if err != nil {
			uc.log.Error("failed to get files", zap.Error(err))
			return err
		}

		// 4. Удаление файлов
		for _, file := range files {
			if err := uc.storage.DeleteFile(file); err != nil {
				uc.log.Error("failed to delete file", zap.Error(err))
				return err
			}
		}

		// 5. Удаление картинок
		if err := uc.storage.DeletePictures(lecturer.Login); err != nil {
			uc.log.Error("failed to delete pictures", zap.Error(err))
			return err
		}

		// 6. Удаление роли лектора
		if err := r.User.DeleteLecturer(ctx, lecturer.ID); err != nil {
			uc.log.Error("failed to delete lecturer role", zap.Error(err))
			return err
		}

		response.Success = true
		return nil
	})

	if err != nil {
		return DeleteLecturerResponse{}, err
	}

	return response, nil
}

func isAdmin(user userrepo.User) bool {
	return user.Login == "olbgvl" || user.Login == "lector" || user.Login == "vasilenk"
}
