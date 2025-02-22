package auth

import (
	"application/ports"
	"context"
	"domain/delivery/models/auth"
)

type AuthUseCase struct {
	authService ports.AuthService
}

func NewAuthUseCase(authService ports.AuthService) *AuthUseCase {
	return &AuthUseCase{
		authService: authService,
	}
}

// TODO: Aqui irán algunas funciones que se encargarán de manejar la logica de negocio
// antes de llamar a los servicios correspondientes
// Por ejemplo, en el caso de  la autenticación, podríamos añadir lógica adicional
// antes o después de llamar al servicio de autenticación.
// Por poner un ejemplo puede ser eventos de dominio, metricas, etc etc xd

func (uc *AuthUseCase) Authenticate(ctx context.Context, credentials *auth.Credentials) (string, error) {
	// 1. Validar credenciales
	authUser, err := uc.authService.ValidateCredentials(ctx, credentials.Email, credentials.Password)
	if err != nil {
		return "", err
	}

	// 2. Crear sesion y obtener token
	token, err := uc.authService.CreateSession(ctx, authUser, credentials.DeviceInfo, credentials.IPAddress)
	if err != nil {
		return "", err
	}

	return token, nil
}

// TODO: Aqui irán algunas funciones que se encargarán de manejar la logica de negocio
// antes de llamar a los servicios correspondientes
// Por ejemplo, en el caso de  la autenticación, podríamos añadir lógica adicional
// antes o después de llamar al servicio de autenticación.
// Por poner un ejemplo puede ser eventos de dominio, metricas, etc etc xd
func (uc *AuthUseCase) SignOut(ctx context.Context, token string) error {
	return uc.authService.InvalidateSession(ctx, token)
}
