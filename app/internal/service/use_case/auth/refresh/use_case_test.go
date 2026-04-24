package refresh

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	db "testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

/*
========================
TEST ENV
========================
*/

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
		&mockAuthService{},
		zap.NewNop(),
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

/*
========================
MOCK AUTH SERVICE
========================
*/

type mockAuthService struct{}

func (m *mockAuthService) GenerateAccess(userID int) (string, error) {
	if userID == 2 {
		return "", assert.AnError
	}
	return "access_token", nil
}

func (m *mockAuthService) GenerateRefresh(userID int) (string, error) {
	if userID == 3 {
		return "", assert.AnError
	}
	return "refresh_token", nil
}

/*
========================
TESTS
========================
*/

func TestUseCase_Execute_Success(t *testing.T) {
	env := setup(t)

	req := AuthRefreshRequest{
		UserID: 1,
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	assert.Equal(t, "access_token", res.AccessToken)
	assert.Equal(t, "refresh_token", res.RefreshToken)
}

func TestUseCase_Execute_InvalidUserID(t *testing.T) {
	env := setup(t)

	req := AuthRefreshRequest{
		UserID: 0,
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrInvalidUserID, err)
	assert.Empty(t, res.AccessToken)
	assert.Empty(t, res.RefreshToken)
}

func TestUseCase_Execute_UserNotFound(t *testing.T) {
	env := setup(t)

	req := AuthRefreshRequest{
		UserID: 99999,
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Empty(t, res.AccessToken)
	assert.Empty(t, res.RefreshToken)
}

func TestUseCase_Execute_AccessTokenError(t *testing.T) {
	env := setup(t)

	req := AuthRefreshRequest{
		UserID: 2,
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrAuthFailed, err)
	assert.Empty(t, res.AccessToken)
	assert.Empty(t, res.RefreshToken)
}

func TestUseCase_Execute_RefreshTokenError(t *testing.T) {
	env := setup(t)

	req := AuthRefreshRequest{
		UserID: 3,
	}

	res, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	assert.Equal(t, ErrAuthFailed, err)
	assert.Empty(t, res.AccessToken)
	assert.Empty(t, res.RefreshToken)
}
