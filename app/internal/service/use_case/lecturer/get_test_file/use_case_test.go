package gettestfile

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

type mockStorage struct{}

func (m *mockStorage) GetFile(fileName string) (*os.File, error) {
	// эмулируем существующий файл
	if fileName == "math_basics.json" ||
		fileName == "linear_algebra.json" ||
		fileName == "physics_intro.json" {
		f, _ := os.CreateTemp("", "testfile")
		return f, nil
	}

	return nil, nil
}

type mockStorageError struct{}

func (m *mockStorageError) GetFile(fileName string) (*os.File, error) {
	return nil, assert.AnError
}

type testEnv struct {
	uc  *UseCase
	ctx context.Context
}

func setup(t *testing.T, storage StorageAdapter) *testEnv {
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
		storage,
		zap.NewNop(),
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

func TestGetTestFile_Success(t *testing.T) {
	env := setup(t, &mockStorage{})

	req := GetTestFileRequest{
		UserID: 6, // lecturer1
		TestID: 2, // Linear Algebra (owner = 6)
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp.File)
}

func TestGetTestFile_NoAccess(t *testing.T) {
	env := setup(t, &mockStorage{})

	req := GetTestFileRequest{
		UserID: 7, // другой лектор
		TestID: 2, // тест lecturer1
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Equal(t, GetTestFileResponse{}, resp)
}

func TestGetTestFile_FileNameEmpty(t *testing.T) {
	env := setup(t, &mockStorage{})

	// в фикстурах нет теста без file_name → используем несуществующий ID
	req := GetTestFileRequest{
		UserID: 6,
		TestID: 999,
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	assert.Equal(t, GetTestFileResponse{}, resp)
}

func TestGetTestFile_StorageFail(t *testing.T) {
	env := setup(t, &mockStorageError{})

	req := GetTestFileRequest{
		UserID: 6,
		TestID: 2,
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	assert.Equal(t, ErrFileNotFound, err)
	assert.Equal(t, GetTestFileResponse{}, resp)
}

func TestGetTestFile_InvalidInput(t *testing.T) {
	env := setup(t, &mockStorage{})

	req := GetTestFileRequest{
		UserID: 0,
		TestID: 0,
	}

	resp, err := env.uc.Execute(env.ctx, req)
	require.Error(t, err)
	assert.Equal(t, ErrForbidden, err)
	assert.Equal(t, GetTestFileResponse{}, resp)
}
