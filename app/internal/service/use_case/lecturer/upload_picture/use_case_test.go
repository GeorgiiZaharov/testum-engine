package uploadpicture

import (
	"context"
	"io"
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
		&realStorage{},
		zap.NewNop(),
		"http://localhost/",
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

// =========================
// STORAGE MOCK
// =========================

type realStorage struct{}

func (r *realStorage) UploadPicture(file io.Reader, fileName string, login string) (string, error) {
	// Для интеграции просто возвращаем имя файла с login префиксом
	return login + "/" + fileName, nil
}

// =========================
// HELPERS
// =========================

func fakeImage() []byte {
	return []byte("fake image content")
}

// =========================
// TESTS
// =========================

func Test_UploadPicture_InvalidInput(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadPictureRequest{})

	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.False(t, resp.Success)
}

func Test_UploadPicture_User_NotLecturer(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadPictureRequest{
		UserID:   2, // student
		File:     fakeImage(),
		FileName: "test.png",
	})

	assert.ErrorIs(t, err, ErrAccessDenied)
	assert.False(t, resp.Success)
}

func Test_UploadPicture_User_NotFound(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadPictureRequest{
		UserID:   9999,
		File:     fakeImage(),
		FileName: "test.png",
	})

	assert.Error(t, err)
	assert.False(t, resp.Success)
}

func Test_UploadPicture_StorageFails(t *testing.T) {
	env := setup(t)

	// используем мока с ошибкой
	env.uc.storage = &storageFailMock{}

	resp, err := env.uc.Execute(env.ctx, UploadPictureRequest{
		UserID:   6, // lecturer
		File:     fakeImage(),
		FileName: "fail.png",
	})

	assert.ErrorIs(t, err, ErrStorageFailed)
	assert.False(t, resp.Success)
}

func Test_UploadPicture_Success(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadPictureRequest{
		UserID:   6, // lecturer
		File:     fakeImage(),
		FileName: "picture.png",
	})

	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "http://localhost/lecturer1/picture.png", resp.URL)
}

// =========================
// STORAGE FAIL MOCK
// =========================

type storageFailMock struct{}

func (s *storageFailMock) UploadPicture(file io.Reader, fileName string, login string) (string, error) {
	return "", ErrStorageFailed
}
