package read

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/lecturer/read/dto"
	"testum-engine/app/internal/handler/middleware"
	gettests "testum-engine/app/internal/service/use_case/lecturer/get_tests"
)

type getTestsUCStub struct {
	executeFn func(ctx context.Context, req gettests.GetTestsRequest) (gettests.GetTestsResponse, error)
}

func (s *getTestsUCStub) Execute(
	ctx context.Context,
	req gettests.GetTestsRequest,
) (gettests.GetTestsResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGetTests_Success(t *testing.T) {
	now := time.Now()

	uc := &getTestsUCStub{
		executeFn: func(
			ctx context.Context,
			req gettests.GetTestsRequest,
		) (gettests.GetTestsResponse, error) {
			require.Equal(t, 10, req.UserID)

			return gettests.GetTestsResponse{
				Tests: []gettests.TestInfo{
					{
						ID:               1,
						Name:             "Math Test",
						CntQuestions:     20,
						CntHardQuestions: 5,
						Groups:           []string{"A-01"},
						DateCreated:      now,
					},
				},
			}, nil
		},
	}

	h := New(uc, nil, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests",
		nil,
	)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetTests(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.GetTestsResponse

	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Len(t, resp.Tests, 1)
	assert.Equal(t, "Math Test", resp.Tests[0].Name)
}

func TestGetTests_Unauthorized(t *testing.T) {
	h := New(&getTestsUCStub{}, nil, nil, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/lecturer/tests",
		nil,
	)

	rec := httptest.NewRecorder()

	h.GetTests(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

