package auth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"

	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	domainPorts "github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	errPackage "github.com/MarlonG1/delivery-backend/internal/infrastructure/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

type authService struct {
	userRepo     domainPorts.UserRepository
	tokenService ports.TokenProvider
}

func NewAuthService(userRepo domainPorts.UserRepository, tokenService ports.TokenProvider) ports.Authenticator {
	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *authService) CreateSession(ctx context.Context, authUser *entities.User, deviceInfo map[string]interface{}, ipAddress string) (string, error) {

	// 1. Obtener el rol principal del usuario
	roles, err := s.userRepo.GetUserRoles(ctx, authUser.ID)
	if err != nil {
		logs.Error("Failed to get users roles", map[string]interface{}{
			"error": err.Error(),
		})
		return "", errPackage.NewGeneralServiceError("Authenticator", "CreateSession", err)
	}
	var roleName string
	if len(roles) > 0 {
		roleName = roles[0].Name
	}

	// 2. Generar token
	claims := &auth.AuthClaims{
		UserID: authUser.ID,
		Role:   roleName,
	}

	token, err := s.tokenService.GenerateToken(claims)
	if err != nil {
		logs.Error("Failed to generate token", map[string]interface{}{
			"error": err.Error(),
		})
		return "", errPackage.NewGeneralServiceError("Authenticator", "CreateSession", err)
	}

	// 3. Crear sesion en base de datos
	deviceInfoJSON, _ := json.Marshal(deviceInfo)
	session := &entities.UserSession{
		ID:           uuid.NewString(),
		UserID:       authUser.ID,
		Token:        token,
		DeviceInfo:   string(deviceInfoJSON),
		IPAddress:    ipAddress,
		ExpiresAt:    time.Now().Add(s.tokenService.GetTokenTTL()),
		LastActivity: time.Now(),
	}
	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		// Si falla la creacion de la sesion, se revoca el token
		_ = s.tokenService.RevokeToken(token)
		logs.Error("Failed to create session", map[string]interface{}{
			"error": err.Error(),
		})
		return "", errPackage.NewGeneralServiceError("Authenticator", "CreateSession", err)
	}

	logs.Info("User logged in successfully", map[string]interface{}{
		"email": authUser.Email,
		"role":  roleName,
	})

	return token, nil
}

func (s *authService) ValidateCredentials(ctx context.Context, email, password string) (*entities.User, error) {
	// 1. Buscar usuario por email
	authUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		logs.Error("Failed to get user by email", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errPackage.NewGeneralServiceError("Authenticator", "ValidateCredentials", errPackage.ErrInvalidCredentials)
	}

	// 2. Verificar si el usuario esta activo
	if !authUser.IsActive {
		logs.Error("User is inactive", map[string]interface{}{
			"email": email,
		})
		return nil, errPackage.NewGeneralServiceError("Authenticator", "ValidateCredentials", errPackage.ErrInactiveUser)
	}

	// 3. Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(authUser.PasswordHash), []byte(password)); err != nil {
		logs.Error("Failed to compare password", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errPackage.NewGeneralServiceError("Authenticator", "ValidateCredentials", errPackage.ErrInvalidCredentials)
	}

	return authUser, nil
}

func (s *authService) InvalidateSession(ctx context.Context, token string) error {
	// 1. Buscar la sesion
	session, err := s.userRepo.GetSessionByToken(ctx, token)
	if err != nil {
		logs.Error("Failed to get session by token", map[string]interface{}{
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError("Authenticator", "InvalidateSession", errPackage.ErrSessionNotFound)
	}

	//2. Eliminar la sesion de la cache
	if err = s.tokenService.RevokeToken(token); err != nil {
		logs.Error("Failed to revoke token", map[string]interface{}{
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError("Authenticator", "InvalidateSession", err)
	}

	// 3. Eliminar la sesion de la base de datos
	if err := s.userRepo.DeleteSession(ctx, session.ID); err != nil {
		logs.Error("Failed to delete session", map[string]interface{}{
			"error": err.Error(),
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errPackage.NewGeneralServiceError("Authenticator", "InvalidateSession", errPackage.ErrSessionDBNotFound)
		}
	}

	logs.Info("User logged out successfully", map[string]interface{}{
		"token": token,
	})

	return nil
}
