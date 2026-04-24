package gettestfile

import (
	"context"
	"os"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
)

type LecturerTestRepository interface {
	GetByID(ctx context.Context, testID int) (lecturertestrepo.TestInfo, error)
}

type AccessRepository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type StorageAdapter interface {
	GetFile(fileName string) (*os.File, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	LecturerTest LecturerTestRepository
	Access       AccessRepository
}
