# UserRepository

**Зона ответственности:**  
Управление пользователями и ролями (лекторы).

**Используемые таблицы:** `User`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `Upsert(params)` | `CreateUserParams` | int | Создает или обновляет пользователя |
| `GetByID(user_id)` | int | `User` | Получает данные пользователя |
| `GetLecturers()` | — | `[]User` | Получает список всех лекторов |
| `CreateLecturer(userID)` | `int` | bool | Назначает пользователя лектором |
| `DeleteLecturer(userID)` | `int` | bool | Удаляет роль лектора |

## CreateUserParams
```go
type CreateUserParams struct {
	Login string
	Mail  string
	Name  string
	Group *string
}
```

## User
```go
type User struct {
	ID           int       `db:"id"`
	Login        string    `db:"login"`
	Mail         string    `db:"mail"`
	Name         string    `db:"name"`
	Group        *string   `db:"group"`
	IsLecturer   bool      `db:"is_lecturer"`
	DateCreated  time.Time `db:"date_created"`
	DateModified time.Time `db:"date_modified"`
}
```
