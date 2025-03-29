package ports

import (
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cacher interface {
	Set(key string, claims []byte, ttl time.Duration) error // Set guarda un token en el cache
	Get(key string) (string, error)                         // Get obtiene un token del cache
	Delete(token string) error                              // Delete elimina un token del cache
	GetRedisClient() *redis.Client                          // GetRedisClient retorna el cliente de Redis
	CacherListService
}

// CacherListService define el comportamiento para la gestión de listas en caché
type CacherListService interface {
	RPush(key string, value []byte) error                   // Añade elemento al final de la lista
	LPush(key string, value []byte) error                   // Añade elemento al inicio de la lista
	LRange(key string, start, stop int64) ([]string, error) // Obtiene rango de elementos de la lista
	LLen(key string) (int64, error)                         // Obtiene longitud de la lista
	LTrim(key string, start, stop int64) error              // Mantiene solo el rango especificado
}

type TokenProvider interface {
	GenerateToken(claims *auth.AuthClaims) (string, error)
	ValidateToken(token string) (*auth.AuthClaims, error)
	RevokeToken(token string) error
	GetTokenTTL() time.Duration
}
