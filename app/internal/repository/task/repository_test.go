package task

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
// SETUP REAL DB (FIXTURES)
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

func TestGetBaseTasks_HappyPath(t *testing.T) {
	env := setup(t)

	tasks, err := env.repo.GetBaseTasks(env.ctx, 1)

	require.NoError(t, err)
	require.NotEmpty(t, tasks)

	// TEST 1 base task (is_hard = false)
	assert.Equal(t, "Solve equation: 2 + 2 = ?", tasks[0].Text)
	assert.False(t, tasks[0].IsHard)

	require.Len(t, tasks[0].Answers, 2)
}

func TestGetHardTasks_HappyPath(t *testing.T) {
	env := setup(t)

	tasks, err := env.repo.GetHardTasks(env.ctx, 1)

	require.NoError(t, err)
	require.NotEmpty(t, tasks)

	assert.Equal(t, "Prove derivative of x^2", tasks[0].Text)
	assert.True(t, tasks[0].IsHard)

	require.Len(t, tasks[0].Answers, 2)
}

// =========================
// NOT FOUND (EMPTY RESULT)
// =========================

func TestGetBaseTasks_NotFound(t *testing.T) {
	env := setup(t)

	// test_id=999 нет в фикстурах
	_, err := env.repo.GetBaseTasks(env.ctx, 999)

	require.Error(t, err)
	assert.Equal(t, ErrTaskNotFound, err)
}

// =========================
// QUERY ERROR (SQL MOCK)
// =========================

func TestGetBaseTasks_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("SELECT").
		WithArgs(1, false).
		WillReturnError(errors.New("db error"))

	_, err = r.GetBaseTasks(context.Background(), 1)

	require.Error(t, err)
}

// =========================
// ROWS ERROR (scan error path)
// =========================

func TestScanTasks_RowsError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	rows := sqlmock.NewRows([]string{
		"id", "text", "image_url", "is_hard",
		"answer_text", "answer_image", "is_correct",
	}).AddRow(
		1, "Task", nil, false,
		"Answer", nil, true,
	).RowError(0, errors.New("row error"))

	mock.ExpectQuery("SELECT").
		WithArgs(1, false).
		WillReturnRows(rows)

	_, err = r.GetBaseTasks(context.Background(), 1)

	require.Error(t, err)
}
