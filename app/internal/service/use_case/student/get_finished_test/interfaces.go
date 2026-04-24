package getfinishedtest

import (
	"context"

	studenttestrepo "testum-engine/app/internal/repository/student_test"
)

type studentTestRepository interface {
	GetFinishedTests(ctx context.Context, userID int) ([]studenttestrepo.StudentFinishedTestInfo, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	StudentTest studentTestRepository
}
