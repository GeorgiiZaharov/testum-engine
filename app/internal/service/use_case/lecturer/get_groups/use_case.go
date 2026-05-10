package getgroups

import (
	"context"
	"errors"

	"go.uber.org/zap"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
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

func (uc *UseCase) Execute(ctx context.Context, req GetGroupsRequest) (GetGroupsResponse, error) {
	if req.UserID <= 0 || req.TestID <= 0 || req.Year < 0 {
		return GetGroupsResponse{}, ErrInvalidInput
	}

	var response GetGroupsResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка доступа лектора
		ok, err := r.Access.HasLecturerAccess(ctx, req.UserID, req.TestID)
		if err != nil {
			uc.log.Error("failed to check lecturer access",
				zap.Error(err),
				zap.Int("user_id", req.UserID),
				zap.Int("test_id", req.TestID),
			)
			return err
		}

		if !ok {
			return ErrAccessDenied
		}

		// 2. Получение групп
		groups, err := r.Test.GetGroups(ctx, req.TestID, req.Year)
		if errors.Is(err, lecturertestrepo.ErrTestNotFound) {
			groups = []lecturertestrepo.GroupInfo{}
		} else if err != nil {
			uc.log.Error("failed to get groups",
				zap.Error(err),
				zap.Int("test_id", req.TestID),
				zap.Int("year", req.Year),
			)
			return err
		}

		// 3. Маппинг repo → use case
		result := make([]GroupInfo, 0, len(groups))

		for _, g := range groups {
			result = append(result, GroupInfo{
				GroupName:    g.GroupName,
				MembersCount: g.MembersCount,
			})
		}

		response.Groups = result
		return nil
	})

	if err != nil {
		return GetGroupsResponse{}, err
	}

	return response, nil
}
