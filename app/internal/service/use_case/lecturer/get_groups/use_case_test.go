package getgroups

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

// =========================
// TESTS
// =========================

func Test_GetGroups_InvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 0,
		TestID: 1,
		Year:   0,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)

	_, err = env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 6,
		TestID: 0,
		Year:   0,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)

	_, err = env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 6,
		TestID: 1,
		Year:   -1,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
}

func Test_GetGroups_AccessDenied(t *testing.T) {
	env := setup(t)

	// lecturer1 (id=6) НЕ владеет test 1 (owner=5)
	_, err := env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 6,
		TestID: 1,
		Year:   0,
	})

	assert.ErrorIs(t, err, ErrAccessDenied)
}

func Test_GetGroups_Success_CurrentYear(t *testing.T) {
	env := setup(t)

	// lecturer1 владеет test 2
	resp, err := env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 6,
		TestID: 2,
		Year:   0,
	})

	require.NoError(t, err)
	require.NotNil(t, resp.Groups)

	// В фикстурах:
	// test_permissions: test 2 → A-101
	assert.Len(t, resp.Groups, 1)
	assert.Equal(t, "A-101", resp.Groups[0].GroupName)
	assert.GreaterOrEqual(t, resp.Groups[0].MembersCount, 1)
}

func Test_GetGroups_Success_WithYearOffset(t *testing.T) {
	env := setup(t)

	// test 1 принадлежит user 5 (magistr)
	resp, err := env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 5,
		TestID: 1,
		Year:   1, // предыдущий год
	})

	require.NoError(t, err)

	// В фикстурах есть old_student (A-101, старый год)
	assert.NotNil(t, resp.Groups)
	assert.GreaterOrEqual(t, len(resp.Groups), 1)
}

func Test_GetGroups_NoGroups(t *testing.T) {
	env := setup(t)

	// test 4 (lecturer2) не имеет permissions в фикстурах
	resp, err := env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 7,
		TestID: 4,
		Year:   0,
	})

	require.Error(t, err)
	assert.Empty(t, resp.Groups)
}

func Test_GetGroups_RepositoryError(t *testing.T) {
	env := setup(t)

	// несуществующий test → ошибка из repo
	_, err := env.uc.Execute(env.ctx, GetGroupsRequest{
		UserID: 6,
		TestID: 9999,
		Year:   0,
	})

	assert.Error(t, err)
}

