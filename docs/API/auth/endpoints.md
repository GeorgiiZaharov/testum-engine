## Endpoint summary

* Method: `POST`
* URL: `/auth/login`
* Auth: `Not required`
* Description: Authenticates user via LDAP and returns access and refresh tokens.

---

## Request

### Request body

```json
{
  "login": "string",
  "password": "string"
}
```

---

## Response

### Success (HTTP 200)

```json
{
  "access_token": "string",
  "refresh_token": "string"
}
```

---

## Errors

* 400 — invalid request body
* 400 — login and password are required
* 400 — invalid input
* 401 — invalid credentials
* 500 — internal server error

---

## cURL example

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "john.doe",
    "password": "secret"
  }'
```

---

## Endpoint summary

* Method: `POST`
* URL: `/auth/refresh`
* Auth: `Required`
* Description: Refreshes JWT tokens for authenticated user.

---

## Request

---

## Response

### Success (HTTP 200)

```json
{
  "access_token": "string",
  "refresh_token": "string"
}
```

---

## Errors

* 401 — unauthorized
* 401 — invalid user id
* 500 — internal server error

---

## cURL example

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Authorization: Bearer <token>"
```

---

## Endpoint summary

* Method: `GET`
* URL: `/auth/me`
* Auth: `Required`
* Description: Returns current authenticated user profile.

---

## Response

### Success (HTTP 200)

```json
{
  "id": 0,
  "login": "string",
  "mail": "string",
  "name": "string",
  "group": "string",
  "is_lecturer": true
}
```

---

## Errors

* 400 — invalid input
* 401 — unauthorized
* 404 — not found
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/auth/me \
  -H "Authorization: Bearer <token>"
```
