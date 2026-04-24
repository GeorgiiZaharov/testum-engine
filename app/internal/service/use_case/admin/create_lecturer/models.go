package createlecturer

type CreateLecturerRequest struct {
	AdminID int
	Login   string
}

type CreateLecturerResponse struct {
	Success bool
}
