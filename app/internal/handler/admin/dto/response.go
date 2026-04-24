package dto

type CreateLecturerResponse struct {
	Success bool `json:"success"`
}

type DeleteLecturerResponse struct {
	Success bool `json:"success"`
}

type GetLecturersResponse struct {
	Lecturers []Lecturer `json:"lecturers"`
}

type Lecturer struct {
	ID           int     `json:"id"`
	Login        string  `json:"login"`
	Mail         string  `json:"mail"`
	Name         string  `json:"name"`
	Group        *string `json:"group,omitempty"`
	DateCreated  string  `json:"date_created"`
	DateModified string  `json:"date_modified"`
}
