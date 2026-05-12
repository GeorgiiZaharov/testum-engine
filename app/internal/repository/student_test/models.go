package studenttest

import "time"

// ==================== INPUT ====================

type StartTestParams struct {
	UserID int
	TestID int
}

type FinishTestParams struct {
	UserID int
	TestID int
	Result TestResult
}

// ==================== DOMAIN ====================

type TestResult struct {
	Mark        *int
	SuccessRate *float64
}

// ==================== ACTIVE TEST ====================
type StudentActiveTestInfo struct {
	ID               int
	Name             string
	LecturerName     string
	CntQuestions     int
	CntHardQuestions int
	DateStart        *time.Time // nil если тест не начинали
}

// ==================== FINISHED TEST ====================
type StudentFinishedTestInfo struct {
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
