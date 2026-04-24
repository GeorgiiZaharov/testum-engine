# GetFinishedTestUseCase

## **Назначение:**

Возвращает список завершённых тестов студента, включая информацию о результатах (оценка, процент успеха и время прохождения).

---

## **Возвращает:**

```go
GetFinishedTestResponse {
  FinishedTests []StudentFinishTestInfo
}
```

---

## **Репозитории:**

* `StudentTestRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (студент)

2. Получает список завершённых тестов студента, для которых есть оценка и дата завершения:

```go
StudentTestRepository.GetFinishedTests(ctx, userID) -> []StudentFinishTestInfo
```

3. Формирует ответ, который содержит список завершённых тестов с результатами.

---

## **Интерфейсы репозиториев:**

### StudentTestRepository

```go
StudentTestRepository.GetFinishedTests(ctx context.Context, userID int) ([]StudentFinishTestInfo, error)
```

---

## **Модели**

### StudentFinishTestInfo

```go
type StudentFinishTestInfo struct {
	ID               int
	Name             string
	LecturerName     string
	CntQuestions     int
	CntHardQuestions int
	Mark             int
	SuccessRate      float64
	DateStart        time.Time
	DateEnd          time.Time
}
```
