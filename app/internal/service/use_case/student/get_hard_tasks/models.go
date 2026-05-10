package gethardtasks

type GetHardTasksRequest struct {
	UserID int
	TestID int
}

type GetHardTasksResponse struct {
	HardTasks []Task
}

type Task struct {
	ID       int
	Text     string
	ImageURL *string
	IsHard   bool
	Answers  []Answer
}

type Answer struct {
	ID       int
	Text     string
	ImageURL *string
}
