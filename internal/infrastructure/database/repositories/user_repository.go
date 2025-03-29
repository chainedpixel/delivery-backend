package repositories

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	"gorm.io/gorm"
	"time"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create inserta un nuevo usuario y su perfil si existe
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) IsUserDeleted(ctx context.Context, userID string) (bool, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Select("deleted_at").
		Where("deleted_at IS NULL").
		First(&user, "id = ?", userID).Error
	if err != nil {
		return true, err
	}
	return false, nil
}

func (r *userRepository) IsUserActive(ctx context.Context, userID string) (bool, error) {
	var user entities.User
	err := r.db.WithContext(ctx).
		Select("is_active").
		First(&user, "id = ?", userID).Error
	if err != nil {
		return false, err
	}
	return user.IsActive, nil
}

// GetByID obtiene un usuario por ID incluyendo su perfil y roles activos
func (r *userRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var usr entities.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Roles", "is_active = ?", true).
		Preload("Roles.Role").
		Preload("Sessions", "expires_at > NOW()").
		First(&usr, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// GetByEmail obtiene un usuario por email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var usr entities.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Roles", "is_active = ?", true).
		Preload("Roles.Role").
		Preload("Sessions", "expires_at > NOW()").
		Where("email = ?", email).
		First(&usr).Error
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// Update actualiza la información del usuario
func (r *userRepository) Update(ctx context.Context, id string, user *entities.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Actualizar usuario
		if err := tx.Model(&entities.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
			return err
		}

		// 2. Actualizar perfil
		if user.Profile != nil {
			if err := tx.Model(&entities.Profile{}).Where("user_id = ?", id).Updates(user.Profile).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete realiza un soft delete del usuario
func (r *userRepository) Delete(ctx context.Context, id string) error {
	now := time.Now()

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.User{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"is_active":  false,
				"deleted_at": now,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) Recover(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.User{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"is_active":  true,
				"deleted_at": nil,
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetProfileByUserID obtiene el perfil de un usuario
func (r *userRepository) GetProfileByUserID(ctx context.Context, userID string) (*entities.Profile, error) {
	var profile entities.Profile
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateProfile actualiza el perfil del usuario
func (r *userRepository) UpdateProfile(ctx context.Context, profile *entities.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// CreateSession crea una nueva sesión
func (r *userRepository) CreateSession(ctx context.Context, session *entities.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

// GetSessionByToken obtiene una sesión por su token
func (r *userRepository) GetSessionByToken(ctx context.Context, token string) (*entities.UserSession, error) {
	var session entities.UserSession
	err := r.db.WithContext(ctx).
		Where("token = ? AND expires_at > NOW()", token).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetActiveSessionsByUserID obtiene todas las sesiones activas de un usuario
func (r *userRepository) GetActiveSessionsByUserID(ctx context.Context, userID string) ([]entities.UserSession, error) {
	var sessions []entities.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > NOW()", userID).
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// DeleteSession elimina una sesión específica
func (r *userRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).
		Delete(&entities.UserSession{}, "id = ?", sessionID).Error
}

// CleanExpiredSessions elimina todas las sesiones expiradas
func (r *userRepository) CleanExpiredSessions(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserSession{}).
		Where("user_id = ? AND expires_at < NOW()", id).
		Update("expires_at", time.Now()).Error
}

// AssignRoleToUser asigna un rol a un usuario
func (r *userRepository) AssignRoleToUser(ctx context.Context, userID string, roleID string, assignedBy string) error {
	userRole := entities.UserRole{
		UserID:     userID,
		RoleID:     roleID,
		AssignedBy: assignedBy,
		AssignedAt: time.Now(),
		IsActive:   true,
	}
	return r.db.WithContext(ctx).Create(&userRole).Error
}

// UpdateRolesToUser actualiza los roles de un usuario
func (r *userRepository) UpdateRolesToUser(ctx context.Context, userID string, loggedUserID string, roles []entities.Role) error {
	// Obtenemos TODOS los roles del usuario (tanto activos como inactivos)
	var allUserRoles []entities.UserRole
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&allUserRoles).Error; err != nil {
		return err
	}

	// Creamos mapas para procesamiento eficiente
	activeRoleIDs := make(map[string]bool)      // Roles actualmente activos
	inactiveRoleIDs := make(map[string]bool)    // Roles actualmente inactivos
	allExistingRoleIDs := make(map[string]bool) // Todos los roles (activos e inactivos)

	for _, userRole := range allUserRoles {
		allExistingRoleIDs[userRole.RoleID] = true
		if userRole.IsActive {
			activeRoleIDs[userRole.RoleID] = true
		} else {
			inactiveRoleIDs[userRole.RoleID] = true
		}
	}

	// Roles que queremos mantener según el array de entrada
	rolesToKeep := make(map[string]bool)
	for _, role := range roles {
		rolesToKeep[role.ID] = true
	}

	// Iniciamos una transacción
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Desactivar roles que ya no están en la lista
	for roleID := range activeRoleIDs {
		if !rolesToKeep[roleID] {
			// El rol ya no debe estar asignado, lo desactivamos
			if err := tx.Model(&entities.UserRole{}).
				Where("user_id = ? AND role_id = ?", userID, roleID).
				Updates(map[string]interface{}{"is_active": false}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 2. Reactivar roles que existían pero estaban inactivos
	for _, role := range roles {
		if inactiveRoleIDs[role.ID] {
			// Reactivamos el rol
			if err := tx.Model(&entities.UserRole{}).
				Where("user_id = ? AND role_id = ?", userID, role.ID).
				Updates(map[string]interface{}{
					"is_active":   true,
					"assigned_at": time.Now(),
					"assigned_by": loggedUserID,
				}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 3. Agregar roles completamente nuevos
	for _, role := range roles {
		// Solo si el rol no existe en absoluto (ni activo ni inactivo)
		if !allExistingRoleIDs[role.ID] {
			// Creamos un nuevo UserRole
			newUserRole := entities.UserRole{
				UserID:     userID,
				RoleID:     role.ID,
				AssignedAt: time.Now(),
				AssignedBy: loggedUserID,
				IsActive:   true,
				CreatedAt:  time.Now(),
			}

			// Insertamos el nuevo registro
			if err := tx.Create(&newUserRole).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Confirmamos la transacción
	return tx.Commit().Error
}

func (r *userRepository) UnassignRole(ctx context.Context, userID string, roleID string) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Update("is_active", false).Error
}

// ActivateOrDeactivate activa o desactiva un usuario
func (r *userRepository) ActivateOrDeactivate(ctx context.Context, id string, active bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.User{}).
			Where("id = ?", id).
			Update("is_active", active).Error; err != nil {
			return err
		}

		return nil
	})
}

// RemoveRoleFromUser remueve un rol de un usuario
func (r *userRepository) RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error {
	return r.db.WithContext(ctx).
		Model(&entities.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Update("is_active", false).Error
}

// GetAllUsersFromCompany obtiene todos los usuarios de una empresa
func (r *userRepository) GetAllUsersFromCompany(ctx context.Context, companyID string, params *entities.UserQueryParams) ([]entities.User, int64, error) {
	var users []entities.User
	query := r.db.WithContext(ctx).
		Preload("Profile").
		Preload("Roles", "is_active = ?", true).
		Preload("Roles.Role").
		Where("company_id = ?", companyID)

	if params.IncludeDeleted {
		query = query.Unscoped()
	} else {
		query = query.Where("deleted_at IS NULL")
	}

	if params.Status {
		query = query.Where("is_active = ?", params.Status)
	}

	if params.CreationDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *params.CreationDate, params.CreationDate.AddDate(0, 0, 1))
	}

	if params.Phone != "" {
		query = query.Where("phone = ?", params.Phone)
	}

	if params.Name != "" {
		query = query.Where("full_name LIKE ?", "%"+params.Name+"%")
	}

	if params.Email != "" {
		query = query.Where("email = ?", params.Email)
	}

	// Paginación
	var total int64
	if err := query.Model(&entities.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("Profile").
		Preload("Roles", "is_active = ?", true).
		Preload("Roles.Role").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserRoles obtiene todos los roles de un usuario
func (r *userRepository) GetUserRoles(ctx context.Context, userID string) ([]entities.Role, error) {
	var roles []entities.Role
	err := r.db.WithContext(ctx).
		Table("roles").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND user_roles.is_active = ?", userID, true).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetUserPermissions obtiene todos los permisos de un usuario a través de sus roles
func (r *userRepository) GetUserPermissions(ctx context.Context, userID string) ([]entities.Permission, error) {
	var permissions []entities.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND user_roles.is_active = ?", userID, true).
		Distinct().
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// MarkEmailAsVerified marca el email como verificado
func (r *userRepository) MarkEmailAsVerified(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", userID).
		Update("email_verified_at", r.db.NowFunc()).Error
}

// MarkPhoneAsVerified marca el teléfono como verificado
func (r *userRepository) MarkPhoneAsVerified(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", userID).
		Update("phone_verified_at", r.db.NowFunc()).Error
}
