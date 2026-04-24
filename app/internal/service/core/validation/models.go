package validation

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

type TestFromFile struct {
	Name      string
	HardCount int
	Tasks     []Task
}

type FormatError struct {
	Error string
}
