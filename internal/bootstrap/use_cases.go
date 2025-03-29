package bootstrap

import (
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/auth"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/order"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/role"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/user"
)

type UseCaseContainer struct {
	services *ServiceContainer

	authUseCase  ports.AuthenticatorUseCase
	userUseCase  ports.UserUseCase
	orderUseCase ports.OrdererUseCase
	roleUseCase  ports.RolerUseCase
}

func NewUseCaseContainer(services *ServiceContainer) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
	}
}

func (c *UseCaseContainer) Initialize() error {
	c.authUseCase = auth.NewAuthUseCase(c.services.GetAuthService())
	c.userUseCase = user.NewUserProfileUseCase(c.services.GetUserService(),
		c.services.GetRoleService(),
		c.services.GetCompanyService(),
		c.services.GetTokenService(),
	)
	c.orderUseCase = order.NewOrderUseCase(c.services.GetOrderService(), c.services.GetCompanyService())
	c.roleUseCase = role.NewRolerUseCase(c.services.GetRoleService())

	return nil
}

func (c *UseCaseContainer) GetAuthUseCase() ports.AuthenticatorUseCase {
	return c.authUseCase
}

func (c *UseCaseContainer) GetUserUseCase() ports.UserUseCase {
	return c.userUseCase
}

func (c *UseCaseContainer) GetOrderUseCase() ports.OrdererUseCase {
	return c.orderUseCase
}

func (c *UseCaseContainer) GetRoleUseCase() ports.RolerUseCase {
	return c.roleUseCase
}
