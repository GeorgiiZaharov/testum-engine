package postbaseanswers

import "errors"

var (
	ErrAccessDenied        = errors.New("access denied")
	ErrAlreadySubmitted    = errors.New("base answers already submitted")
	ErrHardBlockNotPassed  = errors.New("hard block not passed")
	ErrTestAlreadyFinished = errors.New("test already finished")
	ErrFailedToSaveAnswers = errors.New("failed to save base answers")
	ErrFailedToFinishTest  = errors.New("failed to finish test")
	ErrAnswerCheckFailed   = errors.New("answer check failed")
	ErrCalculationFailed   = errors.New("result calculation failed")
	ErrInvalidInput        = errors.New("invalid input")
)
