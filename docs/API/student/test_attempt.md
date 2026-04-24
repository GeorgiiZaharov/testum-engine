## Endpoint summary

* Method: `GET`
* URL: `/student/tests/{test_id}/hard`
* Auth: `Required`
* Description: Returns hard tasks for a student test attempt.

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
  "tasks": [
    {
      "text": "string",
      "image_url": "string",
      "is_hard": true,
      "answers": [
        {
          "text": "string",
          "image_url": "string"
        }
      ]
    }
  ]
}
```

---

## Errors

* 400 — test_id is required
* 400 — invalid test_id
* 401 — unauthorized
* 403 — forbidden / access denied
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/student/tests/123/hard \
  -H "Authorization: Bearer <token>"
```

---

## Endpoint summary

* Method: `GET`
* URL: `/student/tests/{test_id}/base`
* Auth: `Required`
* Description: Returns base tasks for a student test attempt.

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
  "tasks": [
    {
      "text": "string",
      "image_url": "string",
      "is_hard": false,
      "answers": [
        {
          "text": "string",
          "image_url": "string"
        }
      ]
    }
  ]
}
```

---

## Errors

* 400 — test_id is required
* 400 — invalid test_id
* 401 — unauthorized
* 403 — access denied
* 409 — test already completed
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/student/tests/123/base \
  -H "Authorization: Bearer <token>"
```

---

## Endpoint summary

* Method: `POST`
* URL: `/student/tests/{test_id}/hard`
* Auth: `Required`
* Description: Submits answers for hard tasks in a student test.

---

## Request

### Path params

| Name    | Type | Required | Description    |
| ------- | ---- | -------- | -------------- |
| test_id | int  | yes      | Unique test ID |

---

### Request body

```json
{
  "answers": [
    {
      "task_id": 0,
      "options": [
        0
      ]
    }
  ]
}
```

---

## Response

### Success (HTTP 200)

```json
{
  "success": true,
  "is_all_correct": true
}
```

---

## Errors

* 400 — test_id is required
* 400 — invalid test_id
* 400 — invalid request body
* 400 — answers are required
* 401 — unauthorized
* 403 — access denied
* 409 — hard answers already submitted
* 409 — test already finished
* 500 — internal server error

---

## cURL example

```bash
curl -X POST http://localhost:8080/student/tests/123/hard \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "answers": [
      {
        "task_id": 1,
        "options": [2, 3]
      }
    ]
  }'
```

---

## Endpoint summary

* Method: `POST`
* URL: `/student/tests/{test_id}/base`
* Auth: `Required`
* Description: Submits answers for base tasks in a student test.

---

## Request

### Path params

| Name    | Type | Required | Description    |
| ------- | ---- | -------- | -------------- |
| test_id | int  | yes      | Unique test ID |

---

### Request body

```json
{
  "answers": [
    {
      "task_id": 0,
      "options": [
        0
      ]
    }
  ]
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

* 400 — test_id is required
* 400 — invalid test_id
* 400 — invalid request body
* 400 — answers are required
* 401 — unauthorized
* 403 — access denied
* 409 — base answers already submitted
* 409 — hard block not passed
* 409 — test already finished
* 500 — internal server error

---

## cURL example

```bash
curl -X POST http://localhost:8080/student/tests/123/base \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "answers": [
      {
        "task_id": 1,
        "options": [0]
      }
    ]
  }'
```
