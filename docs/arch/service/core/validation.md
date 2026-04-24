# Validation Module

## Назначение

Модуль используется для:
- парсинга теста из текстового файла
- валидации структуры теста
- проверки вопросов, ответов и изображений
- преобразования файла в доменную модель `TestFromFile`

Используется в `FileParserService`.

---

## Входная модель результата

### TestFromFile

```go
type TestFromFile struct {
	Name      string
	HardCount int
	Tasks     []Task
}
````

---

### Task

```go
type Task struct {
	Text     string
	ImageURL string
	IsHard   bool
	Answers  []Answer
}
```

---

### Answer

```go
type Answer struct {
	Text      string
	ImageURL  string
	IsCorrect bool
}
```

---

## Интерфейс (неявный)

```go
type Validator interface {
	Validate(lines []string) (*TestFromFile, []ValidationError, error)
}
```

---

## Основной метод

### Validate

```go
Validate(lines []string) (*TestFromFile, []ValidationError, error)
```

---

### Принимает:

* `[]string` — строки файла теста

---

### Возвращает:

* `*TestFromFile` — распарсенный тест (если нет ошибок)
* `[]ValidationError` — список ошибок валидации
* `error` — системная ошибка (например, чтение файла)

---

### Логика работы:

1. Инициализирует состояние парсера
2. Обрабатывает файл построчно через state machine
3. Поддерживает состояния:

   * Title
   * HardCount
   * Question
   * Answer
4. Строит структуру `TestFromFile`
5. Валидирует:

   * наличие имени теста
   * наличие вопросов
   * наличие ответов
   * наличие правильных ответов
6. Возвращает ошибки или результат

---

## State Machine

```go
StateStart
StateTitle
StateHardCount
StateQuestionStart
StateQuestionBody
StateAnswerStart
StateAnswerBody
```

---

### Поведение:

* `Title` — парсит название теста
* `HardCount` — число сложных вопросов
* `QuestionStart` — начало нового вопроса
* `QuestionBody` — текст вопроса
* `AnswerStart` — начало ответа
* `AnswerBody` — продолжение ответа

---

## Формат входного файла

### Общая структура:

```
Название теста
Количество сложных вопросов
# Вопрос 1
Ответы...
# Вопрос 2
...
```

---

### Вопрос

```
# Текст вопроса
```

---

### Ответы

```
+ правильный ответ
- неправильный ответ
```

---

### Изображения

```
https://site.com/image.png
```

Поддерживаемые форматы:

* png
* jpg
* jpeg
* gif
* webp

---

## Правила валидации

### Ошибки:

* два вопроса подряд без ответов
* вопрос без ответов
* отсутствие правильного ответа
* у вопроса более одного изображения
* у ответа более одного изображения
* отсутствие названия теста
* отсутствие вопросов
* отсутствие количества сложных вопросов

---

## ValidationError

```go
type ValidationError struct {
	Error string
}
```

* содержит текст ошибки с номером строки

---

## Изображения

### extractImage

* принимает строку
* возвращает `*string` если строка — URL изображения
* иначе `nil`

---

## Особенности

* реализован state machine парсер
* поддерживает перенос текста между строками
* автоматически разделяет hard/base вопросы
* собирает ошибки, но не падает
* не использует panic

---

## Ограничения

* строгий формат входного файла
* нет восстановления после критических ошибок структуры
* изображения только через URL (локальные файлы не поддерживаются)
