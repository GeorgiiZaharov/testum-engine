# GetTestUseCase

## **Назначение:**

Возвращает полную информацию о тесте лектора:

* мета-информацию
* список групп
* блоки вопросов (hard / base)

---

## **Возвращает:**

```go
GetTestResponse {
  ID               int
  Name             string
  CntQuestions     int
  CntHardQuestions int
  Groups           []GroupInfo

  HardTasks []Task
  BaseTasks []Task
}
```

---

## **Репозитории:**

* `LecturerTestRepository`
* `TaskRepository`
* `AccessRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (лектор)
   * `testID`

---

2. Проверяет доступ лектора к тесту:

```go
AccessRepository.HasLecturerAccess(ctx, userID, testID) (bool, error)
```

Если доступа нет → возвращается ошибка `ErrAccessDenied`.

---

3. Получает основную информацию о тесте:

```go
LecturerTestRepository.GetByID(ctx, testID) (TestInfo, error)
```

---

4. Получает список групп с доступом / прохождением:

```go
LecturerTestRepository.GetGroups(ctx, testID, 0) ([]GroupInfo, error)
```

> `year = 0` — текущий год (активные группы)

---

5. Получает сложные задания:

```go
TaskRepository.GetHardTasks(ctx, testID) ([]Task, error)
```

---

6. Получает базовые задания:

```go
TaskRepository.GetBaseTasks(ctx, testID) ([]Task, error)
```

---

7. Формирует ответ:

* объединяет `TestInfo`
* добавляет `Groups`
* добавляет `HardTasks` и `BaseTasks`

---

8. Возвращает результат.

---

## **Интерфейсы репозиториев:**

```go
AccessRepository.HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)

LecturerTestRepository.GetByID(ctx context.Context, testID int) (TestInfo, error)

LecturerTestRepository.GetGroups(ctx context.Context, testID int, year int) ([]GroupInfo, error)

TaskRepository.GetHardTasks(ctx context.Context, testID int) ([]Task, error)

TaskRepository.GetBaseTasks(ctx context.Context, testID int) ([]Task, error)
```

---

## **Модели**

### TestInfo

```go
type TestInfo struct {
	ID               int
	Name             string
	CntQuestions     int
	CntHardQuestions int
	Groups           []int
	DateCreated      time.Time
}
```

---

### GroupInfo

```go
type GroupInfo struct {
	GroupName    string
	MembersCount int
}
```

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

