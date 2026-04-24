package getme

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"

	ldapadap "testum-engine/app/internal/adapter/ldap"
	userrepo "testum-engine/app/internal/repository/user"
)

const refreshThreshold = 14 * 24 * time.Hour

type UseCase struct {
	factory repoFactory
	ldap    ldapAdapter
	log     *zap.Logger
}

func NewUseCase(factory repoFactory, ldap ldapAdapter, log *zap.Logger) *UseCase {
	return &UseCase{
		factory: factory,
		ldap:    ldap,
		log:     log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req GetMeRequest) (GetMeResponse, error) {
	if req.UserID <= 0 {
		return GetMeResponse{}, ErrInvalidInput
	}

	var response GetMeResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {
		user, err := r.User.GetByID(ctx, req.UserID)
		if err != nil {
			if errors.Is(err, userrepo.ErrUserNotFound) {
				return ErrNotFound
			}
			uc.log.Error("failed to get user", zap.Error(err))
			return err
		}

		// если данные устарели — выполняем обновление из LDAP
		if time.Since(user.DateModified) > refreshThreshold {
			ldapUser, err := uc.ldap.GetInfo(ctx, user.Login)
			if err != nil {
				if errors.Is(err, ldapadap.ErrUserNotFound) {
					return ErrNotFound
				}
				uc.log.Error("ldap failed", zap.Error(err))
				return ErrLDAPFailed
			}

			_, err = r.User.Upsert(ctx, userrepo.CreateUserParams{
				Login: ldapUser.Login,
				Name:  ldapUser.Name,
				Mail:  ldapUser.Mail,
				Group: ldapUser.Group,
			})
			if err != nil {
				uc.log.Error("upsert failed", zap.Error(err))
				return err
			}

			user.Name = ldapUser.Name
			user.Mail = ldapUser.Mail
			user.Group = ldapUser.Group
		}

		// заполняем ответ
		response = GetMeResponse{
			ID:         user.ID,
			Login:      user.Login,
			Mail:       user.Mail,
			Name:       user.Name,
			Group:      user.Group,
			IsLecturer: user.IsLecturer,
		}

		return nil
	})

	if err != nil {
		return GetMeResponse{}, err
	}

	return response, nil
}
