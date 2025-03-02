package database

import (
	"domain/delivery/models/roles"
	"domain/delivery/models/users"
	"fmt"
	"gorm.io/gorm"
	"shared/logs"
)

// RunMigrations ejecuta todas las migraciones de la base de datos
func RunMigrations(db *gorm.DB) error {
	modelsToMigrate := []interface{}{
		&users.User{},
		&users.Profile{},
		&users.Role{},
		&users.UserSession{},
		&roles.Role{},
		&roles.Permission{},
	}

	logs.Info("Starting database migrations")

	for i, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			logs.Error("Failed to migrate model", map[string]interface{}{
				"model number": i,
				"error":        err.Error(),
			})
			return err
		}

		logs.Info(fmt.Sprintf("Successfully migrated model %d", i))
	}

	logs.Info("All migrations completed successfully")
	return nil
}
