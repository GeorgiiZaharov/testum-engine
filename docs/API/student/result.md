## Endpoint summary

* Method: `GET`
* URL: `/student/test/{test_id}/result`
* Auth: `Required`
* Description: Returns detailed result of a student’s test attempt.

---

## Request

### Path params

| Name    | Type | Required | Description    |
| ------- | ---- | -------- | -------------- |
| test_id | int  | yes      | Unique test ID |

---

## Response

### Success (HTTP 200)

```json
{
  "mark": 0,
  "success_rate": 0,
  "date_start": "string",
  "date_end": "string"
}
```

---

## Errors

* 400 — test_id is required
* 400 — invalid test_id
* 400 — invalid input (use case validation error)
* 401 — unauthorized
* 403 — access denied
* 404 — result not found
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/student/test/123/result \
  -H "Authorization: Bearer <token>"
```
