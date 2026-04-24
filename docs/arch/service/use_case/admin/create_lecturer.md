# CreateLecturerUseCase

**Назначение:**
Выдает пользователю роль лектора, если пользователь еще не существуте берет его с LDAP.

**Возвращает:**
CreateLecturerResponse

```
CreateLecturerResponse {
  success: bool
}
```

**Репозитории / сервисы:**

* **UserRepository**
* **LdapAdapter** 

**Что происходит внутри:**

1. Принимает `userID` пользователя (админа), `login` пользователя (будующий лектор).
2. Проверка, что пользователь администратор через `UserRepository.GetByID(user_id)`.
3. Получает данные пользователя из LDAP: `LdapAdapter.GetInfo(login)`.
4. Вызывает `UserRepository.Upsert(params)` — создает или обновляет пользователя.
5. Назначает роль лектора через `UserRepository.CreateLecturer(userID)`.
6. Возвращает результат выполнения.

**Интерфейс репозитория:**

```go
UserRepository.GetByID(ctx context.Context, userID int) (User, error)
UserRepository.Upsert(ctx context.Context, params CreateUserParams) (int, error)
UserRepository.CreateLecturer(ctx context.Context, userID int) error
```

**Интерфейс LDAP адаптера:**

```go
LdapAdapter.GetInfo(ctx context.Context, login string) (*LdapUserInfo, error)
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

## LdapUserInfo
```go
type LdapUserInfo struct {
	Login string
	Name  string
	Mail  string
	Group *string
}
```

