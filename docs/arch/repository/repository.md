## **UserRepository**

**Зона ответственности:**
Работа с пользователями — получение информации, создание/обновление и проверка роли.

**Используемые таблицы:** `User`, `Lecturer`

**Методы:**

| Метод                    | Аргументы          | Возврат      | Описание                                            |
| ------------------------ | ------------------ | ------------ | --------------------------------------------------- |
| `Upsert(user)`           | `CreateUserParams` | `User`       | Создает или обновляет пользователя в таблице `User` |
| `GetUserById(id)`        | `int`              | `User`       | Получает пользователя по `id`                       |
| `IsLecturer(id)`         | `int`              | `bool`       | Проверяет, является ли пользователь лектором        |
| `GetLecturers()`         | —                  | `list[User]` | Получает список всех лекторов                       |
| `CreateLecturer(id)`     | `int`              | `bool`       | Ставит флаг `is_lecturer` у пользователя            |
| `DeleteLecturer(id)`     | `int`              | `bool`       | Удаляет флаг `is_lecturer` у пользователя           |
| `GetAllImgs(lector_id)`  | `int`              | `list[str]`  | Получает все загруженные изображения лектора        |
| `GetAllFiles(lector_id)` | `int`              | `list[str]`  | Получает все загруженные файлы тестов лектора       |
---

## **TestRepository**

**Зона ответственности:**
Управление тестами, вопросами, ответами и правами на доступ.

**Используемые таблицы:** `Test`, `Task`, `Answer`, `TestPermission`

**Методы:**

| Метод                                        | Аргументы    | Возврат           | Описание                                      |
| -------------------------------------------- | ------------ | ----------------- | --------------------------------------------- |
| `CreateTest(test)`                           | `CreateTestParams`       | int              | Создает новый тест, возвращает его id                            |
| `DeleteTest(test_id)`                        | `int`        | bool              | Удаляет тест по `id`                          |
| `GetById(test_id)`                           | `int`        | `TestInfo`        | Получает тест по `id`                         |
| `GetHardTasks(test_id)`                      | `int`        | list[`Task`]      | Получает сложные вопросы для теста            |
| `GetBaseTasks(test_id)`                      | `int`        | list[`Task`]      | Получает базовые вопросы для теста            |
| `GetHardAnswers(test_id)`                    | `int`        | list[`TaskAnswer`]| Получает правильные ответы на сложные вопросы |
| `GetBaseAnswers(test_id)`                    | `int`        | list[`TaskAnswer`]| Получает правильные ответы на базовые вопросы |
| `GiveAccess(test_id, group_name)`            | `int`, `str` | bool              | Дает группе доступ к тесту                    |
| `GetGroups(test_id)`                         | `int`        | list[`GroupInfo`] | Получает список групп с доступом к тесту      |
| `GetLecturerTests(user_id)`                  | `str`        | list[`TestInfo`]  | Получает список тестов, назначенных лектору   |
| `HasLecturerAccess(user_id, test_id)`        | `int`, `int` | bool              | Проверяет доступ лектора к тесту              |



## Answer

Структура одного варианта ответ из вопроса

```
Answer: {
  text: str,
  image_url: str,
  is_correct: bool
}
```

---

## Task

Струкутра одного вопроса из теста

```
Task: {
  text: str,
  image_url: str,
  is_hard: bool,
  answers: list[Answer]
}

```

## TestInfo

Информация о тесте лектора.

```text
TestInfo: {
    id: int,
    lecturer_id: int,
    name: str,
    cnt_questions: int,
    cnt_hard_questions: int,
    groups: list[str],
    date_created: str,
}
```
---

## CreateTestParams


```
CreateTestParams: {
  name: str,
  lecturer_id: int,
  file_name: str,
  tasks: list[Task]
}
``` 

## TaskAnswer

Ответ на один вопрос.

```text
TaskAnswer: {
    task_id: int,
    options: list[int]
}
```

## GroupInfo

Информация о группе, имеющей доступ к тесту.

```text
GroupInfo: {
    group_name: str,
    cnt: int, // сколько студентов в этой группе
    solve: int, // сколько прошло тест (время конца тестирования не null)
    mean_mark: float, // средний балл среди тех кто прошел
    mean_percent: float, // средний процент среди тех кто прошел
}
```
---

## **StudentTestRepository**

**Зона ответственности:**
Хранение и управление прохождением тестов студентами, их ответами и результатами.

**Используемые таблицы:** `StudentTest`, `StudentAnswer`, `User`

**Методы:**

| Метод                                    | Аргументы                    | Возврат                 | Описание                                       |
| ---------------------------------------- | ---------------------------- | ----------------------- | ---------------------------------------------- |
| `Create(user_id, test_id)`               | `int`, `int`                 | bool                    | Создает запись о прохождении теста студентом   |
| `HasAccess(user_id, test_id)`            | `int`, `int`                 | bool                    | Проверяет, имеет ли студент доступ к тесту     |
| `GetTests(user_id)`                      | `int`                        | list[`StudentTestInfo`] | Получает список тестов, доступных студенту     |
| `SetHardAnswers(user_id, answers)`       | `int`, list[`StudentAnswer`] | bool                    | Сохраняет ответы студента на сложные вопросы   |
| `SetBaseAnswers(user_id, answers)`       | `int`, list[`StudentAnswer`] | bool                    | Сохраняет ответы студента на базовые вопросы   |
| `GetHardAnswers(user_id, test_id)`       | `int`, `int`                 | list[`TaskAnswer`]      | Получает ответы студента на сложные вопросы    |
| `GetBaseAnswers(user_id, test_id)`       | `int`, `int`                 | list[`TaskAnswer`]      | Получает ответы студента на базовые вопросы    |
| `SetDateStart(user_id, test_id, date)`   | `int`, `int`, `str`          | bool                    | Устанавливает дату начала теста студенту       |
| `SetDateEnd(user_id, test_id, date)`     | `int`, `int`, `str`          | bool                    | Устанавливает дату окончания теста студенту    |
| `SetMark(user_id, test_id, mark)`        | `int`, `int`, `float`        | bool                    | Устанавливает оценку студента по тесту         |
| `SetSuccessRate(user_id, test_id, rate)` | `int`, `int`, `float`        | bool                    | Устанавливает процент правильных ответов       |
| `GetResult(user_id, test_id)`            | `int`, `int`                 | `TestResult`            | Получает результат прохождения теста студентом |
| `GetGroupResult(group_name)`             | `str`                        | list[`StudentResult`]   | Получает результаты студентов группы по тесту  |

---
