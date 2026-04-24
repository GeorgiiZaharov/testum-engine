# LDAP Adapter

## Назначение

Адаптер используется для:
- проверки учетных данных пользователя через LDAP
- получения информации о пользователе из LDAP

Используется в AuthService.

---

## Интерфейс

```go
type LdapAdapter interface {
	ValidatePassword(ctx context.Context, login, password string) error
	GetInfo(ctx context.Context, login string) (*LdapUserInfo, error)
}
````

---

## Модели

### LdapUserInfo

```go
type LdapUserInfo struct {
	Login string
	Name  string
	Mail  string
	Group *string
}
```

---

## Методы

### ValidatePassword

```go
ValidatePassword(ctx context.Context, login, password string) error
```

#### Принимает:

* `context.Context`
* `login string`
* `password string`

#### Возвращает:

* `error`

#### Логика:

1. Проверяет, что login и password не пустые
2. Устанавливает LDAP соединение
3. Выполняет поиск пользователя по `uid`
4. Если пользователь найден — получает его `DN`
5. Выполняет `Bind` с DN и паролем
6. Если bind успешен — пароль корректный

---

### GetInfo

```go
GetInfo(ctx context.Context, login string) (*LdapUserInfo, error)
```

#### Принимает:

* `context.Context`
* `login string`

#### Возвращает:

* `*LdapUserInfo`
* `error`

#### Логика:

1. Проверяет, что login не пустой
2. Устанавливает LDAP соединение
3. Выполняет поиск пользователя по `uid`
4. Извлекает атрибуты:

   * `uid`
   * `cn`
   * `mail`
5. Дополнительно извлекает группу из DN
6. Формирует и возвращает `LdapUserInfo`

---

## Ошибки

```go
var (
	ErrConnectionFailed
	ErrSearchFailed
	ErrUserNotFound
	ErrInvalidPassword
	ErrEmptyCredentials
	ErrEmptyLogin
)
```
