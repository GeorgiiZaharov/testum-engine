package app

import (
	authhandler "testum-engine/app/internal/handler/auth"

	getmeuc "testum-engine/app/internal/service/use_case/auth/get_me"
	loginuc "testum-engine/app/internal/service/use_case/auth/login"
	refreshuc "testum-engine/app/internal/service/use_case/auth/refresh"
)

func buildAuth(container *Container) *authhandler.Handler {
	loginUC := loginuc.NewUseCase(
		loginuc.NewFactory(container.DB, container.Logger),
		*container.LDAP,
		*container.AuthService,
		container.Logger,
	)

	refreshUC := refreshuc.NewUseCase(
		refreshuc.NewFactory(container.DB, container.Logger),
		*container.AuthService,
		container.Logger,
	)

	getMeUC := getmeuc.NewUseCase(
		getmeuc.NewFactory(container.DB, container.Logger),
		*container.LDAP,
		container.Logger,
	)

	return authhandler.New(
		loginUC,
		refreshUC,
		getMeUC,
	)
}
