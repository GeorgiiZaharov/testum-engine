package gethardtasks

import (
	"context"
	"testum-engine/app/internal/repository/task"
)

type accessRepository interface {
	HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type taskRepository interface {
	GetHardTasks(ctx context.Context, testID int) ([]task.Task, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access accessRepository
	Task   taskRepository
}
