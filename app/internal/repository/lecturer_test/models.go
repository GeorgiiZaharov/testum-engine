package lecturertest

import "time"

// ==================== INPUT ====================
type Test struct {
	Name      string
	FileName  string
	HardCount int
	Tasks     []Task
	Groups    []string
}

type Task struct {
	Text     string
	ImageURL *string
	IsHard   bool
	Answers  []Answer
}

type Answer struct {
	Text      string
	ImageURL  *string
	IsCorrect bool
}

// ==================== OUTPUT ====================
type TestInfo struct {
	ID               int
	Name             string
	CntQuestions     int
	CntHardQuestions int
	FileName         string
	Groups           []string
	DateCreated      time.Time
}

type GroupInfo struct {
	GroupName    string
	MembersCount int
}
