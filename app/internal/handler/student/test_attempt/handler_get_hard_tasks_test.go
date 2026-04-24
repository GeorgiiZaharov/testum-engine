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

	handlermiddleware "testum-engine/app/internal/handler/middleware"
	handlerdto "testum-engine/app/internal/handler/student/test_attempt/dto"
	gethardtasks "testum-engine/app/internal/service/use_case/student/get_hard_tasks"
)

type getHardTasksUCStub struct {
	executeFn func(
		ctx context.Context,
		req gethardtasks.GetHardTasksRequest,
	) (gethardtasks.GetHardTasksResponse, error)
}

func (s *getHardTasksUCStub) Execute(
	ctx context.Context,
	req gethardtasks.GetHardTasksRequest,
) (gethardtasks.GetHardTasksResponse, error) {
	return s.executeFn(ctx, req)
}

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func TestGetHardTasks_Success(t *testing.T) {
	imageURL := "https://example.com/image.png"

	uc := &getHardTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req gethardtasks.GetHardTasksRequest,
		) (gethardtasks.GetHardTasksResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 15, req.TestID)

			return gethardtasks.GetHardTasksResponse{
				HardTasks: []gethardtasks.Task{
					{
						Text:     "Solve equation",
						ImageURL: &imageURL,
						IsHard:   true,
						Answers: []gethardtasks.Answer{
							{
								Text: "Answer 1",
							},
							{
								Text: "Answer 2",
							},
						},
					},
				},
			}, nil
		},
	}

	h := New(
		uc,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/test/15/hard",
		nil,
	)

	req.SetPathValue("test_id", "15")
	req = handlermiddleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetHardTasks(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp handlerdto.GetTasksResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	require.Len(t, resp.Tasks, 1)

	task := resp.Tasks[0]

	assert.Equal(t, "Solve equation", task.Text)
	assert.True(t, task.IsHard)

	require.NotNil(t, task.ImageURL)
	assert.Equal(t, imageURL, *task.ImageURL)

	require.Len(t, task.Answers, 2)
	assert.Equal(t, "Answer 1", task.Answers[0].Text)
	assert.Equal(t, "Answer 2", task.Answers[1].Text)
}

func TestGetHardTasks_Unauthorized(t *testing.T) {
	called := false

	uc := &getHardTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req gethardtasks.GetHardTasksRequest,
		) (gethardtasks.GetHardTasksResponse, error) {
			called = true
			return gethardtasks.GetHardTasksResponse{}, nil
		},
	}

	h := New(
		uc,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/test/15/hard",
		nil,
	)

	req.SetPathValue("test_id", "15")

	rec := httptest.NewRecorder()

	h.GetHardTasks(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.False(t, called)

	var resp errorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "unauthorized", resp.Error)
}

func TestGetHardTasks_MissingTestID(t *testing.T) {
	called := false

	uc := &getHardTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req gethardtasks.GetHardTasksRequest,
		) (gethardtasks.GetHardTasksResponse, error) {
			called = true
			return gethardtasks.GetHardTasksResponse{}, nil
		},
	}

	h := New(
		uc,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/test/hard",
		nil,
	)

	req = handlermiddleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetHardTasks(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)

	var resp errorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "test_id is required", resp.Error)
}

func TestGetHardTasks_InvalidTestID(t *testing.T) {
	called := false

	uc := &getHardTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req gethardtasks.GetHardTasksRequest,
		) (gethardtasks.GetHardTasksResponse, error) {
			called = true
			return gethardtasks.GetHardTasksResponse{}, nil
		},
	}

	h := New(
		uc,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/test/abc/hard",
		nil,
	)

	req.SetPathValue("test_id", "abc")
	req = handlermiddleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetHardTasks(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)

	var resp errorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "invalid test_id", resp.Error)
}

func TestGetHardTasks_AccessDenied(t *testing.T) {
	uc := &getHardTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req gethardtasks.GetHardTasksRequest,
		) (gethardtasks.GetHardTasksResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 25, req.TestID)

			return gethardtasks.GetHardTasksResponse{},
				gethardtasks.ErrAccessDenied
		},
	}

	h := New(
		uc,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/test/25/hard",
		nil,
	)

	req.SetPathValue("test_id", "25")
	req = handlermiddleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetHardTasks(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)

	var resp errorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "access denied", resp.Error)
}

func TestGetHardTasks_InternalError(t *testing.T) {
	uc := &getHardTasksUCStub{
		executeFn: func(
			ctx context.Context,
			req gethardtasks.GetHardTasksRequest,
		) (gethardtasks.GetHardTasksResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 99, req.TestID)

			return gethardtasks.GetHardTasksResponse{},
				errors.New("unexpected failure")
		},
	}

	h := New(
		uc,
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest(
		http.MethodGet,
		"/student/test/99/hard",
		nil,
	)

	req.SetPathValue("test_id", "99")
	req = handlermiddleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetHardTasks(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp errorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "internal server error", resp.Error)
}
