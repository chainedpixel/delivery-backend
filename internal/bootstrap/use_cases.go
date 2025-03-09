package bootstrap

import (
	"application/ports"
	"application/usecases/auth"
	"application/usecases/order"
	"application/usecases/user"
)

type UseCaseContainer struct {
	services *ServiceContainer

	authUseCase  ports.AuthUseCase
	userUseCase  ports.UserUseCase
	orderUseCase ports.OrdererUseCase
}

func NewUseCaseContainer(services *ServiceContainer) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
	}
}

func (c *UseCaseContainer) Initialize() error {
	c.authUseCase = auth.NewAuthUseCase(c.services.GetAuthService())
	c.userUseCase = user.NewUserProfileUseCase(c.services.GetUserService())
	c.orderUseCase = order.NewOrderUseCase(c.services.GetOrderService(), c.services.GetCompanyService())

	return nil
}

func (c *UseCaseContainer) GetAuthUseCase() ports.AuthUseCase {
	return c.authUseCase
}

func (c *UseCaseContainer) GetUserUseCase() ports.UserUseCase {
	return c.userUseCase
}

func (c *UseCaseContainer) GetOrderUseCase() ports.OrdererUseCase {
	return c.orderUseCase
}
