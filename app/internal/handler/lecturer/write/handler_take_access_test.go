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
	takeaccess "testum-engine/app/internal/service/use_case/lecturer/take_access"
)

type takeAccessUCStub struct {
	executeFn func(
		ctx context.Context,
		req takeaccess.TakeAccessRequest,
	) (takeaccess.TakeAccessResponse, error)
}

func (s *takeAccessUCStub) Execute(
	ctx context.Context,
	req takeaccess.TakeAccessRequest,
) (takeaccess.TakeAccessResponse, error) {
	return s.executeFn(ctx, req)
}

func TestTakeAccess_Success(t *testing.T) {
	uc := &takeAccessUCStub{
		executeFn: func(
			ctx context.Context,
			req takeaccess.TakeAccessRequest,
		) (takeaccess.TakeAccessResponse, error) {
			require.Equal(t, 10, req.UserID)
			require.Equal(t, 55, req.TestID)
			require.Equal(t, "A-01", req.Group)

			return takeaccess.TakeAccessResponse{
				Success: true,
			}, nil
		},
	}

	h := New(nil, nil, uc, nil)

	body := dto.AccessRequest{
		Group: "A-01",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/lecturer/tests/55/access",
		bytes.NewReader(b),
	)
	req.SetPathValue("test_id", "55")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.TakeAccess(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.AccessResponse

	err = json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.True(t, resp.Success)
}

func TestTakeAccess_InvalidBody(t *testing.T) {
	h := New(nil, nil, &takeAccessUCStub{}, nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/lecturer/tests/55/access",
		bytes.NewReader([]byte("{invalid")),
	)
	req.SetPathValue("test_id", "55")
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.TakeAccess(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTakeAccess_ErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "invalid input",
			err:        takeaccess.ErrInvalidInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "forbidden",
			err:        takeaccess.ErrAccessDenied,
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
			uc := &takeAccessUCStub{
				executeFn: func(
					ctx context.Context,
					req takeaccess.TakeAccessRequest,
				) (takeaccess.TakeAccessResponse, error) {
					return takeaccess.TakeAccessResponse{}, tt.err
				},
			}

			h := New(nil, nil, uc, nil)

			body := dto.AccessRequest{
				Group: "A-01",
			}

			b, _ := json.Marshal(body)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/lecturer/tests/55/access",
				bytes.NewReader(b),
			)
			req.SetPathValue("test_id", "55")
			req = middleware.WithUserID(req, 10)

			rec := httptest.NewRecorder()

			h.TakeAccess(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
