package gettestresult

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
			Host: "localhost",
			Port: "3306",
			User: "testum_user",
			Pass: "testum_pass",
			Name: "testum",
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

func Test_GetTestResult_Success(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 3, // student3 (B-202)
		TestID: 1, // Math Basics
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	assert.Equal(t, 4, *resp.Mark)           // предполагаемый балл
	assert.Equal(t, 80.0, *resp.SuccessRate) // предполагаемый процент успешных ответов
	assert.NotEmpty(t, resp.DateStart)
	assert.NotEmpty(t, *resp.DateEnd)
}

func Test_GetTestResult_AccessDenied(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 3, // student3 (B-202)
		TestID: 2, // Linear Algebra (доступ только A-101)
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrAccessDenied, err)
	assert.Empty(t, resp.Mark)
}

func Test_GetTestResult_InvalidInput(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 0, // invalid UserID
		TestID: 0, // invalid TestID
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
	assert.Empty(t, resp.Mark)
}

func Test_GetTestResult_ResultNotFound(t *testing.T) {
	env := setup(t)

	req := GetTestResultRequest{
		UserID: 2, // student2 (A-101)
		TestID: 2, // (студент не проходил)
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrResultNotFound, err)
	assert.Empty(t, resp.Mark)
}

