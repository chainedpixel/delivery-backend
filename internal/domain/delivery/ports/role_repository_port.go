package ports

import (
	"context"
	"domain/delivery/models/role"
)

// RoleRepository define las operaciones disponibles para la gestión de roles y permisos
type RoleRepository interface {
	// Operaciones de Roles
	CreateRole(ctx context.Context, role *role.Role) error
	GetRoleByID(ctx context.Context, id string) (*role.Role, error)
	GetRoleByName(ctx context.Context, name string) (*role.Role, error)
	UpdateRole(ctx context.Context, role *role.Role) error
	DeleteRole(ctx context.Context, id string) error
	ListRoles(ctx context.Context) ([]role.Role, error)

	// Operaciones de Permisos
	CreatePermission(ctx context.Context, permission *role.Permission) error
	GetPermissionByID(ctx context.Context, id string) (*role.Permission, error)
	UpdatePermission(ctx context.Context, permission *role.Permission) error
	DeletePermission(ctx context.Context, id string) error
	ListPermissions(ctx context.Context) ([]role.Permission, error)

	// Operaciones de Asignación
	AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]role.Permission, error)
}
