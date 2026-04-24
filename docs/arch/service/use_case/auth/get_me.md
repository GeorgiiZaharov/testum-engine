# GetMeUseCase

**Назначение:**
Возвращает информацию о текущем пользователе.
Если с момента последнего обновления данных прошло больше 2 недель — выполняется синхронизация с LDAP перед возвратом данных.

---

## **Возвращает:**

```go
GetMeResponse {
  ID         int
  Login      string
  Mail       string
  Name       string
  Group      *string
  IsLecturer bool
}
```

---

## **Репозитории / сервисы:**

* **UserRepository**
* **LdapAdapter**

---

## **Что происходит внутри:**

1. Принимает `userID` (из JWT middleware / контекста).

2. Получает пользователя из базы:

```go
UserRepository.GetByID(ctx, userID) (User, error)
```

---

3. Проверяет актуальность данных:

```
if time.Since(user.DateModified) > 14 * 24h
```

---

### Если данные устарели:

4. Получает актуальные данные из LDAP:

```go
LdapAdapter.GetInfo(ctx, user.Login) (*LdapUserInfo, error)
```

5. Обновляет пользователя в системе:

```go
UserRepository.Upsert(ctx, CreateUserParams) (userID int, error)
```

---

### Если данные актуальны:

* LDAP не вызывается
* используется текущий пользователь из БД

---

6. Возвращает актуальные данные пользователя.

---

## **Интерфейс репозитория:**

```go
UserRepository.GetByID(ctx context.Context, userID int) (User, error)

UserRepository.Upsert(ctx context.Context, params CreateUserParams) (int, error)
```

---

## **Интерфейс LDAP адаптера:**

```go
LdapAdapter.GetInfo(ctx context.Context, login string) (*LdapUserInfo, error)
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
