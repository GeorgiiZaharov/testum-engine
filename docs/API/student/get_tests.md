## Endpoint summary

* Method: `GET`
* URL: `/student/tests`
* Auth: `Required`
* Description: Returns active tests assigned to the authenticated student.

---

## Response

### Success (HTTP 200)

```json
{
  "active_tests": [
    {
      "id": 0,
      "name": "string",
      "lecturer_name": "string",
      "cnt_questions": 0,
      "cnt_hard_questions": 0,
      "date_start": "string"
    }
  ]
}
```

---

## Errors

* 401 — unauthorized
* 400 — invalid input
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/student/tests \
  -H "Authorization: Bearer <token>"
```

---

---

## Endpoint summary

* Method: `GET`
* URL: `/student/tests/finished`
* Auth: `Required`
* Description: Returns finished tests for the authenticated student.

---

## Response

### Success (HTTP 200)

```json
{
  "finished_tests": [
    {
      "id": 0,
      "name": "string",
      "lecturer_name": "string",
      "cnt_questions": 0,
      "cnt_hard_questions": 0,
      "mark": 0,
      "success_rate": 0,
      "date_start": "string",
      "date_end": "string"
    }
  ]
}
```

---

## Errors

* 401 — unauthorized
* 400 — invalid input
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/student/tests/finished \
  -H "Authorization: Bearer <token>"
```
