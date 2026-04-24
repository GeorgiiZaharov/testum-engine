# **StudentHandler**

Эндпоинты для студента по доступу к тестам, прохождению тестов и получению результатов.
Использует UseCases: **GetTestsUseCase**
-------------------------------------------

## **GET /student/tests + jwt**

**query**: `offset`, `limit`
**resp**: `list[StudentTestInfo]`

Возвращает список доступных студенту тестов.
Приоритет выдачи:
1) Незаконченные тесты
2) Не начатые тесты
3) Завершенные тесты

```
GetTestsUseCase(StudentTestRepository)
```

## **GET /student/test/{test_id} + jwt**

**resp**: `StudentTestInfo`

Возвращает информацию о конкретном тесте.

```
GetTestUseCase(StudentTestRepository)
```

## **GET /student/test/{test_id}/hard + jwt**

**resp**: `list[Task]`

Возвращает сложные вопросы для теста.

```
GetHardTasksUseCase(StudentTestRepository)
```

## **GET /student/test/{test_id}/base + jwt**

**resp**: `list[Task]`

Возвращает базовые вопросы для теста.

```
GetBaseTasksUseCase(StudentTestRepository)
```

## **POST /student/test/{test_id}/hard + jwt**

**query**: `test_id`
**body**: `list[TaskAnswer]`
**resp**: `HardTestResult`

Сохраняет результаты сложных вопросов теста.

```
PostHardAnswersUseCase(StudentTestRepository, TestRepository, AnswerCheckService)
```

## **POST /student/test/{test_id}/base + jwt**

**body**: `list[TaskAnswer]`
**resp**: `ok`

Сохраняет результаты базовых вопросов теста.

```
PostBaseAnswersUseCase(ResultCalculationService, AnswerCheckService, StudentTestRepository, TestRepository)
```

## **GET /student/test/{test_id}/result + jwt**

**resp**: `TestResult`

Возвращает результат прохождения теста.

```
GetTestResultUseCase(ResultCalculationService, AnswerCheckService, StudentTestRepository, TestRepository)
```

