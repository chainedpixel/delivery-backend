package database

import (
	"fmt"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// RunMigrations ejecuta todas las migraciones de la base de datos
func RunMigrations(db *gorm.DB) error {
	return HandleCircularDependencies(db, func() error {
		return migrateAllEntities(db)
	})
}

// migrateAllEntities migra todas las entidades usando un enfoque por fases
// Fase 1 - Modelos base (usuarios, roles, permisos, zonas, empresas)
// Fase 2 - Modelos de conductores
// Fase 3 - Modelos de almacén
// Fase 4 - Modelos de órdenes
// Fase 5 - Modelos de inventario (dependientes de órdenes)
// Fase 6 - Modelos de notificaciones y eventos
func migrateAllEntities(db *gorm.DB) error {
	// PASO 1: Migrar primero las tablas base (sin relaciones complejas)
	logs.Info("FASE 1: Migrando modelos base...")
	baseModels := []schema.Tabler{
		// Modelos base de usuarios
		&entities.User{},
		&entities.Profile{},
		&entities.Role{},
		&entities.Permission{},
		&entities.RolePermission{},
		&entities.UserRole{},
		&entities.UserSession{},

		// Modelos base geográficos
		&entities.Zone{},
		&entities.Coverage{},
		&entities.AdjacentZone{},

		// Modelos base de empresas
		&entities.Company{},
		&entities.CompanyAddress{},
		&entities.Branch{},
		&entities.CompanyUser{},
	}

	if err := migrateModels(db, baseModels, "base"); err != nil {
		return err
	}

	// PASO 2: Migrar los modelos de conductores
	logs.Info("FASE 2: Migrando modelos de conductores...")
	driverModels := []schema.Tabler{
		&entities.Driver{},
		&entities.DriverZone{},
		&entities.Availability{},
	}

	if err := migrateModels(db, driverModels, "conductores"); err != nil {
		return err
	}

	// PASO 3: Migrar modelos de almacén
	logs.Info("FASE 3: Migrando modelos de almacén...")
	warehouseModels := []schema.Tabler{
		&entities.Warehouse{},
	}

	if err := migrateModels(db, warehouseModels, "almacén"); err != nil {
		return err
	}

	// PASO 4: Migrar modelos de órdenes
	logs.Info("FASE 4: Migrando modelos de órdenes...")
	orderModels := []schema.Tabler{
		&entities.Order{},
		&entities.Details{},
		&entities.PackageDetail{},
		&entities.DeliveryAddress{},
		&entities.PickupAddress{},
		&entities.Tracking{},
		&entities.QRCode{},
		&entities.StatusHistory{},
	}

	if err := migrateModels(db, orderModels, "órdenes"); err != nil {
		return err
	}

	// PASO 5: Migrar modelos de inventario que dependen de órdenes
	logs.Info("FASE 5: Migrando modelos de inventario relacionados con órdenes...")
	inventoryModels := []schema.Tabler{
		&entities.Inventory{},
		&entities.PackageTracking{},
	}

	if err := migrateModels(db, inventoryModels, "inventario"); err != nil {
		return err
	}

	// PASO 6: Migrar modelos de notificaciones y eventos
	logs.Info("FASE 6: Migrando modelos de notificaciones y eventos...")
	notificationModels := []schema.Tabler{
		&entities.Notification{},
		&entities.NotificationTemplate{},
		&entities.NotificationDevice{},
		&entities.NotificationPreference{},
		&entities.AuditLog{},
		&entities.SystemEvent{},
		&entities.EventLog{},
	}

	return migrateModels(db, notificationModels, "notificaciones")
}

// Función auxiliar para migrar un conjunto de modelos
func migrateModels(db *gorm.DB, models []schema.Tabler, phase string) error {
	for i, model := range models {
		tableName := model.TableName()

		logs.Info(fmt.Sprintf("Migrando modelo %s #%d: %s", phase, i, tableName))

		if err := db.AutoMigrate(model); err != nil {
			logs.Error("Error al migrar modelo", map[string]interface{}{
				"fase":   phase,
				"modelo": tableName,
				"índice": i,
				"error":  err.Error(),
			})
			return err
		}

		logs.Info(fmt.Sprintf("Migración exitosa del modelo %s #%d: %s", phase, i, tableName))
	}

	logs.Info(fmt.Sprintf("Completada la migración de modelos de %s", phase))
	return nil
}

// HandleCircularDependencies maneja las dependencias circulares durante la migración
// desactivando temporalmente las comprobaciones de claves foráneas
func HandleCircularDependencies(db *gorm.DB, migrationFunc func() error) error {
	// PASO 1: Desactivar temporalmente la comprobación de claves foráneas
	logs.Info("Desactivando temporalmente la comprobación de claves foráneas...")

	// Para MySQL/MariaDB
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		logs.Warn("No se pudo desactivar la comprobación de claves foráneas en MySQL, intentando con PostgreSQL...")

		// Para PostgreSQL
		if err := db.Exec("SET CONSTRAINTS ALL DEFERRED").Error; err != nil {
			logs.Warn("No se pudo desactivar la comprobación en PostgreSQL, intentando con SQLite...")

			// Para SQLite
			if err := db.Exec("PRAGMA foreign_keys = OFF").Error; err != nil {
				logs.Warn("No se pudo desactivar la comprobación de claves foráneas", map[string]interface{}{
					"error": err.Error(),
				})
				// Seguimos adelante de todas formas, ya que algunos dialectos no soportan esta operación
			}
		}
	}

	// PASO 2: Ejecutar las migraciones
	logs.Info("Ejecutando migraciones con comprobación de claves foráneas desactivada...")
	err := migrationFunc()

	// PASO 3: Reactivar la comprobación de claves foráneas (independientemente del resultado)
	logs.Info("Reactivando la comprobación de claves foráneas...")

	// Para MySQL/MariaDB
	if dbErr := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; dbErr != nil {
		// Para PostgreSQL (no es necesario reactivar en PostgreSQL, se reactiva automáticamente)

		// Para SQLite
		if dbErr := db.Exec("PRAGMA foreign_keys = ON").Error; dbErr != nil {
			logs.Warn("No se pudo reactivar la comprobación de claves foráneas", map[string]interface{}{
				"error": dbErr.Error(),
			})
		}
	}

	return err
}
