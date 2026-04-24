package giveaccess

type GiveAccessRequest struct {
	UserID int
	TestID int
	Group  string
}

type GiveAccessResponse struct {
	Success bool
}
