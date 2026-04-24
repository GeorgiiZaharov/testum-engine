# GetTestResultUseCase

## **Назначение:**

Возвращает результат прохождения теста студентом.

---

## **Возвращает:**

```go
GetTestResultResponse {
  Mark        *int
  SuccessRate *float64
  DateStart   time.Time
  DateEnd     *time.Time
}
```

---

## **Репозитории:**

* `AccessRepository`
* `ResultRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (студент)
   * `testID`

---

2. Проверяет доступ студента к тесту:

```go
AccessRepository.HasStudentAccess(ctx, userID, testID) -> bool
```

---

3. Если доступ разрешён, получает результат теста:

```go
ResultRepository.GetStudentResult(ctx, userID, testID) -> TestResult
```

---

4. Формирует и возвращает ответ.

---

## **Интерфейсы репозиториев:**

### **AccessRepository**

```go
AccessRepository.HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
```

---

### **ResultRepository**

```go
ResultRepository.GetStudentResult(ctx context.Context, userID int, testID int) (TestResult, error)
```

---

## **Модели**

### **TestResult**

```go
type TestResult struct {
	Mark        *int
	SuccessRate *float64
	DateStart   time.Time
	DateEnd     *time.Time
}
```
