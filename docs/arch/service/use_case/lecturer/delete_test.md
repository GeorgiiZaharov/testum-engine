# DeleteTestUseCase

## **Назначение:**

Удаляет тест лектора из системы после проверки прав доступа.

---

## **Возвращает:**

```go
DeleteTestResponse {
  Success bool
}
```

---

## **Репозитории:**

* `LecturerTestRepository`
* `AccessRepository`
* `StorageAdapter`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (лектор)
   * `testID`

---

2. Проверяет доступ лектора к тесту:

```go
AccessRepository.HasLecturerAccess(ctx, userID, testID) (bool, error)
```

Если доступа нет → возвращается ошибка `ErrAccessDenied`.

---
3. Получает информацию о тесте:
```go 
LecturerTestRepository.GetByID(ctx context.Context, testID int) (TestInfo, error)
```

4. Удаляет файл теста с платформы:
```go
StorageAdapter.DeleteFile(fileName string) error 
```

5. Удаляет тест:

```go
LecturerTestRepository.Delete(ctx, testID) (bool, error)
```

---

6. Возвращает результат операции.

---

## **Интерфейс репозитория доступа:**

```go
AccessRepository.HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
```

---

## **Интерфейс репозитория тестов:**

```go
LecturerTestRepository.GetByID(ctx context.Context, testID int) (TestInfo, error)
LecturerTestRepository.Delete(ctx context.Context, testID int) (bool, error)
```

## **Интерфейс адаптера:**

```go
StorageAdapter.DeleteFile(fileName string) error 
```

---

## **Бизнес-логика:**

* удаление возможно только владельцем теста
* каскадное удаление затрагивает:

  * задания
  * ответы
  * доступы

