# AccessRepository

**Зона ответственности:**  
Управление доступами к тестам.

**Используемые таблицы:** `TestPermission`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `HasLecturerAccess(user_id, test_id)` | int, int | bool | Проверяет доступ лектора к тесту |
| `HasStudentAccess(user_id, test_id)` | int, int | bool | Проверяет доступ студента к тесту |
| `GiveAccess(test_id, group)` | int, string | bool | Выдает доступ группе |
| `TakeAccess(test_id, group)` | int, string | bool | Забирает доступ у группы |

