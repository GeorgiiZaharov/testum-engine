package file

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
	repo FileRepository
	ctx  context.Context
}

// =========================
// SETUP (REAL SQLITE + FIXTURES)
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

	repo := NewRepository(database.DB, zap.NewNop())

	return &testEnv{
		repo: repo,
		ctx:  ctx,
	}
}

// =========================
// HAPPY PATH (FIXTURES)
// =========================

func TestGetAllTestFiles_HappyPath(t *testing.T) {
	env := setup(t)

	files, err := env.repo.GetAllTestFiles(context.Background(), 6)

	require.NoError(t, err)
	require.NotEmpty(t, files)

	assert.Contains(t, files, "linear_algebra.json")
	assert.Contains(t, files, "physics_intro.json")
}

// =========================
// NOT FOUND
// =========================

func TestGetAllTestFiles_NotFound(t *testing.T) {
	env := setup(t)

	files, err := env.repo.GetAllTestFiles(context.Background(), 999)

	assert.NoError(t, err)
	assert.Empty(t, files)
}

// =========================
// SQL MOCK TESTS
// =========================

func TestGetAllTestFiles_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("SELECT file_name").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	_, err = r.GetAllTestFiles(context.Background(), 1)

	require.Error(t, err)
}

func TestGetAllTestFiles_RowsError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	rows := sqlmock.NewRows([]string{"file_name"}).
		AddRow("test.json").
		RowError(0, errors.New("row error"))

	mock.ExpectQuery("SELECT file_name").
		WithArgs(1).
		WillReturnRows(rows)

	_, err = r.GetAllTestFiles(context.Background(), 1)

	require.Error(t, err)
}
