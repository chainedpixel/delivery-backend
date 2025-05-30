package middleware

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"

	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/responser"
	errPackage "github.com/MarlonG1/delivery-backend/internal/infrastructure/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

type AuthMiddleware struct {
	tokenService ports.TokenProvider
	respWriter   *responser.ResponseWriter
}

func NewAuthMiddleware(tokenService ports.TokenProvider) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		respWriter:   responser.NewResponseWriter(),
	}
}

// Handle del middleware permite validar el token de autorización.
// Si el token es válido, se agrega al contexto de la petición.
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// Verificar si es una conexión WebSocket
		if websocket.IsWebSocketUpgrade(r) {
			// Para WebSocket, obtener el token del query string
			tokenString = r.URL.Query().Get("token")
			if tokenString == "" {
				logs.Warn("Missing token in query string for WebSocket connection", map[string]interface{}{
					"path":   r.URL.Path,
					"method": r.Method,
				})
				m.respWriter.Error(w, http.StatusUnauthorized, errPackage.ErrAuthorizationHeaderNotFound.Error(), nil)
				return
			}
		} else {
			// Para peticiones HTTP normales, obtener el token del header
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
			tokenString = parts[1]
		}

		// Validar el token y obtener los claims
		claims, err := m.tokenService.ValidateToken(tokenString)
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
