package getme

type GetMeRequest struct {
	UserID int
}

type GetMeResponse struct {
	ID         int
	Login      string
	Mail       string
	Name       string
	Group      *string
	IsLecturer bool
}
