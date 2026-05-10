package app

import (
	"net/http"

	"testum-engine/app/internal/config"

	"testum-engine/app/internal/handler/middleware"

	"github.com/rs/cors"
)

func Build(cfg config.Config) (*App, error) {
	container, err := NewContainer(cfg)
	if err != nil {
		return nil, err
	}

	authHandler := buildAuth(container)
	adminHandler := buildAdmin(container)
	lecturerFileHandler := buildLecturerFile(container)
	lecturerReadHandler := buildLecturerRead(container)
	lecturerWriteHandler := buildLecturerWrite(container)
	studentGetTestHandler := buildStudentGetTests(container)
	studentResultHandler := buildStudentResult(container)
	studentTestAttemptHandler := buildStudentTestAttempt(container)

	jwtMiddleware := middleware.JWT(*container.AuthService)
	mux := http.NewServeMux()

	// Auth endpoints
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	mux.Handle(
		"GET /auth/me",
		jwtMiddleware(http.HandlerFunc(authHandler.GetMe)),
	)

	// Admin endpoints
	mux.Handle(
		"POST /admin/lecturers",
		jwtMiddleware(http.HandlerFunc(adminHandler.CreateLecturer)),
	)

	mux.Handle(
		"DELETE /admin/lecturers/{lecturer_id}",
		jwtMiddleware(http.HandlerFunc(adminHandler.DeleteLecturer)),
	)

	mux.Handle(
		"GET /admin/lecturers",
		jwtMiddleware(http.HandlerFunc(adminHandler.GetLecturers)),
	)

	// Lecturer endpoints
	// File
	mux.Handle(
		"POST /lecturer/tests",
		jwtMiddleware(http.HandlerFunc(lecturerFileHandler.UploadTest)),
	)

	mux.Handle(
		"GET /lecturer/tests/{test_id}/file",
		jwtMiddleware(http.HandlerFunc(lecturerFileHandler.GetTestFile)),
	)

	// Read
	mux.Handle(
		"GET /lecturer/tests",
		jwtMiddleware(http.HandlerFunc(lecturerReadHandler.GetTests)),
	)

	mux.Handle(
		"GET /lecturer/tests/{test_id}",
		jwtMiddleware(http.HandlerFunc(lecturerReadHandler.GetTest)),
	)

	mux.Handle(
		"GET /lecturer/tests/{test_id}/groups",
		jwtMiddleware(http.HandlerFunc(lecturerReadHandler.GetGroups)),
	)

	mux.Handle(
		"GET /lecturer/tests/{test_id}/result",
		jwtMiddleware(http.HandlerFunc(lecturerReadHandler.GetTestResult)),
	)

	// Write
	mux.Handle(
		"DELETE /lecturer/tests/{test_id}",
		jwtMiddleware(http.HandlerFunc(lecturerWriteHandler.DeleteTest)),
	)

	mux.Handle(
		"POST /lecturer/tests/{test_id}/access",
		jwtMiddleware(http.HandlerFunc(lecturerWriteHandler.GiveAccess)),
	)

	mux.Handle(
		"DELETE /lecturer/tests/{test_id}/access",
		jwtMiddleware(http.HandlerFunc(lecturerWriteHandler.TakeAccess)),
	)

	// Student
	// Get test
	mux.Handle(
		"GET /student/tests",
		jwtMiddleware(http.HandlerFunc(studentGetTestHandler.GetActiveTests)),
	)

	mux.Handle(
		"GET /student/tests/finished",
		jwtMiddleware(http.HandlerFunc(studentGetTestHandler.GetFinishedTests)),
	)

	// Result
	mux.Handle(
		"GET /student/test/{test_id}/result",
		jwtMiddleware(http.HandlerFunc(studentResultHandler.GetTestResult)),
	)

	// Test attempt
	mux.Handle(
		"GET /student/tests/{test_id}/hard",
		jwtMiddleware(http.HandlerFunc(studentTestAttemptHandler.GetHardTasks)),
	)

	mux.Handle(
		"GET /student/tests/{test_id}/base",
		jwtMiddleware(http.HandlerFunc(studentTestAttemptHandler.GetBaseTasks)),
	)

	mux.Handle(
		"POST /student/tests/{test_id}/hard",
		jwtMiddleware(http.HandlerFunc(studentTestAttemptHandler.PostHardAnswers)),
	)

	mux.Handle(
		"POST /student/tests/{test_id}/base",
		jwtMiddleware(http.HandlerFunc(studentTestAttemptHandler.PostBaseAnswers)),
	)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:4173",
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
	}).Handler(mux)

	server := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: corsHandler,
	}

	return &App{
		server: server,
		logger: container.Logger,
	}, nil
}
