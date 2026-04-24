package postbaseanswers

import (
	"context"

	"go.uber.org/zap"

	answercheckserv "testum-engine/app/internal/service/core/answer"

	answerrepo "testum-engine/app/internal/repository/answer"
	studenttestrepo "testum-engine/app/internal/repository/student_test"
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
) (PostBaseAnswersResponse, error) {

	if userID <= 0 || testID <= 0 || len(answers) == 0 {
		return PostBaseAnswersResponse{}, ErrInvalidInput
	}

	var response PostBaseAnswersResponse

	repoAnswers := convertModelToRepoAnswers(answers)

	err := uc.factory.WithTx(ctx, func(r repositories) error {

		// 1. Access check
		ok, err := r.Access.HasStudentAccess(ctx, userID, testID)
		if err != nil {
			uc.log.Error("access check failed", zap.Error(err))
			return err
		}
		if !ok {
			return ErrAccessDenied
		}

		// 2. hard block must exist
		hardStudentAnswers, err := r.Answer.GetHardAnswers(ctx, userID, testID)
		if err != nil {
			uc.log.Error("get hard answers failed", zap.Error(err))
			return err
		}
		if len(hardStudentAnswers) == 0 {
			return ErrHardBlockNotPassed
		}

		// 3. already submitted base answers?
		existingBase, err := r.Answer.GetBaseAnswers(ctx, userID, testID)
		if err != nil {
			uc.log.Error("get base answers failed", zap.Error(err))
			return err
		}
		if len(existingBase) > 0 {
			return ErrAlreadySubmitted
		}

		// 4. test already finished?
		resp, err := r.Result.GetStudentResult(ctx, userID, testID)
		if err != nil {
			return err
		}
		if resp.DateEnd != nil {
			return ErrTestAlreadyFinished
		}

		// 5. correct answers
		baseTrue, err := r.Answer.GetBaseAnswersByTest(ctx, testID)
		if err != nil {
			return err
		}

		hardTrue, err := r.Answer.GetHardAnswersByTest(ctx, testID)
		if err != nil {
			return err
		}

		// 6. check answers
		baseCheck, err := uc.answerCheckService.Check(
			convertAnswerRepoToService(hardStudentAnswers),
			convertAnswerRepoToService(hardTrue),
		)
		if err != nil {
			return ErrAnswerCheckFailed
		}

		hardCheck, err := uc.answerCheckService.Check(
			convertAnswerRepoToService(repoAnswers),
			convertAnswerRepoToService(baseTrue),
		)
		if err != nil {
			return ErrAnswerCheckFailed
		}

		// 7. aggregate
		total := answercheckserv.CheckResult{
			TrueCnt: baseCheck.TrueCnt + hardCheck.TrueCnt,
			Total:   baseCheck.Total + hardCheck.Total,
		}

		// 8. calc result
		final := uc.resultCalculationService.Calc(answercheckserv.CheckResult{
			TrueCnt: total.TrueCnt,
			Total:   total.Total,
		})

		// 9. save base answers
		ok, err = r.Answer.SaveBaseAnswers(ctx, userID, extractRepoAnswerIDs(repoAnswers))
		if err != nil || !ok {
			return ErrFailedToSaveAnswers
		}

		// 10. finish test
		_, err = r.StudentTest.FinishTest(ctx, studenttestrepo.FinishTestParams{
			UserID: userID,
			TestID: testID,
			Result: studenttestrepo.TestResult{
				Mark:        &final.Mark,
				SuccessRate: &final.SuccessRate,
			},
		})
		if err != nil {
			uc.log.Error("finish test failed", zap.Error(err))
			return ErrFailedToFinishTest
		}

		response.Success = true
		return nil
	})

	if err != nil {
		return PostBaseAnswersResponse{}, err
	}

	return response, nil
}

// --------------------
// adapters (repo -> service DTO)
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
