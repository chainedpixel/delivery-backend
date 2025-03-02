package ports

import (
	"context"
	"domain/delivery/models/roles"
	"domain/delivery/models/users"
)

// UserRepository define las operaciones disponibles para la persistencia de usuarios y sus relaciones
type UserRepository interface {
	// Operaciones de Usuario
	Create(ctx context.Context, user *users.User) error
	GetByID(ctx context.Context, id string) (*users.User, error)
	GetByEmail(ctx context.Context, email string) (*users.User, error)
	Update(ctx context.Context, user *users.User) error
	Delete(ctx context.Context, id string) error

	// Operaciones de Perfil
	GetProfileByUserID(ctx context.Context, userID string) (*users.Profile, error)
	UpdateProfile(ctx context.Context, profile *users.Profile) error

	// Operaciones de Sesión
	CreateSession(ctx context.Context, session *users.UserSession) error
	GetSessionByToken(ctx context.Context, token string) (*users.UserSession, error)
	GetActiveSessionsByUserID(ctx context.Context, userID string) ([]users.UserSession, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CleanExpiredSessions(ctx context.Context) error

	// Operaciones de Roles y Permisos
	AssignRoleToUser(ctx context.Context, userID string, roleID string, assignedBy string) error
	RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]roles.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]roles.Permission, error)

	// Operaciones de verificación
	MarkEmailAsVerified(ctx context.Context, userID string) error
	MarkPhoneAsVerified(ctx context.Context, userID string) error
}
