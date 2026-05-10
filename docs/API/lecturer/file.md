## Endpoint summary

* Method: `POST`
* URL: `/lecturer/tests/upload`
* Auth: `Required`
* Description: Uploads a test file for lecturer with optional validation step.

---

## Request

### Request body

`multipart/form-data`

| Field             | Type | Required | Description              |
| ----------------- | ---- | -------- | ------------------------ |
| file              | file | yes      | Test file content        |
| ignore_validation | bool | no       | Skips validation if true |

---

## Response

### Success (HTTP 200)

```json
{
  "format_errors": [
    {
      "error": "string"
    }
  ],
  "validation_errors": [
    {
      "line": 0,
      "error": "string"
    }
  ],
  "test_id": 0,
  "success": true
}
```

---

## Errors

* 400 — invalid multipart data
* 400 — file is required
* 400 — invalid file
* 401 — unauthorized
* 403 — access denied
* 500 — internal server error

---

## cURL example

```bash
curl -X POST http://localhost:8080/lecturer/tests \
  -H "Authorization: Bearer <token>" \
  -F "file=@test.txt" \
  -F "ignore_validation=true"
```

---

## Endpoint summary

* Method: `GET`
* URL: `/lecturer/tests/{test_id}/file`
* Auth: `Required`
* Description: Downloads test file by test ID for lecturer.

---

## Request

### Path params

| Name    | Type | Required | Description    |
| ------- | ---- | -------- | -------------- |
| test_id | int  | yes      | Test unique ID |

---

## Response

### Success (HTTP 200)

Binary file stream:

* Content-Type: `application/octet-stream`
* Content-Disposition: `attachment`

---

## Errors

* 400 — invalid test_id
* 401 — unauthorized
* 403 — forbidden
* 404 — file not found
* 500 — internal server error

---

## cURL example

```bash
curl -X GET http://localhost:8080/lecturer/tests/123/file \
  -H "Authorization: Bearer <token>" \
  --output test_file
```

## Endpoint summary

*Method: POST`
* URL: `/lecturer/picture`*
**Auth:** Required  
**Description:** Uploads a profile picture for a lecturer.

### Request

`multipart/form-data`

| Field   | Type | Required | Description       |
| ------- | ---- | -------- | ---------------- |
| picture | file | yes      | Image file content |

### Response

**Success (HTTP 200)**

```json
{
  "url": "http://localhost/lecturer1/picture.png",
  "success": true
}
````

### Errors

* 400 — invalid multipart data
* 400 — file is required
* 400 — invalid file
* 401 — unauthorized
* 403 — access denied
* 500 — internal server error

### cURL Example

```bash
curl -X POST http://localhost:8080/lecturer/picture \
  -H "Authorization: Bearer <token>" \
  -F "picture=@profile.png"
```
