# Answer Check Service

## Назначение

Сервис используется для:
- проверки ответов студента
- сравнения с эталонными ответами
- подсчёта количества правильных ответов

Используется в логике прохождения теста.

---

## Интерфейс

```go
type CheckService interface {
	Check(studentAnswers []TaskAnswer, trueAnswers []TaskAnswer) (CheckResult, error)
}
````

---

## Модели

### Входные данные

```go
type TaskAnswer struct {
	TaskID          int
	SelectedOptions []int
}
```

### Результат

```go
type CheckResult struct {
	TrueCnt int
	Total   int
}
```

---

## Метод

### Check

```go
Check(studentAnswers []TaskAnswer, trueAnswers []TaskAnswer) (CheckResult, error)
```

#### Принимает:

* `[]TaskAnswer` — ответы студента
* `[]TaskAnswer` — правильные ответы

#### Возвращает:

* `CheckResult`
* `error`

#### Логика:

1. Проверяет совпадение количества задач
2. Строит map правильных ответов по `TaskID`
3. Для каждого ответа студента:

   * ищет соответствующий правильный ответ
   * сортирует варианты ответов (нормализация)
   * сравнивает списки
4. Считает количество совпадений
5. Возвращает результат

---

## Ошибки

```go
ErrTaskMismatch
```

Возникает если:

* количество задач не совпадает
* отсутствует соответствующий `TaskID`

