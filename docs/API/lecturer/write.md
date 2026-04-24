## Endpoint summary

* Method: `DELETE`
* URL: `/lecturer/tests/{test_id}`
* Auth: `Required`
* Description: Deletes a test owned by lecturer.

---

## Request

### Path params

| Name    | Type | Required | Description            |
| ------- | ---- | -------- | ---------------------- |
| test_id | int  | yes      | Unique test identifier |

### Success (HTTP 200)

```json
{
  "success": true
}
```

---

## Errors

* 400 — invalid test_id
* 400 — invalid input
* 401 — unauthorized
* 403 — access denied
* 404 — test not found
* 500 — internal server error

---

## cURL example

```bash
curl -X DELETE http://localhost:8080/lecturer/tests/123 \
  -H "Authorization: Bearer <token>"
```

---

---

## Endpoint summary

* Method: `POST`
* URL: `/lecturer/tests/{test_id}/access`
* Auth: `Required`
* Description: Grants access to a test for a specific group.

---

## Request

### Path params

| Name    | Type | Required | Description            |
| ------- | ---- | -------- | ---------------------- |
| test_id | int  | yes      | Unique test identifier |

---

### Request body

```json
{
  "group": "string"
}
```

---

## Response

### Success (HTTP 200)

```json
{
  "success": true
}
```

---

## Errors

* 400 — invalid test_id
* 400 — invalid body
* 400 — invalid input
* 401 — unauthorized
* 403 — access denied
* 500 — internal server error

---

## cURL example

```bash
curl -X POST http://localhost:8080/lecturer/tests/123/access \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "group": "A-01"
  }'
```

---

---

## Endpoint summary

* Method: `DELETE`
* URL: `/lecturer/tests/{test_id}/access`
* Auth: `Required`
* Description: Revokes access to a test for a specific group.

---

## Request

### Path params

| Name    | Type | Required | Description            |
| ------- | ---- | -------- | ---------------------- |
| test_id | int  | yes      | Unique test identifier |

---

### Request body

```json
{
  "group": "string"
}
```

---

## Response

### Success (HTTP 200)

```json
{
  "success": true
}
```

---

## Errors

* 400 — invalid test_id
* 400 — invalid body
* 400 — invalid input
* 401 — unauthorized
* 403 — access denied
* 500 — internal server error

---

## cURL example

```bash
curl -X DELETE http://localhost:8080/lecturer/tests/123/access \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "group": "A-01"
  }'
```
