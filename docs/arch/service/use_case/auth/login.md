# AuthLoginUseCase

**Назначение:**
Аутентифицирует пользователя через LDAP, создаёт или обновляет пользователя в системе и выдаёт JWT access и refresh токены.

---

## **Возвращает:**

```go
AuthLoginResponse {
  AccessToken  string
  RefreshToken string
}
```

---

## **Репозитории / сервисы:**

* **UserRepository**
* **LdapAdapter**
* **AuthService**

---

## **Что происходит внутри:**

1. Принимает `login` и `password`.
2. Проверка учетных данных через LDAP:

   ```go
   LdapAdapter.ValidatePassword(ctx, login, password) error
   ```
3. Получает данные пользователя из LDAP:

   ```go
   LdapAdapter.GetInfo(ctx, login) (*LdapUserInfo, error)
   ```
4. Создаёт или обновляет пользователя в системе:

   ```go
   UserRepository.Upsert(ctx, params) (userID int, error)
   ```
5. Генерирует JWT access токен:

   ```go
   AuthService.GenerateAccess(userID int) (string, error)
   ```
6. Генерирует JWT refresh токен:

   ```go
   AuthService.GenerateRefresh(userID int) (string, error)
   ```
7. Возвращает пару токенов клиенту.

---

## **Интерфейс репозитория:**

```go
UserRepository.Upsert(ctx context.Context, params CreateUserParams) (int, error)
```

---

## **Интерфейс LDAP адаптера:**

```go
LdapAdapter.ValidatePassword(ctx context.Context, login, password string) error

LdapAdapter.GetInfo(ctx context.Context, login string) (*LdapUserInfo, error)
```

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

---

## **LdapUserInfo**

```go
type LdapUserInfo struct {
	Login string
	Name  string
	Mail  string
	Group *string
}
```
