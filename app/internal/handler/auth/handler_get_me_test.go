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
	getme "testum-engine/app/internal/service/use_case/auth/get_me"
)

type getMeUCStub struct {
	executeFn func(ctx context.Context, req getme.GetMeRequest) (getme.GetMeResponse, error)
}

func (s *getMeUCStub) Execute(ctx context.Context, req getme.GetMeRequest) (getme.GetMeResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGetMe_Success(t *testing.T) {
	uc := &getMeUCStub{
		executeFn: func(ctx context.Context, req getme.GetMeRequest) (getme.GetMeResponse, error) {

			require.Equal(t, 42, req.UserID)

			group := "admin"

			return getme.GetMeResponse{
				ID:         42,
				Login:      "john",
				Mail:       "john@mail.com",
				Name:       "John Doe",
				Group:      &group,
				IsLecturer: true,
			}, nil
		},
	}

	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = middleware.WithUserID(req, 42)

	rec := httptest.NewRecorder()

	h.GetMe(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.GetMeResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Equal(t, 42, resp.ID)
	assert.Equal(t, "john", resp.Login)
	assert.Equal(t, "john@mail.com", resp.Mail)
	assert.Equal(t, "John Doe", resp.Name)
	assert.NotNil(t, resp.Group)
	assert.Equal(t, "admin", *resp.Group)
	assert.True(t, resp.IsLecturer)
}

func TestGetMe_Unauthorized(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	h.GetMe(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetMe_NotFound(t *testing.T) {
	uc := &getMeUCStub{
		executeFn: func(ctx context.Context, req getme.GetMeRequest) (getme.GetMeResponse, error) {
			return getme.GetMeResponse{}, getme.ErrNotFound
		},
	}

	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = middleware.WithUserID(req, 42)

	rec := httptest.NewRecorder()

	h.GetMe(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetMe_InvalidInput(t *testing.T) {
	uc := &getMeUCStub{
		executeFn: func(ctx context.Context, req getme.GetMeRequest) (getme.GetMeResponse, error) {
			return getme.GetMeResponse{}, getme.ErrInvalidInput
		},
	}

	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = middleware.WithUserID(req, 42)

	rec := httptest.NewRecorder()

	h.GetMe(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetMe_InternalError(t *testing.T) {
	uc := &getMeUCStub{
		executeFn: func(ctx context.Context, req getme.GetMeRequest) (getme.GetMeResponse, error) {
			return getme.GetMeResponse{}, getme.ErrLDAPFailed
		},
	}

	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = middleware.WithUserID(req, 42)

	rec := httptest.NewRecorder()

	h.GetMe(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
