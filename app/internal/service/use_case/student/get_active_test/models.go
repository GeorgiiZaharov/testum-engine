package getactivetest

import "time"

type GetActiveTestRequest struct {
	UserID int
}

type GetActiveTestResponse struct {
	ActiveTests []StudentActiveTest
}

type StudentActiveTest struct {
	ID               int
	Name             string
	LecturerName     string
	CntQuestions     int
	CntHardQuestions int
	DateStart        *time.Time
}
