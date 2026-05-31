package gettestresult

import (
	"context"
	"fmt"
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

func TestGetTestResult_Success_CurrentYear(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 6,       // lecturer1
		TestID: 2,       // принадлежит lecturer1
		Group:  "A-101", // есть доступ
		Year:   0,       // текущий
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp.Results)

	// Проверяем что структура корректная
	for _, r := range resp.Results {
		assert.NotZero(t, r.UserID)
		assert.NotEmpty(t, r.Login)
	}
	fmt.Println(resp)
	assert.Equal(t, len(resp.Results), 2)
}

func TestGetTestResult_Success_PreviousYear(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 5, // magistr владелец test 1
		TestID: 1,
		Group:  "A-101",
		Year:   1, // предыдущий год (есть old_student)
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp.Results)

	// В фикстурах есть old_student → проверяем что не пусто
	assert.GreaterOrEqual(t, len(resp.Results), 1)
}

func TestGetTestResult_NoAccess(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 7, // lecturer2
		TestID: 2, // принадлежит lecturer1
		Group:  "A-101",
		Year:   0,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Equal(t, GetTestResultResponse{}, resp)
}

func TestGetTestResult_InvalidInput(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 0,
		TestID: 0,
		Group:  "",
		Year:   0,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
	assert.Equal(t, GetTestResultResponse{}, resp)
}

func TestGetTestResult_EmptyGroup(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 6,
		TestID: 2,
		Group:  "NON_EXISTENT",
		Year:   0,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp.Results)
	assert.Len(t, resp.Results, 0)
}
