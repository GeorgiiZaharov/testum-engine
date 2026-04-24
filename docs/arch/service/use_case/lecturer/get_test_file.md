# GetTestFileUseCase

## **Назначение:**

Возвращает файл теста лектора из хранилища.

---

## **Возвращает:**

```go
GetTestFileResponse {
  File *os.File
}
```

---

## **Репозитории / адаптеры:**

* **LecturerTestRepository**
* **AccessRepository**
* **StorageAdapter**

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

* если доступа нет → ошибка (forbidden)

---

3. Получает информацию о тесте:

```go
LecturerTestRepository.GetByID(ctx, testID) (TestInfo, error)
```

---

4. Проверяет наличие файла:

* если `file_name == ""` → ошибка (файл не найден)

---

5. Получает файл из storage:

```go
StorageAdapter.GetFile(fileName string) (*os.File, error)
```

---

6. Проверяет результат:

* если файл не найден (`nil`) → ошибка (file not found)

---

7. Возвращает файл клиенту.

---

## **Интерфейсы:**

### LecturerTestRepository

```go
LecturerTestRepository.GetByID(ctx context.Context, testID int) (TestInfo, error)
```

---

### AccessRepository

```go
AccessRepository.HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
```

---

### StorageAdapter

```go
StorageAdapter.GetFile(fileName string) (*os.File, error)
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
	FileName         string
	Groups           []int
	DateCreated      time.Time
}
```

---

## **Ошибки:**

```go
var (
	ErrForbidden    = errors.New("no access to test")
	ErrFileNotFound = errors.New("test file not found")
)
```
