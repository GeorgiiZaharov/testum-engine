package deletetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

type testEnv struct {
	uc  *UseCase
	ctx context.Context
}

func setup(t *testing.T, storage storageAdapter) *testEnv {
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
		storage,
		zap.NewNop(),
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

// -------------------- STORAGE MOCKS --------------------

type mockStorage struct{}

func (m *mockStorage) DeleteFile(fileName string) error {
	return nil
}

type failingStorage struct{}

func (f *failingStorage) DeleteFile(fileName string) error {
	return ErrStorageFailed
}

// -------------------- TESTS --------------------

func TestDeleteTest_Success(t *testing.T) {
	env := setup(t, &mockStorage{})

	req := DeleteTestRequest{
		UserID: 5,
		TestID: 1,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	require.True(t, resp.Success)
}

func TestDeleteTest_AccessDenied(t *testing.T) {
	env := setup(t, &mockStorage{})

	req := DeleteTestRequest{
		UserID: 7, // lecturer2
		TestID: 1, // belongs to another lecturer
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.Equal(t, ErrAccessDenied, err)
	require.False(t, resp.Success)
}

func TestDeleteTest_InvalidInput(t *testing.T) {
	env := setup(t, &mockStorage{})

	req := DeleteTestRequest{
		UserID: 0,
		TestID: 0,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.Equal(t, ErrInvalidInput, err)
	require.False(t, resp.Success)
}

func TestDeleteTest_TestNotFound(t *testing.T) {
	env := setup(t, &mockStorage{})

	req := DeleteTestRequest{
		UserID: 5,
		TestID: 9999,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.False(t, resp.Success)
}

func TestDeleteTest_StorageFailure(t *testing.T) {
	env := setup(t, &failingStorage{})

	req := DeleteTestRequest{
		UserID: 5,
		TestID: 1,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.Error(t, err)
	require.Equal(t, ErrStorageFailed, err)
	require.False(t, resp.Success)
}
