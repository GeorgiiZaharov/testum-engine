package postbaseanswers

import "time"

type TaskAnswer struct {
	TaskID  int
	Options []int
}

type CheckResult struct {
	TrueCnt int
	Total   int
}

type TestResult struct {
	Mark        int
	SuccessRate float64
	DateStart   time.Time
	DateEnd     time.Time
}

type PostBaseAnswersResponse struct {
	Success bool
}
