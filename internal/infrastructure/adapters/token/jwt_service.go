package token

import (
	"application/ports"
	"domain/delivery/models/auth"
	errPackage "domain/error"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"shared/logs"
	"time"
)

type JWTService struct {
	secretKey    string
	tokenTTL     time.Duration
	cacheService ports.CacheService
}

func NewJWTService(secretKey string, cache ports.CacheService) *JWTService {
	return &JWTService{
		secretKey:    secretKey,
		tokenTTL:     15 * 24 * time.Hour,
		cacheService: cache,
	}
}

// GenerateToken genera un token JWT para un usuario
func (s *JWTService) GenerateToken(claims *auth.AuthClaims) (string, error) {
	now := time.Now()
	exp := now.Add(s.tokenTTL)

	claims.ExpiresAt = exp

	// Creamos los claims del JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  claims.UserID,
		"role": claims.Role,
		"exp":  exp.Unix(),
		"iat":  now.Unix(),
	})

	// Firmamos el token
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		logs.Error("Failed to sign token", map[string]interface{}{
			"error": err.Error(),
		})
		return "", errPackage.ErrFailedToSignToken
	}

	// Guardamos el token en caché
	key := "token:" + signedToken
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		logs.Error("Failed to marshal claims", map[string]interface{}{
			"error": err.Error(),
		})
		return "", errPackage.ErrFailedToMarshalClaims
	}

	if err := s.cacheService.Set(key, jsonClaims, s.tokenTTL); err != nil {
		logs.Error("Failed to store token in cache", map[string]interface{}{
			"error": err.Error(),
		})
		return "", err
	}

	logs.Info("Token generated successfully", map[string]interface{}{
		"user_id": claims.UserID,
	})

	return signedToken, nil
}

// ValidateToken valida un token JWT y retorna los claims si es válido
func (s *JWTService) ValidateToken(tokenString string) (*auth.AuthClaims, error) {
	// Primero verificamos en caché
	key := "token:" + tokenString
	cachedClaims, err := s.cacheService.Get(key)
	if err != nil {
		logs.Error("Token not found in cache", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var userClaims auth.AuthClaims
	if err := json.Unmarshal([]byte(cachedClaims), &userClaims); err != nil {
		logs.Error("Failed to unmarshal claims", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errPackage.ErrFailedToUnmarshalClaims
	}

	// Validamos el token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", errPackage.ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		logs.Error("Invalid token", map[string]interface{}{
			"error": err,
		})
		return nil, errPackage.ErrInvalidToken
	}

	if time.Now().After(userClaims.ExpiresAt) {
		return nil, errPackage.ErrTokenExpired
	}

	logs.Info("Token validated successfully", map[string]interface{}{
		"user_id": userClaims.UserID,
	})

	return &userClaims, nil
}

// RevokeToken revoca un token eliminándolo de la caché
func (s *JWTService) RevokeToken(token string) error {
	key := "token:" + token
	if err := s.cacheService.Delete(key); err != nil {
		logs.Error("Failed to revoke token", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	logs.Info("Token revoked successfully")
	return nil
}

// GetTokenTTL retorna el tiempo de vida del token
func (s *JWTService) GetTokenTTL() time.Duration {
	return s.tokenTTL
}
