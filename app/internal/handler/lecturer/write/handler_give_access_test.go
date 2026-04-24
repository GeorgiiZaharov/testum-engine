package write

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

	"testum-engine/app/internal/handler/lecturer/write/dto"
	"testum-engine/app/internal/handler/middleware"
	giveaccess "testum-engine/app/internal/service/use_case/lecturer/give_access"
)

type giveAccessUCStub struct {
	executeFn func(
		ctx context.Context,
		req giveaccess.GiveAccessRequest,
	) (giveaccess.GiveAccessResponse, error)
}

func (s *giveAccessUCStub) Execute(
	ctx context.Context,
	req giveaccess.GiveAccessRequest,
) (giveaccess.GiveAccessResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGiveAccess_Success(t *testing.T) {
	uc := &giveAccessUCStub{
		executeFn: func(
			ctx context.Context,
			req giveaccess.GiveAccessRequest,
		) (giveaccess.GiveAccessResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 55, req.TestID)
			require.Equal(t, "A-01", req.Group)

			return giveaccess.GiveAccessResponse{
				Success: true,
			}, nil
		},
	}

	h := New(nil, uc, nil)

	body := dto.AccessRequest{
		Group: "A-01",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/tests/55/access",
		bytes.NewReader(b),
	)
	req.SetPathValue("test_id", "55")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GiveAccess(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.AccessResponse

	err = json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
}

func TestGiveAccess_InvalidBody(t *testing.T) {
	h := New(nil, &giveAccessUCStub{}, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/lecturer/tests/55/access",
		bytes.NewReader([]byte("{invalid")),
	)
	req.SetPathValue("test_id", "55")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GiveAccess(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGiveAccess_ErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "invalid input",
			err:        giveaccess.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "forbidden",
			err:        giveaccess.ErrAccessDenied,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "unknown",
			err:        errors.New("unexpected"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &giveAccessUCStub{
				executeFn: func(
					ctx context.Context,
					req giveaccess.GiveAccessRequest,
				) (giveaccess.GiveAccessResponse, error) {
					return giveaccess.GiveAccessResponse{}, tt.err
				},
			}

			h := New(nil, uc, nil)

			body := dto.AccessRequest{
				Group: "A-01",
			}

			b, _ := json.Marshal(body)

			req := httptest.NewRequest(
				http.MethodPost,
				"/lecturer/tests/55/access",
				bytes.NewReader(b),
			)
			req.SetPathValue("test_id", "55")
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()

			h.GiveAccess(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
