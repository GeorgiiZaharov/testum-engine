package getbasetasks

import (
	"context"

	resultrepo "testum-engine/app/internal/repository/result"
	studenttestrepo "testum-engine/app/internal/repository/student_test"
	taskrepo "testum-engine/app/internal/repository/task"
)

type studentTestRepository interface {
	StartTest(ctx context.Context, params studenttestrepo.StartTestParams) (bool, error)
}

type taskRepository interface {
	GetBaseTasks(ctx context.Context, testID int) ([]taskrepo.Task, error)
}

type resultRepository interface {
	GetStudentResult(ctx context.Context, userID int, testID int) (resultrepo.TestResult, error)
}

type accessRepository interface {
	HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access      accessRepository
	StudentTest studentTestRepository
	Result      resultRepository
	Task        taskRepository
}
