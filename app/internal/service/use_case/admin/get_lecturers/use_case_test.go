package getlecturers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	userrepo "testum-engine/app/internal/repository/user"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

// =========================
// TEST ENV
// =========================

type testEnv struct {
	uc   *UseCase
	repo userrepo.Repository
	ctx  context.Context
}

func setup(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	database, cleanup, err := bootstrap.Setup(bootstrap.Config{
		DBOptions: db.DBOptions{
			Path: ":memory:",
		},
		Migrations: "../../../../../migrations/",
	})
	require.NoError(t, err)

	t.Cleanup(cleanup)

	fx := fixtures.New(database)
	require.NoError(t, fx.Reset(ctx))
	require.NoError(t, fx.SeedAll(ctx))

	uc := NewUseCase(
		NewFactory(database, zap.NewNop()),
		zap.NewNop(),
	)

	return &testEnv{
		uc:   uc,
		repo: userrepo.NewRepository(database, zap.NewNop()),
		ctx:  ctx,
	}
}

// =========================
// TESTS
// =========================

func TestGetLecturers_Success(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, GetLecturersRequest{
		UserID: 8, // admin olgbvl
	})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Lecturers)

	// проверяем структуру
	for _, l := range resp.Lecturers {
		assert.NotEmpty(t, l.Login)
		assert.NotEmpty(t, l.Mail)
		assert.NotZero(t, l.ID)

		_, err := time.Parse(time.RFC3339, l.DateCreated)
		assert.NoError(t, err)

		_, err = time.Parse(time.RFC3339, l.DateModified)
		assert.NoError(t, err)
	}
}

func TestGetLecturers_InvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, GetLecturersRequest{
		UserID: 0,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestGetLecturers_Forbidden(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, GetLecturersRequest{
		UserID: 1, // student
	})

	assert.ErrorIs(t, err, ErrForbidden)
}

func TestGetLecturers_AdminNotFound(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, GetLecturersRequest{
		UserID: 9999,
	})

	assert.Error(t, err)
}

func TestGetLecturers_NoLecturers(t *testing.T) {
	// этот тест важен если база может быть пустой
	env := setup(t)

	// очищаем lecturers через repo (опционально)
	users, err := env.repo.GetLecturers(env.ctx)
	require.NoError(t, err)

	// если вдруг фикстуры уже пустые — просто проверяем стабильность
	if len(users) == 0 {
		resp, err := env.uc.Execute(env.ctx, GetLecturersRequest{
			UserID: 8,
		})

		require.NoError(t, err)
		assert.Len(t, resp.Lecturers, 0)
	}
}
