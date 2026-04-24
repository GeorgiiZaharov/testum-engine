package admin

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

	"testum-engine/app/internal/handler/admin/dto"
	"testum-engine/app/internal/handler/middleware"

	createlecturer "testum-engine/app/internal/service/use_case/admin/create_lecturer"
)

type createLecturerUCStub struct {
	executeFn func(ctx context.Context, req createlecturer.CreateLecturerRequest) (createlecturer.CreateLecturerResponse, error)
}

func (s *createLecturerUCStub) Execute(ctx context.Context, req createlecturer.CreateLecturerRequest) (createlecturer.CreateLecturerResponse, error) {
	return s.executeFn(ctx, req)
}

func TestCreateLecturer_Success(t *testing.T) {
	uc := &createLecturerUCStub{
		executeFn: func(ctx context.Context, req createlecturer.CreateLecturerRequest) (createlecturer.CreateLecturerResponse, error) {

			assert.Equal(t, 10, req.AdminID)
			assert.Equal(t, "john", req.Login)

			return createlecturer.CreateLecturerResponse{
				Success: true,
			}, nil
		},
	}

	h := New(uc, nil, nil)

	body := dto.CreateLecturerRequest{
		Login: " john ",
	}

	b, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/admin/lecturers", bytes.NewReader(b))
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.CreateLecturer(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp dto.CreateLecturerResponse
	err = json.NewDecoder(rec.Body).Decode(&resp)

	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestCreateLecturer_Unauthorized(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/admin/lecturers", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()

	h.CreateLecturer(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateLecturer_InvalidJSON(t *testing.T) {
	h := New(nil, nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/admin/lecturers", bytes.NewBufferString(`{invalid-json}`))
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.CreateLecturer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateLecturer_EmptyLogin(t *testing.T) {
	h := New(nil, nil, nil)

	body := dto.CreateLecturerRequest{
		Login: "   ",
	}

	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/admin/lecturers", bytes.NewReader(b))
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.CreateLecturer(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateLecturer_ForbiddenError(t *testing.T) {
	uc := &createLecturerUCStub{
		executeFn: func(ctx context.Context, req createlecturer.CreateLecturerRequest) (createlecturer.CreateLecturerResponse, error) {
			return createlecturer.CreateLecturerResponse{}, createlecturer.ErrForbidden
		},
	}

	h := New(uc, nil, nil)

	body, _ := json.Marshal(dto.CreateLecturerRequest{Login: "john"})

	req := httptest.NewRequest(http.MethodPost, "/admin/lecturers", bytes.NewReader(body))
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.CreateLecturer(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCreateLecturer_DefaultError(t *testing.T) {
	uc := &createLecturerUCStub{
		executeFn: func(ctx context.Context, req createlecturer.CreateLecturerRequest) (createlecturer.CreateLecturerResponse, error) {
			return createlecturer.CreateLecturerResponse{}, errors.New("boom")
		},
	}

	h := New(uc, nil, nil)

	body, _ := json.Marshal(dto.CreateLecturerRequest{Login: "john"})

	req := httptest.NewRequest(http.MethodPost, "/admin/lecturers", bytes.NewReader(body))
	req = middleware.WithUserID(req, 10)

	rec := httptest.NewRecorder()

	h.CreateLecturer(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

