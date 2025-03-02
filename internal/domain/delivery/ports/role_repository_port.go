package ports

import (
	"context"
	"domain/delivery/models/roles"
)

// RoleRepository define las operaciones disponibles para la gestión de roles y permisos
type RoleRepository interface {
	// Operaciones de Roles
	CreateRole(ctx context.Context, role *roles.Role) error
	GetRoleByID(ctx context.Context, id string) (*roles.Role, error)
	GetRoleByName(ctx context.Context, name string) (*roles.Role, error)
	UpdateRole(ctx context.Context, role *roles.Role) error
	DeleteRole(ctx context.Context, id string) error
	ListRoles(ctx context.Context) ([]roles.Role, error)

	// Operaciones de Permisos
	CreatePermission(ctx context.Context, permission *roles.Permission) error
	GetPermissionByID(ctx context.Context, id string) (*roles.Permission, error)
	UpdatePermission(ctx context.Context, permission *roles.Permission) error
	DeletePermission(ctx context.Context, id string) error
	ListPermissions(ctx context.Context) ([]roles.Permission, error)

	// Operaciones de Asignación
	AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]roles.Permission, error)
}
