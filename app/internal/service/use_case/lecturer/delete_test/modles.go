package deletetest

type DeleteTestRequest struct {
	UserID int
	TestID int
}

type DeleteTestResponse struct {
	Success bool
}
