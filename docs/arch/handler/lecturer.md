# **LecturerHandler**

Эндпоинты для работы с тестами, управления доступом и просмотра результатов студентов.
Использует UseCases: **UploadTestUseCase**, **GetTestFileUseCase**, **DeleteTestUseCase**, **GiveAccessUseCase**, **GetTestsUseCase**, **GetGroupsUseCase**, **GetTestResultUseCase**.
-------------------------------------------

## **POST /lecturer/test + jwt**

**query**: `ignore_parser`
**file**: `test_file`
**resp**: `CreateResponse`

Загружает новый тест, валидирует и сохраняет его. Если `ignore_parser = true`, игнорирует предупреждения парсера.

```
UploadTestUseCase(UserRepository, TestRepository, StorageAdapter, FileValidationService, FileParserService)
```

## **GET /lecturer/test/{test_id} + jwt**

**resp**: `TestFile`

Возвращает файл для указанного теста.

```
GetTestFileUseCase(TestRepository)
```

## **DELETE /lecturer/test/{test_id} + jwt**

**resp**: `ok`

Удаляет указанный тест.

```
DeleteTestUseCase(TestRepository)
```

## **POST /lecturer/test/{test_id}/access + jwt**

**query**: `group`
**resp**: `ok`

Дает доступ к тесту для группы.

```
GiveAccessUseCase(TestRepository)
```

## **POST /lecturer/test/{test_id}/take + jwt**

**query**: `group`
**resp**: `ok`

Дает доступ к тесту для группы.

```
TakeAccessUseCase(TestRepository)
```

## **GET /lecturer/tests + jwt**

**query**: `offset`, `limit`
**resp**: `list[TestInfo]`

Возвращает список тестов текущего лектора.

```
GetTestsUseCase(TestRepository)
```

## **GET /lecturer/test/{test_id}/group/{group_name} + jwt**

**resp**: `list[StudentResult]`

Возвращает результаты студентов группы по тесту.

```
GetTestResultUseCase(StudentTestRepository)
```

## **DELETE /lecturer/attempt/test/{test_id}/user/{user_id} + jwt**

**resp**: `ok`

Удаляет попытку в тесте указанного студента.

```
DeleteAttemptUseCase(TestRepository)
```

## **POST /lecturer/picture + jwt**

**file**: `picture`
**resp**: `createpictureresponse`

загружает изображение на сервер.

```
uploadpictureusecase(userrepository, storageadapter)
```

