package getbasetasks

type GetBaseTasksRequest struct {
	UserID int
	TestID int
}

type GetBaseTasksResponse struct {
	BaseTasks []Task
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
