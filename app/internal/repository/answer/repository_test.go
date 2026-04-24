package answer_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/testing/bootstrap"
	"testum-engine/app/internal/testing/fixtures"

	answer "testum-engine/app/internal/repository/answer"
)

// =========================
// TEST ENV
// =========================

type testEnv struct {
	repo answer.Repository
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
		repo: answer.NewRepository(database.DB, zap.NewNop()),
		ctx:  ctx,
	}
}

// =========================
// REAL DB TESTS
// =========================

func Test_GetBaseAnswers_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetBaseAnswers(env.ctx, 3, 1)

	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_GetHardAnswers_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetHardAnswers(env.ctx, 3, 1)

	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_GetBaseAnswersByTest_HappyPath(t *testing.T) {
	env := setup(t)

	res, err := env.repo.GetBaseAnswersByTest(env.ctx, 1)

	require.NoError(t, err)
	assert.NotNil(t, res)
}

func Test_DeleteAttempt_HappyPath(t *testing.T) {
	env := setup(t)

	ok, err := env.repo.DeleteAttempt(env.ctx, 3, 1)

	require.NoError(t, err)
	assert.True(t, ok)
}

// =========================
// SQL MOCK ERROR TESTS
// =========================

func mockRepo(t *testing.T) (*answer.Repository, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()

	dbConn, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := answer.NewRepository(dbConn, zap.NewNop())

	return &repo, mock, dbConn
}

func Test_GetBaseAnswers_DBError(t *testing.T) {
	repo, mock, dbConn := mockRepo(t)
	defer func() { _ = dbConn.Close() }()

	mock.ExpectQuery("select").
		WillReturnError(errors.New("db error"))

	_, err := (*repo).GetBaseAnswers(context.Background(), 1, 1)

	require.Error(t, err)
}

func Test_GetHardAnswers_DBError(t *testing.T) {
	repo, mock, dbConn := mockRepo(t)
	defer func() { _ = dbConn.Close() }()

	mock.ExpectQuery("select").
		WillReturnError(errors.New("db error"))

	_, err := (*repo).GetHardAnswers(context.Background(), 1, 1)

	require.Error(t, err)
}

func Test_SaveBaseAnswers_EmptyInput(t *testing.T) {
	repo, _, dbConn := mockRepo(t)
	defer func() { _ = dbConn.Close() }()

	ok, err := (*repo).SaveBaseAnswers(context.Background(), 1, []int{})

	require.NoError(t, err)
	assert.True(t, ok)
}

func Test_SaveBaseAnswers_DBError(t *testing.T) {
	repo, mock, dbConn := mockRepo(t)
	defer func() { _ = dbConn.Close() }()

	mock.ExpectExec("INSERT").
		WillReturnError(errors.New("insert failed"))

	ok, err := (*repo).SaveBaseAnswers(context.Background(), 1, []int{1, 2})

	require.Error(t, err)
	assert.False(t, ok)
}

func Test_DeleteAttempt_DBError(t *testing.T) {
	repo, mock, dbConn := mockRepo(t)
	defer func() { _ = dbConn.Close() }()

	mock.ExpectExec("delete").
		WillReturnError(errors.New("delete failed"))

	ok, err := (*repo).DeleteAttempt(context.Background(), 1, 1)

	require.Error(t, err)
	assert.False(t, ok)
}
