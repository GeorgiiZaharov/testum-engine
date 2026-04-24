# AnswerRepository

**Зона ответственности:**  
Работа с ответами студентов и правильными ответами.

**Используемые таблицы:** `StudentAnswer`, `Answer`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `GetHardAnswers(user_id, test_id)` | int, int | `[]TaskAnswer` | Ответы студента (сложный блок) |
| `GetBaseAnswers(user_id, test_id)` | int, int | `[]TaskAnswer` | Ответы студента (базовый блок) |
| `GetHardAnswersByTest(test_id)` | int | `[]TaskAnswer` | **Правильные ответы** (сложный блок) |
| `GetBaseAnswersByTest(test_id)` | int | `[]TaskAnswer` | **Правильные ответы** (базовый блок) |
| `SaveHardAnswers(user_id, answers)` | int, [] | bool | Сохраняет ответы на сложные вопросы |
| `SaveBaseAnswers(user_id, answers)` | int, [] | bool | Сохраняет ответы на базовые вопросы |
| `DeleteAttempt(user_id, test_id)` | int, int | bool | Удаляет попытку прохождения теста студентом |

## TaskAnswer
TaskAnswer {
  task_id: int,
	options: list[int]
}
