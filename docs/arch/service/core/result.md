# Result Calculation Service

## Назначение

Сервис используется для:
- расчёта итоговой оценки за тест
- вычисления процента успешных ответов
- преобразования результата проверки в систему оценок (2–5)

Используется после `Answer Check Service`.

---

## Модель оценки

### Mark

```go
type Mark int
````

```go
const (
	MarkBad          Mark = iota + 2 // 2
	MarkSatisfactory                 // 3
	MarkGood                         // 4
	MarkExcellent                    // 5
)
```

---

## Результат расчёта

### CalcResult

```go id="r3"
type CalcResult struct {
	Mark        Mark
	SuccessRate float64
}
```

* `Mark` — итоговая оценка
* `SuccessRate` — процент правильных ответов

---

## Интерфейс

```go id="r4"
type CalculationService interface {
	Calc(res answer.CheckResult) CalcResult
}
```

---

## Метод

### Calc

```go
Calc(res answer.CheckResult) CalcResult
```

#### Принимает:

* `answer.CheckResult`

  * `TrueCnt`
  * `Total`

#### Возвращает:

* `CalcResult`

#### Логика:

1. Если `Total == 0`:

   * возвращает `Mark = 2`
   * `SuccessRate = 0`
2. Вычисляет процент:

   ```
   successRate = TrueCnt / Total * 100
   ```
3. Определяет оценку через `calcMark`
4. Возвращает результат

---

## Логика оценки

### calcMark

```go
func calcMark(rate float64) Mark
```

#### Правила:

* `rate == 100` → `5 (Excellent)`
* `rate >= 80` → `4 (Good)`
* `rate >= 50` → `3 (Satisfactory)`
* `< 50` → `2 (Bad)`

---

## Особенности

* не имеет состояния (stateless)
* не обращается к БД
* вся логика детерминирована
* легко тестируется через unit tests

---

## Ограничения

* фиксированная шкала оценок (2–5)
* нет конфигурации порогов (захардкожены)

