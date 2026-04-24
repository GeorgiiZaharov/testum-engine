package deletelecturer

type DeleteLecturerRequest struct {
	AdminID    int
	LecturerID int
}

type DeleteLecturerResponse struct {
	Success bool
}
