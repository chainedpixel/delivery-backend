package bootstrap

import "github.com/MarlonG1/delivery-backend/internal/infrastructure/api/middleware"

type MiddlewareContainer struct {
	services *ServiceContainer

	errMiddleware  *middleware.ErrorMiddleware
	authMiddleware *middleware.AuthMiddleware
	tokenExtractor *middleware.TokenExtractor
	corsMiddleware *middleware.CorsMiddleware
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
	c.corsMiddleware = middleware.NewCorsMiddleware(
		[]string{"*"},
		nil,
		nil,
	)

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

func (c *MiddlewareContainer) GetCorsMiddleware() *middleware.CorsMiddleware {
	return c.corsMiddleware
}
