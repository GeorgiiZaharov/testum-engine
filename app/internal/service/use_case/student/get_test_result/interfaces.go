package gettestresult

import (
	"context"

	"testum-engine/app/internal/repository/result"
)

type accessRepository interface {
	HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type resultRepository interface {
	GetStudentResult(ctx context.Context, userID int, testID int) (result.TestResult, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access accessRepository
	Result resultRepository
}
