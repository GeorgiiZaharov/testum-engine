package gettests

import (
	"context"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
)

type lecturerTestRepository interface {
	GetByLecturer(ctx context.Context, userID int) ([]lecturertestrepo.TestInfo, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	LecturerTest lecturerTestRepository
}
