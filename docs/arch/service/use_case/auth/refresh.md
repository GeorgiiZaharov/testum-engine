# AuthRefreshUseCase

**Назначение:**
Обновляет пару JWT токенов (access + refresh) по действующему userID из контекста.

---

## **Возвращает:**

```go
AuthRefreshResponse {
  AccessToken  string
  RefreshToken string
}
```

---

## **Репозитории / сервисы:**

* **AuthService**

---

## **Что происходит внутри:**

1. Принимает `userID` (из JWT middleware / контекста).

2. Генерирует новый access токен:

```go
AuthService.GenerateAccess(userID int) (string, error)
```

3. Генерирует новый refresh токен:

```go
AuthService.GenerateRefresh(userID int) (string, error)
```

4. Возвращает новую пару токенов клиенту.

---

## **Интерфейс сервиса:**

```go
AuthService.GenerateAccess(userID int) (string, error)

AuthService.GenerateRefresh(userID int) (string, error)
```

---

## **User**

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
