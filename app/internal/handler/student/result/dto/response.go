package dto

type GetTestResultResponse struct {
	Mark        *int     `json:"mark"`
	SuccessRate *float64 `json:"success_rate"`
	DateStart   string   `json:"date_start"`
	DateEnd     *string  `json:"date_end"`
}
