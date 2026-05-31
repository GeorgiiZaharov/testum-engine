package access

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

type testEnv struct {
	repo Repository
	ctx  context.Context
}

// =========================
// SETUP (REAL DB)
// =========================
func setup(t *testing.T) *testEnv {
	t.Helper()

	ctx := context.Background()

	database, cleanup, err := bootstrap.Setup(bootstrap.Config{
		DBOptions: db.DBOptions{
			Path: ":memory:",
		},
		Migrations: "../../../migrations",
	})
	require.NoError(t, err)

	t.Cleanup(cleanup)

	_, err = database.Exec(`PRAGMA foreign_keys = ON`)
	require.NoError(t, err)

	fx := fixtures.New(database)

	require.NoError(t, fx.Reset(ctx))
	require.NoError(t, fx.SeedAll(ctx))

	return &testEnv{
		repo: NewRepository(database.DB, zap.NewNop()),
		ctx:  ctx,
	}
}

//
// =========================
// HasLecturerAccess
// =========================
//

func TestHasLecturerAccess_HappyPath(t *testing.T) {
	env := setup(t)

	// owner_id = 5 (magistr owns test 1)
	ok, err := env.repo.HasLecturerAccess(env.ctx, 5, 1)

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestHasLecturerAccess_NotFound(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.HasLecturerAccess(env.ctx, 999, 999)

	require.NoError(t, err)
	assert.False(t, ok)
}

func TestHasLecturerAccess_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("SELECT EXISTS").
		WillReturnError(errors.New("db error"))

	ok, err := r.HasLecturerAccess(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.False(t, ok)
}

//
// =========================
// HasStudentAccess
// =========================
//

func TestHasStudentAccess_HappyPath(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.HasStudentAccess(env.ctx, 2, 1)

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestHasStudentAccess_NotFound(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.HasStudentAccess(env.ctx, 999, 999)

	require.NoError(t, err)
	assert.False(t, ok)
}

func TestHasStudentAccess_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("SELECT EXISTS").
		WillReturnError(errors.New("db error"))

	ok, err := r.HasStudentAccess(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.False(t, ok)
}

//
// =========================
// GiveAccess
// =========================
//

func TestGiveAccess_HappyPath(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.GiveAccess(env.ctx, 3, "A-101")

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestGiveAccess_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectExec("INSERT INTO test_permissions").
		WillReturnError(errors.New("db error"))

	ok, err := r.GiveAccess(context.Background(), 1, "A-101")

	assert.Error(t, err)
	assert.False(t, ok)
}

//
// =========================
// TakeAccess
// =========================
//

func TestTakeAccess_HappyPath(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.TakeAccess(env.ctx, 2, "A-101")

	require.NoError(t, err)
	assert.True(t, ok)
}

func TestTakeAccess_NotFound(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.TakeAccess(env.ctx, 1, "GGGGG")

	assert.Error(t, err)
	assert.Equal(t, ErrAccessNotFound, err)
	assert.False(t, ok)
}

func TestTakeAccess_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectExec("DELETE FROM test_permissions").
		WillReturnError(errors.New("db error"))

	ok, err := r.TakeAccess(context.Background(), 1, "A-101")

	assert.Error(t, err)
	assert.False(t, ok)
}
