package test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/middleware"
	getfile "testum-engine/app/internal/service/use_case/lecturer/get_test_file"
)

//
// STUB
//

type getFileUCStub struct {
	executeFn func(ctx context.Context, req getfile.GetTestFileRequest) (getfile.GetTestFileResponse, error)
}

func (s *getFileUCStub) Execute(ctx context.Context, req getfile.GetTestFileRequest) (getfile.GetTestFileResponse, error) {
	return s.executeFn(ctx, req)
}

//
// TESTS
//

func TestGetTestFile_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-file-*.txt")
	require.NoError(t, err)

	defer func() { _ = os.Remove(tmpFile.Name()) }()
	defer func() { _ = tmpFile.Close() }()

	content := []byte("hello file")
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	_, err = tmpFile.Seek(0, 0)
	require.NoError(t, err)

	uc := &getFileUCStub{
		executeFn: func(ctx context.Context, req getfile.GetTestFileRequest) (getfile.GetTestFileResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 99, req.TestID)

			return getfile.GetTestFileResponse{
				File: tmpFile,
			}, nil
		},
	}

	// =========================
	// 3. HANDLER
	// =========================
	h := New(nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/lecturer/test/99/file", nil)
	req.SetPathValue("test_id", "99")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	// =========================
	// 4. EXEC
	// =========================
	h.GetTestFile(rec, req)

	// =========================
	// 5. ASSERT
	// =========================
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "attachment; filename="+tmpFile.Name(), rec.Header().Get("Content-Disposition"))
	assert.Equal(t, "application/octet-stream", rec.Header().Get("Content-Type"))

	body, err := io.ReadAll(rec.Body)
	require.NoError(t, err)

	assert.Equal(t, "hello file", string(body))
}

func TestGetTestFile_Unauthorized(t *testing.T) {
	h := New(nil, &getFileUCStub{})

	req := httptest.NewRequest(http.MethodGet, "/lecturer/test/99/file", nil)

	rec := httptest.NewRecorder()

	h.GetTestFile(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetTestFile_InvalidID(t *testing.T) {
	h := New(nil, &getFileUCStub{})

	req := httptest.NewRequest(http.MethodGet, "/lecturer/test/xx/file", nil)
	req.SetPathValue("test_id", "xx")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestFile(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetTestFile_UseCaseErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"forbidden", getfile.ErrForbidden, http.StatusForbidden},
		{"not found", getfile.ErrFileNotFound, http.StatusNotFound},
		{"unknown", errors.New("x"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &getFileUCStub{
				executeFn: func(ctx context.Context, req getfile.GetTestFileRequest) (getfile.GetTestFileResponse, error) {
					return getfile.GetTestFileResponse{}, tt.err
				},
			}

			h := New(nil, uc)

			req := httptest.NewRequest(http.MethodGet, "/lecturer/test/99/file", nil)
			req.SetPathValue("test_id", "99")
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()

			h.GetTestFile(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestGetTestFile_FileStatError(t *testing.T) {
	// use real temp file
	tmp, err := os.CreateTemp("", "test")
	require.NoError(t, err)

	_, _ = tmp.Write([]byte("data"))
	_, _ = tmp.Seek(0, 0)

	// close file to simulate stat failure
	_ = tmp.Close()

	uc := &getFileUCStub{
		executeFn: func(ctx context.Context, req getfile.GetTestFileRequest) (getfile.GetTestFileResponse, error) {
			return getfile.GetTestFileResponse{
				File: tmp,
			}, nil
		},
	}

	h := New(nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/lecturer/test/99/file", nil)
	req.SetPathValue("test_id", "99")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestFile(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
