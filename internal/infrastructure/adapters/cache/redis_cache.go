package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"

	"github.com/MarlonG1/delivery-backend/configs"
	errPackage "github.com/MarlonG1/delivery-backend/internal/infrastructure/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

type RedisTokenCache struct {
	client *redis.Client
	config *config.EnvConfig
	ctx    context.Context
}

// NewRedisTokenCache crea una nueva instancia de RedisTokenCache
func NewRedisTokenCache(redisConf *config.RedisConfig) (*RedisTokenCache, error) {
	opt, err := redis.ParseURL(redisConf.GetURL())
	if err != nil {
		logs.Error("Failed to parse Redis URL", map[string]interface{}{
			"url":   redisConf.GetURL(),
			"error": err.Error(),
		})
		return nil, errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"NewRedisTokenCache",
			errPackage.ErrFailedToConnectRedis,
		)
	}

	client := redis.NewClient(opt)

	// Verificar conexión
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		logs.Error("Failed to connect to Redis", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"NewRedisTokenCache",
			errPackage.ErrFailedToPingRedis,
		)
	}

	logs.Info("Successfully connected to Redis")
	return &RedisTokenCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set guarda un token en Redis con un tiempo de vida determinado
func (c *RedisTokenCache) Set(key string, saveInfo []byte, ttl time.Duration) error {
	err := c.client.Set(c.ctx, key, saveInfo, ttl).Err()
	if err != nil {
		logs.Error("Failed to set value in Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"Set",
			errPackage.ErrFailedToSetKeyRedis,
		)
	}

	logs.Info("Value set successfully in Redis", map[string]interface{}{
		"key": key,
		"ttl": ttl.Seconds(),
	})
	return nil
}

// Get obtiene un token de Redis y lo convierte en un AuthClaims
func (c *RedisTokenCache) Get(key string) (string, error) {
	cacheInfo, err := c.client.Get(c.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		logs.Error("Token not found in Redis", map[string]interface{}{
			"key": key,
		})
		return "", errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"Get",
			errPackage.ErrTokenNotFound,
		)
	}
	if err != nil {
		logs.Error("Failed to get value from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return "", errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"Get",
			errPackage.ErrFailedToGetKey,
		)
	}

	logs.Info("Value retrieved successfully from Redis")
	return cacheInfo, nil
}

// Delete elimina un token de Redis
func (c *RedisTokenCache) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		logs.Error("Failed to delete token from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"Delete",
			errPackage.ErrFailedToDeleteKey,
		)
	}

	logs.Info("Token deleted successfully from Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}

// Close cierra la conexión con Redis
func (c *RedisTokenCache) Close() error {
	err := c.client.Close()
	if err != nil {
		logs.Error("Failed to close Redis client", map[string]interface{}{
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"Close",
			errPackage.ErrFailedToCloseRedis,
		)
	}

	logs.Info("Redis client closed successfully", map[string]interface{}{})
	return nil
}

func (c *RedisTokenCache) GetRedisClient() *redis.Client {
	return c.client
}

func (c *RedisTokenCache) RPush(key string, value []byte) error {
	err := c.client.RPush(c.ctx, key, value).Err()
	if err != nil {
		logs.Error("Failed to RPush value to Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"RPush",
			errPackage.ErrFailedRPush,
		)
	}

	logs.Info("Value RPushed successfully to Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}

func (c *RedisTokenCache) LPush(key string, value []byte) error {
	err := c.client.LPush(c.ctx, key, value).Err()
	if err != nil {
		logs.Error("Failed to LPush value to Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"LPush",
			errPackage.ErrFailedLPush,
		)
	}

	logs.Info("Value LPushed successfully to Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}

func (c *RedisTokenCache) LRange(key string, start, stop int64) ([]string, error) {
	result, err := c.client.LRange(c.ctx, key, start, stop).Result()
	if err != nil {
		logs.Error("Failed to get range from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return nil, errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"LRange",
			errPackage.ErrFailedLRange,
		)
	}

	return result, nil
}

func (c *RedisTokenCache) LLen(key string) (int64, error) {
	length, err := c.client.LLen(c.ctx, key).Result()
	if err != nil {
		logs.Error("Failed to get list length from Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return 0, errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"LLen",
			errPackage.ErrFailedLLen,
		)
	}

	return length, nil
}

func (c *RedisTokenCache) LTrim(key string, start, stop int64) error {
	err := c.client.LTrim(c.ctx, key, start, stop).Err()
	if err != nil {
		logs.Error("Failed to trim list in Redis", map[string]interface{}{
			"key":   key,
			"error": err.Error(),
		})
		return errPackage.NewGeneralServiceError(
			"RedisTokenCache",
			"LTrim",
			errPackage.ErrFailedLTrim,
		)
	}

	logs.Info("List trimmed successfully in Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}
