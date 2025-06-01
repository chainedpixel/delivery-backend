package bootstrap

import (
	"github.com/MarlonG1/delivery-backend/configs"
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	domainPorts "github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/services"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/adapters/auth"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/adapters/cache"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/adapters/email"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/adapters/token"
)

type ServiceContainer struct {
	repositories *RepositoryContainer
	config       *config.EnvConfig

	jwtService        ports.TokenProvider
	cacheService      ports.Cacher
	authService       ports.Authenticator
	userService       domainPorts.Userer
	orderService      domainPorts.Orderer
	companyService    domainPorts.Companyrer
	metricsService    domainPorts.MetricsService
	trackerService    domainPorts.OrderTracker
	emailService      ports.EmailService
	roleService       domainPorts.Roler
	simulationService *services.OrderSimulationService
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
	c.trackerService = services.NewTrackerService(c.repositories.GetTrackerRepository())
	c.emailService = email.NewEmailService(c.config)
	c.orderService = services.NewOrderService(c.repositories.GetOrderRepository(), c.trackerService, c.emailService, c.repositories.GetCompanyRepository())
	c.simulationService = services.NewOrderSimulationService(c.orderService, c.repositories.GetCompanyRepository(), c.repositories.GetUserRepository(), c.trackerService)
	c.metricsService = services.NewCompanyMetricsService(c.repositories.GetCompanyRepository(), c.repositories.GetMetricsRepository())
	c.companyService = services.NewCompanyService(c.repositories.GetCompanyRepository(), c.metricsService)
	c.roleService = services.NewRoleService(c.repositories.GetRoleRepository())

	return nil
}

func (c *ServiceContainer) GetSimulationService() *services.OrderSimulationService {
	return c.simulationService
}

func (c *ServiceContainer) GetEmailService() ports.EmailService {
	return c.emailService
}

func (c *ServiceContainer) GetMetricsService() domainPorts.MetricsService {
	return c.metricsService
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

func (c *ServiceContainer) GetRoleService() domainPorts.Roler {
	return c.roleService
}

func (c *ServiceContainer) GetTrackerService() domainPorts.OrderTracker {
	return c.trackerService
}
