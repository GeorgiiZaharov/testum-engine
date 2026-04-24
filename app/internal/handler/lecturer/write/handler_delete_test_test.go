package write

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/lecturer/write/dto"
	"testum-engine/app/internal/handler/middleware"
	deletetest "testum-engine/app/internal/service/use_case/lecturer/delete_test"
)

type deleteTestUCStub struct {
	executeFn func(
		ctx context.Context,
		req deletetest.DeleteTestRequest,
	) (deletetest.DeleteTestResponse, error)
}

func (s *deleteTestUCStub) Execute(
	ctx context.Context,
	req deletetest.DeleteTestRequest,
) (deletetest.DeleteTestResponse, error) {
	return s.executeFn(ctx, req)
}

func TestDeleteTest_Success(t *testing.T) {
	uc := &deleteTestUCStub{
		executeFn: func(
			ctx context.Context,
			req deletetest.DeleteTestRequest,
		) (deletetest.DeleteTestResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 42, req.TestID)

			return deletetest.DeleteTestResponse{
				Success: true,
			}, nil
		},
	}

	h := New(uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/lecturer/tests/42",
		nil,
	)
	req.SetPathValue("test_id", "42")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.DeleteTest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.DeleteTestResponse

	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
}

func TestDeleteTest_Unauthorized(t *testing.T) {
	h := New(&deleteTestUCStub{}, nil, nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/lecturer/tests/42",
		nil,
	)

	rec := httptest.NewRecorder()

	h.DeleteTest(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteTest_InvalidTestID(t *testing.T) {
	h := New(&deleteTestUCStub{}, nil, nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/lecturer/tests/abc",
		nil,
	)
	req.SetPathValue("test_id", "abc")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.DeleteTest(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteTest_ErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "invalid input",
			err:        deletetest.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "forbidden",
			err:        deletetest.ErrAccessDenied,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "not found",
			err:        deletetest.ErrTestNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "unknown",
			err:        errors.New("unexpected"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &deleteTestUCStub{
				executeFn: func(
					ctx context.Context,
					req deletetest.DeleteTestRequest,
				) (deletetest.DeleteTestResponse, error) {
					return deletetest.DeleteTestResponse{}, tt.err
				},
			}

			h := New(uc, nil, nil)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/lecturer/tests/42",
				nil,
			)
			req.SetPathValue("test_id", "42")
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()

			h.DeleteTest(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
