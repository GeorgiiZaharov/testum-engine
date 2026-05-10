# ResultRepository

**Зона ответственности:**  
Получение результатов тестов (без логики подсчета).

**Используемые таблицы:** `StudentTest`, `StudentAnswer`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `GetGroupResult(test_id, group, year)` | int, string, int | `[]StudentResult` | Получает результаты группы year лет назад|
| `GetStudentResult(user_id, test_id)` | int, int | `TestResult` | Получает результат студента |
| `DeleteAttempt(test_id, user_id)` | int, int | error | Удаляет попытку студента в тесте |

## TestResult
TestResult {
  mark: int,
  success_rate: float,
  date_start: time,
  date_end: time,
}
 
## StudentResult
StudentResult {
  user_id: int,
  name: str,
  login: str,
  mail: str,
  mark: int,
  success_rate: float,
  date_start: time,
  date_end: time,
}
