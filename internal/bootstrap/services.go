package bootstrap

import (
	"application/ports"
	"config"
	domainPorts "domain/delivery/interfaces"
	"domain/delivery/services"
	"infrastructure/adapters/auth"
	"infrastructure/adapters/cache"
	"infrastructure/adapters/token"
)

type ServiceContainer struct {
	repositories *RepositoryContainer
	config       *config.EnvConfig

	jwtService     ports.TokenProvider
	cacheService   ports.Cacher
	authService    ports.Authenticator
	userService    domainPorts.Userer
	orderService   domainPorts.Orderer
	companyService domainPorts.Companyrer
}

func NewServiceContainer(repositories *RepositoryContainer, config *config.EnvConfig) *ServiceContainer {
	return &ServiceContainer{
		config:       config,
		repositories: repositories,
	}
}

func (c *ServiceContainer) Initialize() error {
	var err error

	c.cacheService, err = cache.NewRedisTokenCache(config.NewRedisConfig(c.config))
	if err != nil {
		return err
	}

	c.jwtService = token.NewJWTService(c.config.Server.JWTSecret, c.cacheService)
	c.authService = auth.NewAuthService(c.repositories.GetUserRepository(), c.jwtService)
	c.userService = services.NewUserService(c.repositories.GetUserRepository())
	c.orderService = services.NewOrderService(c.repositories.GetOrderRepository())
	c.companyService = services.NewCompanyService(c.repositories.GetCompanyAddressRepository())

	return nil
}

func (c *ServiceContainer) GetTokenService() ports.TokenProvider {
	return c.jwtService
}

func (c *ServiceContainer) GetCacheService() ports.Cacher {
	return c.cacheService
}

func (c *ServiceContainer) GetAuthService() ports.Authenticator {
	return c.authService
}

func (c *ServiceContainer) GetUserService() domainPorts.Userer {
	return c.userService
}

func (c *ServiceContainer) GetOrderService() domainPorts.Orderer {
	return c.orderService
}

func (c *ServiceContainer) GetCompanyService() domainPorts.Companyrer {
	return c.companyService
}
