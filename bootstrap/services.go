package bootstrap

import (
	"application/ports"
	"config"
	"infrastructure/adapters/auth"
	"infrastructure/adapters/cache"
	"infrastructure/adapters/token"
)

type ServiceContainer struct {
	repositories *RepositoryContainer
	config       *config.EnvConfig

	jwtService   ports.TokenService
	cacheService ports.CacheService
	authService  ports.AuthService
}

func NewServiceContainer(repositories *RepositoryContainer) *ServiceContainer {
	return &ServiceContainer{

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

	return nil
}

func (c *ServiceContainer) GetTokenService() ports.TokenService {
	return c.jwtService
}

func (c *ServiceContainer) GetCacheService() ports.CacheService {
	return c.cacheService
}

func (c *ServiceContainer) GetAuthService() ports.AuthService {
	return c.authService
}
