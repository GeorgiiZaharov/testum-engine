package testattempt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"testum-engine/app/internal/handler/middleware"
	"testum-engine/app/internal/handler/student/test_attempt/dto"

	getbasetasks "testum-engine/app/internal/service/use_case/student/get_base_tasks"
	gethardtasks "testum-engine/app/internal/service/use_case/student/get_hard_tasks"
	postbaseanswers "testum-engine/app/internal/service/use_case/student/post_base_answers"
	posthardanswers "testum-engine/app/internal/service/use_case/student/post_hard_answers"
)

type GetHardTasksUseCase interface {
	Execute(
		ctx context.Context,
		req gethardtasks.GetHardTasksRequest,
	) (gethardtasks.GetHardTasksResponse, error)
}

type GetBaseTasksUseCase interface {
	Execute(
		ctx context.Context,
		req getbasetasks.GetBaseTasksRequest,
	) (getbasetasks.GetBaseTasksResponse, error)
}

type PostHardAnswersUseCase interface {
	Execute(
		ctx context.Context,
		userID int,
		testID int,
		answers []posthardanswers.TaskAnswer,
	) (posthardanswers.PostHardAnswersResponse, error)
}

type PostBaseAnswersUseCase interface {
	Execute(
		ctx context.Context,
		userID int,
		testID int,
		answers []postbaseanswers.TaskAnswer,
	) (postbaseanswers.PostBaseAnswersResponse, error)
}

type Handler struct {
	getHardTasksUC    GetHardTasksUseCase
	getBaseTasksUC    GetBaseTasksUseCase
	postHardAnswersUC PostHardAnswersUseCase
	postBaseAnswersUC PostBaseAnswersUseCase
}

func New(
	getHardTasksUC GetHardTasksUseCase,
	getBaseTasksUC GetBaseTasksUseCase,
	postHardAnswersUC PostHardAnswersUseCase,
	postBaseAnswersUC PostBaseAnswersUseCase,
) *Handler {
	return &Handler{
		getHardTasksUC:    getHardTasksUC,
		getBaseTasksUC:    getBaseTasksUC,
		postHardAnswersUC: postHardAnswersUC,
		postBaseAnswersUC: postBaseAnswersUC,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"error": message,
	})
}

func parseTestID(r *http.Request) (int, error) {
	raw := r.PathValue("test_id")
	if raw == "" {
		return 0, errors.New("test_id is required")
	}

	testID, err := strconv.Atoi(raw)
	if err != nil || testID <= 0 {
		return 0, errors.New("invalid test_id")
	}

	return testID, nil
}

func mapGetHardTasksError(err error) (int, string) {
	switch {
	case errors.Is(err, gethardtasks.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, gethardtasks.ErrForbidden):
		return http.StatusForbidden, "forbidden"

	case errors.Is(err, gethardtasks.ErrAccessDenied):
		return http.StatusForbidden, "access denied"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapGetBaseTasksError(err error) (int, string) {
	switch {
	case errors.Is(err, getbasetasks.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, getbasetasks.ErrAccessDenied):
		return http.StatusForbidden, "access denied"

	case errors.Is(err, getbasetasks.ErrTestCompleted):
		return http.StatusConflict, "test already completed"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapPostHardAnswersError(err error) (int, string) {
	switch {
	case errors.Is(err, posthardanswers.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, posthardanswers.ErrAccessDenied):
		return http.StatusForbidden, "access denied"

	case errors.Is(err, posthardanswers.ErrAlreadySubmitted):
		return http.StatusConflict, "hard answers already submitted"

	case errors.Is(err, posthardanswers.ErrTestAlreadyFinished):
		return http.StatusConflict, "test already finished"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapPostBaseAnswersError(err error) (int, string) {
	switch {
	case errors.Is(err, postbaseanswers.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, postbaseanswers.ErrAccessDenied):
		return http.StatusForbidden, "access denied"

	case errors.Is(err, postbaseanswers.ErrAlreadySubmitted):
		return http.StatusConflict, "base answers already submitted"

	case errors.Is(err, postbaseanswers.ErrHardBlockNotPassed):
		return http.StatusConflict, "hard block not passed"

	case errors.Is(err, postbaseanswers.ErrTestAlreadyFinished):
		return http.StatusConflict, "test already finished"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

//
// GET /student/tests/{test_id}/hard
//

func (h *Handler) GetHardTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseTestID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.getHardTasksUC.Execute(
		r.Context(),
		gethardtasks.GetHardTasksRequest{
			UserID: userID,
			TestID: testID,
		},
	)
	if err != nil {
		code, message := mapGetHardTasksError(err)
		writeError(w, code, message)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToGetHardTasksResponse(res),
	)
}

//
// GET /student/tests/{test_id}/base
//

func (h *Handler) GetBaseTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseTestID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.getBaseTasksUC.Execute(
		r.Context(),
		getbasetasks.GetBaseTasksRequest{
			UserID: userID,
			TestID: testID,
		},
	)
	if err != nil {
		code, message := mapGetBaseTasksError(err)
		writeError(w, code, message)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToGetBaseTasksResponse(res),
	)
}

//
// POST /student/tests/{test_id}/hard
//

func (h *Handler) PostHardAnswers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseTestID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.PostAnswersRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Answers) == 0 {
		writeError(w, http.StatusBadRequest, "answers are required")
		return
	}

	res, err := h.postHardAnswersUC.Execute(
		r.Context(),
		userID,
		testID,
		dto.ToPostHardAnswersModel(req),
	)
	if err != nil {
		code, message := mapPostHardAnswersError(err)
		writeError(w, code, message)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToPostHardAnswersResponse(res),
	)
}

//
// POST /student/tests/{test_id}/base
//

func (h *Handler) PostBaseAnswers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	testID, err := parseTestID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req dto.PostAnswersRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Answers) == 0 {
		writeError(w, http.StatusBadRequest, "answers are required")
		return
	}

	res, err := h.postBaseAnswersUC.Execute(
		r.Context(),
		userID,
		testID,
		dto.ToPostBaseAnswersModel(req),
	)
	if err != nil {
		code, message := mapPostBaseAnswersError(err)
		writeError(w, code, message)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToPostBaseAnswersResponse(res),
	)
}
