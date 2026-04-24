package posthardanswers

import "time"

type TaskAnswer struct {
	TaskID  int
	Options []int
}

type CheckResult struct {
	TrueCnt int
	Total   int
}

type CalcResult struct {
	Mark        int
	SuccessRate float64
}

type PostHardAnswersResponse struct {
	IsAllCorrect bool
}

type TestResult struct {
	Mark        int
	SuccessRate float64
	DateStart   time.Time
	DateEnd     time.Time
}
