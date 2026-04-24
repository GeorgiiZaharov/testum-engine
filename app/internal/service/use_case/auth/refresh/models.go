package refresh

type AuthRefreshRequest struct {
	UserID int
}

type AuthRefreshResponse struct {
	AccessToken  string
	RefreshToken string
}
