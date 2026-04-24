реализуй данный репозиторий
напиши содержимое файлов 
- errors.go
- models.go
- repostitory.go

используй zap для логирования

вот пример:
## repository.go
```go
package user

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

type repository interface {
	upsert(ctx context.context, params createuserparams) (int, error)
	getbyid(ctx context.context, userid int) (user, error)
	getlecturers(ctx context.context) ([]user, error)
	createlecturer(ctx context.context, login string) (bool, error)
	deletelecturer(ctx context.context, login string) (bool, error)
}

type repository struct {
	db  *sql.db
	log *zap.logger
}

func newrepository(db *sql.db, log *zap.logger) repository {
	return &repository{
		db:  db,
		log: log,
	}
}

// ================= upsert =================
func (r *repository) upsert(ctx context.context, params createuserparams) (int, error) {
	query := `
		insert into users (login, mail, name, ` + "`group`" + `, is_lecturer, date_created, date_modified)
		values (?, ?, ?, ?, false, now(), now())
		on duplicate key update
			mail = values(mail),
			name = values(name),
			` + "`group`" + ` = values(` + "`group`" + `),
			date_modified = now()
	`

	res, err := r.db.execcontext(ctx,
		query,
		params.login,
		params.mail,
		params.name,
		params.group,
	)
	if err != nil {
		r.log.error("upsert failed",
			zap.error(err),
			zap.string("login", params.login),
		)
		return 0, err
	}

	id, err := res.lastinsertid()
	if err != nil || id == 0 {
		// fallback — получить id по login
		err = r.db.queryrowcontext(ctx,
			`select id from users where login = ?`,
			params.login,
		).scan(&id)

		if err != nil {
			r.log.error("upsert get id failed",
				zap.error(err),
				zap.string("login", params.login),
			)
			return 0, err
		}
	}

	return int(id), nil
}
```

## errors.go
```go
package user

import "errors"

var (
	errusernotfound = errors.new("user not found")
)
```
 
## models.go
```go
package user

import "time"

type createuserparams struct {
	login string
	mail  string
	name  string
	group *string
}

type user struct {
	id           int       `db:"id"`
	login        string    `db:"login"`
	mail         string    `db:"mail"`
	name         string    `db:"name"`
	group        *string   `db:"group"`
	islecturer   bool      `db:"is_lecturer"`
	datecreated  time.time `db:"date_created"`
	datemodified time.time `db:"date_modified"`
}
```
вот схема базы данных:

# модели базы данных
---

## **users**

| поле              | тип данных     | ограничения   | описание                              |
| ----------------- | -------------- | ------------- | ------------------------------------- |
| **id**            | `int`          | `primary key` | уникальный идентификатор пользователя |
| **login**         | `varchar(255)` | `not null`    | логин пользователя                    |
| **group**         | `varchar(100)` | `null`        | группа студента                       |
| **mail**          | `varchar(100)` | `not null`    | почта пользователя                    |
| **name**          | `varchar(255)` | `not null`    | имя пользователя                      |
| **is_lecturer**   | `boolean`      | `not null`    | флаг, указывающий является ли пользователь лектором |
| **date_modified** | `timestamp`    | `not null`    | дата последнего изменения записи      |
| **date_created**  | `timestamp`    | `not null`    | дата создания пользователя            |

---

## **tests**

| поле             | тип данных     | ограничения   | описание                       |
| ---------------- | -------------- | ------------- | ------------------------------ |
| **id**           | `int`          | `primary key` | уникальный идентификатор теста |
| **owner_id**     | `int`          | `not null`    | внешний ключ на `user.id`  |
| **name**         | `varchar(255)` | `not null`    | название теста                 |
| **file_name**    | `varchar(255)` | `not null`    | имя файла теста                |
| **date_created** | `timestamp`    | `not null`    | дата создания теста            |

---

## **tasks**

| поле                   | тип данных     | ограничения   | описание                             |
| ---------------------- | -------------- | ------------- | ------------------------------------ |
| **id**                 | `int`          | `primary key` | уникальный идентификатор вопроса     |
| **test_id**            | `int`          | `not null`    | внешний ключ на `test.id`            |
| **text**               | `text`         | `not null`    | текст вопроса                        |
| **image_url**          | `varchar(255)` | `null`        | url изображения для вопроса          |
| **is_hard**            | `boolean`      | `not null`    | флаг, определяющий сложность вопроса |

---

## **answers**

| поле           | тип данных     | ограничения   | описание                                 |
| -------------- | -------------- | ------------- | ---------------------------------------- |
| **id**         | `int`          | `primary key` | уникальный идентификатор ответа          |
| **task_id**    | `int`          | `not null`    | внешний ключ на `task.id`                |
| **text**       | `text`         | `not null`    | текст ответа                             |
| **image_url**  | `varchar(255)` | `null`        | url изображения для ответа               |
| **is_correct** | `boolean`      | `not null`    | флаг, указывающий на правильность ответа |

---

## **test_permissions**

| поле        | тип данных     | ограничения   | описание                                |
| ----------- | -------------- | ------------- | --------------------------------------- |
| **id**      | `int`          | `primary key` | уникальный идентификатор записи доступа |
| **test_id** | `int`          | `not null`    | внешний ключ на `test.id`               |
| **group**   | `varchar(100)` | `not null`    | группа, которая имеет доступ к тесту    |

---

## **student_answers**

| поле             | тип данных  | ограничения   | описание                                 |
| ---------------- | ----------- | ------------- | ---------------------------------------- |
| **id**           | `int`       | `primary key` | уникальный идентификатор ответа студента |
| **student_id**   | `int`       | `not null`    | внешний ключ на `user.id`             |
| **answer_id**    | `int`       | `not null`    | внешний ключ на `answer.id`              |
| **date_created** | `timestamp` | `not null`    | дата создания ответа                     |

---

## **student_tests**

| поле             | тип данных  | ограничения   | описание                                          |
| ---------------- | ----------- | ------------- | ------------------------------------------------- |
| **id**           | `int`       | `primary key` | уникальный идентификатор записи прохождения теста |
| **student_id**   | `int`       | `not null`    | внешний ключ на `user.id`                         |
| **test_id**      | `int`       | `not null`    | внешний ключ на `test.id`                         |
| **mark**         | `int`       | `null`        | оценка студента                                   |
| **group*         | `varchar(100)` | `not null` | группа студента на момент прохождения теста       |
| **success_rate** | `float`     | `null`        | процент правильных ответов                        |
| **date_start**   | `timestamp` | `null`        | дата начала теста                                 |
| **date_end**     | `timestamp` | `null`        | дата завершения теста                             |

---
