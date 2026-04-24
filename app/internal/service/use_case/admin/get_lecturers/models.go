package getlecturers

type GetLecturersRequest struct {
	UserID int
}

type GetLecturersResponse struct {
	Lecturers []Lecturer
}

type Lecturer struct {
	ID           int
	Login        string
	Mail         string
	Name         string
	Group        *string
	DateCreated  string
	DateModified string
}
