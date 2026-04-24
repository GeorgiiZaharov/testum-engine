package answer

type TaskAnswer struct {
	TaskID  int   `db:"task_id"`
	Options []int `db:"options"`
}

type SaveAnswersParams struct {
	UserID  int
	TestID  int
	Answers []TaskAnswer
}
