package config

import (
	errPackage "config/error"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
	"strings"
)

type EnvConfig struct {
	Server struct {
		Port      string
		JWTSecret string
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

	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "..")

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(projectRoot)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, errPackage.ErrEnvFileNotFound
		}
		return nil, err
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
		return fmt.Errorf("DB_HOST is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("DB_USERNAME is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}

func MapEnvKeys(v *viper.Viper) {
	// .env keys for database connection
	v.Set("database.host", v.GetString("db_host"))
	v.Set("database.port", v.GetString("db_port"))
	v.Set("database.name", v.GetString("db_name"))
	v.Set("database.user", v.GetString("db_user"))
	v.Set("database.password", v.GetString("db_password"))
	v.Set("database.charset", v.GetString("db_charset"))
	v.Set("database.driver", v.GetString("db_driver"))

	// .env keys for redis connection
	v.Set("redis.host", v.GetString("redis_host"))
	v.Set("redis.port", v.GetString("redis_port"))
	v.Set("redis.password", v.GetString("redis_password"))

	// .env keys for server configuration
	v.Set("server.port", v.GetString("server_port"))
	v.Set("server.jwtSecret", v.GetString("jwt_secret"))

	// .env keys for log configuration
	v.Set("log.level", v.GetString("log_level"))
	v.Set("log.fileLogging", v.GetString("log_file_logging"))

}
