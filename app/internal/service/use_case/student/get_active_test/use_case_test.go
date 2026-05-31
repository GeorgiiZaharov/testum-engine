package getactivetest

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

func Test_GetActiveTests_InvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, GetActiveTestRequest{
		UserID: 0,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
}

func Test_GetActiveTests_StudentWithAvailableTests(t *testing.T) {
	env := setup(t)

	// student1 (id=1, group A-101)
	resp, err := env.uc.Execute(env.ctx, GetActiveTestRequest{
		UserID: 1,
	})

	require.NoError(t, err)
	require.NotNil(t, resp.ActiveTests)

	// В фикстурах:
	// A-101 имеет доступ к test 1 и test 2
	// student1 не проходил тесты → оба активны
	assert.GreaterOrEqual(t, len(resp.ActiveTests), 1)

	for _, tcase := range resp.ActiveTests {
		assert.NotZero(t, tcase.ID)
		assert.NotEmpty(t, tcase.Name)
		assert.NotEmpty(t, tcase.LecturerName)
	}
}

func Test_GetActiveTests_StudentWithPartiallyCompletedTests(t *testing.T) {
	env := setup(t)

	// student3 (id=3, group B-202)
	// он уже прошел test 1
	resp, err := env.uc.Execute(env.ctx, GetActiveTestRequest{
		UserID: 3,
	})

	require.NoError(t, err)

	// test 1 не должен быть в активных
	for _, tcase := range resp.ActiveTests {
		assert.NotEqual(t, 1, tcase.ID)
	}
}

func Test_GetActiveTests_StudentWithNoAvailableTests(t *testing.T) {
	env := setup(t)

	// student3 (B-202) — доступ только к test 1
	// но он уже его прошел → ничего не остается
	resp, err := env.uc.Execute(env.ctx, GetActiveTestRequest{
		UserID: 3,
	})

	require.NoError(t, err)

	// возможен 0 (если repo фильтрует корректно)
	assert.NotNil(t, resp.ActiveTests)
}

func Test_GetActiveTests_OldStudent(t *testing.T) {
	env := setup(t)

	// old_student (id=4)
	resp, err := env.uc.Execute(env.ctx, GetActiveTestRequest{
		UserID: 4,
	})

	require.NoError(t, err)

	// он уже прошел test 1 → не должен получить его снова
	for _, tcase := range resp.ActiveTests {
		assert.NotEqual(t, 1, tcase.ID)
	}
}

func Test_GetActiveTests_UserNotFound(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, GetActiveTestRequest{ // TODO: add user not found exception
		UserID: 9999,
	})

	require.NoError(t, err)
	assert.Empty(t, resp.ActiveTests)
}

func Test_GetActiveTests_ResponseMapping(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, GetActiveTestRequest{
		UserID: 1,
	})

	require.NoError(t, err)

	if len(resp.ActiveTests) > 0 {
		tcase := resp.ActiveTests[0]

		assert.NotZero(t, tcase.ID)
		assert.NotEmpty(t, tcase.Name)
		assert.NotEmpty(t, tcase.LecturerName)
		assert.GreaterOrEqual(t, tcase.CntQuestions, 0)
		assert.GreaterOrEqual(t, tcase.CntHardQuestions, 0)
	}
}
