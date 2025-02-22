package config

import (
	errPackage "config/error"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type EnvConfig struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		Charset  string
		Driver   string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	Log struct {
		Level       string
		FileLogging bool
	}
}

func NewEnvConfig() (*EnvConfig, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath("../")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, errPackage.ErrEnvFileNotFound
		}
	}

	v.AutomaticEnv()
	MapEnvKeys(v)

	var config EnvConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, errPackage.ErrFailedToLoadEnv
	}

	return &config, validateConfig(&config)
}

func validateConfig(config *EnvConfig) error {
	if config.Database.Host == "" {
		return fmt.Errorf("DB_HOST es requerido")
	}
	if config.Database.User == "" {
		return fmt.Errorf("DB_USERNAME es requerido")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD es requerido")
	}
	return nil
}

func MapEnvKeys(v *viper.Viper) {
	v.BindEnv("server.port", "SERVER_PORT")
	v.BindEnv("database.host", "DB_HOST")
	v.BindEnv("database.port", "DB_PORT")
	v.BindEnv("database.name", "DB_DATABASE")
	v.BindEnv("database.user", "DB_USERNAME")
	v.BindEnv("database.driver", "DB_DRIVER")
	v.BindEnv("database.password", "DB_PASSWORD")
	v.BindEnv("database.charset", "DB_CHARSET")
	v.BindEnv("redis.host", "REDIS_HOST")
	v.BindEnv("redis.port", "REDIS_PORT")
	v.BindEnv("redis.password", "REDIS_PASSWORD")
	v.BindEnv("log.level", "LOG_LEVEL")
	v.BindEnv("log.fileLogging", "LOG_FILE_LOGGING")
}
