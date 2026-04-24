# Auth DTO (AuthRouter)

DTO используемые в эндпоинтах авторизации.

### AuthLoginRequest

```
AuthLoginRequest: {
    login: str,
    password: str
}
```

### AuthRefreshRequest

```
AuthRefreshRequest: {
    refresh_token: str
}
```

### Tokens

```
Tokens: {
    access_token: str,
    refresh_token: str
}
```

### UserInfo

```
UserInfo: {
    login: str,
    name: str,
    group: str,
    date_modified: str,
    date_created: str,
    roles: list[str]
}
```

---
