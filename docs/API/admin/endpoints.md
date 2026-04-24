## Endpoint summary

* Method: `POST`
* URL: `/admin/lecturers`
* Auth: `Required`
* Description: Creates a new lecturer using admin privileges.

---

## Request

### Request body

```json
{
  "login": "string"
}
```

---

## Response

### Success (HTTP 201)

```json
{
  "success": true
}
```

---

## Errors

* 400 — invalid request body
* 400 — login is required
* 400 — invalid input
* 401 — unauthorized
* 403 — forbidden
* 500 — internal server error

---

## cURL example

```bash
curl -X POST http://localhost:8080/admin/lecturers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "login": "john.doe"
  }'
```

---

## Endpoint summary

* Method: `DELETE`
* URL: `/admin/lecturers/{lecturer_id}`
* Auth: `Required`
* Description: Deletes a lecturer by ID.

---

## Request

### Path params

| Name        | Type | Required | Description        |
| ----------- | ---- | -------- | ------------------ |
| lecturer_id | int  | yes      | Lecturer unique ID |

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

* 400 — lecturer_id is required
* 400 — invalid lecturer_id
* 401 — unauthorized
* 403 — forbidden
* 404 — not found
* 409 — user is not lecturer
* 500 — internal server error

---

## cURL example

```bash
curl -X DELETE http://localhost:8080/admin/lecturers/123 \
  -H "Authorization: Bearer <token>"
```

---

## Endpoint summary

* Method: `GET`
* URL: `/admin/lecturers`
* Auth: `Required`
* Description: Returns list of all lecturers accessible to admin.

---

## Response

### Success (HTTP 200)

```json
{
  "lecturers": [
    {
      "id": 0,
      "login": "string",
      "mail": "string",
      "name": "string",
      "group": "string",
      "date_created": "string",
      "date_modified": "string"
    }
  ]
}
```

---

## Errors

* 400 — invalid input
* 401 — unauthorized
* 403 — forbidden
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/admin/lecturers \
  -H "Authorization: Bearer <token>"
```
