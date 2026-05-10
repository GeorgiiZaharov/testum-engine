# UploadPictureUseCase

## **Назначение:**

Загружает изображение на сервер. 

---

## **Возвращает:**

```go
UploadPictureResponse {
  url: str
}
```

---

## **Репозитории / сервисы:**

* **UserRepository**

---

## **Адаптеры:**

* **StorageAdapter**

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


3. Сохраняет файл в storage:

```go
StorageAdapter.UploadPicture(file, fileName, login) -> (string, error)
```


4. Возвращает результат загрузки.

---


## **Интерфейс адаптера:**

```go
StorageAdapter.UploadPicture(file, fileName, login) -> (string, error)
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



