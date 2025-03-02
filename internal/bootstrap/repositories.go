package bootstrap

import (
	"domain/delivery/ports"
	"gorm.io/gorm"
	"infrastructure/database/repositories"
)

type RepositoryContainer struct {
	db *gorm.DB

	roleRepo ports.RoleRepository
	userRepo ports.UserRepository
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		db: db,
	}
}

func (c *RepositoryContainer) Initialize() error {
	c.roleRepo = repositories.NewRoleRepository(c.db)
	c.userRepo = repositories.NewUserRepository(c.db)

	return nil
}

func (c *RepositoryContainer) GetRoleRepository() ports.RoleRepository {
	return c.roleRepo
}

func (c *RepositoryContainer) GetUserRepository() ports.UserRepository {
	return c.userRepo
}
