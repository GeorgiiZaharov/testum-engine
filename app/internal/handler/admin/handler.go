package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"testum-engine/app/internal/handler/admin/dto"
	"testum-engine/app/internal/handler/middleware"

	createlecturer "testum-engine/app/internal/service/use_case/admin/create_lecturer"
	deletelecturer "testum-engine/app/internal/service/use_case/admin/delete_lecturer"
	getlecturers "testum-engine/app/internal/service/use_case/admin/get_lecturers"
)

type CreateLecturerUseCase interface {
	Execute(
		ctx context.Context,
		req createlecturer.CreateLecturerRequest,
	) (createlecturer.CreateLecturerResponse, error)
}

type DeleteLecturerUseCase interface {
	Execute(
		ctx context.Context,
		req deletelecturer.DeleteLecturerRequest,
	) (deletelecturer.DeleteLecturerResponse, error)
}

type GetLecturersUseCase interface {
	Execute(
		ctx context.Context,
		req getlecturers.GetLecturersRequest,
	) (getlecturers.GetLecturersResponse, error)
}

type Handler struct {
	createLecturerUC CreateLecturerUseCase
	deleteLecturerUC DeleteLecturerUseCase
	getLecturersUC   GetLecturersUseCase
}

func New(
	createLecturerUC CreateLecturerUseCase,
	deleteLecturerUC DeleteLecturerUseCase,
	getLecturersUC GetLecturersUseCase,
) *Handler {
	return &Handler{
		createLecturerUC: createLecturerUC,
		deleteLecturerUC: deleteLecturerUC,
		getLecturersUC:   getLecturersUC,
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

func mapCreateLecturerError(err error) (int, string) {
	switch {
	case errors.Is(err, createlecturer.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, createlecturer.ErrForbidden):
		return http.StatusForbidden, "forbidden"

	case errors.Is(err, createlecturer.ErrLDAPFailed):
		return http.StatusInternalServerError, "internal server error"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapDeleteLecturerError(err error) (int, string) {
	switch {
	case errors.Is(err, deletelecturer.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, deletelecturer.ErrForbidden):
		return http.StatusForbidden, "forbidden"

	case errors.Is(err, deletelecturer.ErrNotFound):
		return http.StatusNotFound, "not found"

	case errors.Is(err, deletelecturer.ErrNotLecturer):
		return http.StatusConflict, "user is not lecturer"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func mapGetLecturersError(err error) (int, string) {
	switch {
	case errors.Is(err, getlecturers.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"

	case errors.Is(err, getlecturers.ErrForbidden):
		return http.StatusForbidden, "forbidden"

	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

//
// POST /admin/lecturers
//

func (h *Handler) CreateLecturer(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateLecturerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Login = strings.TrimSpace(req.Login)

	if req.Login == "" {
		writeError(w, http.StatusBadRequest, "login is required")
		return
	}

	res, err := h.createLecturerUC.Execute(
		r.Context(),
		createlecturer.CreateLecturerRequest{
			AdminID: adminID,
			Login:   req.Login,
		},
	)
	if err != nil {
		code, message := mapCreateLecturerError(err)
		writeError(w, code, message)
		return
	}

	writeJSON(w, http.StatusCreated, dto.CreateLecturerResponse{
		Success: res.Success,
	})
}

//
// DELETE /admin/lecturers/{lecturer_id}
//

func (h *Handler) DeleteLecturer(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	lecturerIDRaw := r.PathValue("lecturer_id")
	if lecturerIDRaw == "" {
		writeError(w, http.StatusBadRequest, "lecturer_id is required")
		return
	}

	lecturerID, err := strconv.Atoi(lecturerIDRaw)
	if err != nil || lecturerID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid lecturer_id")
		return
	}

	res, err := h.deleteLecturerUC.Execute(
		r.Context(),
		deletelecturer.DeleteLecturerRequest{
			AdminID:    adminID,
			LecturerID: lecturerID,
		},
	)
	if err != nil {
		code, message := mapDeleteLecturerError(err)
		writeError(w, code, message)
		return
	}

	writeJSON(w, http.StatusOK, dto.DeleteLecturerResponse{
		Success: res.Success,
	})
}

//
// GET /admin/lecturers
//

func (h *Handler) GetLecturers(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	res, err := h.getLecturersUC.Execute(
		r.Context(),
		getlecturers.GetLecturersRequest{
			UserID: userID,
		},
	)
	if err != nil {
		code, message := mapGetLecturersError(err)
		writeError(w, code, message)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToGetLecturersResponse(res),
	)
}
