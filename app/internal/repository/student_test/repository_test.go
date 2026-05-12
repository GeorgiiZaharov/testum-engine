package studenttest

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

// =========================
// REAL DB SETUP
// =========================
type testEnv struct {
	repo Repository
	ctx  context.Context
}

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

	return &testEnv{
		repo: NewRepository(database.DB, zap.NewNop()),
		ctx:  ctx,
	}
}

// ======================================================
// START TEST
// ======================================================

func TestStartTest_HappyPath(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.StartTest(env.ctx, StartTestParams{
		UserID: 1,
		TestID: 1,
	})

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestStartTest_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectExec("insert into student_tests").
		WillReturnError(errors.New("db error"))

	ok, err := r.StartTest(context.Background(), StartTestParams{})

	assert.Error(t, err)
	assert.False(t, ok)
}

// ======================================================
// FINISH TEST
// ======================================================

func TestFinishTest_HappyPath(t *testing.T) {
	env := setup(t)

	mark := 5
	rate := 100.0

	ok, err := env.repo.FinishTest(env.ctx, FinishTestParams{
		UserID: 3,
		TestID: 1,
		Result: TestResult{
			Mark:        &mark,
			SuccessRate: &rate,
		},
	})

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestFinishTest_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectExec("update student_tests").
		WillReturnError(errors.New("db error"))

	ok, err := r.FinishTest(context.Background(), FinishTestParams{})

	assert.Error(t, err)
	assert.False(t, ok)
}

func TestFinishTest_NotFound(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectExec("update student_tests").
		WillReturnResult(sqlmock.NewResult(0, 0))

	ok, err := r.FinishTest(context.Background(), FinishTestParams{})

	assert.ErrorIs(t, err, ErrStudentTestNotFound)
	assert.False(t, ok)
}

// ======================================================
// GET ACTIVE TESTS
// ======================================================
func TestGetActiveTests_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetActiveTests(env.ctx, 2)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.GreaterOrEqual(t, len(res), 1)

	found := false
	for _, tinfo := range res {
		if tinfo.ID == 1 {
			found = true

			assert.Equal(t, "Math Basics", tinfo.Name)
			assert.Equal(t, "Magistr Magistr", tinfo.LecturerName)
			assert.GreaterOrEqual(t, tinfo.CntQuestions, 1)
			assert.GreaterOrEqual(t, tinfo.CntHardQuestions, 0)
			assert.NotNil(t, tinfo.DateStart)
		}
	}

	assert.True(t, found, "test 1 should be available for student2 via A-101 group")
}

func TestGetActiveTests_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("select").
		WillReturnError(errors.New("db error"))

	_, err = r.GetActiveTests(context.Background(), 1)

	assert.Error(t, err)
}

func TestGetActiveTests_RowsError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "lecturer_name", "cnt", "hard", "date_start",
	}).AddRow(1, "test", "lecturer", 10, 2, &now)

	rows.CloseError(errors.New("rows error"))

	mock.ExpectQuery("select").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop())

	_, err = r.GetActiveTests(context.Background(), 1)

	assert.Error(t, err)
}

// ======================================================
// GET FINISHED TESTS
// ======================================================
func TestGetFinishedTests_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetFinishedTests(env.ctx, 3)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.GreaterOrEqual(t, len(res), 1)

	found := false

	for _, r := range res {
		if r.ID == 1 {
			found = true

			assert.Equal(t, "Math Basics", r.Name)
			assert.Equal(t, "Magistr Magistr", r.LecturerName)

			assert.Equal(t, 4, r.Mark)
			assert.Equal(t, 80.0, r.SuccessRate)

			assert.False(t, r.DateStart.IsZero())
			assert.False(t, r.DateEnd.IsZero())
		}
	}

	assert.True(t, found, "finished test for student3 should exist")
}

func TestGetFinishedTests_VoidResult(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetFinishedTests(env.ctx, 1)

	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestGetFinishedTests_QueryError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	r := NewRepository(dbMock, zap.NewNop())

	mock.ExpectQuery("select").
		WillReturnError(errors.New("db error"))

	_, err = r.GetFinishedTests(context.Background(), 1)

	assert.Error(t, err)
}

func TestGetFinishedTests_RowsError(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = dbMock.Close() }()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "lecturer_name", "cnt", "hard", "mark", "rate", "start", "end",
	}).AddRow(1, "test", "lecturer", 10, 2, 5, 100.0, now, now)

	rows.CloseError(errors.New("rows error"))

	mock.ExpectQuery("select").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop())

	_, err = r.GetFinishedTests(context.Background(), 1)

	assert.Error(t, err)
}
