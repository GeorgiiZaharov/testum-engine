package dto

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetMeResponse struct {
	ID         int     `json:"id"`
	Login      string  `json:"login"`
	Mail       string  `json:"mail"`
	Name       string  `json:"name"`
	Group      *string `json:"group,omitempty"`
	IsLecturer bool    `json:"is_lecturer"`
}
