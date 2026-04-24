package gettestresult

import (
	"context"
	"time"

	resultrepo "testum-engine/app/internal/repository/result"
)

type AccessRepository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type ResultRepository interface {
	GetGroupResult(ctx context.Context, testID int, group string, year int) ([]resultrepo.StudentResult, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access AccessRepository
	Result ResultRepository
}

// DTO из repo слоя (используем напрямую)
type StudentResult struct {
	UserID int
	Name   string
	Login  string
	Mail   string

	Mark        *int
	SuccessRate *float64
	DateStart   *time.Time
	DateEnd     *time.Time
}
