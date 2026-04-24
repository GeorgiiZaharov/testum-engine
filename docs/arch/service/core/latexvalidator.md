# LaTeX Validator Module

## Назначение

Модуль используется для:
- валидации LaTeX-строк в тестах
- проверки корректности скобок
- проверки LaTeX-команд
- подсказок при ошибках написания команд

Используется в `FileParserService` и `ValidationService`.

---

## Основная структура

### Validator

```go
type Validator struct {
	dict *Dictionary
}
````

---

## Интерфейс (неявный)

Явного интерфейса нет, используется структура:

```go
type Validator struct
```

Рекомендуемый интерфейс:

```go
type LatexValidator interface {
	Validate(text string) []ValidationError
}
```

---

## Модель ошибок

### ValidationError

```go
type ValidationError struct {
	Line  int
	Error string
}
```

* `Line` — номер строки
* `Error` — текст ошибки

---

## Метод

### Validate

```go
Validate(text string) []ValidationError
```

#### Принимает:

* `text string` — LaTeX текст

#### Возвращает:

* `[]ValidationError`

#### Логика:

1. Разделяет текст на строки
2. Для каждой строки:

   * проверяет скобки (`checkBrackets`)
   * проверяет команды (`checkCommands`)
3. Собирает все ошибки

---

## Проверка скобок

### checkBrackets

```go
checkBrackets(line string, lineNum int) []ValidationError
```

#### Логика:

* Используется стек
* Проверяются `()` и `{}`
* Ошибки:

  * лишняя закрывающая скобка
  * неправильный порядок
  * незакрытая скобка

---

## Проверка команд

### checkCommands

```go
checkCommands(line string, lineNum int, dict *Dictionary) []ValidationError
```

#### Логика:

* Ищет LaTeX команды через regex: `\\[a-zA-Z_]+`
* Проверяет наличие в словаре
* Если команда неверная:

  * ищет похожую через `Suggest`
  * возвращает предупреждение с подсказкой

---

## Dictionary

### Назначение

Хранит список допустимых LaTeX команд и функций.

---

### Структура

```go
type Dictionary struct {
	functions []string
}
```

---

### Инициализация

```go
NewDictionary()
```

Содержит список стандартных LaTeX-команд:

* математические функции (`frac`, `sum`, `lim`)
* операторы (`leq`, `geq`)
* греческие буквы (`alpha`, `beta`, `omega`)

---

## Методы Dictionary

### IsValid

```go
IsValid(cmd string) bool
```

#### Логика:

* бинарный поиск по отсортированному списку
* проверка существования команды

---

### Suggest

```go
Suggest(word string) string
```

#### Логика:

* сравнивает слово со всеми командами
* использует similarity (Левенштейн)
* возвращает:

  * ближайшее совпадение
  * или несколько через `, `
* если сходство < 0.5 → пустая строка

---

## Алгоритмы

### similarity

```go
similarity(a, b string) float64
```

* основан на расстоянии Левенштейна
* нормализуется по длине строки

---

### distance

```go
distance(a, b string) int
```

* классический DP алгоритм Левенштейна
* учитывает вставку / удаление / замену

---

## Проверка команд (regex)

```go
\\([a-zA-Z_]+)
```

Извлекает LaTeX команды вида:

* `\frac`
* `\alpha`
* `\sum`

---

## Ошибки

Модуль не использует `error`, только:

```go
[]ValidationError
```

---

## Особенности

* не падает при ошибках — всегда возвращает список
* поддерживает подсказки при опечатках
* работает построчно
* полностью stateless (кроме словаря)
