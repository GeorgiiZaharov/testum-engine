package admin

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"testum-engine/app/internal/handler/middleware"

	deletelecturer "testum-engine/app/internal/service/use_case/admin/delete_lecturer"
)

type deleteLecturerUCStub struct {
	executeFn func(ctx context.Context, req deletelecturer.DeleteLecturerRequest) (deletelecturer.DeleteLecturerResponse, error)
}

func (s *deleteLecturerUCStub) Execute(ctx context.Context, req deletelecturer.DeleteLecturerRequest) (deletelecturer.DeleteLecturerResponse, error) {
	return s.executeFn(ctx, req)
}

func TestDeleteLecturer_Success(t *testing.T) {
	uc := &deleteLecturerUCStub{
		executeFn: func(ctx context.Context, req deletelecturer.DeleteLecturerRequest) (deletelecturer.DeleteLecturerResponse, error) {

			assert.Equal(t, 10, req.AdminID)
			assert.Equal(t, 5, req.LecturerID)

			return deletelecturer.DeleteLecturerResponse{
				Success: true,
			}, nil
		},
	}

	h := New(nil, uc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/lecturers/5", nil)
	req = middleware.WithUserID(req, 10)
	req.SetPathValue("lecturer_id", "5")

	rec := httptest.NewRecorder()

	h.DeleteLecturer(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestDeleteLecturer_Unauthorized(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/lecturers/5", nil)
	rec := httptest.NewRecorder()

	h.DeleteLecturer(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteLecturer_InvalidID(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/lecturers/abc", nil)
	req = middleware.WithUserID(req, 10)
	req.SetPathValue("lecturer_id", "abc")

	rec := httptest.NewRecorder()

	h.DeleteLecturer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteLecturer_Forbidden(t *testing.T) {
	uc := &deleteLecturerUCStub{
		executeFn: func(ctx context.Context, req deletelecturer.DeleteLecturerRequest) (deletelecturer.DeleteLecturerResponse, error) {
			return deletelecturer.DeleteLecturerResponse{}, deletelecturer.ErrForbidden
		},
	}

	h := New(nil, uc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/lecturers/5", nil)
	req = middleware.WithUserID(req, 10)
	req.SetPathValue("lecturer_id", "5")

	rec := httptest.NewRecorder()

	h.DeleteLecturer(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteLecturer_NotFound(t *testing.T) {
	uc := &deleteLecturerUCStub{
		executeFn: func(ctx context.Context, req deletelecturer.DeleteLecturerRequest) (deletelecturer.DeleteLecturerResponse, error) {
			return deletelecturer.DeleteLecturerResponse{}, deletelecturer.ErrNotFound
		},
	}

	h := New(nil, uc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/lecturers/5", nil)
	req = middleware.WithUserID(req, 10)
	req.SetPathValue("lecturer_id", "5")

	rec := httptest.NewRecorder()

	h.DeleteLecturer(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteLecturer_Conflict(t *testing.T) {
	uc := &deleteLecturerUCStub{
		executeFn: func(ctx context.Context, req deletelecturer.DeleteLecturerRequest) (deletelecturer.DeleteLecturerResponse, error) {
			return deletelecturer.DeleteLecturerResponse{}, deletelecturer.ErrNotLecturer
		},
	}

	h := New(nil, uc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/lecturers/5", nil)
	req = middleware.WithUserID(req, 10)
	req.SetPathValue("lecturer_id", "5")

	rec := httptest.NewRecorder()

	h.DeleteLecturer(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestDeleteLecturer_DefaultError(t *testing.T) {
	uc := &deleteLecturerUCStub{
		executeFn: func(ctx context.Context, req deletelecturer.DeleteLecturerRequest) (deletelecturer.DeleteLecturerResponse, error) {
			return deletelecturer.DeleteLecturerResponse{}, errors.New("boom")
		},
	}

	h := New(nil, uc, nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/lecturers/5", nil)
	req = middleware.WithUserID(req, 10)
	req.SetPathValue("lecturer_id", "5")

	rec := httptest.NewRecorder()

	h.DeleteLecturer(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
