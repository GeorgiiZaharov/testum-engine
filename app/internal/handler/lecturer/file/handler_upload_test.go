package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/lecturer/file/dto"
	"testum-engine/app/internal/handler/middleware"
	uploadtest "testum-engine/app/internal/service/use_case/lecturer/upload_test"
)

type uploadTestUCStub struct {
	executeFn func(ctx context.Context, req uploadtest.UploadTestRequest) (uploadtest.UploadTestResponse, error)
}

func (s *uploadTestUCStub) Execute(ctx context.Context, req uploadtest.UploadTestRequest) (uploadtest.UploadTestResponse, error) {
	return s.executeFn(ctx, req)
}

func createMultipartRequest(t *testing.T, fileContent []byte, ignore bool) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)

	_, err = part.Write(fileContent)
	require.NoError(t, err)

	require.NoError(t, writer.WriteField("ignore_validation", map[bool]string{true: "true", false: "false"}[ignore]))
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/lecturer/tests/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestUploadTest_Success(t *testing.T) {
	uc := &uploadTestUCStub{
		executeFn: func(ctx context.Context, req uploadtest.UploadTestRequest) (uploadtest.UploadTestResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, "test.txt", req.FileName)
			require.True(t, req.IgnoreValidation)

			return uploadtest.UploadTestResponse{
				Success: true,
				TestID:  ptr(42),
			}, nil
		},
	}

	h := New(uc, nil, nil)

	req := createMultipartRequest(t, []byte("hello world"), true)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.UploadTest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.UploadTestResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.Empty(t, resp.FormatErrors)
	assert.Empty(t, resp.ValidationErrors)
	assert.Equal(t, 42, *resp.TestID)
}

func TestUploadTest_Unauthorized(t *testing.T) {
	h := New(&uploadTestUCStub{}, nil, nil)

	req := createMultipartRequest(t, []byte("data"), false)

	rec := httptest.NewRecorder()

	h.UploadTest(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUploadTest_InvalidMultipart(t *testing.T) {
	h := New(&uploadTestUCStub{}, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/lecturer/tests/upload", bytes.NewReader([]byte("broken")))
	req.Header.Set("Content-Type", "multipart/form-data")

	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.UploadTest(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadTest_FileRequired(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("ignore_validation", "false")
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/lecturer/tests/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = middleware.WithUserID(req, 10)

	h := New(&uploadTestUCStub{}, nil, nil)

	rec := httptest.NewRecorder()

	h.UploadTest(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadTest_UseCaseErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "invalid input",
			err:        uploadtest.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "forbidden",
			err:        uploadtest.ErrAccessDenied,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "storage error",
			err:        uploadtest.ErrStorageFailed,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "unknown",
			err:        errors.New("unknown"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &uploadTestUCStub{
				executeFn: func(ctx context.Context, req uploadtest.UploadTestRequest) (uploadtest.UploadTestResponse, error) {
					return uploadtest.UploadTestResponse{}, tt.err
				},
			}

			h := New(uc, nil, nil)

			req := createMultipartRequest(t, []byte("data"), false)
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()

			h.UploadTest(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func ptr(i int) *int {
	return &i
}
