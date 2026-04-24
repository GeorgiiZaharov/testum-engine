package posthardanswers

import (
	"context"

	answercheckserv "testum-engine/app/internal/service/core/answer"
	resultcalcserv "testum-engine/app/internal/service/core/result"

	answerrepo "testum-engine/app/internal/repository/answer"
	resultrepo "testum-engine/app/internal/repository/result"
	studenttestrepo "testum-engine/app/internal/repository/student_test"
	userrepo "testum-engine/app/internal/repository/user"
)

type UserRepository interface {
	GetByID(ctx context.Context, userID int) (userrepo.User, error)
}

type AccessRepository interface {
	HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
}

type AnswerRepository interface {
	GetHardAnswers(ctx context.Context, userID int, testID int) ([]answerrepo.TaskAnswer, error)
	GetHardAnswersByTest(ctx context.Context, testID int) ([]answerrepo.TaskAnswer, error)

	SaveHardAnswers(ctx context.Context, userID int, answers []int) (bool, error)
}

type StudentTestRepository interface {
	FinishTest(ctx context.Context, result studenttestrepo.FinishTestParams) (bool, error)
	StartTest(ctx context.Context, params studenttestrepo.StartTestParams) (bool, error)
}

type ResultRepository interface {
	GetStudentResult(ctx context.Context, userID int, testID int) (resultrepo.TestResult, error)
}

type AnswerCheckService interface {
	Check(
		studentAnswers []answercheckserv.TaskAnswer,
		trueAnswers []answercheckserv.TaskAnswer,
	) (answercheckserv.CheckResult, error)
}

type ResultCalculationService interface {
	Calc(res answercheckserv.CheckResult) resultcalcserv.CalcResult
}

type repoFactory interface {
	WithTx(ctx context.Context, fn func(r repositories) error) error
}

type repositories struct {
	User        UserRepository
	Access      AccessRepository
	Answer      AnswerRepository
	Result      ResultRepository
	StudentTest StudentTestRepository
}
