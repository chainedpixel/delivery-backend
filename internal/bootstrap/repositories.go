package bootstrap

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/database/repositories"
	"gorm.io/gorm"
)

type RepositoryContainer struct {
	db *gorm.DB

	roleRepo    ports.RolerRepository
	userRepo    ports.UserRepository
	orderRepo   ports.OrdererRepository
	companyRepo ports.CompanyRepository
	metricsRepo ports.MetricsRepository
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		db: db,
	}
}

func (c *RepositoryContainer) Initialize() error {
	c.roleRepo = repositories.NewRoleRepository(c.db)
	c.userRepo = repositories.NewUserRepository(c.db)
	c.orderRepo = repositories.NewOrderRepository(c.db)
	c.companyRepo = repositories.NewCompanyRepository(c.db)
	c.metricsRepo = repositories.NewMetricsRepository(c.db)

	return nil
}

func (c *RepositoryContainer) GetMetricsRepository() ports.MetricsRepository {
	return c.metricsRepo
}

func (c *RepositoryContainer) GetRoleRepository() ports.RolerRepository {
	return c.roleRepo
}

func (c *RepositoryContainer) GetUserRepository() ports.UserRepository {
	return c.userRepo
}

func (c *RepositoryContainer) GetOrderRepository() ports.OrdererRepository {
	return c.orderRepo
}

func (c *RepositoryContainer) GetCompanyRepository() ports.CompanyRepository {
	return c.companyRepo
}
