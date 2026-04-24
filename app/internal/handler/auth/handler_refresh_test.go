package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testum-engine/app/internal/handler/auth/dto"
	"testum-engine/app/internal/handler/middleware"
	refresh "testum-engine/app/internal/service/use_case/auth/refresh"
)

type refreshUCStub struct {
	executeFn func(ctx context.Context, req refresh.AuthRefreshRequest) (refresh.AuthRefreshResponse, error)
}

func (s *refreshUCStub) Execute(ctx context.Context, req refresh.AuthRefreshRequest) (refresh.AuthRefreshResponse, error) {
	return s.executeFn(ctx, req)
}
func TestRefresh_Success(t *testing.T) {
	uc := &refreshUCStub{
		executeFn: func(ctx context.Context, req refresh.AuthRefreshRequest) (refresh.AuthRefreshResponse, error) {

			assert.Equal(t, 42, req.UserID)

			return refresh.AuthRefreshResponse{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}, nil
		},
	}

	h := New(nil, uc, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req = middleware.WithUserID(req, 42)

	rec := httptest.NewRecorder()

	h.Refresh(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.RefreshResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, "access", resp.AccessToken)
	assert.Equal(t, "refresh", resp.RefreshToken)
}
func TestRefresh_Unauthorized(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	rec := httptest.NewRecorder()

	h.Refresh(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
func TestRefresh_AuthFailed(t *testing.T) {
	uc := &refreshUCStub{
		executeFn: func(ctx context.Context, req refresh.AuthRefreshRequest) (refresh.AuthRefreshResponse, error) {
			return refresh.AuthRefreshResponse{}, refresh.ErrAuthFailed
		},
	}

	h := New(nil, uc, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	req = middleware.WithUserID(req, 42)

	rec := httptest.NewRecorder()

	h.Refresh(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
