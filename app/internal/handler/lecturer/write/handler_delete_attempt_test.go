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
	deleteattempt "testum-engine/app/internal/service/use_case/lecturer/delete_attempt"
)

// -------------------- STUB --------------------

type deleteAttemptUCStub struct {
	executeFn func(ctx context.Context, req deleteattempt.DeleteAttemptRequest) (deleteattempt.DeleteAttemptResponse, error)
}

func (s *deleteAttemptUCStub) Execute(ctx context.Context, req deleteattempt.DeleteAttemptRequest) (deleteattempt.DeleteAttemptResponse, error) {
	return s.executeFn(ctx, req)
}

// -------------------- TESTS --------------------

func TestDeleteAttempt_Success(t *testing.T) {
	uc := &deleteAttemptUCStub{
		executeFn: func(ctx context.Context, req deleteattempt.DeleteAttemptRequest) (deleteattempt.DeleteAttemptResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 42, req.TestID)
			return deleteattempt.DeleteAttemptResponse{Success: true}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/lecturer/attempt/test/42/user/10",
		nil,
	)
	req.SetPathValue("test_id", "42")
	req.SetPathValue("user_id", "10")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.DeleteAttempt(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.DeleteAttemptResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestDeleteAttempt_Unauthorized(t *testing.T) {
	h := New(nil, nil, nil, &deleteAttemptUCStub{})

	req := httptest.NewRequest(
		http.MethodDelete,
		"/lecturer/attempt/test/42/user/10",
		nil,
	)
	rec := httptest.NewRecorder()

	h.DeleteAttempt(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteAttempt_InvalidIDs(t *testing.T) {
	h := New(nil, nil, nil, &deleteAttemptUCStub{})

	tests := []struct {
		name     string
		testID   string
		userID   string
		wantCode int
	}{
		{"invalid test_id", "abc", "10", http.StatusBadRequest},
		{"invalid user_id", "42", "xyz", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/dummy", nil)
			req.SetPathValue("test_id", tt.testID)
			req.SetPathValue("user_id", tt.userID)
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()
			h.DeleteAttempt(rec, req)

			assert.Equal(t, tt.wantCode, rec.Code)
		})
	}
}

func TestDeleteAttempt_ErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"invalid input", deleteattempt.ErrInvalidInput, http.StatusBadRequest},
		{"access denied", deleteattempt.ErrAccessDenied, http.StatusForbidden},
		{"delete failed", deleteattempt.ErrDeleteFailed, http.StatusInternalServerError},
		{"unknown error", errors.New("unexpected"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &deleteAttemptUCStub{
				executeFn: func(ctx context.Context, req deleteattempt.DeleteAttemptRequest) (deleteattempt.DeleteAttemptResponse, error) {
					return deleteattempt.DeleteAttemptResponse{}, tt.err
				},
			}

			h := New(nil, nil, nil, uc)

			req := httptest.NewRequest(http.MethodDelete, "/dummy", nil)
			req.SetPathValue("test_id", "42")
			req.SetPathValue("user_id", "10")
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()
			h.DeleteAttempt(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
