package giveaccess

import "context"

type accessRepository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
	GiveAccess(ctx context.Context, testID int, group string) (bool, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	Access accessRepository
}
