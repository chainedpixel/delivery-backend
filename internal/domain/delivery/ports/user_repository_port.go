package ports

import (
	"context"
	"domain/delivery/models/role"
	"domain/delivery/models/user"
)

// UserRepository define las operaciones disponibles para la persistencia de usuarios y sus relaciones
type UserRepository interface {
	// Operaciones de Usuario
	Create(ctx context.Context, user *user.User) error
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id string) error

	// Operaciones de Perfil
	GetProfileByUserID(ctx context.Context, userID string) (*user.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *user.UserProfile) error

	// Operaciones de Sesión
	CreateSession(ctx context.Context, session *user.UserSession) error
	GetSessionByToken(ctx context.Context, token string) (*user.UserSession, error)
	GetActiveSessionsByUserID(ctx context.Context, userID string) ([]user.UserSession, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CleanExpiredSessions(ctx context.Context) error

	// Operaciones de Roles y Permisos
	AssignRoleToUser(ctx context.Context, userID string, roleID string, assignedBy string) error
	RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]role.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]role.Permission, error)

	// Operaciones de verificación
	MarkEmailAsVerified(ctx context.Context, userID string) error
	MarkPhoneAsVerified(ctx context.Context, userID string) error
}
