package getbasetasks

import (
	"context"
	"fmt"
	"testing"
	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"

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

func Test_GetBaseTasks_Success(t *testing.T) {
	env := setup(t)

	req := GetBaseTasksRequest{
		UserID: 2, // student1 (A-101)
		TestID: 1, // Math Basics
	}

	resp, err := env.uc.Execute(env.ctx, req)

	fmt.Println(resp, err)
	require.NoError(t, err)
	require.NotEmpty(t, resp.BaseTasks)

	// assert.Equal(t, false, resp.BaseTasks[0].IsHard)
}

func Test_GetBaseTasks_InvalidInput(t *testing.T) {
	env := setup(t)

	// Test with invalid UserID
	req := GetBaseTasksRequest{
		UserID: 0, // invalid UserID
		TestID: 1,
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	require.Empty(t, resp.BaseTasks)

	// Test with invalid TestID
	req = GetBaseTasksRequest{
		UserID: 2,
		TestID: 0, // invalid TestID
	}

	resp, err = env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	require.Empty(t, resp.BaseTasks)
}

func Test_GetBaseTasks_AccessDenied(t *testing.T) {
	env := setup(t)

	req := GetBaseTasksRequest{
		UserID: 3, // student3 (B-202) - no access to Test 2
		TestID: 2, // Math Basics
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	require.Equal(t, ErrAccessDenied, err)
	require.Empty(t, resp.BaseTasks)
}

func Test_GetBaseTasks_TestAlreadyCompleted(t *testing.T) {
	env := setup(t)

	req := GetBaseTasksRequest{
		UserID: 4, // old_student
		TestID: 1, // Math Basics
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	require.Equal(t, ErrTestCompleted, err)
	require.Empty(t, resp.BaseTasks)
}

func Test_GetBaseTasks_ErrorOnGetTasks(t *testing.T) {
	env := setup(t)

	req := GetBaseTasksRequest{
		UserID: 2,    // student1 (A-101)
		TestID: 9999, // non-existing TestID
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	require.Empty(t, resp.BaseTasks)
}
