package result

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"
)

type testEnv struct {
	repo *repository
	fx   *fixtures.Manager
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

	r := NewRepository(database.DB, zap.NewNop()).(*repository)

	return &testEnv{
		repo: r,
		fx:   fx,
		ctx:  ctx,
	}
}

// helper: фиксируем время чтобы академические тесты были стабильны
func (e *testEnv) freezeTime(t time.Time) {
	e.repo.now = func() time.Time {
		return t
	}
}

// =========================
// GetStudentResult
// =========================
func TestGetStudentResult_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetStudentResult(env.ctx, 3, 1)

	require.NoError(t, err)
	assert.Equal(t, 4, *res.Mark)
	assert.Equal(t, 80.0, *res.SuccessRate)
	assert.NotNil(t, res.DateStart)
	assert.NotNil(t, *res.DateEnd)
}

func TestGetStudentResult_HappyPath_WithNullDateEnd(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetStudentResult(env.ctx, 2, 1)

	require.NoError(t, err)
	assert.Nil(t, res.Mark)
	assert.Nil(t, res.SuccessRate)
	assert.NotNil(t, res.DateStart)
	assert.Nil(t, res.DateEnd)
}

func TestGetStudentResult_NotFound(t *testing.T) {
	env := setup(t)

	_, err := env.repo.GetStudentResult(env.ctx, 1, 1)

	require.Error(t, err)
	assert.Equal(t, ErrResultNotFound, err)
}

func TestGetStudentResult_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("db error"))

	_, err = r.GetStudentResult(context.Background(), 1, 1)

	assert.Error(t, err)
}

// =========================
// GetGroupResult
// =========================

func TestGetGroupResult_CurrentAcademicYear_NullsAndValues(t *testing.T) {
	env := setup(t)

	env.freezeTime(time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC))

	results, err := env.repo.GetGroupResult(env.ctx, 1, "A-101", 0)

	require.NoError(t, err)
	require.NotEmpty(t, results)

	var student1, student2 bool

	for _, r := range results {
		switch r.Login {

		case "student1":
			student1 = true
			assert.Nil(t, r.Mark)
			assert.Nil(t, r.SuccessRate)
			assert.Nil(t, r.DateStart)
			assert.Nil(t, r.DateEnd)

		case "student2":
			student2 = true
			assert.Nil(t, r.Mark)
			assert.NotNil(t, r.DateStart)
		}
	}

	assert.True(t, student1)
	assert.True(t, student2)
	assert.Equal(t, len(results), 2)
}

func TestGetGroupResult_PreviousAcademicYear_FullData(t *testing.T) {
	env := setup(t)

	env.freezeTime(time.Date(2027, 9, 10, 0, 0, 0, 0, time.UTC))

	results, err := env.repo.GetGroupResult(env.ctx, 1, "A-101", 3)

	require.NoError(t, err)
	require.NotEmpty(t, results)

	found := false

	for _, r := range results {
		if r.Login == "old_student" {
			found = true

			assert.NotNil(t, r.Mark)
			assert.NotNil(t, r.SuccessRate)
			assert.NotNil(t, r.DateStart)
			assert.NotNil(t, r.DateEnd)
		}
	}

	assert.True(t, found)
}

func TestGetGroupResult_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("db error"))

	_, err = r.GetGroupResult(context.Background(), 1, "A-101", 0)

	assert.Error(t, err)
}

func TestGetGroupResult_RowsError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	rows := sqlmock.NewRows([]string{
		"id", "name", "login", "mail",
		"mark", "success_rate", "date_start", "date_end",
	}).AddRow(
		1, "Test", "login", "mail",
		nil, nil, nil, nil,
	).RowError(0, errors.New("row error"))

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	_, err = r.GetGroupResult(context.Background(), 1, "A-101", 0)

	assert.Error(t, err)
}

// =========================
// DeleteAttempt
// =========================
func TestDeleteAttempt_HappyPath(t *testing.T) {
	env := setup(t)

	err := env.repo.DeleteAttempt(env.ctx, 1, 3)

	require.NoError(t, err)
}

func TestDeleteAttempt_NotFound(t *testing.T) {
	env := setup(t)

	err := env.repo.DeleteAttempt(env.ctx, 999, 999)

	require.NoError(t, err)
}

func TestDeleteAttempt_StudentAnswersError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	r := NewRepository(db, zap.NewNop())

	mock.ExpectExec("DELETE sa").
		WithArgs(3, 1).
		WillReturnError(errors.New("db error"))

	err = r.DeleteAttempt(context.Background(), 1, 3)

	require.Error(t, err)
	require.Equal(t, ErrDeleteFailed, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteAttempt_StudentTestsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	r := NewRepository(db, zap.NewNop())

	// first query OK
	mock.ExpectExec("DELETE sa").
		WithArgs(3, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// second query FAILS
	mock.ExpectExec("DELETE FROM student_tests").
		WithArgs(3, 1).
		WillReturnError(errors.New("db error"))

	err = r.DeleteAttempt(context.Background(), 1, 3)

	require.Error(t, err)
	require.Equal(t, ErrDeleteFailed, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
