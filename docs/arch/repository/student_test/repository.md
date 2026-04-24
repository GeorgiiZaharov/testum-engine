# StudentTestRepository

**Зона ответственности:**  
Управление прохождением теста студентом.

**Используемые таблицы:** `StudentTest`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `StartTest(user_id, test_id, group)` | int, int, str | bool | Создает попытку прохождения теста |
| `FinishTest(user_id, test_id, result)` | int, int, TestResult | bool | Завершает тест устанавливая результат (дату окончания проставляет БД)|
| `GetActiveTests(user_id)` | int | `[]StudentActiveTestInfo` | Получает список тестов, которые доступны пользователю, но не решенные |
| `GetFinishedTests(user_id)` | int | `[]StudentFinishTestInfo` | Получает список решенных тестов (по которым есть оценка)|

## TestResult
TestResult {
  mark: int | nil,
  success_rate: float | nil,
}

## StudentActiveTestInfo
StudentTestInfo {
  id: int,
  name: str,
  lecturer_name: str,
  cnt_questions: int,
  cnt_hard_questions: int,
  date_start: time | nil,
}

## StudentFinishTestInfo
StudentTestInfo {
  id: int,
  name: str,
  lecturer_name: str,
  cnt_questions: int,
  cnt_hard_questions: int,
  mark: int,
  success_rate: float,
  date_start: time,
  date_end: time,
}
