package gettests

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/middleware"
	handlerdto "testum-engine/app/internal/handler/student/get_tests/dto"
	getactivetest "testum-engine/app/internal/service/use_case/student/get_active_test"
)

type getActiveTestUCStub struct {
	executeFn func(
		ctx context.Context,
		req getactivetest.GetActiveTestRequest,
	) (getactivetest.GetActiveTestResponse, error)
}

func (s *getActiveTestUCStub) Execute(
	ctx context.Context,
	req getactivetest.GetActiveTestRequest,
) (getactivetest.GetActiveTestResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGetActiveTests_Success(t *testing.T) {
	uc := &getActiveTestUCStub{
		executeFn: func(ctx context.Context, req getactivetest.GetActiveTestRequest) (getactivetest.GetActiveTestResponse, error) {
			require.Equal(t, 10, req.UserID)

			return getactivetest.GetActiveTestResponse{
				ActiveTests: []getactivetest.StudentActiveTest{
					{
						ID:               1,
						Name:             "test",
						LecturerName:     "lecturer",
						CntQuestions:     10,
						CntHardQuestions: 5,
					},
				},
			}, nil
		},
	}

	h := New(uc, nil)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/active", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetActiveTests(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []handlerdto.StudentActiveTest
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Len(t, resp, 1)
	assert.Equal(t, 1, resp[0].ID)
	assert.Equal(t, "test", resp[0].Name)
}

func TestGetActiveTests_Unauthorized(t *testing.T) {
	uc := &getActiveTestUCStub{
		executeFn: func(ctx context.Context, req getactivetest.GetActiveTestRequest) (getactivetest.GetActiveTestResponse, error) {
			t.Fatal("should not be called")
			return getactivetest.GetActiveTestResponse{}, nil
		},
	}

	h := New(uc, nil)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/active", nil)
	rec := httptest.NewRecorder()

	h.GetActiveTests(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetActiveTests_InvalidInput(t *testing.T) {
	uc := &getActiveTestUCStub{
		executeFn: func(ctx context.Context, req getactivetest.GetActiveTestRequest) (getactivetest.GetActiveTestResponse, error) {
			return getactivetest.GetActiveTestResponse{}, getactivetest.ErrInvalidInput
		},
	}

	h := New(uc, nil)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/active", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetActiveTests(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetActiveTests_InternalError(t *testing.T) {
	uc := &getActiveTestUCStub{
		executeFn: func(ctx context.Context, req getactivetest.GetActiveTestRequest) (getactivetest.GetActiveTestResponse, error) {
			return getactivetest.GetActiveTestResponse{}, errors.New("boom")
		},
	}

	h := New(uc, nil)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/active", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetActiveTests(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
