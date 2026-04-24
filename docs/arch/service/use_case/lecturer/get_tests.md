# GetTestsUseCase

## **Назначение:**

Возвращает список тестов, созданных конкретным лектором.

---

## **Возвращает:**

```go
GetTestsResponse {
  Tests []TestInfo
}
```

---

## **Репозитории:**

* `LecturerTestRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (лектор)

---

2. Получает список тестов лектора:

```go
LecturerTestRepository.GetByLecturer(ctx, userID) ([]TestInfo, error)
```

---

3. Возвращает результат клиенту.

---

## **Интерфейс репозитория:**

```go
LecturerTestRepository.GetByLecturer(ctx context.Context, userID int) ([]TestInfo, error)
```

---

## **Модель TestInfo:**

```go
type TestInfo struct {
	ID                int
	Name              string
	CntQuestions      int
	CntHardQuestions  int
	Groups            []string
	DateCreated       time.Time
}
```
