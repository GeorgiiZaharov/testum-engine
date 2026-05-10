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

	"testum-engine/app/internal/handler/middleware"
	"testum-engine/app/internal/service/use_case/lecturer/upload_picture"
)

// =========================
// UC STUB
// =========================

type uploadPictureUCStub struct {
	executeFn func(ctx context.Context, req uploadpicture.UploadPictureRequest) (uploadpicture.UploadPictureResponse, error)
}

func (s *uploadPictureUCStub) Execute(ctx context.Context, req uploadpicture.UploadPictureRequest) (uploadpicture.UploadPictureResponse, error) {
	return s.executeFn(ctx, req)
}

// =========================
// HELPERS
// =========================

func createPictureMultipartRequest(t *testing.T, fileContent []byte, fileName string) *http.Request {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	require.NoError(t, err)
	_, err = part.Write(fileContent)
	require.NoError(t, err)

	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/lecturer/picture", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// =========================
// TESTS
// =========================

func TestUploadPicture_Success(t *testing.T) {
	uc := &uploadPictureUCStub{
		executeFn: func(ctx context.Context, req uploadpicture.UploadPictureRequest) (uploadpicture.UploadPictureResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, "avatar.png", req.FileName)
			assert.Equal(t, []byte("fake image"), req.File)

			return uploadpicture.UploadPictureResponse{
				Success: true,
				URL:     "http://localhost/lecturer10/avatar.png",
			}, nil
		},
	}

	h := New(nil, nil, uc)

	req := createPictureMultipartRequest(t, []byte("fake image"), "avatar.png")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()
	h.UploadPicture(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp uploadpicture.UploadPictureResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
	assert.Equal(t, "http://localhost/lecturer10/avatar.png", resp.URL)
}

func TestUploadPicture_Unauthorized(t *testing.T) {
	uc := &uploadPictureUCStub{}
	h := New(nil, nil, uc)

	req := createPictureMultipartRequest(t, []byte("data"), "avatar.png")
	rec := httptest.NewRecorder()

	h.UploadPicture(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUploadPicture_InvalidMultipart(t *testing.T) {
	uc := &uploadPictureUCStub{}
	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodPost, "/lecturer/picture", bytes.NewReader([]byte("broken")))
	req.Header.Set("Content-Type", "multipart/form-data")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()
	h.UploadPicture(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadPicture_FileRequired(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/lecturer/picture", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = middleware.WithUserID(req, 10)

	uc := &uploadPictureUCStub{}
	h := New(nil, nil, uc)

	rec := httptest.NewRecorder()
	h.UploadPicture(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUploadPicture_UseCaseErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"invalid input", uploadpicture.ErrInvalidInput, http.StatusBadRequest},
		{"forbidden", uploadpicture.ErrAccessDenied, http.StatusForbidden},
		{"storage error", uploadpicture.ErrStorageFailed, http.StatusInternalServerError},
		{"unknown", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &uploadPictureUCStub{
				executeFn: func(ctx context.Context, req uploadpicture.UploadPictureRequest) (uploadpicture.UploadPictureResponse, error) {
					return uploadpicture.UploadPictureResponse{}, tt.err
				},
			}

			h := New(nil, nil, uc)

			req := createPictureMultipartRequest(t, []byte("data"), "avatar.png")
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()
			h.UploadPicture(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
