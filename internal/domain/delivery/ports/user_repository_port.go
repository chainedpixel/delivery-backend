package ports

import (
	"context"
	"domain/delivery/models/entities"
)

// UserRepository define las operaciones disponibles para la persistencia de usuarios y sus relaciones
type UserRepository interface {
	// Operaciones de Usuario
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Update(ctx context.Context, id string, user *entities.User) error
	Delete(ctx context.Context, id string) error
	ActivateOrDeactivate(ctx context.Context, id string, active bool) error

	// Operaciones de Perfil
	GetProfileByUserID(ctx context.Context, userID string) (*entities.Profile, error)
	UpdateProfile(ctx context.Context, profile *entities.Profile) error
	Recover(ctx context.Context, id string) error

	// Operaciones de Verificación
	IsUserDeleted(ctx context.Context, userID string) (bool, error)
	IsUserActive(ctx context.Context, userID string) (bool, error)

	// Operaciones de Sesión
	CreateSession(ctx context.Context, session *entities.UserSession) error
	GetSessionByToken(ctx context.Context, token string) (*entities.UserSession, error)
	GetActiveSessionsByUserID(ctx context.Context, userID string) ([]entities.UserSession, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CleanExpiredSessions(ctx context.Context) error

	// Operaciones de Roles y Permisos
	AssignRoleToUser(ctx context.Context, userID string, roleID string, assignedBy string) error
	RemoveRoleFromUser(ctx context.Context, userID string, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]entities.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]entities.Permission, error)

	// Operaciones de verificación
	MarkEmailAsVerified(ctx context.Context, userID string) error
	MarkPhoneAsVerified(ctx context.Context, userID string) error
}
