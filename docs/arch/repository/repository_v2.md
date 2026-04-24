# Repository Structure

# UserRepository

**Зона ответственности:**  
Управление пользователями и ролями (лекторы).

**Используемые таблицы:** `User`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `Upsert(params)` | `CreateUserParams` | `User` | Создает или обновляет пользователя |
| `GetLecturers()` | — | `[]LecturerInfo` | Получает список всех лекторов |
| `CreateLecturer(login)` | `string` | bool | Назначает пользователя лектором |
| `DeleteLecturer(login)` | `string` | bool | Удаляет роль лектора |

---

# FileRepository

**Зона ответственности:**  
Работа с файлами тестов и изображениями.

**Используемые компоненты:** `FileStorage`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `GetAllTestFiles(lecturer_id)` | int | `[]string` | Получает все файлы тестов лектора |
| `GetAllImages(lecturer_id)` | int | `[]string` | Получает все изображения лектора |

---

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

---

# AccessRepository

**Зона ответственности:**  
Управление доступами к тестам (RBAC на уровне тестов).

**Используемые таблицы:** `TestPermission`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `HasLecturerAccess(user_id, test_id)` | int, int | bool | Проверяет доступ лектора к тесту |
| `HasStudentAccess(user_id, test_id)` | int, int | bool | Проверяет доступ студента к тесту |
| `GiveAccess(test_id, group)` | int, string | bool | Выдает доступ группе |
| `TakeAccess(test_id, group)` | int, string | bool | Забирает доступ у группы |

---

# ResultRepository

**Зона ответственности:**  
Получение результатов тестов (без логики подсчета).

**Используемые таблицы:** `StudentTest`, `StudentAnswer`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `GetGroupResult(test_id, group)` | int, string | `[]Result` | Получает результаты группы |
| `GetStudentResult(user_id, test_id)` | int, int | `TestResult` | Получает результат студента |

---

# StudentTestRepository

**Зона ответственности:**  
Управление прохождением теста студентом.

**Используемые таблицы:** `StudentTest`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `StartTest(user_id, test_id)` | int, int | bool | Создает попытку прохождения теста |
| `SetDateStart(user_id, test_id)` | int, int | bool | Устанавливает время начала |
| `SetDateEnd(user_id, test_id)` | int, int | bool | Устанавливает время окончания |
| `GetResult(user_id, test_id)` | int, int | `TestResult` | Получает результат |
| `GetTests(user_id)` | int | `[]StudentTestInfo` | Получает список тестов |

---

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

---

# TaskRepository

**Зона ответственности:**  
Получение заданий теста.

**Используемые таблицы:** `Task`

| Метод | Аргументы | Возврат | Описание |
|------|----------|--------|----------|
| `GetHardTasks(test_id)` | int | `[]Task` | Получает сложные задания |
| `GetBaseTasks(test_id)` | int | `[]Task` | Получает базовые задания |

