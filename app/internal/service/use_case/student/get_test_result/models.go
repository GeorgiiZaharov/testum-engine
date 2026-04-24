package gettestresult

import "time"

type GetTestResultRequest struct {
	UserID int
	TestID int
}

type GetTestResultResponse struct {
	Mark        *int
	SuccessRate *float64
	DateStart   time.Time
	DateEnd     *time.Time
}
