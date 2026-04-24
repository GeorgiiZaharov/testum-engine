package result

import "time"

type TestResult struct {
	Mark        *int       `db:"mark"`
	SuccessRate *float64   `db:"success_rate"`
	DateStart   time.Time  `db:"date_start"`
	DateEnd     *time.Time `db:"date_end"`
}

type StudentResult struct {
	UserID int    `db:"user_id"`
	Name   string `db:"name"`
	Login  string `db:"login"`
	Mail   string `db:"mail"`

	Mark        *int       `db:"mark"`
	SuccessRate *float64   `db:"success_rate"`
	DateStart   *time.Time `db:"date_start"`
	DateEnd     *time.Time `db:"date_end"`
}
