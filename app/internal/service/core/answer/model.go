package answer

type TaskAnswer struct {
	TaskID          int
	SelectedOptions []int
}

type CheckResult struct {
	TrueCnt int
	Total   int
}
