# LecturerTestRepository

**Зона ответственности:**  
Управление тестами лектора.

**Используемые таблицы:** `Test`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `Create(test)` | `Test` | bool | Создает тест |
| `Delete(test_id)` | int | bool | Удаляет тест |
| `GetByID(test_id)` | int | `TestInfo` | Получает тест |
| `GetByLecturer(user_id)` | int | `[]TestInfo` | Получает тесты лектора |
| `GetGroups(test_id, year)` | int, int | `[]GroupInfo` | Получает группы, проходившие тест в прошлом или группы имеющие сейчас к нему доступ |


## TestInfo

```
TestInfo {
  id: int,
  name: str,
  cnt_questions: int,
  cnt_hard_questions: int,
  file_name: str,
  groups: list[int],
  date_created: str
}
```

## Test
```
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

Test {
  name: str
  hard_count: int,
  file_name: str,
  tasks: list[Task]
}
```

## GroupInfo
```
GroupInfo {
  group_name: str,
  members_count: int
}
```
