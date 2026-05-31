package gettest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

	uc := NewUseCase(
		NewFactory(database, zap.NewNop()),
		zap.NewNop(),
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

func TestGetTest_Success(t *testing.T) {
	env := setup(t)

	req := GetTestRequest{
		UserID: 5, // lecturer (magistr)
		TestID: 1, // Math Basics
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)

	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "Math Basics", resp.Name)

	assert.Greater(t, len(resp.Groups), 0)
	assert.GreaterOrEqual(t, len(resp.HardTasks), 0)
	assert.GreaterOrEqual(t, len(resp.BaseTasks), 0)
}

func TestGetTest_AccessDenied(t *testing.T) {
	env := setup(t)

	req := GetTestRequest{
		UserID: 7, // lecturer2
		TestID: 1, // owned by another lecturer
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrAccessDenied, err)
	assert.Equal(t, GetTestResponse{}, resp)
}

func TestGetTest_InvalidInput(t *testing.T) {
	env := setup(t)

	req := GetTestRequest{
		UserID: 0,
		TestID: 0,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
	assert.Equal(t, GetTestResponse{}, resp)
}

func TestGetTest_NotFound(t *testing.T) {
	env := setup(t)

	req := GetTestRequest{
		UserID: 5,
		TestID: 9999,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, GetTestResponse{}, resp)
}

func TestGetTest_TaskSplitValidation(t *testing.T) {
	env := setup(t)

	req := GetTestRequest{
		UserID: 6, // lecturer1
		TestID: 3, // Physics Intro
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)

	var hasHard bool
	var hasBase bool

	for _, task := range resp.HardTasks {
		if task.IsHard {
			hasHard = true
		}
	}

	for _, task := range resp.BaseTasks {
		if !task.IsHard {
			hasBase = true
		}
	}

	assert.True(t, hasHard)
	assert.True(t, hasBase)
}
