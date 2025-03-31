package bootstrap

import (
	"github.com/MarlonG1/delivery-backend/configs"
	"gorm.io/gorm"
	"sync"
)

type ContainerDependency interface {
	Initialize() error
}

type Container struct {
	db     *gorm.DB
	config *config.EnvConfig

	repositories *RepositoryContainer
	services     *ServiceContainer
	useCases     *UseCaseContainer
	handlers     *HandlerContainer
	middleware   *MiddlewareContainer

	mu sync.RWMutex
}

func NewContainer(db *gorm.DB, config *config.EnvConfig) *Container {
	return &Container{
		db:     db,
		config: config,
	}
}

// Initialize Inicializa todos los contenedores de la aplicación en orden de dependencia
// El orden de inicialización es importante para evitar errores de dependencia
// 1 - Repositories, 2 - Services, 3 - UseCases, 4 - Middleware, 5 - Handlers
// Especificamente en ese orde
func (c *Container) Initialize() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.repositories = NewRepositoryContainer(c.db)
	if err := c.repositories.Initialize(); err != nil {
		return err
	}

	c.services = NewServiceContainer(c.repositories, c.config)
	if err := c.services.Initialize(); err != nil {
		return err
	}

	c.useCases = NewUseCaseContainer(c.services)
	if err := c.useCases.Initialize(); err != nil {
		return err
	}

	c.middleware = NewMiddlewareContainer(c.services)
	if err := c.middleware.Initialize(); err != nil {
		return err
	}

	c.handlers = NewHandlerContainer(c.useCases, c.services)
	if err := c.handlers.Initialize(); err != nil {
		return err
	}

	return nil
}

func (c *Container) GetRepositoryContainer() *RepositoryContainer {
	return c.repositories
}

func (c *Container) GetServiceContainer() *ServiceContainer {
	return c.services
}

func (c *Container) GetUseCaseContainer() *UseCaseContainer {
	return c.useCases
}

func (c *Container) GetHandlerContainer() *HandlerContainer {
	return c.handlers
}

func (c *Container) GetMiddlewareContainer() *MiddlewareContainer {
	return c.middleware
}
