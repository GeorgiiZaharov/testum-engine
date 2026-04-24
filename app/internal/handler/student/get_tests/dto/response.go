package dto

import "time"

// Модель для активного теста
type StudentActiveTest struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	LecturerName     string     `json:"lecturer_name"`
	CntQuestions     int        `json:"cnt_questions"`
	CntHardQuestions int        `json:"cnt_hard_questions"`
	DateStart        *time.Time `json:"date_start,omitempty"`
}

type GetActiveTestResponse struct {
	ActiveTests []StudentActiveTest `json:"active_tests"`
}

// Модель для завершённого теста
type StudentFinishTest struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	LecturerName     string    `json:"lecturer_name"`
	CntQuestions     int       `json:"cnt_questions"`
	CntHardQuestions int       `json:"cnt_hard_questions"`
	Mark             int       `json:"mark"`
	SuccessRate      float64   `json:"success_rate"`
	DateStart        time.Time `json:"date_start"`
	DateEnd          time.Time `json:"date_end"`
}

type GetFinishedTestResponse struct {
	FinishedTests []StudentFinishTest `json:"finished_tests"`
}
