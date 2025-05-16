package bootstrap

import "github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"

type HandlerContainer struct {
	usesCases *UseCaseContainer
	services  *ServiceContainer

	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	orderHandler   *handlers.OrderHandler
	roleHandler    *handlers.RoleHandler
	companyHandler *handlers.CompanyHandler
	branchHandler  *handlers.BranchHandler
	trackerHandler *handlers.TrackerHandler
}

func NewHandlerContainer(userCases *UseCaseContainer, services *ServiceContainer) *HandlerContainer {
	return &HandlerContainer{
		usesCases: userCases,
		services:  services,
	}
}

func (c *HandlerContainer) Initialize() error {
	c.authHandler = handlers.NewAuthHandler(c.usesCases.GetAuthUseCase())
	c.userHandler = handlers.NewUserHandler(c.usesCases.GetUserUseCase())
	c.orderHandler = handlers.NewOrderHandler(c.usesCases.GetOrderUseCase())
	c.roleHandler = handlers.NewRoleHandler(c.usesCases.GetRoleUseCase())
	c.companyHandler = handlers.NewCompanyHandler(c.usesCases.GetCompanyUseCase())
	c.branchHandler = handlers.NewBranchHandler(c.usesCases.GetBranchUseCase())
	c.trackerHandler = handlers.NewTrackerHandler(c.usesCases.GetTrackerUseCase())

	return nil
}

func (c *HandlerContainer) GetBranchHandler() *handlers.BranchHandler {
	return c.branchHandler
}

func (c *HandlerContainer) GetCompanyHandler() *handlers.CompanyHandler {
	return c.companyHandler
}

func (c *HandlerContainer) GetAuthHandler() *handlers.AuthHandler {
	return c.authHandler
}

func (c *HandlerContainer) GetUserHandler() *handlers.UserHandler {
	return c.userHandler
}

func (c *HandlerContainer) GetOrderHandler() *handlers.OrderHandler {
	return c.orderHandler
}

func (c *HandlerContainer) GetRoleHandler() *handlers.RoleHandler {
	return c.roleHandler
}

func (c *HandlerContainer) GetTrackerHandler() *handlers.TrackerHandler {
	return c.trackerHandler
}
