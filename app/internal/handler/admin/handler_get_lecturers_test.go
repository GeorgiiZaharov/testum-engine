package admin

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"testum-engine/app/internal/handler/middleware"
	getlecturers "testum-engine/app/internal/service/use_case/admin/get_lecturers"
)

type getLecturersUCStub struct {
	executeFn func(ctx context.Context, req getlecturers.GetLecturersRequest) (getlecturers.GetLecturersResponse, error)
}

func (s *getLecturersUCStub) Execute(ctx context.Context, req getlecturers.GetLecturersRequest) (getlecturers.GetLecturersResponse, error) {
	return s.executeFn(ctx, req)
}

func TestGetLecturers_Success(t *testing.T) {
	uc := &getLecturersUCStub{
		executeFn: func(ctx context.Context, req getlecturers.GetLecturersRequest) (getlecturers.GetLecturersResponse, error) {

			assert.Equal(t, 10, req.UserID)

			return getlecturers.GetLecturersResponse{
				Lecturers: []getlecturers.Lecturer{
					{ID: 1, Login: "john"},
				},
			}, nil
		},
	}

	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/admin/lecturers", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetLecturers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetLecturers_Unauthorized(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/lecturers", nil)
	rec := httptest.NewRecorder()

	h.GetLecturers(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetLecturers_Forbidden(t *testing.T) {
	uc := &getLecturersUCStub{
		executeFn: func(ctx context.Context, req getlecturers.GetLecturersRequest) (getlecturers.GetLecturersResponse, error) {
			return getlecturers.GetLecturersResponse{}, getlecturers.ErrForbidden
		},
	}

	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/admin/lecturers", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetLecturers(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGetLecturers_DefaultError(t *testing.T) {
	uc := &getLecturersUCStub{
		executeFn: func(ctx context.Context, req getlecturers.GetLecturersRequest) (getlecturers.GetLecturersResponse, error) {
			return getlecturers.GetLecturersResponse{}, errors.New("boom")
		},
	}

	h := New(nil, nil, uc)

	req := httptest.NewRequest(http.MethodGet, "/admin/lecturers", nil)
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.GetLecturers(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
