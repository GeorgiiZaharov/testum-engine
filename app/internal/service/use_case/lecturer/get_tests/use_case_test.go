package gettests

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	testrepo "testum-engine/app/internal/repository/lecturer_test"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

type testEnv struct {
	uc   *UseCase
	ctx  context.Context
	repo testrepo.Repository
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
	t.Cleanup(
		func() {
			cleanup()
			fmt.Println(123234)
		},
	)

	fx := fixtures.New(database)
	require.NoError(t, fx.Reset(ctx))
	require.NoError(t, fx.SeedAll(ctx))

	repo := testrepo.NewRepository(database, zap.NewNop())

	uc := NewUseCase(
		NewFactory(database, zap.NewNop()),
		zap.NewNop(),
	)

	return &testEnv{
		uc:   uc,
		ctx:  ctx,
		repo: repo,
	}
}

func TestGetTests_Success(t *testing.T) {
	env := setup(t)

	// lecturer1 (id=6) имеет тесты 2 и 3
	resp, err := env.uc.Execute(env.ctx, GetTestsRequest{
		UserID: 6,
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Tests)

	assert.Len(t, resp.Tests, 2)

	ids := []int{resp.Tests[0].ID, resp.Tests[1].ID}
	assert.Contains(t, ids, 2)
	assert.Contains(t, ids, 3)
}

func TestGetTests_Empty(t *testing.T) {
	env := setup(t)

	// студент не имеет тестов
	resp, err := env.uc.Execute(env.ctx, GetTestsRequest{
		UserID: 1,
	})

	require.NoError(t, err)
	assert.Len(t, resp.Tests, 0)
}

func TestGetTests_InvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, GetTestsRequest{
		UserID: 0,
	})

	require.Error(t, err)
	assert.Equal(t, ErrInvalidInput, err)
}

func TestGetTests_AnotherLecturer(t *testing.T) {
	env := setup(t)

	// lecturer2 (id=7) имеет один тест (id=4)
	resp, err := env.uc.Execute(env.ctx, GetTestsRequest{
		UserID: 7,
	})

	require.NoError(t, err)
	require.Len(t, resp.Tests, 1)

	assert.Equal(t, 4, resp.Tests[0].ID)
	assert.Equal(t, "Programming Basics", resp.Tests[0].Name)
}
