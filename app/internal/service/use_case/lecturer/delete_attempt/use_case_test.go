package deleteattempt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

// =========================
// Test environment
// =========================
type testEnv struct {
	uc  *UseCase
	ctx context.Context
	fx  *fixtures.Manager
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
		fx:  fx,
	}
}

// =========================
// Tests
// =========================

// -------------------- 1. Успешное удаление попытки --------------------
func TestDeleteAttempt_Success(t *testing.T) {
	env := setup(t)

	req := DeleteAttemptRequest{
		LecturerID: 5, // владельцем теста 1 является magistr (id=5)
		UserID:     2, // student2
		TestID:     1, // Math Basics
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.True(t, resp.Success)
}

// -------------------- 2. Попытка доступа без прав --------------------
func TestDeleteAttempt_AccessDenied(t *testing.T) {
	env := setup(t)

	req := DeleteAttemptRequest{
		LecturerID: 7, // lecturer2 не владеет тестом 1
		UserID:     2,
		TestID:     1,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.Equal(t, ErrAccessDenied, err)
	require.False(t, resp.Success)
}

// -------------------- 3. Некорректный ввод --------------------
func TestDeleteAttempt_InvalidInput(t *testing.T) {
	env := setup(t)

	req := DeleteAttemptRequest{
		LecturerID: 0,
		UserID:     0,
		TestID:     0,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.Equal(t, ErrInvalidInput, err)
	require.False(t, resp.Success)
}

// -------------------- 4. Попытка удалить несуществующую запись --------------------
func TestDeleteAttempt_NotFound(t *testing.T) {
	env := setup(t)

	req := DeleteAttemptRequest{
		LecturerID: 5,
		UserID:     999, // несуществующий студент
		TestID:     1,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.True(t, resp.Success)
}
