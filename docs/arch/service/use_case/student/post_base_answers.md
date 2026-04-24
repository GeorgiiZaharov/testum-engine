# PostBaseAnswersUseCase

## **Назначение:**

Сохраняет ответы студента на базовые задания, проверяет их корректность с учётом сложного блока и завершает тест с расчётом итогового результата.

---

## **Возвращает:**

```go
PostBaseAnswersResponse {
  Success bool
}
```

---

## **Репозитории / сервисы:**

* `AccessRepository`
* `AnswerRepository`
* `ResultRepository`
* `AnswerCheckService`
* `ResultCalculationService`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (студент)
   * `testID`
   * `answers []TaskAnswer`

---

2. Проверяет доступ студента к тесту:

```go
AccessRepository.HasStudentAccess(ctx, userID, testID) -> bool
```

---

3. Проверяет, не были ли уже отправлены базовые ответы:

```go
AnswerRepository.GetBaseAnswers(ctx, userID, testID) -> []TaskAnswer
```

Если ответы уже существуют → прекращает выполнение.

---

4. Проверяет, что сложный блок уже пройден:

```go
AnswerRepository.GetHardAnswers(ctx, userID, testID) -> []TaskAnswer
```

Если сложные ответы отсутствуют → выполнение запрещено.

---

5. Проверяет, что тест ещё не завершён:

```go
ResultRepository.GetStudentResult(ctx, userID, testID) -> TestResult
```

Если результат уже существует → выполнение запрещено.

---

6. Получает правильные ответы:

```go
AnswerRepository.GetHardAnswersByTest(ctx, testID) -> []TaskAnswer

AnswerRepository.GetBaseAnswersByTest(ctx, testID) -> []TaskAnswer
```

---

7. Проверяет ответы:

```go
AnswerCheckService.Check(hardStudentAnswers, hardTrueAnswers) -> CheckResult

AnswerCheckService.Check(baseStudentAnswers, baseTrueAnswers) -> CheckResult
```

---

8. Объединяет результаты:

```go
totalTrue  = hard.TrueCnt + base.TrueCnt
totalTasks = hard.Total   + base.Total
```

---

9. Рассчитывает итог:

```go
ResultCalculationService.Calc(CheckResult{TrueCnt: totalTrue, Total: totalTasks}) -> CalcResult
```

---

10. Сохраняет базовые ответы:

```go
AnswerRepository.SaveBaseAnswers(ctx, userID, answers) -> error
```

---

11. Завершает тест и сохраняет результат:

```go
StudentTestRepository.FinishTest(ctx, userID, testID, TestResult) -> bool
```

---

12. Возвращает статус выполнения.

---

## **Интерфейсы репозиториев:**

### AccessRepository

```go
AccessRepository.HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
```

---

### AnswerRepository

```go
AnswerRepository.GetHardAnswers(ctx context.Context, userID int, testID int) ([]TaskAnswer, error)

AnswerRepository.GetBaseAnswers(ctx context.Context, userID int, testID int) ([]TaskAnswer, error)

AnswerRepository.GetHardAnswersByTest(ctx context.Context, testID int) ([]TaskAnswer, error)

AnswerRepository.GetBaseAnswersByTest(ctx context.Context, testID int) ([]TaskAnswer, error)

AnswerRepository.SaveBaseAnswers(ctx context.Context, userID int, answers []TaskAnswer) error
```

---

### ResultRepository

```go
ResultRepository.GetStudentResult(ctx context.Context, userID int, testID int) (TestResult, error)
```

---

### StudentTestRepository

```go
StudentTestRepository.FinishTest(ctx context.Context, userID int, testID int, result TestResult) (bool, error)
```

---

## **Интерфейсы сервисов:**

### AnswerCheckService

```go
AnswerCheckService.Check(studentAnswers []TaskAnswer, trueAnswers []TaskAnswer) (CheckResult, error)
```

---

### ResultCalculationService

```go
ResultCalculationService.Calc(res CheckResult) CalcResult
```

---

## **Модели**

### TaskAnswer

```go
type TaskAnswer struct {
	TaskID  int
	Options []int
}
```

---

### CheckResult

```go
type CheckResult struct {
	TrueCnt int
	Total   int
}
```

---

### TestResult

```go
type TestResult struct {
	Mark        int
	SuccessRate float64
	DateStart   time.Time
	DateEnd     time.Time
}
```
