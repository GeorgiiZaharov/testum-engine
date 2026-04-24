package uploadtest

import (
	"bufio"
	"bytes"
	"context"

	"go.uber.org/zap"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
	latexvalidator "testum-engine/app/internal/service/core/latexvalidator"
	validation "testum-engine/app/internal/service/core/validation"
)

type UseCase struct {
	factory        repoFactory
	storage        storageAdapter
	validator      fileValidationService
	latexValidator fileLatexValidationService
	log            *zap.Logger
}

func NewUseCase(
	factory repoFactory,
	storage storageAdapter,
	validator fileValidationService,
	latexValidator fileLatexValidationService,
	log *zap.Logger,
) *UseCase {
	return &UseCase{
		factory:        factory,
		storage:        storage,
		validator:      validator,
		latexValidator: latexValidator,
		log:            log,
	}
}

func (uc *UseCase) Execute(ctx context.Context, req UploadTestRequest) (UploadTestResponse, error) {
	if req.UserID <= 0 || len(req.File) == 0 || req.FileName == "" {
		return UploadTestResponse{}, ErrInvalidInput
	}

	var response UploadTestResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Проверка пользователя
		user, err := r.User.GetByID(ctx, req.UserID)
		if err != nil {
			uc.log.Error("failed to get user", zap.Error(err))
			return err
		}

		if !user.IsLecturer {
			return ErrAccessDenied
		}

		// 2. Чтение файла → []string
		lines, err := readLines(req.File)
		if err != nil {
			uc.log.Error("failed to read file", zap.Error(err))
			return err
		}

		// 3. Валидация структуры и формата
		testFromFile, formatErrs, err := uc.validator.Validate(lines)
		if err != nil {
			uc.log.Error("validation service failed", zap.Error(err))
			return err
		}

		if len(formatErrs) > 0 && !req.IgnoreValidation {
			response.FormatErrors = mapFormatErrors(formatErrs)
			return nil
		}

		if testFromFile == nil {
			uc.log.Error("validator returned nil test")
			return ErrInvalidInput
		}

		// 4. Latex валидация
		if !req.IgnoreValidation {
			latexErrs := uc.latexValidator.Validate(lines)

			if len(latexErrs) > 0 {
				response.ValidationErrors = mapLatexErrors(latexErrs)
				return nil
			}
		}

		// 5. Upload файла
		if err := uc.storage.UploadFile(bytes.NewReader(req.File), req.FileName); err != nil {
			uc.log.Error("failed to upload file", zap.Error(err))
			return ErrStorageFailed
		}

		// 6. Маппинг в repo DTO
		test := mapToRepoTest(testFromFile, req.FileName)

		// 7. Сохранение теста
		testID, err := r.Test.Create(ctx, req.UserID, test)
		if err != nil {
			uc.log.Error("failed to create test", zap.Error(err))
			return err
		}

		response.TestID = &testID
		response.Success = true

		return nil
	})

	if err != nil {
		return UploadTestResponse{}, err
	}

	return response, nil
}

// =========================
// helpers
// =========================

func readLines(data []byte) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func mapFormatErrors(src []validation.FormatError) []FormatError {
	res := make([]FormatError, 0, len(src))

	for _, e := range src {
		res = append(res, FormatError{
			Error: e.Error,
		})
	}

	return res
}

func mapLatexErrors(src []latexvalidator.ValidationError) []ValidationError {
	res := make([]ValidationError, 0, len(src))

	for _, e := range src {
		res = append(res, ValidationError{
			Line:  e.Line,
			Error: e.Error,
		})
	}

	return res
}

func mapToRepoTest(src *validation.TestFromFile, fileName string) lecturertestrepo.Test {
	tasks := make([]lecturertestrepo.Task, 0, len(src.Tasks))

	for _, t := range src.Tasks {
		answers := make([]lecturertestrepo.Answer, 0, len(t.Answers))

		for _, a := range t.Answers {
			answers = append(answers, lecturertestrepo.Answer{
				Text:      a.Text,
				ImageURL:  a.ImageURL,
				IsCorrect: a.IsCorrect,
			})
		}

		tasks = append(tasks, lecturertestrepo.Task{
			Text:     t.Text,
			ImageURL: t.ImageURL,
			IsHard:   t.IsHard,
			Answers:  answers,
		})
	}

	return lecturertestrepo.Test{
		Name:      src.Name,
		FileName:  fileName,
		Tasks:     tasks,
		HardCount: src.HardCount,
	}
}
