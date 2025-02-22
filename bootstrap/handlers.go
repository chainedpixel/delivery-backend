package bootstrap

import "infrastructure/api/handlers"

type HandlerContainer struct {
	userCases *UseCaseContainer

	authHandler *handlers.AuthHandler
}

func NewHandlerContainer(userCases *UseCaseContainer) *HandlerContainer {
	return &HandlerContainer{
		userCases: userCases,
	}
}

func (c *HandlerContainer) Initialize() error {
	c.authHandler = handlers.NewAuthHandler(c.userCases.GetAuthUseCase())

	return nil
}

func (c *HandlerContainer) GetAuthHandler() *handlers.AuthHandler {
	return c.authHandler
}
