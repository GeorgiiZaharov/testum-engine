package uploadtest

import (
	"context"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"

	latexvalidator "testum-engine/app/internal/service/core/latexvalidator"
	validation "testum-engine/app/internal/service/core/validation"
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
		&fakeStorage{},
		validation.NewParser(zap.NewNop()),
		latexvalidator.New(),
		zap.NewNop(),
	)

	return &testEnv{
		uc:  uc,
		ctx: ctx,
	}
}

// =========================
// STORAGE MOCK
// =========================

type fakeStorage struct{}

func (s *fakeStorage) UploadFile(file io.Reader, fileName string) (string, error) {
	salt := "abc123"

	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	return name + "_" + salt + ext, nil
}

// =========================
// HELPERS
// =========================

func validTestFile() []byte {
	return []byte(`
МАТЕМАТИКА

1
# 2 + 2 = ?

+ 4
- 5

# 3 + 3 = ?

+ 6
- 7
`)
}

// ❌ неверный формат (нет названия + сломанная структура)
func invalidFormatFile() []byte {
	return []byte(`
# broken file

+ no question here
`)
}

// ❌ latex ошибка (невалидный command)
func latexErrorFile() []byte {
	return []byte(`
TEST

1
# \fracc{2}{2} = ?

+ 1
- 2
`)
}

// =========================
// TESTS
// =========================

func Test_UploadTest_InvalidInput(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadTestRequest{})

	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.False(t, resp.Success)
}

func Test_UploadTest_User_NotLecturer(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadTestRequest{
		UserID:   2,
		File:     validTestFile(),
		FileName: "test.txt",
	})

	assert.ErrorIs(t, err, ErrAccessDenied)
	assert.False(t, resp.Success)
}

func Test_UploadTest_User_NotFound(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadTestRequest{
		UserID:   9999,
		File:     validTestFile(),
		FileName: "test.txt",
	})

	assert.Error(t, err)
	assert.False(t, resp.Success)
}

func Test_UploadTest_FormatValidation_Error(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadTestRequest{
		UserID:   6,
		File:     invalidFormatFile(),
		FileName: "broken.txt",
	})

	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.NotEmpty(t, resp.FormatErrors)
	assert.Nil(t, resp.TestID)
}

func Test_UploadTest_LatexValidation_Error(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadTestRequest{
		UserID:   6,
		File:     latexErrorFile(),
		FileName: "latex.txt",
	})

	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.NotEmpty(t, resp.ValidationErrors)
	assert.Nil(t, resp.TestID)
}

func Test_UploadTest_IgnoreValidation_StillSaves(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadTestRequest{
		UserID:           6,
		File:             latexErrorFile(),
		FileName:         "latex.txt",
		IgnoreValidation: true,
	})

	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.NotNil(t, resp.TestID)
	assert.Empty(t, resp.ValidationErrors)
}

func Test_UploadTest_Success(t *testing.T) {
	env := setup(t)

	resp, err := env.uc.Execute(env.ctx, UploadTestRequest{
		UserID:   6,
		File:     validTestFile(),
		FileName: "math.txt",
	})

	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.NotNil(t, resp.TestID)
}
