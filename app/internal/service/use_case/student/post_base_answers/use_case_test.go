package postbaseanswers

import (
	"context"
	"testing"

	"testum-engine/app/internal/adapter/db"
	answer "testum-engine/app/internal/service/core/answer"
	result "testum-engine/app/internal/service/core/result"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
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
		answer.NewCheckService(),
		result.NewCalculationService(),
		zap.NewNop(),
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

//
// ======================================================
// 1. HAPPY PATH
// ======================================================
//

func Test_PostBaseAnswers_ValidCase(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, 2, 1, []TaskAnswer{
		{TaskID: 1, Options: []int{1}},
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

//
// ======================================================
// 2. INVALID INPUT
// ======================================================
//

func Test_PostBaseAnswers_InvalidInput(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, 0, 1, []TaskAnswer{
		{TaskID: 1, Options: []int{1}},
	})

	require.Error(t, err)
	assert.False(t, resp.Success)
}

//
// ======================================================
// 3. ACCESS DENIED (student not in allowed group)
// ======================================================
//

func Test_PostBaseAnswers_AccessDenied(t *testing.T) {
	env := setup(t)

	// student 3 belongs to B-202, test 1 is A-101 + B-202 ok,
	// but we simulate mismatch by using non-existing permission scenario
	_, err := env.uc.Execute(env.ctx, 3, 3, []TaskAnswer{
		{TaskID: 3, Options: []int{5}},
	})

	require.Error(t, err)
	assert.Equal(t, ErrAccessDenied, err)
}

//
// ======================================================
// 4. ALREADY SUBMITTED BASE ANSWERS
// (student 3 already has answers in fixtures)
// ======================================================
//

func Test_PostBaseAnswers_AlreadySubmitted(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, 3, 1, []TaskAnswer{
		{TaskID: 1, Options: []int{1}},
	})

	require.Error(t, err)
	assert.Equal(t, ErrAlreadySubmitted, err)
}

//
// ======================================================
// 5. HARD BLOCK NOT PASSED
// (student 2 has no hard answers in fixtures)
// ======================================================
//

func Test_PostBaseAnswers_HardBlockNotPassed(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, 1, 1, []TaskAnswer{
		{TaskID: 1, Options: []int{1}},
	})

	require.Error(t, err)
	assert.Equal(t, ErrHardBlockNotPassed, err)
}
