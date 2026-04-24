# DeleteLecturerUseCase

**Назначение:**
Удаляет лектора и связанные с ним файлы.

**Возвращает:**
DeleteLecturerResponse
```
DeleteLecturerResponse {
  success: bool
}
```

**Репозитории / сервисы:**

* **UserRepository**
* **FileRepository**
* **StorageAdapter**

**Что происходит внутри:**

1. Принимает `userID` пользователя (админа), `lectorID` пользователя (текущий лектор).
2. Проверка, что пользователь администратор через `UserRepository.GetByID(user_id)`(userID).
3. Проверка, что собираемя удалить лектора через `UserRepository.GetByID(user_id)`(lectorID).
4. Получает список файлов с тестами `FileRepository.GetAllTestFiles(lecturerID)`.
5. Удаляет файлы через `StorageAdapter.DeleteFile(fileName string)`, `StorageAdapter.DeletePictures(login)`
6. Удаляет роль лектора `DeleteLecturer(userID int)`.
7. Возвращает результат.

**Интерфейс репозиториев:**

```go
UserRepository.GetByID(userID int) (User, error)
UserRepository.GetAllTestFiles(lecturerID int) ([]string, error)
UserRepository.DeleteLecturer(lecturerID int) error
```

**Интерфейс адаптера:**
```go
StorageAdapter.DeleteFile(fileName string) error
StorageAdapter.DeletePictures(login string) error
```

## User

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
