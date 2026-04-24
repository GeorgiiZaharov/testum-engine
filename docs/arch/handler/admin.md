# **AdminHandler**

Эндпоинты для администрирования пользователей и назначения ролей.
Использует use cases: GetLecturersUseCase, CreateLecturerUseCase, DeleteLecturerUseCase
-----------------------------------------------------------------------------

## **GET /admin/lecturers + jwt**

**query**: `offset`, `limit`  
**resp**: `list[LecturerInfo]`

Возвращает список всех лекторов системы.

```
GetLecturersUseCase(lecturerRepository)
```

---

## **POST /admin/{lecturer_name} + jwt**

**resp**: `ok`

Добавляет пользователя в список лекторов.

```
CreateLecturerUseCase(lecturerRepository)
```

---

## **DELETE /admin/{lecturer_name} + jwt**

**resp**: `ok`

Удаляет пользователя из списка лекторов.

```
DeleteLecturerUseCase(lecturerRepository)
```

---
