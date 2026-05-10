# DeleteAttemptUseCase

## **Назначение:**

Удаляет попытку прохождения теста студентом.

---

## **Возвращает:**

```go
DeleteAttemptResponse {
  Success bool
}
```

---

## **Репозитории:**

* `AccessRepository`
* `ResultRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `lecturerID`
   * `userID`
   * `testID`

---

2. Проверяет доступ лектора к тесту:

```go
AccessRepository.HasLecturerAccess(ctx, lecturerID, testID) (bool, error)
```

Если доступа нет → возвращается ошибка `ErrAccessDenied`.

---
3. Удаляет попытку студента:
```go 
ResultRepository.DeleteAttempt(ctx context.Context, testID int, userID int) error 
```

---

## **Интерфейс репозитория доступа:**

```go
AccessRepository.HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
```

---

## **Интерфейс репозитория тестов:**

```go
ResultRepository.DeleteAttempt(ctx context.Context, testID int, userID int) error
```
---

## **Бизнес-логика:**

* удаление попытки возможно только владельцем теста
* удаление затрагивает:

  * ответы
  * результат

