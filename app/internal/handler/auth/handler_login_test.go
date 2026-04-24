package auth

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

	"testum-engine/app/internal/handler/auth/dto"

	loginuc "testum-engine/app/internal/service/use_case/auth/login"
)

type loginUCStub struct {
	executeFn func(ctx context.Context, req loginuc.AuthLoginRequest) (loginuc.AuthLoginResponse, error)
}

func (s *loginUCStub) Execute(ctx context.Context, req loginuc.AuthLoginRequest) (loginuc.AuthLoginResponse, error) {
	return s.executeFn(ctx, req)
}

func TestLogin_Success(t *testing.T) {
	uc := &loginUCStub{
		executeFn: func(ctx context.Context, req loginuc.AuthLoginRequest) (loginuc.AuthLoginResponse, error) {

			assert.Equal(t, "john", req.Login)
			assert.Equal(t, "pass", req.Password)

			return loginuc.AuthLoginResponse{
				AccessToken:  "access",
				RefreshToken: "refresh",
			}, nil
		},
	}

	h := New(uc, nil, nil)

	body := dto.LoginRequest{
		Login:    "john",
		Password: "pass",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.LoginResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)

	require.NoError(t, err)
	assert.Equal(t, "access", resp.AccessToken)
	assert.Equal(t, "refresh", resp.RefreshToken)
}

func TestLogin_InvalidJSON(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString("{invalid}"))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLogin_EmptyFields(t *testing.T) {
	h := New(nil, nil, nil)

	body := dto.LoginRequest{
		Login:    "   ",
		Password: "",
	}

	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLogin_Unauthorized(t *testing.T) {
	uc := &loginUCStub{
		executeFn: func(ctx context.Context, req loginuc.AuthLoginRequest) (loginuc.AuthLoginResponse, error) {
			return loginuc.AuthLoginResponse{}, loginuc.ErrUnauthorized
		},
	}

	h := New(uc, nil, nil)

	body, _ := json.Marshal(dto.LoginRequest{
		Login:    "john",
		Password: "bad",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLogin_DefaultError(t *testing.T) {
	uc := &loginUCStub{
		executeFn: func(ctx context.Context, req loginuc.AuthLoginRequest) (loginuc.AuthLoginResponse, error) {
			return loginuc.AuthLoginResponse{}, errors.New("boom")
		},
	}

	h := New(uc, nil, nil)

	body, _ := json.Marshal(dto.LoginRequest{
		Login:    "john",
		Password: "pass",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
