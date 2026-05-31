package gethardtasks

import (
	"context"
	"testing"
	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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

	uc := NewUseCase(
		NewFactory(database, zap.NewNop()),
		zap.NewNop(),
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

func Test_GetHardTasks_Success(t *testing.T) {
	env := setup(t)

	req := GetHardTasksRequest{
		UserID: 2, // student3 (B-202)
		TestID: 1, // Linear Algebra
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.NotEmpty(t, resp.HardTasks)

	assert.True(t, resp.HardTasks[0].IsHard)
	assert.Len(t, resp.HardTasks[0].Answers, 2)
}

func Test_GetHardTasks_AccessDenied(t *testing.T) {
	env := setup(t)

	req := GetHardTasksRequest{
		UserID: 3, // student3 (B-202)
		TestID: 2,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrAccessDenied, err)
	assert.Empty(t, resp.HardTasks)
}

func Test_GetHardTasks_InvalidInput(t *testing.T) {
	env := setup(t)

	req := GetHardTasksRequest{
		UserID: 0, // invalid UserID
		TestID: 0, // invalid TestID
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
	assert.Empty(t, resp.HardTasks)
}
