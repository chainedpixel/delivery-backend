package middleware

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

type TokenExtractor struct{}

func NewTokenExtractor() *TokenExtractor {
	return &TokenExtractor{}
}

func (te *TokenExtractor) ExtractToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		if websocket.IsWebSocketUpgrade(r) {
			tokenString = r.URL.Query().Get("token")
		} else {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")
			tokenString = parts[1]
		}

		// Almacenar el token en el contexto
		ctx := context.WithValue(r.Context(), "userToken", tokenString)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
