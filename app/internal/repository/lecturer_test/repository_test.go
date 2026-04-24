package lecturertest

import (
	"context"
	"database/sql"
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
// CREATE
// ======================================================

func TestCreate_HappyPath(t *testing.T) {
	env := setup(t)

	test := Test{
		Name:     "New Test",
		FileName: "new.json",
		Groups:   []string{"A-101"},
		Tasks: []Task{
			{
				Text:    "Q1",
				IsHard:  false,
				Answers: []Answer{{Text: "A1", IsCorrect: true}},
			},
		},
	}

	tid, err := env.repo.Create(env.ctx, 5, test)

	assert.NoError(t, err)
	assert.True(t, tid > 0)
}

func TestCreate_InsertTestError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("insert into tests").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	tid, err := r.Create(context.Background(), 1, Test{})

	assert.ErrorIs(t, err, ErrCreateFailed)
	assert.True(t, tid == 0)
}

func TestCreate_TestLastInsertIDError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("insert into tests").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("id error")))

	r := NewRepository(dbMock, zap.NewNop())

	tid, err := r.Create(context.Background(), 1, Test{})

	assert.Error(t, err)
	assert.True(t, tid == 0)
}

func TestCreate_InsertTaskError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("insert into tests").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("insert into tasks").
		WillReturnError(errors.New("task error"))

	r := NewRepository(dbMock, zap.NewNop())

	tid, err := r.Create(context.Background(), 1, Test{
		Tasks: []Task{{}},
	})

	assert.ErrorIs(t, err, ErrCreateFailed)
	assert.True(t, tid == 0)
}

func TestCreate_TaskLastInsertIDError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("insert into tests").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("insert into tasks").
		WillReturnResult(sqlmock.NewErrorResult(errors.New("id error")))

	r := NewRepository(dbMock, zap.NewNop())

	tid, err := r.Create(context.Background(), 1, Test{
		Tasks: []Task{{}},
	})

	assert.Error(t, err)
	assert.True(t, tid == 0)
}

func TestCreate_InsertAnswerError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("insert into tests").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("insert into tasks").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("insert into answers").
		WillReturnError(errors.New("answer error"))

	r := NewRepository(dbMock, zap.NewNop())

	tid, err := r.Create(context.Background(), 1, Test{
		Tasks: []Task{
			{
				Answers: []Answer{{}},
			},
		},
	})

	assert.ErrorIs(t, err, ErrCreateFailed)
	assert.True(t, tid == 0)
}

func TestCreate_InsertPermissionError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("insert into tests").
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := NewRepository(dbMock, zap.NewNop())

	tid, err := r.Create(context.Background(), 1, Test{
		Groups: []string{"A"},
	})

	assert.ErrorIs(t, err, ErrCreateFailed)
	assert.True(t, tid == 0)
}

// ======================================================
// DELETE
// ======================================================

func TestDelete_HappyPath(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.Delete(env.ctx, 1)

	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestDelete_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("delete from tests").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	ok, err := r.Delete(context.Background(), 1)

	assert.ErrorIs(t, err, ErrDeleteFailed)
	assert.False(t, ok)
}

func TestDelete_NotFound(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("delete from tests").
		WillReturnResult(sqlmock.NewResult(0, 0))

	r := NewRepository(dbMock, zap.NewNop())

	ok, err := r.Delete(context.Background(), 1)

	assert.ErrorIs(t, err, ErrTestNotFound)
	assert.False(t, ok)
}

// ======================================================
// GET BY ID
// ======================================================

func TestGetByID_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetByID(env.ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.ID)
	assert.Equal(t, "Math Basics", res.Name)
	assert.GreaterOrEqual(t, res.CntQuestions, 1)
	assert.GreaterOrEqual(t, res.CntHardQuestions, 0)
	assert.Equal(t, "math_basics.json", res.FileName)
	assert.NotEmpty(t, res.Groups)
}

func TestGetByID_NotFound(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectQuery("select").
		WillReturnError(sql.ErrNoRows)

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetByID(context.Background(), 1)

	assert.ErrorIs(t, err, ErrTestNotFound)
}

func TestGetByID_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectQuery("select").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetByID(context.Background(), 1)

	assert.Error(t, err)
}

func TestGetByID_RowsError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	now := time.Now()

	mock.ExpectQuery("select").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "cnt", "hard", "date",
		}).AddRow(1, "test", 1, 0, now))

	rows := sqlmock.NewRows([]string{"group"}).
		AddRow("A").
		CloseError(errors.New("rows error"))

	mock.ExpectQuery("select").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetByID(context.Background(), 1)

	assert.Error(t, err)
}

// ======================================================
// GET BY LECTURER
// ======================================================

func TestGetByLecturer_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetByLecturer(env.ctx, 6)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 1)

	assert.NotEmpty(t, res[0].Name)
}

func TestGetByLecturer_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectQuery("select").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetByLecturer(context.Background(), 1)

	assert.Error(t, err)
}

func TestGetByLecturer_RowsError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "name", "cnt", "hard", "date",
	}).AddRow(1, "test", 1, 0, now)

	rows.CloseError(errors.New("rows error"))

	mock.ExpectQuery("select").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetByLecturer(context.Background(), 1)

	assert.Error(t, err)
}

// ======================================================
// GET GROUPS
// ======================================================

func TestGetGroups_CurrentYear_HappyPath(t *testing.T) {
	env := setup(t)

	fixedNow := time.Date(2026, time.May, 10, 0, 0, 0, 0, time.UTC)

	repo := env.repo.(*repository)
	repo.now = func() time.Time { return fixedNow }

	res, err := repo.GetGroups(env.ctx, 1, 0)

	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.True(t, len(res) == 2)

	foundA := false
	foundB := false
	for _, g := range res {
		if g.GroupName == "A-101" {
			foundA = true
			assert.GreaterOrEqual(t, g.MembersCount, 1)
		}
		if g.GroupName == "B-202" {
			foundB = true
			assert.GreaterOrEqual(t, g.MembersCount, 1)
		}
	}
	assert.True(t, foundA)
	assert.True(t, foundB)
}

func TestGetGroups_PastYear_HappyPath(t *testing.T) {
	env := setup(t)

	fixedNow := time.Date(2026, time.September, 10, 0, 0, 0, 0, time.UTC)

	repo := env.repo.(*repository)
	repo.now = func() time.Time { return fixedNow }

	res, err := repo.GetGroups(env.ctx, 1, 1)

	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.True(t, len(res) == 1)
}

func TestGetGroups_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop()).(*repository)
	r.now = func() time.Time { return time.Now() }

	_, err := r.GetGroups(context.Background(), 1, 0)

	assert.Error(t, err)
}

func TestGetGroups_ScanError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	rows := sqlmock.NewRows([]string{"group", "count"}).
		AddRow("A-101", "invalid_int") // ломаем scan

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop()).(*repository)
	r.now = func() time.Time { return time.Now() }

	_, err := r.GetGroups(context.Background(), 1, 0)

	assert.Error(t, err)
}

func TestGetGroups_RowsError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	rows := sqlmock.NewRows([]string{"group", "count"}).
		AddRow("A-101", 1).
		CloseError(errors.New("rows error"))

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop()).(*repository)
	r.now = func() time.Time { return time.Now() }

	_, err := r.GetGroups(context.Background(), 1, 0)

	assert.Error(t, err)
}

func TestGetGroups_NotFound(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	rows := sqlmock.NewRows([]string{"group", "count"})

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop()).(*repository)
	r.now = func() time.Time { return time.Now() }

	_, err := r.GetGroups(context.Background(), 1, 0)

	assert.ErrorIs(t, err, ErrTestNotFound)
}
