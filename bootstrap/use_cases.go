package bootstrap

import (
	"application/ports"
	"application/usecases/auth"
)

type UseCaseContainer struct {
	services *ServiceContainer

	userUseCase ports.AuthUseCase
}

func NewUseCaseContainer(services *ServiceContainer) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
	}
}

func (c *UseCaseContainer) Initialize() error {
	c.userUseCase = auth.NewAuthUseCase(c.services.GetAuthService())

	return nil
}

func (c *UseCaseContainer) GetAuthUseCase() ports.AuthUseCase {
	return c.userUseCase
}
