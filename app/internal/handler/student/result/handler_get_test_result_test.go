package result

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/middleware"
	handlerdto "testum-engine/app/internal/handler/student/result/dto"
	gettestresult "testum-engine/app/internal/service/use_case/student/get_test_result"
)

type getTestResultUCStub struct {
	executeFn func(
		ctx context.Context,
		req gettestresult.GetTestResultRequest,
	) (gettestresult.GetTestResultResponse, error)
}

func (s *getTestResultUCStub) Execute(
	ctx context.Context,
	req gettestresult.GetTestResultRequest,
) (gettestresult.GetTestResultResponse, error) {
	return s.executeFn(ctx, req)
}

/*
-------------------------
SUCCESS
-------------------------
*/

func TestGetTestResult_Success(t *testing.T) {
	mark := 5
	successRate := 92.5

	start := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 24, 14, 30, 0, 0, time.UTC)

	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 15, req.TestID)

			return gettestresult.GetTestResultResponse{
				Mark:        &mark,
				SuccessRate: &successRate,
				DateStart:   start,
				DateEnd:     &end,
			}, nil
		},
	}

	h := New(uc)

	req := httptest.NewRequest(http.MethodGet, "/student/test/15/result", nil)
	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp handlerdto.GetTestResultResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	require.NotNil(t, resp.Mark)
	assert.Equal(t, 5, *resp.Mark)

	require.NotNil(t, resp.SuccessRate)
	assert.Equal(t, 92.5, *resp.SuccessRate)
}

/*
-------------------------
AUTH
-------------------------
*/

func TestGetTestResult_Unauthorized(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			t.Fatal("must not be called")
			return gettestresult.GetTestResultResponse{}, nil
		},
	}

	h := New(uc)

	req := httptest.NewRequest(http.MethodGet, "/student/test/15/result", nil)
	req.SetPathValue("test_id", "15")

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	var errResp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Equal(t, "unauthorized", errResp["error"])
}

/*
-------------------------
VALIDATION
-------------------------
*/

func TestGetTestResult_MissingTestID(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			t.Fatal("must not be called")
			return gettestresult.GetTestResultResponse{}, nil
		},
	}

	h := New(uc)

	req := httptest.NewRequest(http.MethodGet, "/student/test/result", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Equal(t, "test_id is required", errResp["error"])
}

func TestGetTestResult_InvalidTestID(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			t.Fatal("must not be called")
			return gettestresult.GetTestResultResponse{}, nil
		},
	}

	h := New(uc)

	req := httptest.NewRequest(http.MethodGet, "/student/test/abc/result", nil)
	req.SetPathValue("test_id", "abc")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var errResp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Equal(t, "invalid test_id", errResp["error"])
}

/*
-------------------------
USE CASE ERRORS
-------------------------
*/

func TestGetTestResult_AccessDenied(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			return gettestresult.GetTestResultResponse{}, gettestresult.ErrAccessDenied
		},
	}

	h := New(uc)

	req := httptest.NewRequest(http.MethodGet, "/student/test/7/result", nil)
	req.SetPathValue("test_id", "7")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var errResp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Equal(t, "access denied", errResp["error"])
}

func TestGetTestResult_NotFound(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			return gettestresult.GetTestResultResponse{}, gettestresult.ErrResultNotFound
		},
	}

	h := New(uc)

	req := httptest.NewRequest(http.MethodGet, "/student/test/100/result", nil)
	req.SetPathValue("test_id", "100")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var errResp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Equal(t, "result not found", errResp["error"])
}

func TestGetTestResult_InternalError(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			return gettestresult.GetTestResultResponse{}, errors.New("unexpected failure")
		},
	}

	h := New(uc)

	req := httptest.NewRequest(http.MethodGet, "/student/test/55/result", nil)
	req.SetPathValue("test_id", "55")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var errResp map[string]string
	err := json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)

	assert.Equal(t, "internal server error", errResp["error"])
}
