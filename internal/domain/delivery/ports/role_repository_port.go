package ports

import (
	"context"
	"domain/delivery/models/auth"
)

// RoleRepository define las operaciones disponibles para la gestión de roles y permisos
type RoleRepository interface {
	// Operaciones de Roles
	CreateRole(ctx context.Context, role *auth.Role) error
	GetRoleByID(ctx context.Context, id string) (*auth.Role, error)
	GetRoleByName(ctx context.Context, name string) (*auth.Role, error)
	UpdateRole(ctx context.Context, role *auth.Role) error
	DeleteRole(ctx context.Context, id string) error
	ListRoles(ctx context.Context) ([]auth.Role, error)

	// Operaciones de Permisos
	CreatePermission(ctx context.Context, permission *auth.Permission) error
	GetPermissionByID(ctx context.Context, id string) (*auth.Permission, error)
	UpdatePermission(ctx context.Context, permission *auth.Permission) error
	DeletePermission(ctx context.Context, id string) error
	ListPermissions(ctx context.Context) ([]auth.Permission, error)

	// Operaciones de Asignación
	AssignPermissionToRole(ctx context.Context, roleID string, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID string, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]auth.Permission, error)
}
