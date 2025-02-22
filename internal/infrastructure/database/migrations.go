package database

import (
	"domain/delivery/models/auth"
	"domain/delivery/models/user"
	"fmt"
	"gorm.io/gorm"
	"shared/logs"
)

// RunMigrations ejecuta todas las migraciones de la base de datos
func RunMigrations(db *gorm.DB) error {
	modelsToMigrate := []interface{}{
		&user.User{},
		&user.UserProfile{},
		&user.UserRole{},
		&user.UserSession{},
		&auth.Role{},
		&auth.Permission{},
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
