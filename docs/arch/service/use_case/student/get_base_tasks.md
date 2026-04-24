# GetBaseTasksUseCase

## **Назначение:**

Возвращает базовые вопросы для теста студента и фиксирует начало прохождения теста.

---

## **Возвращает:**

```go
GetBaseTasksResponse {
  BaseTasks []Task
}
```

---

## **Репозитории:**

* `StudentTestRepository`
* `ResultRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (студент)
   * `testID` (идентификатор теста)

---

2. Проверяет, имеет ли студент доступ к тесту:

```go
StudentTestRepository.HasAccess(ctx, userID, testID) -> bool
```

---

3. Проверяет, был ли тест уже завершён:

```go
ResultRepository.GetResult(ctx, userID, testID) -> TestResult
```

Если результат существует → доступ к базовым заданиям запрещён.

---

4. Отмечает начало прохождения теста (если ещё не начат):

```go
StudentTestRepository.StartTest(ctx, userID, testID, group) -> bool
```

---

5. Получает базовые задания:

```go
StudentTestRepository.GetBaseTasks(ctx, testID) -> []Task
```

---

6. Возвращает список базовых заданий.

---

## **Интерфейсы репозиториев:**

### **StudentTestRepository**

```go
StudentTestRepository.HasAccess(ctx context.Context, userID int, testID int) (bool, error)

StudentTestRepository.StartTest(ctx context.Context, userID int, testID int, group string) (bool, error)

StudentTestRepository.GetBaseTasks(ctx context.Context, testID int) ([]Task, error)
```

---

### **ResultRepository**

```go
ResultRepository.GetResult(ctx context.Context, userID int, testID int) (TestResult, error)
```

---

## **Модели**

### **Task**

```go
type Task struct {
  Text     string
  ImageURL string
  IsHard   bool
  Answers  []Answer
}
```

---

### **Answer**

```go
type Answer struct {
  Text     string
  ImageURL string
}
```

---

### **TestResult**

```go
type TestResult struct {
	Mark        int
	SuccessRate float64
	DateStart   time.Time
	DateEnd     time.Time
}
```

