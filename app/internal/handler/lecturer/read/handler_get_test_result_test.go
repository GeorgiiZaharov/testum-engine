package read

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	handlerdto "testum-engine/app/internal/handler/lecturer/read/dto"
	"testum-engine/app/internal/handler/middleware"
	gettestresult "testum-engine/app/internal/service/use_case/lecturer/get_test_result"
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
----------------------------------------------------
SUCCESS
----------------------------------------------------
*/

func TestGetTestResult_Success(t *testing.T) {
	group := "A-01"

	score := 87.5
	mark := 5

	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 15, req.TestID)
			require.Equal(t, "A-01", req.Group)
			require.Equal(t, 2026, req.Year)

			return gettestresult.GetTestResultResponse{
				Results: []gettestresult.StudentResult{
					{
						UserID:      1,
						Name:        "John",
						Mail:        "john@mail.com",
						SuccessRate: &score,
						Mark:        &mark,
					},
				},
			}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/15/result?group="+group+"&year=2026",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp handlerdto.GetTestResultResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	require.Len(t, resp.Results, 1)

	assert.Equal(t, 1, resp.Results[0].StudentID)
	assert.Equal(t, "John", resp.Results[0].Name)
	assert.Equal(t, "john@mail.com", resp.Results[0].Email)

	require.NotNil(t, resp.Results[0].Score)
	assert.Equal(t, 87.5, *resp.Results[0].Score)

	require.NotNil(t, resp.Results[0].Mark)
	assert.Equal(t, 5, *resp.Results[0].Mark)
}

/*
----------------------------------------------------
AUTH
----------------------------------------------------
*/

func TestGetTestResult_Unauthorized(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			t.Fatal("use case must not be called")
			return gettestresult.GetTestResultResponse{}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/15/result?group=A-01&year=2026",
		nil,
	)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

/*
----------------------------------------------------
VALIDATION
----------------------------------------------------
*/

func TestGetTestResult_InvalidTestID(t *testing.T) {
	called := false

	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			called = true
			return gettestresult.GetTestResultResponse{}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/abc/result?group=A-01&year=2026",
		nil,
	)

	req.SetPathValue("test_id", "abc")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

func TestGetTestResult_MissingGroupOrYear(t *testing.T) {
	called := false

	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			called = true
			return gettestresult.GetTestResultResponse{}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/15/result",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

func TestGetTestResult_InvalidYear(t *testing.T) {
	called := false

	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			called = true
			return gettestresult.GetTestResultResponse{}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/15/result?group=A-01&year=abc",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

/*
----------------------------------------------------
USE CASE ERROR
----------------------------------------------------
*/

func TestGetTestResult_InternalError(t *testing.T) {
	uc := &getTestResultUCStub{
		executeFn: func(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error) {
			return gettestresult.GetTestResultResponse{}, errors.New("boom")
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/15/result?group=A-01&year=2026",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTestResult(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
