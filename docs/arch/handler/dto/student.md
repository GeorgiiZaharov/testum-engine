
# Student DTO (StudentRouter)

DTO используемые при работе студента с тестами.

---

## StudentTestInfo

Информация о тесте в списке доступных тестов.

```text
StudentTestInfo: {
    id: int,
    name: str,
    cnt_questions: int,
    creator: str,
    date_start: str,
    date_end: str,
    date_created: str
}
```

---

## TaskOption

Вариант ответа.

```text
TaskOption: {
    id: int,
    text: str,
    image_url: str
}
```

---

## Task

Описание вопроса теста.

```text
Task: {
    id: int,
    is_multiple_choice: bool,
    text: str,
    image_url: str,
    options: list[TaskOption]
}
```

---

## TaskAnswer

Ответ студента на один вопрос.

```text
TaskAnswer: {
    task_id: int,
    selected_options: list[int]
}
```

---

## HardTestResult

Результат прохождения адаптивной (сложной) части теста.

```text
HardTestResult: {
    success: bool,
    all_correct: bool
}
```

---

## TestResult

Итоговый результат теста.

```text
TestResult: {
    mark: float,
    success_rate: float,
    date_start: str,
    date_end: str
}
```

---
