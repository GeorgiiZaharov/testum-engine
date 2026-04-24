package refresh

import (
	"context"

	userrepo "testum-engine/app/internal/repository/user"
)

type authService interface {
	GenerateAccess(userID int) (string, error)
	GenerateRefresh(userID int) (string, error)
}

type userRepository interface {
	GetByID(ctx context.Context, userID int) (userrepo.User, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User userRepository
}
