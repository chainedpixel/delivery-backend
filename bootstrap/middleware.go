package bootstrap

import "infrastructure/api/middleware"

type MiddlewareContainer struct {
	services *ServiceContainer

	errMiddleware  *middleware.ErrorMiddleware
	authMiddleware *middleware.AuthMiddleware
	tokenExtractor *middleware.TokenExtractor
}

func NewMiddlewareContainer(services *ServiceContainer) *MiddlewareContainer {
	return &MiddlewareContainer{
		services: services,
	}
}

func (c *MiddlewareContainer) Initialize() error {
	c.errMiddleware = middleware.NewErrorMiddleware()
	c.authMiddleware = middleware.NewAuthMiddleware(c.services.GetTokenService())
	c.tokenExtractor = middleware.NewTokenExtractor()

	return nil
}

func (c *MiddlewareContainer) GetErrorMiddleware() *middleware.ErrorMiddleware {
	return c.errMiddleware
}

func (c *MiddlewareContainer) GetAuthMiddleware() *middleware.AuthMiddleware {
	return c.authMiddleware
}

func (c *MiddlewareContainer) GetTokenExtractor() *middleware.TokenExtractor {
	return c.tokenExtractor
}
