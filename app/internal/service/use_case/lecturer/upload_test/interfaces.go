package uploadtest

import (
	"context"
	"io"

	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
	userrepo "testum-engine/app/internal/repository/user"

	latexvalidator "testum-engine/app/internal/service/core/latexvalidator"
	validation "testum-engine/app/internal/service/core/validation"
)

type userRepository interface {
	GetByID(ctx context.Context, userID int) (userrepo.User, error)
}

type lecturerTestRepository interface {
	Create(ctx context.Context, lecturerID int, test lecturertestrepo.Test) (int, error)
}

type storageAdapter interface {
	UploadFile(file io.Reader, fileName string) error
}

type fileValidationService interface {
	Validate(lines []string) (*validation.TestFromFile, []validation.FormatError, error)
}

type fileLatexValidationService interface {
	Validate(lines []string) []latexvalidator.ValidationError
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User userRepository
	Test lecturerTestRepository
}
