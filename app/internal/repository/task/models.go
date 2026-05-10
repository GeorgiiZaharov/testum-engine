package task

type Answer struct {
	ID        int
	Text      string
	ImageURL  *string
	IsCorrect bool
}

type Task struct {
	ID       int
	Text     string
	ImageURL *string
	IsHard   bool
	Answers  []Answer
}
