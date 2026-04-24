# GetLecturersUseCase

**Назначение:**  
Возвращает список всех лекторов системы.

**Возвращает:**
[]User
```
User {
  id: int,
  login: str,
	mail: str,
	name: str,
  group: str | nil,
	date_created: time,
	date_modified: time,
}
```

**Репозитории / сервисы:**

- Использует **UserRepository**

**Что происходит внутри:**

1. Принимает `userID`.
2. Проверка, что пользователь администратор через `UserRepository.GetByID(user_id)`.
3. Вызывается метод `UserRepository.GetLecturers()`.
4. Возвращает данные клиенту.

**Интерфейс репозитория:**

```go
UserRepository.GetByID(ctx context.Context, userID int) (User, error)
UserRepository.GetLecturers(ctx context.Context) ([]User, error)
```

## User
```go 
type User struct {
	ID           int
	Login        string
	Mail         string
	Name         string
	Group        *string
	IsLecturer   bool
	DateCreated  time.Time
	DateModified time.Time
}
```
