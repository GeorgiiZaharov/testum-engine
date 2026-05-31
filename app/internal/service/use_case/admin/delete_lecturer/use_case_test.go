package deletelecturer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	userrepo "testum-engine/app/internal/repository/user"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

// =========================
// MOCK STORAGE
// =========================

type mockStorage struct {
	DeleteFileFunc     func(fileName string) error
	DeletePicturesFunc func(login string) error
}

func (m *mockStorage) DeleteFile(fileName string) error {
	return m.DeleteFileFunc(fileName)
}

func (m *mockStorage) DeletePictures(login string) error {
	return m.DeletePicturesFunc(login)
}

// =========================
// TEST ENV
// =========================

type testEnv struct {
	uc   *UseCase
	ctx  context.Context
	repo userrepo.Repository
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

	storage := &mockStorage{
		DeleteFileFunc:     func(string) error { return nil },
		DeletePicturesFunc: func(string) error { return nil },
	}

	uc := NewUseCase(
		NewFactory(database, zap.NewNop()),
		storage,
		zap.NewNop(),
	)

	return &testEnv{
		uc:   uc,
		ctx:  ctx,
		repo: userrepo.NewRepository(database, zap.NewNop()),
	}
}

// =========================
// TESTS
// =========================

func TestDeleteLecturer_Success(t *testing.T) {
	env := setup(t)

	req := DeleteLecturerRequest{
		AdminID:    8,
		LecturerID: 6,
	}

	resp, err := env.uc.Execute(env.ctx, req)

	require.NoError(t, err)
	assert.True(t, resp.Success)

	user, err := env.repo.GetByID(env.ctx, 6)
	require.NoError(t, err)
	assert.False(t, user.IsLecturer)
}

func TestDeleteLecturer_InvalidInput(t *testing.T) {
	env := setup(t)

	_, err := env.uc.Execute(env.ctx, DeleteLecturerRequest{
		AdminID:    0,
		LecturerID: 0,
	})

	assert.ErrorIs(t, err, ErrInvalidInput)
}

func TestDeleteLecturer_Forbidden(t *testing.T) {
	env := setup(t)

	req := DeleteLecturerRequest{
		AdminID:    1,
		LecturerID: 6,
	}

	_, err := env.uc.Execute(env.ctx, req)

	assert.ErrorIs(t, err, ErrForbidden)
}

func TestDeleteLecturer_NotFound(t *testing.T) {
	env := setup(t)

	req := DeleteLecturerRequest{
		AdminID:    8,
		LecturerID: 999,
	}

	_, err := env.uc.Execute(env.ctx, req)

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDeleteLecturer_NotLecturer(t *testing.T) {
	env := setup(t)

	req := DeleteLecturerRequest{
		AdminID:    8,
		LecturerID: 1,
	}

	_, err := env.uc.Execute(env.ctx, req)

	assert.ErrorIs(t, err, ErrNotLecturer)
}

func TestDeleteLecturer_DeleteFileError(t *testing.T) {
	env := setup(t)

	storage := &mockStorage{
		DeleteFileFunc: func(string) error {
			return errors.New("delete file error")
		},
		DeletePicturesFunc: func(string) error { return nil },
	}

	env.uc = NewUseCase(env.uc.factory, storage, zap.NewNop())

	req := DeleteLecturerRequest{
		AdminID:    8,
		LecturerID: 6,
	}

	_, err := env.uc.Execute(env.ctx, req)

	assert.Error(t, err)
}

func TestDeleteLecturer_DeletePicturesError(t *testing.T) {
	env := setup(t)

	storage := &mockStorage{
		DeleteFileFunc:     func(string) error { return nil },
		DeletePicturesFunc: func(string) error { return errors.New("delete pictures error") },
	}

	env.uc = NewUseCase(env.uc.factory, storage, zap.NewNop())

	req := DeleteLecturerRequest{
		AdminID:    8,
		LecturerID: 6,
	}

	_, err := env.uc.Execute(env.ctx, req)

	assert.Error(t, err)
}
