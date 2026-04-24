# GetHardTasksUseCase

## **Назначение:**

Возвращает сложные вопросы для теста студента.

---

## **Возвращает:**

```go
GetHardTasksResponse {
  HardTasks []Task
}
```

---

## **Репозитории:**

* `AccessRepository`
* `TaskRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (студент)
   * `testID` (идентификатор теста)

2. Проверяет, имеет ли студент доступ к тесту:

```go
AccessRepository.HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
```

3. Если доступ разрешён, запрашивает у репозитория сложные задания для данного теста:

```go
TaskRepository.GetHardTasks(ctx, testID) -> []Task
```

4. Возвращает список сложных заданий для теста.

---

## **Интерфейсы репозиториев:**

### **AccessRepository**

```go
AccessRepository.HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
```

### **TaskRepository**

```go
TaskRepository.GetHardTasks(ctx context.Context, testID int) ([]Task, error)
```

---

## **Модели**

### **Task**

```go
type Task struct {
  Text     string
  ImageURL string
  IsHard   bool
  Answers  []Answer
}
```

### **Answer**

```go
type Answer struct {
  Text     string
  ImageURL string
}
```
