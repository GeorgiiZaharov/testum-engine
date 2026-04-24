# **AuthRouter**

Эндпоинты логина и получения информации о пользователе.
Используют uses cases: **AuthLoginUseCase**, **AuthRefreshUseCase**, **GetCurrentUserUseCase**
-------------------------------

## **POST /auth/login**

**body**: `AuthLoginRequest`
**resp**: `Tokens`

Получение токенов по логину и паролю.

```
AuthLoginUseCase(AuthService, LdapAdapter, UserRepository) -> Tokens
```

## **POST /auth/refresh**

**body**: `AuthRefreshRequest`
**resp**: `Tokens`

Проверка refresh токена и выдача новой пары токенов.

```
AuthRefreshUseCase(AuthRefreshRequest, UserRepository) -> Tokens
```

## **GET /auth/me + jwt**

**resp**: `UserInfo`

Возвращает информацию о пользователе и его права в системе.

```
GetCurrentUserUseCase(UserRepository) -> UserInfo
```
---
