package testattempt

import (
	"bytes"
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
	postbaseanswers "testum-engine/app/internal/service/use_case/student/post_base_answers"
)

type postBaseAnswersUCStub struct {
	executeFn func(
		ctx context.Context,
		userID int,
		testID int,
		answers []postbaseanswers.TaskAnswer,
	) (postbaseanswers.PostBaseAnswersResponse, error)
}

func (s *postBaseAnswersUCStub) Execute(
	ctx context.Context,
	userID int,
	testID int,
	answers []postbaseanswers.TaskAnswer,
) (postbaseanswers.PostBaseAnswersResponse, error) {
	return s.executeFn(ctx, userID, testID, answers)
}

type postBaseAnswersResponse struct {
	Success bool `json:"success"`
}

func TestPostBaseAnswers_Success(t *testing.T) {
	uc := &postBaseAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []postbaseanswers.TaskAnswer,
		) (postbaseanswers.PostBaseAnswersResponse, error) {
			require.Equal(t, 10, userID)
			require.Equal(t, 15, testID)

			require.Len(t, answers, 1)
			require.Equal(t, 100, answers[0].TaskID)
			require.Equal(t, []int{1, 2}, answers[0].Options)

			return postbaseanswers.PostBaseAnswersResponse{
				Success: true,
			}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	body := handlerdto.PostAnswersRequest{
		Answers: []handlerdto.TaskAnswerRequest{
			{
				TaskID:  100,
				Options: []int{1, 2},
			},
		},
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/base",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostBaseAnswers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp postBaseAnswersResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
}

func TestPostBaseAnswers_Unauthorized(t *testing.T) {
	called := false

	uc := &postBaseAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []postbaseanswers.TaskAnswer,
		) (postbaseanswers.PostBaseAnswersResponse, error) {
			called = true
			return postbaseanswers.PostBaseAnswersResponse{}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/base",
		nil,
	)

	req.SetPathValue("test_id", "15")

	rec := httptest.NewRecorder()

	h.PostBaseAnswers(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.False(t, called)
}

func TestPostBaseAnswers_InvalidJSON(t *testing.T) {
	called := false

	uc := &postBaseAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []postbaseanswers.TaskAnswer,
		) (postbaseanswers.PostBaseAnswersResponse, error) {
			called = true
			return postbaseanswers.PostBaseAnswersResponse{}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/base",
		bytes.NewBufferString("{invalid json"),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostBaseAnswers(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

func TestPostBaseAnswers_EmptyAnswers(t *testing.T) {
	called := false

	uc := &postBaseAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []postbaseanswers.TaskAnswer,
		) (postbaseanswers.PostBaseAnswersResponse, error) {
			called = true
			return postbaseanswers.PostBaseAnswersResponse{}, nil
		},
	}

	h := New(nil, nil, nil, uc)

	body := handlerdto.PostAnswersRequest{
		Answers: []handlerdto.TaskAnswerRequest{},
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/base",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostBaseAnswers(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

func TestPostBaseAnswers_HardBlockNotPassed(t *testing.T) {
	uc := &postBaseAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []postbaseanswers.TaskAnswer,
		) (postbaseanswers.PostBaseAnswersResponse, error) {
			return postbaseanswers.PostBaseAnswersResponse{},
				postbaseanswers.ErrHardBlockNotPassed
		},
	}

	h := New(nil, nil, nil, uc)

	body := handlerdto.PostAnswersRequest{
		Answers: []handlerdto.TaskAnswerRequest{
			{
				TaskID:  1,
				Options: []int{1},
			},
		},
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/base",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostBaseAnswers(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)

	var resp errorResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "hard block not passed", resp.Error)
}

func TestPostBaseAnswers_InternalError(t *testing.T) {
	uc := &postBaseAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []postbaseanswers.TaskAnswer,
		) (postbaseanswers.PostBaseAnswersResponse, error) {
			return postbaseanswers.PostBaseAnswersResponse{},
				errors.New("unexpected failure")
		},
	}

	h := New(nil, nil, nil, uc)

	body := handlerdto.PostAnswersRequest{
		Answers: []handlerdto.TaskAnswerRequest{
			{
				TaskID:  1,
				Options: []int{1},
			},
		},
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/base",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostBaseAnswers(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp errorResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.False(t, resp.Success)
	assert.Equal(t, "internal server error", resp.Error)
}
