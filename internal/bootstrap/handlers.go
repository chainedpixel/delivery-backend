package bootstrap

import "github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"

type HandlerContainer struct {
	usesCases *UseCaseContainer

	authHandler  *handlers.AuthHandler
	userHandler  *handlers.UserHandler
	orderHandler *handlers.OrderHandler
	roleHandler  *handlers.RoleHandler
}

func NewHandlerContainer(userCases *UseCaseContainer) *HandlerContainer {
	return &HandlerContainer{
		usesCases: userCases,
	}
}

func (c *HandlerContainer) Initialize() error {
	c.authHandler = handlers.NewAuthHandler(c.usesCases.GetAuthUseCase())
	c.userHandler = handlers.NewUserHandler(c.usesCases.GetUserUseCase())
	c.orderHandler = handlers.NewOrderHandler(c.usesCases.GetOrderUseCase())
	c.roleHandler = handlers.NewRoleHandler(c.usesCases.GetRoleUseCase())

	return nil
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
