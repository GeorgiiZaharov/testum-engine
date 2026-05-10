package uploadpicture

import (
	"bytes"
	"context"
	"go.uber.org/zap"
)

type UseCase struct {
	factory repoFactory
	storage storageAdapter
	log     *zap.Logger
	baseURL string
}

func NewUseCase(factory repoFactory, storage storageAdapter, log *zap.Logger, baseURL string) *UseCase {
	return &UseCase{
		factory: factory,
		storage: storage,
		log:     log,
		baseURL: baseURL,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req UploadPictureRequest) (UploadPictureResponse, error) {
	if req.UserID <= 0 || len(req.File) == 0 || req.FileName == "" {
		return UploadPictureResponse{}, ErrInvalidInput
	}

	var resp UploadPictureResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {
		// 1. Проверяем пользователя
		user, err := r.User.GetByID(ctx, req.UserID)
		if err != nil {
			uc.log.Error("failed to get user", zap.Error(err))
			return err
		}

		if !user.IsLecturer {
			return ErrAccessDenied
		}

		// 2. Загружаем картинку
		url, err := uc.storage.UploadPicture(bytes.NewReader(req.File), req.FileName, user.Login)
		if err != nil {
			uc.log.Error("failed to upload picture", zap.Error(err))
			return ErrStorageFailed
		}

		resp.URL = uc.baseURL + url
		resp.Success = true

		return nil
	})

	if err != nil {
		return UploadPictureResponse{}, err
	}

	return resp, nil
}
