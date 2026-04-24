package dto

import (
	getme "testum-engine/app/internal/service/use_case/auth/get_me"
	login "testum-engine/app/internal/service/use_case/auth/login"
	refresh "testum-engine/app/internal/service/use_case/auth/refresh"
)

func ToLoginResponse(res login.AuthLoginResponse) LoginResponse {
	return LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
}

func ToRefreshResponse(res refresh.AuthRefreshResponse) RefreshResponse {
	return RefreshResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
}

func ToGetMeResponse(res getme.GetMeResponse) GetMeResponse {
	return GetMeResponse{
		ID:         res.ID,
		Login:      res.Login,
		Mail:       res.Mail,
		Name:       res.Name,
		Group:      res.Group,
		IsLecturer: res.IsLecturer,
	}
}
