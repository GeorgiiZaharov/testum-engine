package deletelecturer

import (
	"context"

	userrepo "testum-engine/app/internal/repository/user"
)

type userRepository interface {
	GetByID(ctx context.Context, userID int) (userrepo.User, error)
	DeleteLecturer(ctx context.Context, lecturerID int) error
}

type fileRepository interface {
	GetAllTestFiles(ctx context.Context, lecturerID int) ([]string, error)
}

type storageAdapter interface {
	DeleteFile(fileName string) error
	DeletePictures(login string) error
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User userRepository
	File fileRepository
}
