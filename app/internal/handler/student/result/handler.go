package result

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"testum-engine/app/internal/handler/middleware"
	"testum-engine/app/internal/handler/student/result/dto"
	gettestresult "testum-engine/app/internal/service/use_case/student/get_test_result"
)

/*
----------------------------
UseCase interface
----------------------------
*/

type GetTestResultUseCase interface {
	Execute(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error)
}

/*
----------------------------
Handler
----------------------------
*/

type Handler struct {
	getTestResultUC GetTestResultUseCase
}

func New(getTestResultUC GetTestResultUseCase) *Handler {
	return &Handler{
		getTestResultUC: getTestResultUC,
	}
}

/*
----------------------------
HTTP handler
----------------------------
*/

// GET /student/test/{test_id}/result
func (h *Handler) GetTestResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testIDRaw := r.PathValue("test_id")
	if testIDRaw == "" {
		writeError(w, http.StatusBadRequest, "test_id is required")
		return
	}

	testID, err := strconv.Atoi(testIDRaw)
	if err != nil || testID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	ucResp, err := h.getTestResultUC.Execute(r.Context(), gettestresult.GetTestResultRequest{
		UserID: userID,
		TestID: testID,
	})
	if err != nil {
		code, msg := mapGetTestResultError(err)
		writeError(w, code, msg)
		return
	}

	// ✅ SUCCESS → DTO напрямую
	writeJSON(w, http.StatusOK, dto.ToResponse(ucResp))
}

/*
----------------------------
Error mapping
----------------------------
*/

func mapGetTestResultError(err error) (int, string) {
	switch {
	case errors.Is(err, gettestresult.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, gettestresult.ErrAccessDenied):
		return http.StatusForbidden, "access denied"

	case errors.Is(err, gettestresult.ErrResultNotFound):
		return http.StatusNotFound, "result not found"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

/*
----------------------------
Helpers
----------------------------
*/

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}
