package dto

import "time"

type GetTestsResponse struct {
	Tests []TestInfo `json:"tests"`
}

type TestInfo struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	CntQuestions     int       `json:"cnt_questions"`
	CntHardQuestions int       `json:"cnt_hard_questions"`
	Groups           []string  `json:"groups"`
	DateCreated      time.Time `json:"date_created"`
}

type GetTestResponse struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	CntQuestions     int     `json:"cnt_questions"`
	CntHardQuestions int     `json:"cnt_hard_questions"`
	Groups           []Group `json:"groups"`
	HardTasks        []Task  `json:"hard_tasks"`
	BaseTasks        []Task  `json:"base_tasks"`
}

type Group struct {
	GroupName    string `json:"group_name"`
	MembersCount int    `json:"members_count"`
}

type Task struct {
	Text    string   `json:"text"`
	Image   *string  `json:"image,omitempty"`
	Answers []Answer `json:"answers"`
	IsHard  bool     `json:"is_hard"`
}

type Answer struct {
	Text      string  `json:"text"`
	Image     *string `json:"image,omitempty"`
	IsCorrect bool    `json:"is_correct"`
}

type GetGroupsResponse struct {
	Groups []Group `json:"groups"`
}

type GetTestResultResponse struct {
	Results []StudentResult `json:"results"`
}

type StudentResult struct {
	StudentID int      `json:"student_id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Score     *float64 `json:"score,omitempty"`
	Mark      *int     `json:"mark,omitempty"`
}
