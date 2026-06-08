package createlecturer

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	userrepo "testum-engine/app/internal/repository/user"
)

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

func (uc *UseCase) Execute(ctx context.Context, req CreateLecturerRequest) (CreateLecturerResponse, error) {
	if req.AdminID <= 0 || req.Login == "" {
		return CreateLecturerResponse{}, ErrInvalidInput
	}

	var response CreateLecturerResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка администратора
		admin, err := r.User.GetByID(ctx, req.AdminID)
		if err != nil {
			uc.log.Error("failed to get admin user", zap.Error(err))
			return err
		}

		if !isAdmin(admin) {
			return ErrForbidden
		}

		// 2. Получение данных из LDAP
		ldapUser, err := uc.ldap.GetInfo(ctx, req.Login)
		if err != nil {
			uc.log.Error("ldap get info failed", zap.Error(err))
			return ErrLDAPFailed
		}

		// 3. Upsert пользователя
		userID, err := r.User.Upsert(ctx, userrepo.CreateUserParams{
			Login: ldapUser.Login,
			Name:  ldapUser.Name,
			Mail:  ldapUser.Mail,
			Group: ldapUser.Group,
		})
		if err != nil {
			uc.log.Error("failed to upsert user", zap.Error(err))
			return err
		}

		// 4. Назначение роли лектора
		if err := r.User.CreateLecturer(ctx, userID); err != nil {
			fmt.Println(11111111)
			uc.log.Error("failed to create lecturer", zap.Error(err))
			return err
		}

		response.Success = true
		return nil
	})

	if err != nil {
		return CreateLecturerResponse{}, err
	}

	return response, nil
}

func isAdmin(user userrepo.User) bool {
	return user.Login == "olbgvl" || user.Login == "lector" || user.Login == "vasilenk"
}
