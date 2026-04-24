# PostHardAnswersUseCase

## **Назначение:**

Сохраняет ответы студента на сложные задания теста, проверяет их корректность и при необходимости завершает тест с расчётом результата.

---

## **Возвращает:**

```go
PostHardAnswersResponse {
  IsAllCorrect bool
}
```

---

## **Репозитории / сервисы:**

* `AccessRepository`
* `AnswerRepository`
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

3. Проверяет, не были ли уже отправлены ответы:

```go
AnswerRepository.GetHardAnswers(ctx, userID, testID) -> []TaskAnswer
```

Если ответы уже существуют → прекращает выполнение (тест уже пройден).

---

4. Получает правильные ответы:

```go
AnswerRepository.GetHardAnswersByTest(ctx, testID) -> []TaskAnswer
```

---

5. Проверяет ответы студента:

```go
AnswerCheckService.Check(studentAnswers, trueAnswers) -> CheckResult
```

---

6. Сохраняет ответы студента:

```go
AnswerRepository.SaveHardAnswers(ctx, userID, answers) -> bool
```

---

7. Если проверка успешная (все ответы верны):

   * рассчитывает результат:

```go
ResultCalculationService.Calc(checkResult) -> CalcResult
```

* сохраняет результат:

```go
AnswerRepository.SaveResult(ctx, userID, testID, CalcResult) -> bool
```

---

8. Возвращает результат проверки:

```go
IsAllCorrect = (checkResult.TrueCnt == checkResult.Total)
```

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

AnswerRepository.GetHardAnswersByTest(ctx context.Context, testID int) ([]TaskAnswer, error)

AnswerRepository.SaveHardAnswers(ctx context.Context, userID int, answers []TaskAnswer) error

AnswerRepository.SaveResult(ctx context.Context, userID int, testID int, result CalcResult) error
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

## **Модели:**

### TaskAnswer

```go
type TaskAnswer struct {
	TaskID          int
	Options         []int
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

### CalcResult

```go
type CalcResult struct {
	Mark        int
	SuccessRate float64
}
```
