package createlecturer

import (
	"context"

	ldapapdap "testum-engine/app/internal/adapter/ldap"
	userrepo "testum-engine/app/internal/repository/user"
)

type userRepository interface {
	GetByID(ctx context.Context, userID int) (userrepo.User, error)
	Upsert(ctx context.Context, params userrepo.CreateUserParams) (int, error)
	CreateLecturer(ctx context.Context, userID int) error
}

type ldapAdapter interface {
	GetInfo(ctx context.Context, login string) (*ldapapdap.LdapUserInfo, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User userRepository
}
