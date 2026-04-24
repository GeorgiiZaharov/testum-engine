package gettest

import (
	"context"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
	taskrepo "testum-engine/app/internal/repository/task"
)

type accessRepository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type lecturerTestRepository interface {
	GetByID(ctx context.Context, testID int) (lecturertestrepo.TestInfo, error)
	GetGroups(ctx context.Context, testID int, year int) ([]lecturertestrepo.GroupInfo, error)
}

type taskRepository interface {
	GetHardTasks(ctx context.Context, testID int) ([]taskrepo.Task, error)
	GetBaseTasks(ctx context.Context, testID int) ([]taskrepo.Task, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access       accessRepository
	LecturerTest lecturerTestRepository
	Task         taskRepository
}
