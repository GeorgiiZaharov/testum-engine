## Endpoint summary

* Method: `GET`
* URL: `/lecturer/tests`
* Auth: `Required`
* Description: Returns list of tests created by the authenticated lecturer.

---

## 3. Response

### Success (HTTP 200)

```json
{
  "tests": [
    {
      "id": 0,
      "name": "string",
      "cnt_questions": 0,
      "cnt_hard_questions": 0,
      "groups": ["string"],
      "date_created": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

## Errors

* 401 — unauthorized
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/lecturer/tests \
  -H "Authorization: Bearer <token>"
```

---

## Endpoint summary

* Method: `GET`
* URL: `/lecturer/tests/{test_id}`
* Auth: `Required`
* Description: Returns full details of a specific test including tasks and groups.

---

## Request

### Path params

| Name    | Type | Required | Description    |
| ------- | ---- | -------- | -------------- |
| test_id | int  | yes      | Test unique ID |

---

## 3. Response

### Success (HTTP 200)

```json
{
  "id": 0,
  "name": "string",
  "cnt_questions": 0,
  "cnt_hard_questions": 0,
  "groups": [
    {
      "group_name": "string",
      "members_count": 0
    }
  ],
  "hard_tasks": [
    {
      "text": "string",
      "image": "string",
      "answers": [
        {
          "text": "string",
          "image": "string",
          "is_correct": true
        }
      ],
      "is_hard": true
    }
  ],
  "base_tasks": [
    {
      "text": "string",
      "image": "string",
      "answers": [
        {
          "text": "string",
          "image": "string",
          "is_correct": true
        }
      ],
      "is_hard": false
    }
  ]
}
```

---

## Errors

* 400 — invalid test_id
* 401 — unauthorized
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/lecturer/tests/123 \
  -H "Authorization: Bearer <token>"
```

---

## Endpoint summary

* Method: `GET`
* URL: `/lecturer/tests/{test_id}/groups`
* Auth: `Required`
* Description: Returns groups assigned to a test filtered by year.

---

## 2. Request

### Path params

| Name    | Type | Required | Description    |
| ------- | ---- | -------- | -------------- |
| test_id | int  | yes      | Test unique ID |

---

### Query params

| Name | Type | Required | Description (source: query) |
| ---- | ---- | -------- | --------------------------- |
| year | int  | yes      | Academic year filter        |

---


## Response

### Success (HTTP 200)

```json
{
  "groups": [
    {
      "group_name": "string",
      "members_count": 0
    }
  ]
}
```

---

## Errors

* 400 — invalid test_id
* 400 — year is required
* 400 — invalid year
* 401 — unauthorized
* 500 — internal server error

---

## cURL example

```bash
curl -X GET "http://localhost:8080/lecturer/tests/123/groups?year=0" \
  -H "Authorization: Bearer <token>"
```

---

## Endpoint summary

* Method: `GET`
* URL: `/lecturer/tests/{test_id}/result`
* Auth: `Required`
* Description: Returns students results for a test filtered by group and year.

---

##  Request

### Path params

| Name    | Type | Required | Description    |
| ------- | ---- | -------- | -------------- |
| test_id | int  | yes      | Test unique ID |

---

### Query params

| Name  | Type   | Required | Description (source: query) |
| ----- | ------ | -------- | --------------------------- |
| group | string | yes      | Student group               |
| year  | int    | yes      | Academic year               |

---

## Response

### Success (HTTP 200)

```json
{
  "results": [
    {
      "student_id": 0,
      "name": "string",
      "email": "string",
      "score": 0,
      "mark": 0
    }
  ]
}
```

---

## Errors

* 400 — invalid test_id
* 400 — group and year required
* 400 — invalid year
* 401 — unauthorized
* 500 — internal server error

---

## cURL example

```bash
curl -X GET "http://localhost:8080/lecturer/tests/123/result?group=A-01&year=0" \
  -H "Authorization: Bearer <token>"
```
