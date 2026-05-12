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
// SETUP (REAL DB + FIXTURES)
// =========================

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
		Migrations: "../../../migrations",
	})
	require.NoError(t, err)

	t.Cleanup(cleanup)

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

	files, err := env.repo.GetAllTestFiles(context.Background(), 6) // lecturer1 has tests

	require.NoError(t, err)
	require.NotEmpty(t, files)

	assert.Contains(t, files, "linear_algebra.json")
	assert.Contains(t, files, "physics_intro.json")
}

// =========================
// NOT FOUND (EMPTY RESULT)
// =========================

func TestGetAllTestFiles_NotFound(t *testing.T) {
	env := setup(t)

	files, err := env.repo.GetAllTestFiles(context.Background(), 999)

	assert.Empty(t, files)
	assert.NoError(t, err)
}

// =========================
// QUERY ERROR (SQL MOCK)
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

// =========================
// ROWS ERROR
// =========================

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
