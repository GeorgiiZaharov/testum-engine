# TaskRepository

**Зона ответственности:**  
Получение заданий теста.

**Используемые таблицы:** `Task`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `GetHardTasks(test_id)` | int | `[]Task` | Получает сложные задания |
| `GetBaseTasks(test_id)` | int | `[]Task` | Получает базовые задания |

## Task
Answer {
  text: str,
  image_url: str,
  is_correct: bool
}

Task {
  text: str,
  image_url: str,
  is_hard: bool,
	answers:  list[Answer]
}
