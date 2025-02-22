package config

import "fmt"

type RedisConfig struct {
	config *EnvConfig
}

func NewRedisConfig(config *EnvConfig) *RedisConfig {
	return &RedisConfig{}
}

func (c *RedisConfig) GetURL() string {
	if c.config.Redis.Password != "" {
		return fmt.Sprintf("redis://%s@%s:%s", c.config.Redis.Password, c.config.Redis.Host, c.config.Redis.Port)
	}
	return fmt.Sprintf("redis://%s:%s", c.config.Redis.Host, c.config.Redis.Password)
}
