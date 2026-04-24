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
	posthardanswers "testum-engine/app/internal/service/use_case/student/post_hard_answers"
)

type postHardAnswersUCStub struct {
	executeFn func(
		ctx context.Context,
		userID int,
		testID int,
		answers []posthardanswers.TaskAnswer,
	) (posthardanswers.PostHardAnswersResponse, error)
}

func (s *postHardAnswersUCStub) Execute(
	ctx context.Context,
	userID int,
	testID int,
	answers []posthardanswers.TaskAnswer,
) (posthardanswers.PostHardAnswersResponse, error) {
	return s.executeFn(ctx, userID, testID, answers)
}

func TestPostHardAnswers_Success(t *testing.T) {
	uc := &postHardAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []posthardanswers.TaskAnswer,
		) (posthardanswers.PostHardAnswersResponse, error) {
			require.Equal(t, 10, userID)
			require.Equal(t, 15, testID)

			require.Len(t, answers, 1)
			require.Equal(t, 100, answers[0].TaskID)
			require.Equal(t, []int{1, 2}, answers[0].Options)

			return posthardanswers.PostHardAnswersResponse{
				IsAllCorrect: true,
			}, nil
		},
	}

	h := New(nil, nil, uc, nil)

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
		"/student/test/15/hard",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostHardAnswers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp handlerdto.PostHardAnswersResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.IsAllCorrect)
}

func TestPostHardAnswers_Unauthorized(t *testing.T) {
	called := false

	uc := &postHardAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []posthardanswers.TaskAnswer,
		) (posthardanswers.PostHardAnswersResponse, error) {
			called = true
			return posthardanswers.PostHardAnswersResponse{}, nil
		},
	}

	h := New(nil, nil, uc, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/hard",
		nil,
	)

	req.SetPathValue("test_id", "15")

	rec := httptest.NewRecorder()

	h.PostHardAnswers(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.False(t, called)
}

func TestPostHardAnswers_InvalidJSON(t *testing.T) {
	called := false

	uc := &postHardAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []posthardanswers.TaskAnswer,
		) (posthardanswers.PostHardAnswersResponse, error) {
			called = true
			return posthardanswers.PostHardAnswersResponse{}, nil
		},
	}

	h := New(nil, nil, uc, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/hard",
		bytes.NewBufferString("{invalid json"),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostHardAnswers(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

func TestPostHardAnswers_EmptyAnswers(t *testing.T) {
	called := false

	uc := &postHardAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []posthardanswers.TaskAnswer,
		) (posthardanswers.PostHardAnswersResponse, error) {
			called = true
			return posthardanswers.PostHardAnswersResponse{}, nil
		},
	}

	h := New(nil, nil, uc, nil)

	body := handlerdto.PostAnswersRequest{
		Answers: []handlerdto.TaskAnswerRequest{},
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/student/test/15/hard",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostHardAnswers(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.False(t, called)
}

func TestPostHardAnswers_AlreadySubmitted(t *testing.T) {
	uc := &postHardAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []posthardanswers.TaskAnswer,
		) (posthardanswers.PostHardAnswersResponse, error) {
			return posthardanswers.PostHardAnswersResponse{},
				posthardanswers.ErrAlreadySubmitted
		},
	}

	h := New(nil, nil, uc, nil)

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
		"/student/test/15/hard",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostHardAnswers(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestPostHardAnswers_InternalError(t *testing.T) {
	uc := &postHardAnswersUCStub{
		executeFn: func(
			ctx context.Context,
			userID int,
			testID int,
			answers []posthardanswers.TaskAnswer,
		) (posthardanswers.PostHardAnswersResponse, error) {
			return posthardanswers.PostHardAnswersResponse{},
				errors.New("unexpected failure")
		},
	}

	h := New(nil, nil, uc, nil)

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
		"/student/test/15/hard",
		bytes.NewReader(b),
	)

	req.SetPathValue("test_id", "15")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.PostHardAnswers(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
