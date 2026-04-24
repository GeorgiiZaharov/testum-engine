package getlecturers

import (
	"context"
	"time"

	"go.uber.org/zap"

	userrepo "testum-engine/app/internal/repository/user"
)

type UseCase struct {
	factory repoFactory
	log     *zap.Logger
}

func NewUseCase(factory repoFactory, log *zap.Logger) *UseCase {
	return &UseCase{
		factory: factory,
		log:     log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req GetLecturersRequest) (GetLecturersResponse, error) {
	if req.UserID <= 0 {
		return GetLecturersResponse{}, ErrInvalidInput
	}

	var response GetLecturersResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка администратора
		admin, err := r.User.GetByID(ctx, req.UserID)
		if err != nil {
			uc.log.Error("failed to get admin", zap.Error(err))
			return err
		}

		if !isAdmin(admin) {
			return ErrForbidden
		}

		// 2. Получение списка лекторов
		users, err := r.User.GetLecturers(ctx)
		if err != nil {
			uc.log.Error("failed to get lecturers", zap.Error(err))
			return err
		}

		// 3. Маппинг в response
		lecturers := make([]Lecturer, 0, len(users))

		for _, u := range users {
			lecturers = append(lecturers, Lecturer{
				ID:           u.ID,
				Login:        u.Login,
				Mail:         u.Mail,
				Name:         u.Name,
				Group:        u.Group,
				DateCreated:  u.DateCreated.Format(time.RFC3339),
				DateModified: u.DateModified.Format(time.RFC3339),
			})
		}

		response.Lecturers = lecturers
		return nil
	})

	if err != nil {
		return GetLecturersResponse{}, err
	}

	return response, nil
}

func isAdmin(user userrepo.User) bool {
	return user.Login == "olbgvl" || user.Login == "lector"
}
