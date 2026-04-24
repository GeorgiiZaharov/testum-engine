# Auth JWT Service

## Назначение

Сервис используется для:
- генерации JWT access и refresh токенов
- парсинга и валидации токенов

Используется в AuthService и middleware.

---

## Интерфейс

```go
type Service interface {
	GenerateAccess(userID int) (string, error)
	GenerateRefresh(userID int) (string, error)
	Parse(tokenStr string) (*Claims, error)
}
````

---

## Модели

### Claims

```go
type Claims struct {
	UserID int
	Type   TokenType
	jwt.RegisteredClaims
}
```

Содержит:

* `UserID` — идентификатор пользователя
* `Type` — тип токена (access / refresh)
* стандартные JWT поля (exp, iat)

---

### TokenType

```go
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)
```

---

## Методы

### GenerateAccess

```go
GenerateAccess(userID int) (string, error)
```

#### Принимает:

* `userID int`

#### Возвращает:

* `string` (JWT токен)
* `error`

#### Логика:

1. Создаёт claims с типом `access`
2. Устанавливает срок жизни: 15 минут
3. Подписывает токен с использованием HS256
4. Возвращает строку токена

---

### GenerateRefresh

```go
GenerateRefresh(userID int) (string, error)
```

#### Принимает:

* `userID int`

#### Возвращает:

* `string`
* `error`

#### Логика:

1. Создаёт claims с типом `refresh`
2. Устанавливает срок жизни: 7 дней
3. Подписывает токен
4. Возвращает строку

---

### Parse

```go
Parse(tokenStr string) (*Claims, error)
```

#### Принимает:

* `tokenStr string`

#### Возвращает:

* `*Claims`
* `error`

#### Логика:

1. Парсит токен с использованием секрета
2. Проверяет срок действия
3. Проверяет валидность claims
4. Проверяет наличие `UserID` и `Type`
5. Возвращает claims

---

## Ошибки

```go
var (
	ErrInvalidToken
	ErrExpiredToken
	ErrWrongType
)
```

* `ErrInvalidToken` — токен некорректен
* `ErrExpiredToken` — токен истёк
* `ErrWrongType` — неверный тип токена

