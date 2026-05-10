package deleteattempt

type DeleteAttemptRequest struct {
	LecturerID int
	UserID     int
	TestID     int
}

type DeleteAttemptResponse struct {
	Success bool
}
