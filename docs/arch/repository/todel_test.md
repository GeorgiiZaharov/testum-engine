напиши тесты для этого кода, которые обеспечат 100% покрытие кода
пиши тесты только для:

- Query errors
- NotFound
- Happy path
- rows.Err()


вот пример:
```go
package user

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
)

type testEnv struct {
	repo Repository
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

	repo := NewRepository(database.DB, zap.NewNop())

	return &testEnv{
		repo: repo,
		fx:   fx,
		ctx:  ctx,
	}
}

// =========================
// REAL DB TESTS
// =========================
func TestGetByID_Success(t *testing.T) {
	env := setup(t)

	u, err := env.repo.GetByID(env.ctx, 1)
	require.NoError(t, err)

	assert.Equal(t, "student1", u.Login)
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
```

вот какие фикстуры заданы в проекте:
```md
# Фикстуры
# users

**Описание:** пользователи системы (студенты + преподаватели).
Используется для авторизации, групп и прав доступа.

```md
| id | login        | mail                | name          | group   | is_lecturer | date_created | date_modified |
|----|-------------|---------------------|---------------|---------|-------------|--------------|----------------|
| 1  | student1    | student1@mail.com   | Student One   | A-101   | false       | now          | now            |
| 2  | student2    | student2@mail.com   | Student Two   | A-101   | false       | now          | now            |
| 3  | student3    | student3@mail.com   | Student Three | B-202   | false       | now          | now            |
| 4  | old_student | old@mail.com        | Old Student   | A-101   | false       | now-1y       | now-1y         |
| 5  | magistr     | magistr@mail.com    | Magistr       | C-000   | true        | now          | now            |
| 6  | lecturer1   | lecturer1@mail.com  | Lecturer One  | NULL    | true        | now          | now            |
| 7  | lecturer2   | lecturer2@mail.com  | Lecturer Two  | NULL    | true        | now          | now            |
```

---

# tests

**Описание:** тесты, созданные лекторами.
Связь: `owner_id → users.id`

```md
| id | owner_id    | name                | file_name                | created_at |
|----|-------------|---------------------|--------------------------|------------|
| 1  | 5           | Math Basics         | math_basics.json         | now        |
| 2  | 6           | Linear Algebra      | linear_algebra.json      | now        |
| 3  | 6           | Physics Intro       | physics_intro.json       | now        |
| 4  | 7           | Programming Basics  | programming_basics.json  | now        |
```

---

# test_permissions

**Описание:** доступ групп к тестам.

```md
| id | test_id | group |
|----|--------|-------|
| 1  | 1      | A-101 |
| 2  | 1      | B-202 |
| 3  | 2      | A-101 |
| 4  | 3      | C-000 |
```

---

# tasks

**Описание:** вопросы тестов.

```
-- TEST 1
| id | test_id | text                          | image_url            | is_hard |
|----|--------|-------------------------------|----------------------|----------|
| 1  | 1      | Solve equation: 2 + 2 = ?     | img/test1_easy.png   | false    |
| 2  | 1      | Prove derivative of x^2       | NULL                 | true     |

-- TEST 3
| id | test_id | text                                      | image_url            | is_hard |
|----|--------|-------------------------------------------|----------------------|----------|
| 3  | 3      | What is Newton's second law?              | img/test3_easy.png   | false    |
| 4  | 3      | Derive kinetic energy formula from work theorem | img/test3_hard.png | true     |
```

---

# answers

**Описание:** варианты ответов

```md
| id | task_id| text                                     | image_url                     | is_correct |
|----|--------|------------------------------------------|-------------------------------|------------|
| 1  | 1      | 4                                        | NULL                          | true       |
| 2  | 1      | 5                                        | NULL                          | false      |

| 3  | 2      | 2x                                       | img/answer_math_hard_1.png    | true       |
| 4  | 2      | x^2                                      | img/answer_math_hard_2.png    | false      |

| 5  | 3      | Force = mass * acceleration              | NULL                          | true       |
| 6  | 3      | Force = mass / acceleration              | NULL                          | false      |

| 7  | 4      | E = (mv^2)/2 derived from work integral  | img/answer_physics_hard_1.png | true       |
| 8  | 4      | E = (mv^2)/2 derived from work integral  | img/answer_physics_hard_2.png | true      |
```

---

# student_tests

```md
| id | student_id| test_id| mark | group | success_rate | date_start | date_end |
|----|-----------|--------|------|-------|--------------|------------|----------|
| 1  | 2         | 1      | NULL | A-101 | NULL         | t          | NULL     |
| 2  | 3         | 1      | 4    | B-202 | 80.0         | t-24h      | t        |
| 3  | 4         | 1      | 5    | A-101 | 100.0        | t-1y-1h    | t-1y     |
```

---

# student_answers

```md
| id | student_id| answer_id | created_at |
|----|-----------|-----------|------------|
| 1  | 3         | 2         | t-1h       |
| 2  | 3         | 3         | t-1h       |
| 3  | 4         | 1         | t-1y       |
```
```
# Используй фикстуры по максимуму
