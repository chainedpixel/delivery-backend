package bootstrap

import (
	"application/ports"
	"application/usecases/auth"
	"application/usecases/user"
)

type UseCaseContainer struct {
	services *ServiceContainer

	authUseCase ports.AuthUseCase
	userUseCase ports.UserUseCase
}

func NewUseCaseContainer(services *ServiceContainer) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
	}
}

func (c *UseCaseContainer) Initialize() error {
	c.authUseCase = auth.NewAuthUseCase(c.services.GetAuthService())
	c.userUseCase = user.NewUserProfileUseCase(c.services.GetUserService())

	return nil
}

func (c *UseCaseContainer) GetAuthUseCase() ports.AuthUseCase {
	return c.authUseCase
}

func (c *UseCaseContainer) GetUserUseCase() ports.UserUseCase {
	return c.userUseCase
}
