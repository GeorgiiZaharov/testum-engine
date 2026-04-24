package gettests

import (
	"context"
	"encoding/json"
	"net/http"

	"testum-engine/app/internal/handler/middleware"
	"testum-engine/app/internal/handler/student/get_tests/dto"
	getactivetest "testum-engine/app/internal/service/use_case/student/get_active_test"
	getfinishedtest "testum-engine/app/internal/service/use_case/student/get_finished_test"
)

/*
----------------------------
UseCases
----------------------------
*/

type GetActiveTestUseCase interface {
	Execute(ctx context.Context, req getactivetest.GetActiveTestRequest) (getactivetest.GetActiveTestResponse, error)
}

type GetFinishedTestUseCase interface {
	Execute(ctx context.Context, req getfinishedtest.GetFinishedTestRequest) (getfinishedtest.GetFinishedTestResponse, error)
}

/*
----------------------------
Handler
----------------------------
*/

type Handler struct {
	activeTestUseCase   GetActiveTestUseCase
	finishedTestUseCase GetFinishedTestUseCase
}

func New(
	activeTestUseCase GetActiveTestUseCase,
	finishedTestUseCase GetFinishedTestUseCase,
) *Handler {
	return &Handler{
		activeTestUseCase:   activeTestUseCase,
		finishedTestUseCase: finishedTestUseCase,
	}
}

/*
----------------------------
GET ACTIVE TESTS
----------------------------
*/

func (h *Handler) GetActiveTests(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req := getactivetest.GetActiveTestRequest{
		UserID: userID,
	}

	resp, err := h.activeTestUseCase.Execute(r.Context(), req)
	if err != nil {
		if err == getactivetest.ErrInvalidInput {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, dto.ConvertActiveTestsToDTO(resp.ActiveTests))
}

/*
----------------------------
GET FINISHED TESTS
----------------------------
*/

func (h *Handler) GetFinishedTests(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req := getfinishedtest.GetFinishedTestRequest{
		UserID: userID,
	}

	resp, err := h.finishedTestUseCase.Execute(r.Context(), req)
	if err != nil {
		if err == getfinishedtest.ErrInvalidInput {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, dto.ConvertFinishedTestsToDTO(resp.FinishedTests))
}

/*
----------------------------
HELPERS
----------------------------
*/

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
