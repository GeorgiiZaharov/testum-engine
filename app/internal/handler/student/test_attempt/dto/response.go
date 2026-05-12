package dto

type GetTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

type TaskResponse struct {
	ID               int              `json:"id"`
	Text             string           `json:"text"`
	ImageURL         *string          `json:"image_url,omitempty"`
	IsHard           bool             `json:"is_hard"`
	IsMultipleChoice bool             `json:"is_multiple_choice"`
	Answers          []AnswerResponse `json:"answers"`
}

type AnswerResponse struct {
	ID       int     `json:"id"`
	Text     string  `json:"text"`
	ImageURL *string `json:"image_url,omitempty"`
}

type PostHardAnswersResponse struct {
	IsAllCorrect bool `json:"is_all_correct"`
}

type PostBaseAnswersResponse struct {
	Success bool `json:"success"`
}
