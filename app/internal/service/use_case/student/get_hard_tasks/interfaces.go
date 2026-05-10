package gethardtasks

import (
	"context"
	"testum-engine/app/internal/repository/task"

	resultrepo "testum-engine/app/internal/repository/result"
	studenttestrepo "testum-engine/app/internal/repository/student_test"
)

type studentTestRepository interface {
	StartTest(ctx context.Context, params studenttestrepo.StartTestParams) (bool, error)
}

type accessRepository interface {
	HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type resultRepository interface {
	GetStudentResult(ctx context.Context, userID int, testID int) (resultrepo.TestResult, error)
}

type taskRepository interface {
	GetHardTasks(ctx context.Context, testID int) ([]task.Task, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access      accessRepository
	Task        taskRepository
	Result      resultRepository
	StudentTest studentTestRepository
}
