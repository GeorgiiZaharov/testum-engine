# GetActiveTestUseCase

## **Назначение:**

Возвращает список активных тестов, доступных студенту, которые он еще не прошел.

---

## **Возвращает:**

```go
GetActiveTestResponse {
  ActiveTests []StudentActiveTestInfo
}
```

---

## **Репозитории:**

* `StudentTestRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (студент)

2. Получает список активных тестов, которые доступны студенту, но еще не решены:

```go
StudentTestRepository.GetActiveTests(ctx, userID) -> []StudentActiveTestInfo
```

3. Формирует ответ, который содержит список активных тестов.

---

## **Интерфейсы репозиториев:**

### StudentTestRepository

```go
StudentTestRepository.GetActiveTests(ctx context.Context, userID int) ([]StudentActiveTestInfo, error)
```

---

## **Модели**

### StudentActiveTestInfo

```go
type StudentActiveTestInfo struct {
	ID               int
	Name             string
	LecturerName     string
	CntQuestions     int
	CntHardQuestions int
	DateStart        *time.Time
}
```
