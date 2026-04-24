package getfinishedtest

import "time"

type GetFinishedTestRequest struct {
	UserID int
}

type GetFinishedTestResponse struct {
	FinishedTests []StudentFinishTest
}

type StudentFinishTest struct {
	ID               int
	Name             string
	LecturerName     string
	CntQuestions     int
	CntHardQuestions int
	Mark             int
	SuccessRate      float64
	DateStart        time.Time
	DateEnd          time.Time
}
