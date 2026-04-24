package login

import (
	"context"

	ldapadapter "testum-engine/app/internal/adapter/ldap"
	userrepo "testum-engine/app/internal/repository/user"
)

type userRepository interface {
	Upsert(ctx context.Context, params userrepo.CreateUserParams) (int, error)
}

type ldapAdapter interface {
	ValidatePassword(ctx context.Context, login, password string) error
	GetInfo(ctx context.Context, login string) (*ldapadapter.LdapUserInfo, error)
}

type authService interface {
	GenerateAccess(userID int) (string, error)
	GenerateRefresh(userID int) (string, error)
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User userRepository
}
