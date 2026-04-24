package gettest

import "time"

type GetTestRequest struct {
	UserID int
	TestID int
}

type GetTestResponse struct {
	ID               int
	Name             string
	CntQuestions     int
	CntHardQuestions int

	Groups    []GroupInfo
	HardTasks []Task
	BaseTasks []Task
}

// ---- repo DTO reuse ----

type TestInfo struct {
	ID               int
	Name             string
	CntQuestions     int
	CntHardQuestions int
	Groups           []int
	DateCreated      time.Time
}

type GroupInfo struct {
	GroupName    string
	MembersCount int
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
