package gettests

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
	handlerdto "testum-engine/app/internal/handler/student/get_tests/dto"
	getfinishedtest "testum-engine/app/internal/service/use_case/student/get_finished_test"
)

type getFinishedTestUCStub struct {
	executeFn func(
		ctx context.Context,
		req getfinishedtest.GetFinishedTestRequest,
	) (getfinishedtest.GetFinishedTestResponse, error)
}

func (s *getFinishedTestUCStub) Execute(
	ctx context.Context,
	req getfinishedtest.GetFinishedTestRequest,
) (getfinishedtest.GetFinishedTestResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGetFinishedTests_Success(t *testing.T) {
	start := time.Date(2026, 4, 24, 10, 0, 0, 0, time.UTC)
	end := time.Date(2026, 4, 24, 12, 0, 0, 0, time.UTC)

	uc := &getFinishedTestUCStub{
		executeFn: func(ctx context.Context, req getfinishedtest.GetFinishedTestRequest) (getfinishedtest.GetFinishedTestResponse, error) {
			require.Equal(t, 10, req.UserID)

			return getfinishedtest.GetFinishedTestResponse{
				FinishedTests: []getfinishedtest.StudentFinishTest{
					{
						ID:           1,
						Name:         "test",
						LecturerName: "lecturer",
						Mark:         5,
						SuccessRate:  95,
						DateStart:    start,
						DateEnd:      end,
					},
				},
			}, nil
		},
	}

	h := New(nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/finished", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetFinishedTests(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []handlerdto.StudentFinishTest
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Len(t, resp, 1)
	assert.Equal(t, 1, resp[0].ID)
	assert.Equal(t, 5, resp[0].Mark)
	assert.Equal(t, 95.0, resp[0].SuccessRate)
}

func TestGetFinishedTests_Unauthorized(t *testing.T) {
	uc := &getFinishedTestUCStub{
		executeFn: func(ctx context.Context, req getfinishedtest.GetFinishedTestRequest) (getfinishedtest.GetFinishedTestResponse, error) {
			t.Fatal("should not be called")
			return getfinishedtest.GetFinishedTestResponse{}, nil
		},
	}

	h := New(nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/finished", nil)
	rec := httptest.NewRecorder()

	h.GetFinishedTests(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetFinishedTests_InvalidInput(t *testing.T) {
	uc := &getFinishedTestUCStub{
		executeFn: func(ctx context.Context, req getfinishedtest.GetFinishedTestRequest) (getfinishedtest.GetFinishedTestResponse, error) {
			return getfinishedtest.GetFinishedTestResponse{}, getfinishedtest.ErrInvalidInput
		},
	}

	h := New(nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/finished", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetFinishedTests(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetFinishedTests_InternalError(t *testing.T) {
	uc := &getFinishedTestUCStub{
		executeFn: func(ctx context.Context, req getfinishedtest.GetFinishedTestRequest) (getfinishedtest.GetFinishedTestResponse, error) {
			return getfinishedtest.GetFinishedTestResponse{}, errors.New("boom")
		},
	}

	h := New(nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/student/tests/finished", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetFinishedTests(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
