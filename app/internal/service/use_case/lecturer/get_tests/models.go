package gettests

import "time"

type GetTestsRequest struct {
	UserID int
}

type GetTestsResponse struct {
	Tests []TestInfo
}

type TestInfo struct {
	ID               int
	Name             string
	CntQuestions     int
	CntHardQuestions int
	Groups           []string
	DateCreated      time.Time
}
