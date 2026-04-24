package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"testum-engine/app/internal/handler/auth/dto"
	"testum-engine/app/internal/handler/middleware"

	getme "testum-engine/app/internal/service/use_case/auth/get_me"
	login "testum-engine/app/internal/service/use_case/auth/login"
	refresh "testum-engine/app/internal/service/use_case/auth/refresh"
)

type LoginUseCase interface {
	Execute(ctx context.Context, req login.AuthLoginRequest) (login.AuthLoginResponse, error)
}

type RefreshUseCase interface {
	Execute(ctx context.Context, req refresh.AuthRefreshRequest) (refresh.AuthRefreshResponse, error)
}

type GetMeUseCase interface {
	Execute(ctx context.Context, req getme.GetMeRequest) (getme.GetMeResponse, error)
}

type Handler struct {
	loginUC   LoginUseCase
	refreshUC RefreshUseCase
	getMeUC   GetMeUseCase
}

func New(
	loginUC LoginUseCase,
	refreshUC RefreshUseCase,
	getMeUC GetMeUseCase,
) *Handler {
	return &Handler{
		loginUC:   loginUC,
		refreshUC: refreshUC,
		getMeUC:   getMeUC,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"success": false,
		"error":   message,
	})
}

func mapLoginError(err error) (int, string) {
	switch {
	case errors.Is(err, login.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, login.ErrUnauthorized):
		return http.StatusUnauthorized, "invalid credentials"

	case errors.Is(err, login.ErrLDAPFailed):
		return http.StatusInternalServerError, "internal server error"

	case errors.Is(err, login.ErrTokenGenerate):
		return http.StatusInternalServerError, "internal server error"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapRefreshError(err error) (int, string) {
	switch {
	case errors.Is(err, refresh.ErrInvalidUserID):
		return http.StatusUnauthorized, "unauthorized"

	case errors.Is(err, refresh.ErrAuthFailed):
		return http.StatusInternalServerError, "internal server error"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapGetMeError(err error) (int, string) {
	switch {
	case errors.Is(err, getme.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, getme.ErrNotFound):
		return http.StatusNotFound, "not found"

	case errors.Is(err, getme.ErrLDAPFailed):
		return http.StatusInternalServerError, "internal server error"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

//
// POST /auth/login
//

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	req.Password = strings.TrimSpace(req.Password)

	if req.Login == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "login and password are required")
		return
	}

	res, err := h.loginUC.Execute(r.Context(), login.AuthLoginRequest{
		Login:    req.Login,
		Password: req.Password,
	})
	if err != nil {
		code, msg := mapLoginError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToLoginResponse(res))
}

//
// POST /auth/refresh
//

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := h.refreshUC.Execute(r.Context(), refresh.AuthRefreshRequest{
		UserID: userID,
	})
	if err != nil {
		code, msg := mapRefreshError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToRefreshResponse(res))
}

//
// GET /auth/me
//

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := h.getMeUC.Execute(r.Context(), getme.GetMeRequest{
		UserID: userID,
	})
	if err != nil {
		code, msg := mapGetMeError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToGetMeResponse(res))
}
