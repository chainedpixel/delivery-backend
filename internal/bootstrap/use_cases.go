package bootstrap

import (
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/auth"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/company"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/order"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/role"
	"github.com/MarlonG1/delivery-backend/internal/application/usecases/user"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/websocket"
)

type UseCaseContainer struct {
	services *ServiceContainer

	authUseCase    ports.AuthenticatorUseCase
	userUseCase    ports.UserUseCase
	orderUseCase   ports.OrdererUseCase
	roleUseCase    ports.RolerUseCase
	companyUseCase ports.CompanyUseCase
	branchUseCase  ports.BranchUseCase
	trackerUseCase ports.TrackerUseCase

	wsHub *websocket.Hub
}

func NewUseCaseContainer(services *ServiceContainer, wsHub *websocket.Hub) *UseCaseContainer {
	return &UseCaseContainer{
		services: services,
		wsHub:    wsHub,
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
	c.companyUseCase = company.NewCompanyUseCase(c.services.GetCompanyService())
	c.branchUseCase = company.NewBranchUseCase(c.services.GetCompanyService())
	c.trackerUseCase = order.NewTrackerUseCase(c.services.GetTrackerService(), c.services.GetOrderService(), c.wsHub)

	return nil
}

func (c *UseCaseContainer) GetBranchUseCase() ports.BranchUseCase {
	return c.branchUseCase
}

func (c *UseCaseContainer) GetCompanyUseCase() ports.CompanyUseCase {
	return c.companyUseCase
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

func (c *UseCaseContainer) GetTrackerUseCase() ports.TrackerUseCase {
	return c.trackerUseCase
}
