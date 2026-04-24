# DB Adapter

## Назначение

DB адаптер отвечает за:
- подключение к базе данных (MySQL)
- предоставление унифицированного интерфейса для выполнения SQL-запросов
- управление пулом соединений
- работу с транзакциями

Адаптер используется в repository-слое как низкоуровневый источник данных.

---

## Конфигурация

### DBOptions

```go
type DBOptions struct {
	Host string
	User string
	Pass string
	Name string
	Port string
}
````

Используется для формирования DSN строки подключения.

---

## Интерфейс Executor

```go
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
```

### Назначение

Позволяет абстрагировать:

* обычное соединение (`*sqlx.DB`)
* транзакцию (`*sql.Tx`)

### Использование

Repository принимает `Executor`, а не конкретную реализацию:

```go
func (r *repository) someMethod(ctx context.Context, exec Executor) error
```

Это позволяет:

* выполнять код как в транзакции, так и без неё
* упрощает тестирование (моки)

---

## Инициализация

### NewDB

```go
func NewDB(opts DBOptions) (*DB, error)
```

### Что делает

1. Формирует DSN строку:

   ```
   user:pass@tcp(host:port)/dbname?parseTime=true
   ```

2. Устанавливает соединение через `sqlx.Connect`

3. Настраивает пул соединений:
   * `MaxOpenConns = 25`
   * `MaxIdleConns = 25`
   * `ConnMaxLifetime = 5 минут`

### Возвращает

* `*DB` — готовый адаптер
* `error` — ошибка подключения

---

## Работа с транзакциями

### BeginTx

```go
func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error)
```

Создаёт новую транзакцию.

---

### WithTx

```go
func (db *DB) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
```

### Назначение

Упрощает работу с транзакциями и гарантирует:

* rollback при любой ошибке
* commit только при успешном выполнении

### Поведение

1. Начинает транзакцию
2. Вызывает переданную функцию `fn`
3. Если `fn` вернула ошибку:

   * транзакция откатывается
4. Если всё успешно:

   * вызывается `Commit`
