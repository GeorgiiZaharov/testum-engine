# GetGroupsUseCase

## **Назначение:**

Возвращает список групп, которым доступен тест и/или, которые его проходили по годам

---

## **Возвращает:**

[]GroupInfo
```go
type GroupInfo struct {
	GroupName    string
	MembersCount int
}
```

---

## **Репозитории:**

* `AccessRepository`
* `LecturerTestRepository`

---

## **Что происходит внутри:**

1. Принимает:

   * `userID` (лектор)
   * `testID`
   * `year` (0 - текущий год, 1 - предыдущий и тд. не может быть отрицательным)

---

2. Проверяет доступ лектора к тесту:

```go
AccessRepository.HasLecturerAccess(ctx, userID, testID) (bool, error)
```

Если доступа нет → возвращается ошибка `ErrAccessDenied`.

---

3. Получает информацию о группе, которой доступен тест или которая его решала:

```go
LecturerTestRepository.GetGroups(ctx context.Context, testID int, year int) ([]GroupInfo, error)
```

---

4. Возвращает результат.

---

## **Интерфейсы репозиториев:**

```go
AccessRepository.HasLecturerAccess(ctx, userID, testID) (bool, error)
LecturerTestRepository.GetGroups(ctx context.Context, testID int, year int) ([]GroupInfo, error)
```

---

## **Модели**

### TestInfo

```go
type TestInfo struct {
	GroupName    string
	MembersCount int
}
```

