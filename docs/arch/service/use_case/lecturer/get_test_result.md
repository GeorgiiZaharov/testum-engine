# GetTestResultUseCase

## **Назначение:**

Возвращает результаты студентов конкретной группы по тесту.

---

## **Возвращает:**

```go
GetTestResultResponse {
  Results []StudentResult
}
```

---

## **Репозитории:**

* **AccessRepository**
* **ResultRepository**

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (лектор)
   * `testID`
   * `group` (название группы)

---

2. Проверяет доступ лектора к тесту:

```go
AccessRepository.HasLecturerAccess(ctx, userID, testID) (bool, error)
```

* если доступа нет → ошибка (forbidden)

---

3. Получает результаты группы:

```go
ResultRepository.GetGroupResult(ctx, testID, group) ([]StudentResult, error)
```

---

4. Возвращает список результатов.

---

## **Интерфейсы:**

### AccessRepository

```go
AccessRepository.HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
```

---

### ResultRepository

```go
ResultRepository.GetGroupResult(ctx context.Context, testID int, group string) ([]StudentResult, error)
```

---

## **Модели**

### Result

```go
type StudentResult struct {
	UserID int
	Name   string
	Login  string
	Mail   string

	Mark        *int
	SuccessRate *float64
	DateStart   *time.Time
	DateEnd     *time.Time
}
```

