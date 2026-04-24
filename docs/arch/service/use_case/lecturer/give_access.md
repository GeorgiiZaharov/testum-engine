# GiveAccessUseCase

## **Назначение:**

Выдаёт доступ к тесту конкретной группе студентов.

---

## **Возвращает:**

```go
GiveAccessResponse {
  Success bool
}
```

---

## **Репозитории:**

* `AccessRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (лектор)
   * `testID`
   * `group string`

---

2. Проверяет, что пользователь имеет доступ к тесту:

```go
AccessRepository.HasLecturerAccess(ctx, userID, testID) -> bool
```

Если доступа нет → возвращается ошибка `ErrAccessDenied`.

---

3. Выдаёт доступ группе:

```go
AccessRepository.GiveAccess(ctx, testID, group) -> error
```

---

4. Возвращает результат операции.

---

## **Интерфейс репозитория:**

```go
AccessRepository.HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)

AccessRepository.GiveAccess(ctx context.Context, testID int, group string) error
```

