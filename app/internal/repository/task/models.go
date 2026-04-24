package task

type Answer struct {
	Text      string
	ImageURL  *string
	IsCorrect bool
}

type Task struct {
	Text     string
	ImageURL *string
	IsHard   bool
	Answers  []Answer
}
