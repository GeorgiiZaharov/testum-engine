package read

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"testum-engine/app/internal/handler/lecturer/read/dto"
	"testum-engine/app/internal/handler/middleware"

	getgroups "testum-engine/app/internal/service/use_case/lecturer/get_groups"
	gettest "testum-engine/app/internal/service/use_case/lecturer/get_test"
	gettestresult "testum-engine/app/internal/service/use_case/lecturer/get_test_result"
	gettests "testum-engine/app/internal/service/use_case/lecturer/get_tests"
)

type GetTestsUseCase interface {
	Execute(ctx context.Context, req gettests.GetTestsRequest) (gettests.GetTestsResponse, error)
}

type GetTestUseCase interface {
	Execute(ctx context.Context, req gettest.GetTestRequest) (gettest.GetTestResponse, error)
}

type GetGroupsUseCase interface {
	Execute(ctx context.Context, req getgroups.GetGroupsRequest) (getgroups.GetGroupsResponse, error)
}

type GetTestResultUseCase interface {
	Execute(ctx context.Context, req gettestresult.GetTestResultRequest) (gettestresult.GetTestResultResponse, error)
}

type Handler struct {
	getTestsUC  GetTestsUseCase
	getTestUC   GetTestUseCase
	getGroupsUC GetGroupsUseCase
	getResultUC GetTestResultUseCase
}

func New(
	getTestsUC GetTestsUseCase,
	getTestUC GetTestUseCase,
	getGroupsUC GetGroupsUseCase,
	getResultUC GetTestResultUseCase,
) *Handler {
	return &Handler{
		getTestsUC:  getTestsUC,
		getTestUC:   getTestUC,
		getGroupsUC: getGroupsUC,
		getResultUC: getResultUC,
	}
}

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

func parseID(r *http.Request) (int, error) {
	id := r.PathValue("test_id")
	return strconv.Atoi(id)
}

//
// GET /lecturer/tests
//

func (h *Handler) GetTests(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := h.getTestsUC.Execute(r.Context(), gettests.GetTestsRequest{
		UserID: userID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToGetTests(res))
}

//
// GET /lecturer/tests/{test_id}
//

func (h *Handler) GetTest(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	res, err := h.getTestUC.Execute(r.Context(), gettest.GetTestRequest{
		UserID: userID,
		TestID: testID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToGetTest(res))
}

//
// GET /lecturer/tests/{test_id}/groups?year=2026
//

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		writeError(w, http.StatusBadRequest, "year is required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid year")
		return
	}

	res, err := h.getGroupsUC.Execute(r.Context(), getgroups.GetGroupsRequest{
		UserID: userID,
		TestID: testID,
		Year:   year,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToGetGroups(res))
}

//
// GET /lecturer/tests/{test_id}/result?group=A-01&year=2026
//

func (h *Handler) GetTestResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid test_id")
		return
	}

	group := r.URL.Query().Get("group")
	yearStr := r.URL.Query().Get("year")

	if group == "" || yearStr == "" {
		writeError(w, http.StatusBadRequest, "group and year required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid year")
		return
	}

	res, err := h.getResultUC.Execute(r.Context(), gettestresult.GetTestResultRequest{
		UserID: userID,
		TestID: testID,
		Group:  group,
		Year:   year,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, dto.ToGetTestResult(res))
}
