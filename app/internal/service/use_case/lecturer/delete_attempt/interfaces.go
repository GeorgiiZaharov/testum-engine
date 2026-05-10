package deleteattempt

import (
	"context"
)

type accessRepository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type resultRepository interface {
	DeleteAttempt(ctx context.Context, testID int, userID int) error
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access accessRepository
	Result resultRepository
}
