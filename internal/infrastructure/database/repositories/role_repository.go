package repositories

import (
	"context"
	"domain/delivery/models/entities"
	"domain/delivery/ports"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) ports.RolerRepository {
	return &roleRepository{
		db: db,
	}
}

// GetRoleByID obtiene un rol por su ID incluyendo sus permisos
func (r *roleRepository) GetRoleByID(ctx context.Context, id string) (*entities.Role, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName obtiene un rol por su nombre
func (r *roleRepository) GetRoleByName(ctx context.Context, name string) (*entities.Role, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// UpdateRole actualiza un rol
func (r *roleRepository) UpdateRole(ctx context.Context, role *entities.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

// DeleteRole elimina un rol (soft delete)
func (r *roleRepository) DeleteRole(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entities.Role{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// ListRoles lista todos los roles activos
func (r *roleRepository) ListRoles(ctx context.Context) ([]entities.Role, error) {
	var roles []entities.Role
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleByIDOrName obtiene un rol por su ID o nombre
func (r *roleRepository) GetRoleByIDOrName(ctx context.Context, param string) (*entities.Role, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).
		Preload("Permissions").
		First(&role, "id = ? OR name = ?", param, param).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// IsRoleExist verifica si un rol existe, ya sea por su ID o nombre
func (r *roleRepository) IsRoleExist(ctx context.Context, param string) (bool, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).
		First(&role, "id = ? OR name = ?", param, param).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeactivateRole desactiva un rol
func (r *roleRepository) DeactivateRole(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entities.Role{}).
		Where("id = ?", id).
		Update("is_active", false).Error
}

// IsRoleActive verifica si un rol está activo
func (r *roleRepository) IsRoleActive(ctx context.Context, param string) (bool, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).
		Select("is_active").
		First(&role, "id = ? OR name = ?", param, param).Error
	if err != nil {
		return false, err
	}
	return role.IsActive, nil
}

// CreatePermission crea un nuevo permiso
func (r *roleRepository) CreatePermission(ctx context.Context, permission *entities.Permission) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

// GetPermissionByID obtiene un permiso por su ID
func (r *roleRepository) GetPermissionByID(ctx context.Context, id string) (*entities.Permission, error) {
	var permission entities.Permission
	err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// UpdatePermission actualiza un permiso
func (r *roleRepository) UpdatePermission(ctx context.Context, permission *entities.Permission) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

// DeletePermission elimina un permiso
func (r *roleRepository) DeletePermission(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entities.Permission{}, "id = ?", id).Error
}

// ListPermissions lista todos los permisos
func (r *roleRepository) ListPermissions(ctx context.Context) ([]entities.Permission, error) {
	var permissions []entities.Permission
	err := r.db.WithContext(ctx).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// AssignPermissionToRole asigna un permiso a un rol
func (r *roleRepository) AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)",
		roleID, permissionID,
	).Error
}

// RemovePermissionFromRole remueve un permiso de un rol
func (r *roleRepository) RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?",
		roleID, permissionID,
	).Error
}

// GetRolePermissions obtiene todos los permisos de un rol
func (r *roleRepository) GetRolePermissions(ctx context.Context, roleID string) ([]entities.Permission, error) {
	var permissions []entities.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
