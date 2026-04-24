# UploadTestUseCase

## **Назначение:**

Загружает тест лектора в систему: валидирует файл, парсит его структуру, сохраняет файл в storage и создаёт запись о тесте в БД.

---

## **Возвращает:**

```go
UploadTestResponse {
  FormatError FormatError[]
  ValidationError ValidationError[]
  TestID *int
  Success bool
}
```

---

## **Репозитории / сервисы:**

* **LecturerTestRepository**
* **UserRepository**
* **AccessRepository**

---

## **Адаптеры:**

* **StorageAdapter**

---

## **Сервисы:**

* **FileValidationService**
* **FileParserService**

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (лектор)
   * `file io.Reader`
   * `fileName string`

---

2. Проверяет пользователя и роль лектора:

```go
UserRepository.GetByID(ctx, userID) -> User
```

---

3. Проверяет файл на корректность формата:

```go
FileValidationService.Validate(file) -> []FormatError
```

Если есть ошибки — возвращает их и прерывает выполнение.

---

4. Парсит файл в структуру теста:

```go
FileParserService.Parse(file) -> ([]ValidationError, Test)
```

---

5. Сохраняет файл в storage:

```go
StorageAdapter.UploadFile(file, fileName) -> error
```

---

6. Создаёт тест в базе данных:

```go
LecturerTestRepository.Create(ctx, test) -> testID int
```

---

7. (Опционально) создаёт доступы к группам:

```go
AccessRepository.CreateDefaultAccess(ctx, testID)
```

---

8. Возвращает результат загрузки.

---

## **Интерфейс репозитория:**

```go
LecturerTestRepository.Create(ctx context.Context, test Test) (int, error)
```

---

## **Интерфейс сервисов:**

```go
FileValidationService.Validate(file []byte) []FormatError

FileParserService.Parse(file []byte) ([]ValidationError, Test)
```

---

## **Интерфейс адаптера:**

```go
StorageAdapter.UploadFile(file io.Reader, fileName string) error
```

---

## **User**

```go
type User struct {
	ID           int
	Login        string
	Mail         string
	Name         string
	Group        *string
	IsLecturer   bool
	DateCreated  time.Time
	DateModified time.Time
}
```

---

## **Test (domain)**

```go
type Test struct {
	Name      string
	HardCount int
	Tasks     []Task
}
```

---

## **Task**

```go
type Task struct {
	Text     string
	ImageURL string
	IsHard   bool
	Answers  []Answer
}
```

---

## **Answer**

```go
type Answer struct {
	Text      string
	ImageURL  string
	IsCorrect bool
}
```

## **ValidationError**
```go
type ValidationError struct {
	Line  int
	Error string
}
```

## **FormatError**
```go
type FormatError struct {
	Error string
}
```


