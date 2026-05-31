package getfinishedtest

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

// =========================
// TEST ENV
// =========================

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

// =========================
// TESTS
// =========================

func Test_GetFinishedTests_InvalidInput(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, GetFinishedTestRequest{
		UserID: 0,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Empty(t, resp.FinishedTests)
}

func Test_GetFinishedTests_NoFinishedTests(t *testing.T) {
	env := setup(t)

	// student1 (id=1) не имеет завершённых тестов
	resp, err := env.uc.Execute(env.ctx, GetFinishedTestRequest{
		UserID: 1,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.FinishedTests)
}

func Test_GetFinishedTests_Success(t *testing.T) {
	env := setup(t)

	// student3 (id=3) имеет завершённый тест (см. fixtures)
	resp, err := env.uc.Execute(env.ctx, GetFinishedTestRequest{
		UserID: 3,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.FinishedTests)

	test := resp.FinishedTests[0]

	assert.Equal(t, 1, test.ID)
	assert.Equal(t, "Math Basics", test.Name)
	assert.Equal(t, "Magistr Magistr", test.LecturerName)
	assert.Equal(t, 2, test.CntQuestions)
	assert.Equal(t, 1, test.CntHardQuestions)
	assert.Equal(t, 4, test.Mark)
	assert.Equal(t, 80.0, test.SuccessRate)

	assert.False(t, test.DateStart.IsZero())
	assert.False(t, test.DateEnd.IsZero())
}

func Test_GetFinishedTests_MultipleResults(t *testing.T) {
	env := setup(t)

	// student4 (id=4) тоже имеет завершённый тест
	resp, err := env.uc.Execute(env.ctx, GetFinishedTestRequest{
		UserID: 4,
	})

	require.NoError(t, err)
	require.Len(t, resp.FinishedTests, 1)

	test := resp.FinishedTests[0]

	assert.Equal(t, 1, test.ID)
	assert.Equal(t, 5, test.Mark)
	assert.Equal(t, 100.0, test.SuccessRate)
}

func Test_GetFinishedTests_UserNotFound(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, GetFinishedTestRequest{
		UserID: 9999,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.FinishedTests)
}
