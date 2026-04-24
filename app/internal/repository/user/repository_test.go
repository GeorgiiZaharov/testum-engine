package user

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

func strPtr(s string) *string {
	return &s
}

// ======================================================
// UPSERT
// ======================================================

func TestUpsert_Insert_HappyPath(t *testing.T) {
	env := setup(t)

	id, err := env.repo.Upsert(env.ctx, CreateUserParams{
		Login: "new_user",
		Mail:  "new@mail.com",
		Name:  "New User",
		Group: strPtr("A-101"),
	})

	assert.NoError(t, err)
	assert.NotZero(t, id)
}

func TestUpsert_Update_HappyPath(t *testing.T) {
	env := setup(t)

	id, err := env.repo.Upsert(env.ctx, CreateUserParams{
		Login: "student1",
		Mail:  "updated@mail.com",
		Name:  "Updated Name",
		Group: strPtr("A-101"),
	})

	assert.NoError(t, err)
	assert.NotZero(t, id)
}

func TestUpsert_Update_NoChanges(t *testing.T) {
	env := setup(t)

	id, err := env.repo.Upsert(env.ctx, CreateUserParams{
		Login: "student1",
		Mail:  "student1@mail.com",
		Name:  "Student One",
		Group: strPtr("A-101"),
	})

	assert.NoError(t, err)
	assert.Equal(t, id, 1)
}

func TestUpsert_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("INSERT INTO users").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.Upsert(context.Background(), CreateUserParams{})

	assert.Error(t, err)
}

func TestUpsert_LastInsertFallback(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(0, 1)) // id = 0 → fallback

	mock.ExpectQuery("SELECT id FROM users").
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	r := NewRepository(dbMock, zap.NewNop())

	id, err := r.Upsert(context.Background(), CreateUserParams{
		Login: "test",
	})

	assert.NoError(t, err)
	assert.Equal(t, 10, id)
}

func TestUpsert_FallbackError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("INSERT INTO users").
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery("SELECT id FROM users").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.Upsert(context.Background(), CreateUserParams{
		Login: "test",
	})

	assert.Error(t, err)
}

// ======================================================
// GET BY ID
// ======================================================

func TestGetByID_HappyPath(t *testing.T) {
	env := setup(t)

	u, err := env.repo.GetByID(env.ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "student1", u.Login)
}

func TestGetByID_NotFound(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectQuery("SELECT").
		WillReturnError(sql.ErrNoRows)

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetByID(context.Background(), 1)

	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestGetByID_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetByID(context.Background(), 1)

	assert.Error(t, err)
}

// ======================================================
// GET LECTURERS
// ======================================================

func TestGetLecturers_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetLecturers(env.ctx)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(res), 1)
}

func TestGetLecturers_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectQuery("SELECT").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetLecturers(context.Background())

	assert.Error(t, err)
}

func TestGetLecturers_ScanError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	rows := sqlmock.NewRows([]string{
		"id", "login", "mail", "name", "group", "is_lecturer", "created", "modified",
	}).AddRow("bad", "login", "mail", "name", "A", true, time.Now(), time.Now())

	mock.ExpectQuery("SELECT").
		WillReturnRows(rows)

	r := NewRepository(dbMock, zap.NewNop())

	_, err := r.GetLecturers(context.Background())

	assert.Error(t, err)
}

// ======================================================
// CREATE LECTURER
// ======================================================

func TestCreateLecturer_HappyPath(t *testing.T) {
	env := setup(t)

	err := env.repo.CreateLecturer(env.ctx, 1)

	assert.NoError(t, err)
}

func TestCreateLecturer_NotFound(t *testing.T) {
	env := setup(t)

	err := env.repo.CreateLecturer(env.ctx, 999)

	assert.Error(t, err)
}
func TestCreateLecturer_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("UPDATE users").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	err := r.CreateLecturer(context.Background(), 999)

	assert.Error(t, err)
}

// ======================================================
// DELETE LECTURER
// ======================================================

func TestDeleteLecturer_HappyPath(t *testing.T) {
	env := setup(t)

	err := env.repo.DeleteLecturer(env.ctx, 5)

	assert.NoError(t, err)
}

func TestDeleteLecturer_NotFound(t *testing.T) {
	env := setup(t)

	err := env.repo.DeleteLecturer(context.Background(), 999)

	assert.Error(t, err)
}

func TestDeleteLecturer_QueryError(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	defer func() { _ = dbMock.Close() }()

	mock.ExpectExec("UPDATE users").
		WillReturnError(errors.New("db error"))

	r := NewRepository(dbMock, zap.NewNop())

	err := r.DeleteLecturer(context.Background(), 999)

	assert.Error(t, err)
}
