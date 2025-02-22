package middleware

import (
	"application/ports"
	"context"
	"infrastructure/api/responser"
	errPackage "infrastructure/error"
	"net/http"
	"shared/logs"
	"strings"
)

type AuthMiddleware struct {
	tokenService ports.TokenService
	respWriter   *responser.ResponseWriter
}

func NewAuthMiddleware(tokenService ports.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		respWriter:   responser.NewResponseWriter(),
	}
}

// Handle del middleware permite validar el token de autorización.
// Si el token es válido, se agrega al contexto de la petición.
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logs.Warn("Missing Authorization header", map[string]interface{}{
				"path":   r.URL.Path,
				"method": r.Method,
			})
			m.respWriter.Error(w, http.StatusUnauthorized, errPackage.ErrAuthorizationHeaderNotFound.Error(), nil)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logs.Warn("Invalid Authorization header format", map[string]interface{}{
				"header": authHeader,
				"path":   r.URL.Path,
				"method": r.Method,
			})
			m.respWriter.Error(w, http.StatusUnauthorized, errPackage.ErrInvalidAuthorizationFormat.Error(), nil)
			return
		}

		claims, err := m.tokenService.ValidateToken(parts[1])
		if err != nil {
			logs.Warn("Invalid token", map[string]interface{}{
				"error":  err.Error(),
				"path":   r.URL.Path,
				"method": r.Method,
			})
			m.respWriter.Error(w, http.StatusUnauthorized, errPackage.ErrTokenExpiredOrTampered.Error(), []string{err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
