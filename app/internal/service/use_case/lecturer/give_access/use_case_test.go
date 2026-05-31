package giveaccess

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

func TestGiveAccess_Success(t *testing.T) {
	env := setup(t)

	req := GiveAccessRequest{
		UserID: 6, // lecturer2 (owner of test_id=2)
		TestID: 2, // Linear Algebra
		Group:  "A-101",
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.True(t, res.Success)
}

func TestGiveAccess_InvalidInput(t *testing.T) {
	env := setup(t)

	reqs := []GiveAccessRequest{
		{UserID: 0, TestID: 2, Group: "A-101"},
		{UserID: 6, TestID: 0, Group: "A-101"},
		{UserID: 6, TestID: 2, Group: ""},
	}

	for _, req := range reqs {
		_, err := env.uc.Execute(env.ctx, req)
		require.Error(t, err)
		require.Equal(t, ErrInvalidInput, err)
	}
}

func TestGiveAccess_AccessDenied(t *testing.T) {
	env := setup(t)

	req := GiveAccessRequest{
		UserID: 6, // lecturer2
		TestID: 1, // owned by lecturer1 (user 5)
		Group:  "A-101",
	}

	_, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.Equal(t, ErrAccessDenied, err)
}

func TestGiveAccess_AlreadyHasAccess(t *testing.T) {
	env := setup(t)

	req := GiveAccessRequest{
		UserID: 6,
		TestID: 2,
		Group:  "A-101",
	}

	// first grant
	res1, err := env.uc.Execute(env.ctx, req)
	require.NoError(t, err)
	require.True(t, res1.Success)

	// second grant (idempotent expectation)
	res2, err := env.uc.Execute(env.ctx, req)
	require.NoError(t, err)
	require.True(t, res2.Success)
}
