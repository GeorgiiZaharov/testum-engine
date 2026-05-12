package posthardanswers

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

func Test_PostHardAnswers_AllCorrect_ShouldFinishTest(t *testing.T) {
	env := setup(t)

	userID := 1 // A-101 → есть доступ к test 1
	testID := 1

	answers := []TaskAnswer{
		{
			TaskID:  2,        // hard задача
			Options: []int{3}, // правильный ответ
		},
	}

	resp, err := env.uc.Execute(env.ctx, userID, testID, answers)

	assert.NoError(t, err)
	assert.True(t, resp.IsAllCorrect)
}

func Test_PostHardAnswers_NotAllCorrect_ShouldNotFinishTest(t *testing.T) {
	env := setup(t)

	userID := 1
	testID := 1

	answers := []TaskAnswer{
		{
			TaskID:  2,
			Options: []int{4}, // неправильный ответ
		},
	}

	resp, err := env.uc.Execute(env.ctx, userID, testID, answers)

	require.NoError(t, err)
	assert.False(t, resp.IsAllCorrect)
}

func Test_PostHardAnswers_AccessDenied(t *testing.T) {
	env := setup(t)

	userID := 3 // B-202
	testID := 3 // доступ только C-000

	answers := []TaskAnswer{
		{
			TaskID:  4,
			Options: []int{7},
		},
	}

	_, err := env.uc.Execute(env.ctx, userID, testID, answers)

	require.Error(t, err)
	assert.Equal(t, ErrAccessDenied, err)
}

func Test_PostHardAnswers_AlreadySubmitted(t *testing.T) {
	env := setup(t)

	userID := 3 // уже есть ответы в student_answers
	testID := 1

	answers := []TaskAnswer{
		{
			TaskID:  2,
			Options: []int{3},
		},
	}

	_, err := env.uc.Execute(env.ctx, userID, testID, answers)

	require.Error(t, err)
	assert.Equal(t, ErrAlreadySubmitted, err)
}
