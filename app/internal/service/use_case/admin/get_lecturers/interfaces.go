package getlecturers

import (
	"context"

	userrepo "testum-engine/app/internal/repository/user"
)

type userRepository interface {
	GetByID(ctx context.Context, userID int) (userrepo.User, error)
	GetLecturers(ctx context.Context) ([]userrepo.User, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User userRepository
}
