# TakeAccessUseCase

## **Назначение:**

Удаляет (отзывает) доступ группы к тесту.

---

## **Возвращает:**

```go
TakeAccessResponse {
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

2. Проверяет права лектора на тест:

```go
AccessRepository.HasLecturerAccess(ctx, userID, testID) -> bool
```

Если доступа нет → возвращается ошибка `ErrAccessDenied`.

---

3. Отзывает доступ у группы:

```go
AccessRepository.TakeAccess(ctx, testID, group) -> error
```

---

4. Возвращает результат операции.

---

## **Интерфейс репозитория:**

```go
AccessRepository.HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)

AccessRepository.TakeAccess(ctx context.Context, testID int, group string) error
```

