package write

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"testum-engine/app/internal/handler/lecturer/write/dto"
	"testum-engine/app/internal/handler/middleware"

	deleteattempt "testum-engine/app/internal/service/use_case/lecturer/delete_attempt"
	deletetest "testum-engine/app/internal/service/use_case/lecturer/delete_test"
	giveaccess "testum-engine/app/internal/service/use_case/lecturer/give_access"
	takeaccess "testum-engine/app/internal/service/use_case/lecturer/take_access"
)

//
// USE CASE INTERFACES
//

type DeleteTestUseCase interface {
	Execute(ctx context.Context, req deletetest.DeleteTestRequest) (deletetest.DeleteTestResponse, error)
}

type GiveAccessUseCase interface {
	Execute(ctx context.Context, req giveaccess.GiveAccessRequest) (giveaccess.GiveAccessResponse, error)
}

type TakeAccessUseCase interface {
	Execute(ctx context.Context, req takeaccess.TakeAccessRequest) (takeaccess.TakeAccessResponse, error)
}

type DeleteAttemptUseCase interface {
	Execute(rCtx context.Context, req deleteattempt.DeleteAttemptRequest) (deleteattempt.DeleteAttemptResponse, error)
}

//
// HANDLER
//

type Handler struct {
	deleteTestUC    DeleteTestUseCase
	giveAccessUC    GiveAccessUseCase
	takeAccessUC    TakeAccessUseCase
	deleteAttemptUC DeleteAttemptUseCase
}

func New(
	deleteTestUC DeleteTestUseCase,
	giveAccessUC GiveAccessUseCase,
	takeAccessUC TakeAccessUseCase,
	deleteAttemptUC DeleteAttemptUseCase,
) *Handler {
	return &Handler{
		deleteTestUC:    deleteTestUC,
		giveAccessUC:    giveAccessUC,
		takeAccessUC:    takeAccessUC,
		deleteAttemptUC: deleteAttemptUC,
	}
}

//
// HELPERS
//

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{
		"success": false,
		"error":   msg,
	})
}

func parseTestID(r *http.Request) (int, error) {
	id := r.PathValue("test_id")
	return strconv.Atoi(id)
}

//
// ERROR MAPPING
//

func mapDeleteError(err error) (int, string) {
	switch {
	case errors.Is(err, deletetest.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	case errors.Is(err, deletetest.ErrAccessDenied):
		return http.StatusForbidden, "access denied"
	case errors.Is(err, deletetest.ErrTestNotFound):
		return http.StatusNotFound, "test not found"
	case errors.Is(err, deletetest.ErrStorageFailed):
		return http.StatusInternalServerError, "internal server error"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapAccessError(err error) (int, string) {
	switch {
	case errors.Is(err, giveaccess.ErrAccessDenied):
		return http.StatusForbidden, "access denied"
	case errors.Is(err, giveaccess.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	case errors.Is(err, takeaccess.ErrAccessDenied):
		return http.StatusForbidden, "access denied"
	case errors.Is(err, takeaccess.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapDeleteAttemptError(err error) (int, string) {
	switch {
	case errors.Is(err, deleteattempt.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	case errors.Is(err, deleteattempt.ErrAccessDenied):
		return http.StatusForbidden, "access denied"
	case errors.Is(err, deleteattempt.ErrDeleteFailed):
		return http.StatusInternalServerError, "failed to delete attempt"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

//
// DELETE /lecturer/tests/{test_id}
//

func (h *Handler) DeleteTest(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseTestID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	res, err := h.deleteTestUC.Execute(r.Context(), deletetest.DeleteTestRequest{
		UserID: userID,
		TestID: testID,
	})
	if err != nil {
		code, msg := mapDeleteError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.DeleteTestResponse{
		Success: res.Success,
	})
}

//
// POST /lecturer/tests/{test_id}/access
//

func (h *Handler) GiveAccess(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseTestID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	var req dto.AccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.giveAccessUC.Execute(r.Context(), giveaccess.GiveAccessRequest{
		UserID: userID,
		TestID: testID,
		Group:  req.Group,
	})
	if err != nil {
		code, msg := mapAccessError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.AccessResponse{
		Success: res.Success,
	})
}

//
// DELETE /lecturer/tests/{test_id}/access
//

func (h *Handler) TakeAccess(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseTestID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	var req dto.AccessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	res, err := h.takeAccessUC.Execute(r.Context(), takeaccess.TakeAccessRequest{
		UserID: userID,
		TestID: testID,
		Group:  req.Group,
	})
	if err != nil {
		code, msg := mapAccessError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.AccessResponse{
		Success: res.Success,
	})
}

//
// DELETE /lecturer/attempt/test/{test_id}/user/{user_id}
//

func (h *Handler) DeleteAttempt(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Парсим test_id и student user_id из URL
	testIDStr := r.PathValue("test_id")
	studentIDStr := r.PathValue("user_id")

	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	studentID, err := strconv.Atoi(studentIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	// Вызов use case
	resp, err := h.deleteAttemptUC.Execute(r.Context(), deleteattempt.DeleteAttemptRequest{
		LecturerID: userID,
		TestID:     testID,
		UserID:     studentID,
	})
	if err != nil {
		code, msg := mapDeleteAttemptError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.FromUseCaseResponse(resp))
}
