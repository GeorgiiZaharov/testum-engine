package login

type AuthLoginRequest struct {
	Login    string
	Password string
}

type AuthLoginResponse struct {
	AccessToken  string
	RefreshToken string
}
