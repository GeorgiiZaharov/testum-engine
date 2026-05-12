package posthardanswers

import (
	"context"

	answercheckserv "testum-engine/app/internal/service/core/answer"

	answerrepo "testum-engine/app/internal/repository/answer"
	studenttestrepo "testum-engine/app/internal/repository/student_test"

	"go.uber.org/zap"
)

type UseCase struct {
	factory                  repoFactory
	answerCheckService       AnswerCheckService
	resultCalculationService ResultCalculationService
	log                      *zap.Logger
}

func NewUseCase(
	factory repoFactory,
	answerCheckService AnswerCheckService,
	resultCalculationService ResultCalculationService,
	log *zap.Logger,
) *UseCase {
	return &UseCase{
		factory:                  factory,
		answerCheckService:       answerCheckService,
		resultCalculationService: resultCalculationService,
		log:                      log,
	}
}

func (uc *UseCase) Execute(
	ctx context.Context,
	userID int,
	testID int,
	answers []TaskAnswer,
) (PostHardAnswersResponse, error) {

	if userID <= 0 || testID <= 0 || len(answers) == 0 {
		return PostHardAnswersResponse{}, ErrInvalidInput
	}

	repoAnswers := convertModelToRepoAnswers(answers)

	var response PostHardAnswersResponse

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. access check
		ok, err := r.Access.HasStudentAccess(ctx, userID, testID)
		if err != nil {
			uc.log.Error("access check failed", zap.Error(err))
			return err
		}
		if !ok {
			return ErrAccessDenied
		}

		_, err = r.Result.GetStudentResult(ctx, userID, testID)

		if err != nil {
			uc.log.Error("get student result failed", zap.Error(err))
			return err
		}

		// 2. already submitted?
		existing, err := r.Answer.GetHardAnswers(ctx, userID, testID)
		if err != nil {
			uc.log.Error("get hard answers failed", zap.Error(err))
			return err
		}
		if len(existing) > 0 {
			return ErrAlreadySubmitted
		}

		// 3. correct answers
		trueAnswers, err := r.Answer.GetHardAnswersByTest(ctx, testID)
		if err != nil {
			uc.log.Error("get hard answers by test failed", zap.Error(err))
			return err
		}

		// 4. check
		check, err := uc.answerCheckService.Check(
			convertAnswerRepoToService(repoAnswers),
			convertAnswerRepoToService(trueAnswers),
		)
		if err != nil {
			uc.log.Error("faild to check answers", zap.Error(err))
			return ErrAnswerCheckFailed
		}

		// 5. save answers
		ok, err = r.Answer.SaveHardAnswers(ctx, userID, extractRepoAnswerIDs(repoAnswers))
		if err != nil || !ok {
			uc.log.Error("faild to save answers", zap.Error(err))
			return ErrFailedToSaveAnswers
		}

		response.IsAllCorrect = (check.TrueCnt == check.Total)

		// 6. if all correct → finish test
		if response.IsAllCorrect {

			res := uc.resultCalculationService.Calc(check)

			_, err = r.StudentTest.FinishTest(ctx, studenttestrepo.FinishTestParams{
				UserID: userID,
				TestID: testID,
				Result: studenttestrepo.TestResult{
					Mark:        &res.Mark,
					SuccessRate: &res.SuccessRate,
				},
			})
			if err != nil {
				uc.log.Error("finish test failed", zap.Error(err))
				return ErrFailedToFinishTest
			}
		}

		return nil
	})

	if err != nil {
		uc.log.Error("faild to post hard answers", zap.Error(err))
		return PostHardAnswersResponse{}, err
	}

	return response, nil
}

// --------------------
// adapters
// --------------------

func convertAnswerRepoToService(in []answerrepo.TaskAnswer) []answercheckserv.TaskAnswer {
	out := make([]answercheckserv.TaskAnswer, 0, len(in))

	for _, a := range in {
		out = append(out, answercheckserv.TaskAnswer{
			TaskID:          a.TaskID,
			SelectedOptions: a.Options,
		})
	}

	return out
}

func convertModelToRepoAnswers(in []TaskAnswer) []answerrepo.TaskAnswer {
	out := make([]answerrepo.TaskAnswer, 0, len(in))

	for _, a := range in {
		out = append(out, answerrepo.TaskAnswer{
			TaskID:  a.TaskID,
			Options: a.Options,
		})
	}

	return out
}

func extractRepoAnswerIDs(in []answerrepo.TaskAnswer) []int {
	out := make([]int, 0)

	for _, t := range in {
		out = append(out, t.Options...)
	}

	return out
}
