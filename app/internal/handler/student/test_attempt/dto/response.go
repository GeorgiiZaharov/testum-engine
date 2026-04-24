package dto

type GetTasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

type TaskResponse struct {
	Text     string           `json:"text"`
	ImageURL *string          `json:"image_url,omitempty"`
	IsHard   bool             `json:"is_hard"`
	Answers  []AnswerResponse `json:"answers"`
}

type AnswerResponse struct {
	Text     string  `json:"text"`
	ImageURL *string `json:"image_url,omitempty"`
}

type PostHardAnswersResponse struct {
	Success      bool `json:"success"`
	IsAllCorrect bool `json:"is_all_correct"`
}

type PostBaseAnswersResponse struct {
	Success bool `json:"success"`
}
