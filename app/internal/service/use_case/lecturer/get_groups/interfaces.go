package getgroups

import (
	"context"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
)

type accessRepository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type lecturerTestRepository interface {
	GetGroups(ctx context.Context, testID int, year int) ([]lecturertestrepo.GroupInfo, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access accessRepository
	Test   lecturerTestRepository
}
