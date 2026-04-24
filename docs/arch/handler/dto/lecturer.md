# Lecturer DTO (LecturerRouter)

DTO используемые преподавателем для управления тестами.

---

## FormatError

Ошибка формата тестового файла.

```
FormatError: {
    error: str,
}
```

---

## ParseError

Ошибка парсинга тестового файла.

```
ParseError: {
    error: str,
    line: int
}
```

---

## CreateResponse

Ответ при создании теста.

```
CreateTestResponse: {
    format_error: list[FormatError],
    parse_error: list[ParseError]
}
```

---

## TestFile

Ответ на запрос файла теста

```
TestFile: {
    file_link: str
}
```

---

## TestInfo

Информация о тесте лектора.

```text
TestInfo: {
    id: int,
    name: str,
    cnt_questions: int,
    date_created: str
}
```

---

## GroupInfo

Информация о группе, имеющей доступ к тесту.

```text
GroupInfo: {
    group_name: str
}
```

---

## StudentResult

Результат студента в тесте.

```text
StudentResult: {
    id: int,
    login: str,
    name: str,
    result: TestResult
}
```

---
