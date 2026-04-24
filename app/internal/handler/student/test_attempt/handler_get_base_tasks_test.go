package testattempt

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
	handlerdto "testum-engine/app/internal/handler/student/test_attempt/dto"
	getbasetasks "testum-engine/app/internal/service/use_case/student/get_base_tasks"
)

type getBaseTasksUCStub struct {
	executeFn func(
		ctx context.Context,
		req getbasetasks.GetBaseTasksRequest,
	) (getbasetasks.GetBaseTasksResponse, error)
}

func (s *getBaseTasksUCStub) Execute(
	ctx context.Context,
	req getbasetasks.GetBaseTasksRequest,
) (getbasetasks.GetBaseTasksResponse, error) {
	return s.executeFn(ctx, req)
}

type getBaseTasksSuccessResponse handlerdto.GetTasksResponse

type getBaseTasksErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func TestGetBaseTasks_Success(t *testing.T) {
	uc := &getBaseTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req getbasetasks.GetBaseTasksRequest,
		) (getbasetasks.GetBaseTasksResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 15, req.TestID)

			img := "https://img.com/1.png"

			return getbasetasks.GetBaseTasksResponse{
				BaseTasks: []getbasetasks.Task{
					{
						Text:   "What is 2+2?",
						IsHard: false,
						Answers: []getbasetasks.Answer{
							{Text: "3"},
							{Text: "4"},
						},
						ImageURL: &img,
					},
				},
			}, nil
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/tests/15/base",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetBaseTasks(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp getBaseTasksSuccessResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	require.Len(t, resp.Tasks, 1)

	task := resp.Tasks[0]
	assert.Equal(t, "What is 2+2?", task.Text)
	assert.False(t, task.IsHard)

	require.Len(t, task.Answers, 2)
	assert.Equal(t, "3", task.Answers[0].Text)
	assert.Equal(t, "4", task.Answers[1].Text)

	require.NotNil(t, task.ImageURL)
	assert.Equal(t, "https://img.com/1.png", *task.ImageURL)
}

func TestGetBaseTasks_Unauthorized(t *testing.T) {
	called := false

	uc := &getBaseTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req getbasetasks.GetBaseTasksRequest,
		) (getbasetasks.GetBaseTasksResponse, error) {
			called = true
			return getbasetasks.GetBaseTasksResponse{}, nil
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/tests/15/base",
		nil,
	)

	req.SetPathValue("test_id", "15")

	rec := httptest.NewRecorder()

	h.GetBaseTasks(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.False(t, called)

	var resp getBaseTasksErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "unauthorized", resp.Error)
}

func TestGetBaseTasks_MissingTestID(t *testing.T) {
	called := false

	uc := &getBaseTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req getbasetasks.GetBaseTasksRequest,
		) (getbasetasks.GetBaseTasksResponse, error) {
			called = true
			return getbasetasks.GetBaseTasksResponse{}, nil
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/tests/base",
		nil,
	)

	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetBaseTasks(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)

	var resp getBaseTasksErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "test_id is required", resp.Error)
}

func TestGetBaseTasks_InvalidTestID(t *testing.T) {
	called := false

	uc := &getBaseTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req getbasetasks.GetBaseTasksRequest,
		) (getbasetasks.GetBaseTasksResponse, error) {
			called = true
			return getbasetasks.GetBaseTasksResponse{}, nil
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/tests/abc/base",
		nil,
	)

	req.SetPathValue("test_id", "abc")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetBaseTasks(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)

	var resp getBaseTasksErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "invalid test_id", resp.Error)
}

func TestGetBaseTasks_AlreadyCompleted(t *testing.T) {
	uc := &getBaseTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req getbasetasks.GetBaseTasksRequest,
		) (getbasetasks.GetBaseTasksResponse, error) {
			return getbasetasks.GetBaseTasksResponse{},
				getbasetasks.ErrTestCompleted
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/tests/15/base",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetBaseTasks(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var resp getBaseTasksErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "test already completed", resp.Error)
}

func TestGetBaseTasks_AccessDenied(t *testing.T) {
	uc := &getBaseTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req getbasetasks.GetBaseTasksRequest,
		) (getbasetasks.GetBaseTasksResponse, error) {
			return getbasetasks.GetBaseTasksResponse{},
				getbasetasks.ErrAccessDenied
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/tests/15/base",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetBaseTasks(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var resp getBaseTasksErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "access denied", resp.Error)
}

func TestGetBaseTasks_InternalError(t *testing.T) {
	uc := &getBaseTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req getbasetasks.GetBaseTasksRequest,
		) (getbasetasks.GetBaseTasksResponse, error) {
			return getbasetasks.GetBaseTasksResponse{},
				errors.New("unexpected failure")
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/tests/15/base",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetBaseTasks(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp getBaseTasksErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "internal server error", resp.Error)
}
