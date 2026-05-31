package takeaccess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

type testEnv struct {
	uc  *UseCase
	ctx context.Context
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

	factory := NewFactory(database, zap.NewNop())
	uc := NewUseCase(factory, zap.NewNop())

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

func TestTakeAccess_Success(t *testing.T) {
	env := setup(t)

	req := TakeAccessRequest{
		UserID: 6, // lecturer2 (owner of test_id=2)
		TestID: 2,
		Group:  "A-101", // есть в фикстурах
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.True(t, res.Success)
}

func TestTakeAccess_AccessDenied(t *testing.T) {
	env := setup(t)

	req := TakeAccessRequest{
		UserID: 6, // lecturer2
		TestID: 1, // принадлежит другому лектору
		Group:  "A-101",
	}

	_, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.Equal(t, ErrAccessDenied, err)
}

func TestTakeAccess_InvalidInput(t *testing.T) {
	env := setup(t)

	cases := []TakeAccessRequest{
		{UserID: 0, TestID: 1, Group: "A-101"},
		{UserID: 6, TestID: 0, Group: "A-101"},
		{UserID: 6, TestID: 1, Group: ""},
	}

	for _, c := range cases {
		_, err := env.uc.Execute(env.ctx, c)
		require.Error(t, err)
		require.Equal(t, ErrInvalidInput, err)
	}
}

func TestTakeAccess_AlreadyRemoved(t *testing.T) {
	env := setup(t)

	req := TakeAccessRequest{
		UserID: 6,
		TestID: 2,
		Group:  "A-101",
	}

	// первый вызов
	res1, err := env.uc.Execute(env.ctx, req)
	require.NoError(t, err)
	require.True(t, res1.Success)

	// второй вызов — запись уже удалена
	res2, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	require.False(t, res2.Success)
}

func TestTakeAccess_GroupNotExists(t *testing.T) {
	env := setup(t)

	req := TakeAccessRequest{
		UserID: 6,
		TestID: 2,
		Group:  "NON_EXISTENT",
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.False(t, res.Success)
}
