package database

import (
	"config"
	"fmt"
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/El_Salvador",
		p.Config.Database.Host,
		p.Config.Database.User,
		p.Config.Database.Password,
		p.Config.Database.Name,
		p.Config.Database.Port)

	return postgres.Open(dsn)
}

func (p *PostgresDriver) GetHost() string {
	return p.Config.Database.Host
}

func (p *PostgresDriver) GetDriverName() string {
	return "PostgreSQL"
}
