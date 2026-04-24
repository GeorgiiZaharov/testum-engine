# Storage Adapter

## Назначение

Адаптер используется для:
- работы с файловым хранилищем
- сохранения и удаления файлов и изображений
- изоляции файловой системы через интерфейс (для тестов и подмены)

---

## Интерфейсы

### SaltGenerator

```go
type SaltGenerator interface {
	Generate() (string, error)
}
```

Используется для генерации уникального суффикса имени файла.

### StorageAdapter

```go
type Storage interface {
	UploadFile(file io.Reader, fileName string) error
	GetFile(fileName string) (*os.File, error)
	DeleteFile(fileName string) error

	UploadPicture(file io.Reader, fileName string, login string) (string, error)
	DeletePictures(login string) error
}
```

---

## Основная структура

### StorageAdapter

```go
type StorageAdapter struct {
	fs        FileSystem
	saltGen   SaltGenerator
	basePath  string
	imagePath string
}
```

### Зависимости:

* `FileSystem` — абстракция файловой системы
* `SaltGenerator` — генератор случайных строк

---

## Методы

### UploadFile

```go
UploadFile(file io.Reader, fileName string) error
```

#### Принимает:

* `io.Reader`
* `fileName string`

#### Возвращает:

* `error`

#### Логика:

1. Валидирует имя файла
2. Проверяет существование файла
3. Создаёт файл
4. Копирует содержимое

---

### GetFile

```go
GetFile(fileName string) (*os.File, error)
```

#### Принимает:

* `fileName string`

#### Возвращает:

* `*os.File`
* `error`

#### Логика:

1. Валидирует имя
2. Проверяет существование
3. Если нет — возвращает `nil`
4. Иначе открывает файл

---

### DeleteFile

```go
DeleteFile(fileName string) error
```

#### Принимает:

* `fileName string`

#### Возвращает:

* `error`

#### Логика:

1. Валидирует имя
2. Проверяет существование
3. Если файла нет — ничего не делает
4. Удаляет файл

---

### UploadPicture

```go
UploadPicture(file io.Reader, fileName string, login string) (string, error)
```

#### Принимает:

* `io.Reader`
* `fileName string`
* `login string`

#### Возвращает:

* `string` (новое имя файла)
* `error`

#### Логика:

1. Валидирует имя файла и login
2. Создаёт директорию `data/images/{login}`
3. Генерирует salt
4. Формирует имя: `original_salt.ext`
5. Сохраняет файл
6. Возвращает новое имя

---

### DeletePictures

```go
DeletePictures(login string) error
```

#### Принимает:

* `login string`

#### Возвращает:

* `error`

#### Логика:

1. Валидирует login
2. Проверяет существование директории пользователя
3. Если нет — ничего не делает
4. Удаляет директорию полностью

---

## Ошибки

```go
var (
	ErrFileExists
	ErrInvalidName
)
```

