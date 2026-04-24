package dto

type PostAnswersRequest struct {
	Answers []TaskAnswerRequest `json:"answers"`
}

type TaskAnswerRequest struct {
	TaskID  int   `json:"task_id"`
	Options []int `json:"options"`
}
