package test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"testum-engine/app/internal/handler/middleware"

	getfile "testum-engine/app/internal/service/use_case/lecturer/get_test_file"
	uploadtest "testum-engine/app/internal/service/use_case/lecturer/upload_test"

	"testum-engine/app/internal/handler/lecturer/file/dto"
)

//
// USE CASE INTERFACES
//

type UploadTestUseCase interface {
	Execute(ctx context.Context, req uploadtest.UploadTestRequest) (uploadtest.UploadTestResponse, error)
}

type GetTestFileUseCase interface {
	Execute(ctx context.Context, req getfile.GetTestFileRequest) (getfile.GetTestFileResponse, error)
}

//
// HANDLER
//

type Handler struct {
	uploadUC UploadTestUseCase
	fileUC   GetTestFileUseCase
}

func New(uploadUC UploadTestUseCase, fileUC GetTestFileUseCase) *Handler {
	return &Handler{
		uploadUC: uploadUC,
		fileUC:   fileUC,
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
	return strconv.Atoi(r.PathValue("test_id"))
}

//
// ERROR MAPPING (UPLOAD)
//

func mapUploadError(err error) (int, string) {
	switch {
	case errors.Is(err, uploadtest.ErrInvalidInput):
		return http.StatusBadRequest, "invalid input"
	case errors.Is(err, uploadtest.ErrAccessDenied):
		return http.StatusForbidden, "access denied"
	case errors.Is(err, uploadtest.ErrStorageFailed):
		return http.StatusInternalServerError, "internal server error"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

//
// ERROR MAPPING (FILE)
//

func mapFileError(err error) (int, string) {
	switch {
	case errors.Is(err, getfile.ErrForbidden):
		return http.StatusForbidden, "forbidden"
	case errors.Is(err, getfile.ErrFileNotFound):
		return http.StatusNotFound, "file not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

//
// POST /lecturer/tests/upload
//

func (h *Handler) UploadTest(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// multipart form
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid multipart data")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer func() { _ = file.Close() }()

	bytes, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid file")
		return
	}

	ignoreValidation := r.FormValue("ignore_validation") == "true"

	res, err := h.uploadUC.Execute(r.Context(), uploadtest.UploadTestRequest{
		UserID:           userID,
		File:             bytes,
		FileName:         header.Filename,
		IgnoreValidation: ignoreValidation,
	})
	if err != nil {
		code, msg := mapUploadError(err)
		writeError(w, code, msg)
		return
	}

	writeJSON(w, http.StatusOK, dto.ToUploadResponse(res))
}

//
// GET /lecturer/tests/{test_id}/file
//

func (h *Handler) GetTestFile(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.fileUC.Execute(r.Context(), getfile.GetTestFileRequest{
		UserID: userID,
		TestID: testID,
	})
	if err != nil {
		code, msg := mapFileError(err)
		writeError(w, code, msg)
		return
	}

	// Get file info
	fileInfo, err := res.File.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "unable to retrieve file info")
		return
	}

	// STREAM FILE
	w.Header().Set("Content-Disposition", "attachment; filename="+res.File.Name())
	w.Header().Set("Content-Type", "application/octet-stream")

	// Serve content with file info
	http.ServeContent(w, r, res.File.Name(), fileInfo.ModTime(), res.File)
}
