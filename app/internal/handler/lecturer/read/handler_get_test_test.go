package read

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/lecturer/read/dto"
	"testum-engine/app/internal/handler/middleware"
	gettest "testum-engine/app/internal/service/use_case/lecturer/get_test"
)

type getTestUCStub struct {
	executeFn func(ctx context.Context, req gettest.GetTestRequest) (gettest.GetTestResponse, error)
}

func (s *getTestUCStub) Execute(
	ctx context.Context,
	req gettest.GetTestRequest,
) (gettest.GetTestResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGetTest_Success(t *testing.T) {
	uc := &getTestUCStub{
		executeFn: func(
			ctx context.Context,
			req gettest.GetTestRequest,
		) (gettest.GetTestResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 42, req.TestID)

			return gettest.GetTestResponse{
				ID:               42,
				Name:             "Physics",
				CntQuestions:     10,
				CntHardQuestions: 3,
			}, nil
		},
	}

	h := New(nil, uc, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/42",
		nil,
	)
	req.SetPathValue("test_id", "42")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTest(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.GetTestResponse

	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, 42, resp.ID)
}

func TestGetTest_InvalidTestID(t *testing.T) {
	h := New(nil, &getTestUCStub{}, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests/abc",
		nil,
	)
	req.SetPathValue("test_id", "abc")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTest(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
