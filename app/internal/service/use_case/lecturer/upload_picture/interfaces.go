package uploadpicture

import (
	"context"
	"io"

	userrepo "testum-engine/app/internal/repository/user"
)

type userRepository interface {
	GetByID(ctx context.Context, userID int) (userrepo.User, error)
}

type storageAdapter interface {
	UploadPicture(file io.Reader, fileName string, login string) (string, error)
}

// RepoFactory нужен для работы с транзакциями
type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User userRepository
}
