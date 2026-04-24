package takeaccess

type TakeAccessRequest struct {
	UserID int
	TestID int
	Group  string
}

type TakeAccessResponse struct {
	Success bool
}
