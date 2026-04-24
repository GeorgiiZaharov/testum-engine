package deletetest

import (
	"context"

	"testum-engine/app/internal/repository/lecturer_test"
)

type lecturerTestRepository interface {
	GetByID(ctx context.Context, testID int) (lecturertest.TestInfo, error)
	Delete(ctx context.Context, testID int) (bool, error)
}

type accessRepository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type storageAdapter interface {
	DeleteFile(fileName string) error
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	LecturerTest lecturerTestRepository
	Access       accessRepository
}
