package database

import (
	"fmt"
	"github.com/MarlonG1/delivery-backend/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDriver struct {
	Config *config.EnvConfig
}

func NewPostgresDriver(config *config.EnvConfig) *PostgresDriver {
	return &PostgresDriver{Config: config}
}

func (p *PostgresDriver) GetDSN() gorm.Dialector {
	return postgres.Open(p.GetStringConnection())
}

func (p *PostgresDriver) GetStringConnection() string {
	return fmt.Sprintf("host=%s users=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/El_Salvador",
		p.Config.Database.Host,
		p.Config.Database.User,
		p.Config.Database.Password,
		p.Config.Database.Name,
		p.Config.Database.Port)
}

func (p *PostgresDriver) GetHost() string {
	return p.Config.Database.Host
}

func (p *PostgresDriver) GetDriverName() string {
	return "PostgreSQL"
}
