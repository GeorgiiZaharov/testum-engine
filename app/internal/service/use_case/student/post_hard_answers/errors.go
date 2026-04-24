package posthardanswers

import "errors"

var (
	ErrAccessDenied        = errors.New("access denied")
	ErrAlreadySubmitted    = errors.New("hard answers already submitted")
	ErrHardBlockNotPassed  = errors.New("hard block not passed")
	ErrTestAlreadyFinished = errors.New("test already finished")
	ErrFailedToSaveAnswers = errors.New("failed to save hard answers")
	ErrFailedToFinishTest  = errors.New("failed to finish test")
	ErrAnswerCheckFailed   = errors.New("answer check failed")
	ErrCalculationFailed   = errors.New("result calculation failed")
	ErrInvalidInput        = errors.New("invalid input")
)
